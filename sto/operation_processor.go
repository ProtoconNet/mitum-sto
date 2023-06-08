package sto

import (
	"context"
	"fmt"
	"io"
	"sync"

	currencybase "github.com/ProtoconNet/mitum-currency/v3/base"
	currency "github.com/ProtoconNet/mitum-currency/v3/operation/currency"
	extensioncurrency "github.com/ProtoconNet/mitum-currency/v3/operation/extension"
	types "github.com/ProtoconNet/mitum-currency/v3/operation/type"
	"github.com/ProtoconNet/mitum-sto/kyc"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
	"github.com/ProtoconNet/mitum2/util/logging"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

var operationProcessorPool = sync.Pool{
	New: func() interface{} {
		return new(OperationProcessor)
	},
}

type GetLastBlockFunc func() (base.BlockMap, bool, error)

type DuplicationType string

const (
	DuplicationTypeSender      DuplicationType = "sender"
	DuplicationTypeCurrency    DuplicationType = "currency"
	DuplicationTypeContractSTO DuplicationType = "contract-sto"
)

type OperationProcessor struct {
	sync.RWMutex
	*logging.Logging
	*base.BaseOperationProcessor
	processorHintSet     *hint.CompatibleSet
	fee                  map[currencybase.CurrencyID]currencybase.Big
	duplicated           map[string]DuplicationType
	duplicatedNewAddress map[string]struct{}
	processorClosers     *sync.Map
	GetStateFunc         base.GetStateFunc
}

func NewOperationProcessor() *OperationProcessor {
	m := sync.Map{}
	return &OperationProcessor{
		Logging: logging.NewLogging(func(c zerolog.Context) zerolog.Context {
			return c.Str("module", "mitum-sto-operations-processor")
		}),
		processorHintSet:     hint.NewCompatibleSet(),
		fee:                  map[currencybase.CurrencyID]currencybase.Big{},
		duplicated:           map[string]DuplicationType{},
		duplicatedNewAddress: map[string]struct{}{},
		processorClosers:     &m,
	}
}

func (opr *OperationProcessor) New(
	height base.Height,
	getStateFunc base.GetStateFunc,
	newPreProcessConstraintFunc base.NewOperationProcessorProcessFunc,
	newProcessConstraintFunc base.NewOperationProcessorProcessFunc) (*OperationProcessor, error) {
	e := util.StringErrorFunc("failed to create new OperationProcessor")

	nopr := operationProcessorPool.Get().(*OperationProcessor)
	if nopr.processorHintSet == nil {
		nopr.processorHintSet = opr.processorHintSet
	}

	if nopr.fee == nil {
		nopr.fee = opr.fee
	}

	if nopr.duplicated == nil {
		nopr.duplicated = make(map[string]DuplicationType)
	}

	if nopr.duplicatedNewAddress == nil {
		nopr.duplicatedNewAddress = make(map[string]struct{})
	}

	if nopr.Logging == nil {
		nopr.Logging = opr.Logging
	}

	b, err := base.NewBaseOperationProcessor(
		height, getStateFunc, newPreProcessConstraintFunc, newProcessConstraintFunc)
	if err != nil {
		return nil, e(err, "")
	}

	nopr.BaseOperationProcessor = b
	nopr.GetStateFunc = getStateFunc
	return nopr, nil
}

func (opr *OperationProcessor) SetProcessor(
	hint hint.Hint,
	newProcessor types.GetNewProcessor,
) (base.OperationProcessor, error) {
	if err := opr.processorHintSet.Add(hint, newProcessor); err != nil {
		if !errors.Is(err, util.ErrFound) {
			return nil, err
		}
	}

	return opr, nil
}

func (opr *OperationProcessor) PreProcess(ctx context.Context, op base.Operation, getStateFunc base.GetStateFunc) (context.Context, base.OperationProcessReasonError, error) {
	e := util.StringErrorFunc("failed to preprocess for OperationProcessor")

	if opr.processorClosers == nil {
		opr.processorClosers = &sync.Map{}
	}

	var sp base.OperationProcessor
	switch i, known, err := opr.getNewProcessor(op); {
	case err != nil:
		return ctx, base.NewBaseOperationProcessReasonError(err.Error()), nil
	case !known:
		return ctx, nil, e(nil, "failed to getNewProcessor, %T", op)
	default:
		sp = i
	}

	switch _, reasonerr, err := sp.PreProcess(ctx, op, getStateFunc); {
	case err != nil:
		return ctx, nil, e(err, "")
	case reasonerr != nil:
		return ctx, reasonerr, nil
	}

	return ctx, nil, nil
}

func (opr *OperationProcessor) Process(ctx context.Context, op base.Operation, getStateFunc base.GetStateFunc) ([]base.StateMergeValue, base.OperationProcessReasonError, error) {
	e := util.StringErrorFunc("failed to process for OperationProcessor")

	if err := opr.checkDuplication(op); err != nil {
		return nil, base.NewBaseOperationProcessReasonError("duplication found: %w", err), nil
	}

	var sp base.OperationProcessor
	switch i, known, err := opr.getNewProcessor(op); {
	case err != nil:
		return nil, nil, e(err, "")
	case !known:
		return nil, nil, e(nil, "failed to getNewProcessor")
	default:
		sp = i
	}

	stateMergeValues, reasonerr, err := sp.Process(ctx, op, getStateFunc)

	return stateMergeValues, reasonerr, err
}

func (opr *OperationProcessor) checkDuplication(op base.Operation) error {
	opr.Lock()
	defer opr.Unlock()

	var did string
	var didtype DuplicationType
	var newAddresses []base.Address

	switch t := op.(type) {
	case currency.CreateAccounts:
		fact, ok := t.Fact().(currency.CreateAccountsFact)
		if !ok {
			return errors.Errorf("expected CreateAccountsFact, not %T", t.Fact())
		}
		as, err := fact.Targets()
		if err != nil {
			return errors.Errorf("failed to get Addresses")
		}
		newAddresses = as
		did = fact.Sender().String()
		didtype = DuplicationTypeSender
	case currency.KeyUpdater:
		fact, ok := t.Fact().(currency.KeyUpdaterFact)
		if !ok {
			return errors.Errorf("expected KeyUpdaterFact, not %T", t.Fact())
		}
		did = fact.Target().String()
		didtype = DuplicationTypeSender
	case currency.Transfers:
		fact, ok := t.Fact().(currency.TransfersFact)
		if !ok {
			return errors.Errorf("expected TransfersFact, not %T", t.Fact())
		}
		did = fact.Sender().String()
		didtype = DuplicationTypeSender
	case extensioncurrency.CreateContractAccounts:
		fact, ok := t.Fact().(extensioncurrency.CreateContractAccountsFact)
		if !ok {
			return errors.Errorf("expected CreateContractAccountsFact, not %T", t.Fact())
		}
		as, err := fact.Targets()
		if err != nil {
			return errors.Errorf("failed to get Addresses")
		}
		newAddresses = as
	case extensioncurrency.Withdraws:
		fact, ok := t.Fact().(extensioncurrency.WithdrawsFact)
		if !ok {
			return errors.Errorf("expected WithdrawsFact, not %T", t.Fact())
		}
		did = fact.Sender().String()
		didtype = DuplicationTypeSender
	case currency.CurrencyRegister:
		fact, ok := t.Fact().(currency.CurrencyRegisterFact)
		if !ok {
			return errors.Errorf("expected CurrencyRegisterFact, not %T", t.Fact())
		}
		did = fact.Currency().Currency().String()
		didtype = DuplicationTypeCurrency
	case currency.CurrencyPolicyUpdater:
		fact, ok := t.Fact().(currency.CurrencyPolicyUpdaterFact)
		if !ok {
			return errors.Errorf("expected CurrencyPolicyUpdaterFact, not %T", t.Fact())
		}
		did = fact.Currency().String()
		didtype = DuplicationTypeCurrency
	case currency.SuffrageInflation:
	case AuthorizeOperators:
		fact, ok := t.Fact().(AuthorizeOperatorsFact)
		if !ok {
			return errors.Errorf("expected AuthorizeOperatorsFact, not %T", t.Fact())
		}
		did = fact.Sender().String()
		didtype = DuplicationTypeSender
	case CreateSecurityTokens:
		fact, ok := t.Fact().(CreateSecurityTokensFact)
		if !ok {
			return errors.Errorf("expected CreateSecurityTokensFact, not %T", t.Fact())
		}
		did = fact.Sender().String()
		didtype = DuplicationTypeSender
	case IssueSecurityTokens:
		fact, ok := t.Fact().(IssueSecurityTokensFact)
		if !ok {
			return errors.Errorf("expected IssueSecurityTokensFact, not %T", t.Fact())
		}
		did = fact.Sender().String()
		didtype = DuplicationTypeSender

	case RedeemTokens:
		fact, ok := t.Fact().(RedeemTokensFact)
		if !ok {
			return errors.Errorf("expected RedeemTokensFact, not %T", t.Fact())
		}
		did = fact.Sender().String()
		didtype = DuplicationTypeSender
	case RevokeOperators:
		fact, ok := t.Fact().(RevokeOperatorsFact)
		if !ok {
			return errors.Errorf("expected RevokeOperatorsFact, not %T", t.Fact())
		}
		did = fact.Sender().String()
		didtype = DuplicationTypeSender
	case SetDocument:
		fact, ok := t.Fact().(SetDocumentFact)
		if !ok {
			return errors.Errorf("expected SetDocumentFact, not %T", t.Fact())
		}
		did = fact.Sender().String()
		didtype = DuplicationTypeSender
	case TransferSecurityTokensPartition:
		fact, ok := t.Fact().(TransferSecurityTokensPartitionFact)
		if !ok {
			return errors.Errorf("expected TransferSecurityTokensPartitionFact, not %T", t.Fact())
		}
		did = fact.Sender().String()
		didtype = DuplicationTypeSender
	case kyc.CreateKYCService:
		fact, ok := t.Fact().(kyc.CreateKYCServiceFact)
		if !ok {
			return errors.Errorf("expected CreateKYCServiceFact, not %T", t.Fact())
		}
		did = fact.Sender().String()
		didtype = DuplicationTypeSender
	case kyc.AddControllers:
		fact, ok := t.Fact().(kyc.AddControllersFact)
		if !ok {
			return errors.Errorf("expected AddControllersFact, not %T", t.Fact())
		}
		did = fact.Sender().String()
		didtype = DuplicationTypeSender
	case kyc.RemoveControllers:
		fact, ok := t.Fact().(kyc.RemoveControllersFact)
		if !ok {
			return errors.Errorf("expected RemoveControllersFact, not %T", t.Fact())
		}
		did = fact.Sender().String()
		didtype = DuplicationTypeSender
	case kyc.AddCustomers:
		fact, ok := t.Fact().(kyc.AddCustomersFact)
		if !ok {
			return errors.Errorf("expected AddCustomersFact, not %T", t.Fact())
		}
		did = fact.Sender().String()
		didtype = DuplicationTypeSender
	case kyc.UpdateCustomers:
		fact, ok := t.Fact().(kyc.UpdateCustomersFact)
		if !ok {
			return errors.Errorf("expected UpdateCustomersFact, not %T", t.Fact())
		}
		did = fact.Sender().String()
		didtype = DuplicationTypeSender
	default:
		return nil
	}

	if len(did) > 0 {
		if _, found := opr.duplicated[did]; found {
			switch didtype {
			case DuplicationTypeSender:
				return errors.Errorf("violates only one sender in proposal")
			case DuplicationTypeCurrency:
				return errors.Errorf("duplicate currency id, %q found in proposal", did)
			default:
				return errors.Errorf("violates duplication in proposal")
			}
		}

		opr.duplicated[did] = didtype
	}

	if len(newAddresses) > 0 {
		if err := opr.checkNewAddressDuplication(newAddresses); err != nil {
			return err
		}
	}

	return nil
}

func (opr *OperationProcessor) checkNewAddressDuplication(as []base.Address) error {
	for i := range as {
		if _, found := opr.duplicatedNewAddress[as[i].String()]; found {
			return errors.Errorf("new address already processed")
		}
	}

	for i := range as {
		opr.duplicatedNewAddress[as[i].String()] = struct{}{}
	}

	return nil
}

func (opr *OperationProcessor) Close() error {
	opr.Lock()

	defer opr.Unlock()
	defer opr.close()

	return nil
}

func (opr *OperationProcessor) Cancel() error {
	opr.Lock()
	defer opr.Unlock()

	defer opr.close()

	return nil
}

func (opr *OperationProcessor) getNewProcessor(op base.Operation) (base.OperationProcessor, bool, error) {
	switch i, err := opr.getNewProcessorFromHintset(op); {
	case err != nil:
		return nil, false, err
	case i != nil:
		return i, true, nil
	}

	switch t := op.(type) {
	case currency.CreateAccounts,
		currency.KeyUpdater,
		currency.Transfers,
		extensioncurrency.CreateContractAccounts,
		extensioncurrency.Withdraws,
		currency.CurrencyRegister,
		currency.CurrencyPolicyUpdater,
		currency.SuffrageInflation,
		AuthorizeOperators,
		CreateSecurityTokens,
		IssueSecurityTokens,
		RedeemTokens,
		RevokeOperators,
		SetDocument,
		TransferSecurityTokensPartition,
		kyc.CreateKYCService,
		kyc.AddControllers,
		kyc.RemoveControllers,
		kyc.AddCustomers,
		kyc.UpdateCustomers:
		return nil, false, errors.Errorf("%T needs SetProcessor", t)
	default:
		return nil, false, nil
	}
}

func (opr *OperationProcessor) getNewProcessorFromHintset(op base.Operation) (base.OperationProcessor, error) {
	var f types.GetNewProcessor

	if hinter, ok := op.(hint.Hinter); !ok {
		return nil, nil
	} else if i := opr.processorHintSet.Find(hinter.Hint()); i == nil {
		return nil, nil
	} else if j, ok := i.(types.GetNewProcessor); !ok {
		return nil, errors.Errorf("invalid GetNewProcessor func, %T", i)
	} else {
		f = j
	}

	opp, err := f(opr.Height(), opr.GetStateFunc, nil, nil)
	if err != nil {
		return nil, err
	}

	h := op.(util.Hasher).Hash().String()
	_, iscloser := opp.(io.Closer)
	if iscloser {
		opr.processorClosers.Store(h, opp)
		iscloser = true
	}

	opr.Log().Debug().
		Str("operation", h).
		Str("processor", fmt.Sprintf("%T", opp)).
		Bool("is_closer", iscloser).
		Msg("operation processor created")

	return opp, nil
}

func (opr *OperationProcessor) close() {
	opr.processorClosers.Range(func(_, v interface{}) bool {
		err := v.(io.Closer).Close()
		if err != nil {
			opr.Log().Error().Err(err).Str("op", fmt.Sprintf("%T", v)).Msg("failed to close operation processor")
		} else {
			opr.Log().Debug().Str("processor", fmt.Sprintf("%T", v)).Msg("operation processor closed")
		}

		return true
	})

	opr.fee = nil
	opr.duplicated = nil
	opr.duplicatedNewAddress = nil
	opr.processorClosers = &sync.Map{}

	operationProcessorPool.Put(opr)

	opr.Log().Debug().Msg("operation processors closed")
}
