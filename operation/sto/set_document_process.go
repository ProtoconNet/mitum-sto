package sto

import (
	"context"
	"sync"

	"github.com/ProtoconNet/mitum-currency/v3/common"
	currencystate "github.com/ProtoconNet/mitum-currency/v3/state"
	currency "github.com/ProtoconNet/mitum-currency/v3/state/currency"
	extensioncurrency "github.com/ProtoconNet/mitum-currency/v3/state/extension"
	currencytypes "github.com/ProtoconNet/mitum-currency/v3/types"
	stostate "github.com/ProtoconNet/mitum-sto/state/sto"
	stotypes "github.com/ProtoconNet/mitum-sto/types/sto"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/pkg/errors"
)

var setDocumentProcessorPool = sync.Pool{
	New: func() interface{} {
		return new(SetDocumentProcessor)
	},
}

func (SetDocument) Process(
	ctx context.Context, getStateFunc base.GetStateFunc,
) ([]base.StateMergeValue, base.OperationProcessReasonError, error) {
	return nil, nil, nil
}

type SetDocumentProcessor struct {
	*base.BaseOperationProcessor
}

func NewSetDocumentProcessor() currencytypes.GetNewProcessor {
	return func(
		height base.Height,
		getStateFunc base.GetStateFunc,
		newPreProcessConstraintFunc base.NewOperationProcessorProcessFunc,
		newProcessConstraintFunc base.NewOperationProcessorProcessFunc,
	) (base.OperationProcessor, error) {
		e := util.StringError("failed to create new SetDocumentProcessor")

		nopp := setDocumentProcessorPool.Get()
		opp, ok := nopp.(*SetDocumentProcessor)
		if !ok {
			return nil, errors.Errorf("expected SetDocumentProcessor, not %T", nopp)
		}

		b, err := base.NewBaseOperationProcessor(
			height, getStateFunc, newPreProcessConstraintFunc, newProcessConstraintFunc)
		if err != nil {
			return nil, e.Wrap(err)
		}

		opp.BaseOperationProcessor = b

		return opp, nil
	}
}

func (opp *SetDocumentProcessor) PreProcess(
	ctx context.Context, op base.Operation, getStateFunc base.GetStateFunc,
) (context.Context, base.OperationProcessReasonError, error) {
	e := util.StringError("failed to preprocess SetDocument")

	fact, ok := op.Fact().(SetDocumentFact)
	if !ok {
		return ctx, nil, e.Wrap(errors.Errorf("not SetDocumentFact, %T", op.Fact()))
	}

	if err := fact.IsValid(nil); err != nil {
		return ctx, nil, e.Wrap(err)
	}

	if err := currencystate.CheckExistsState(currency.StateKeyAccount(fact.Sender()), getStateFunc); err != nil {
		return nil, base.NewBaseOperationProcessReasonError("sender not found, %q: %w", fact.Sender(), err), nil
	}

	if err := currencystate.CheckNotExistsState(extensioncurrency.StateKeyContractAccount(fact.Sender()), getStateFunc); err != nil {
		return nil, base.NewBaseOperationProcessReasonError("contract account cannot update sto documents, %q: %w", fact.Sender(), err), nil
	}

	if err := currencystate.CheckFactSignsByState(fact.Sender(), op.Signs(), getStateFunc); err != nil {
		return ctx, base.NewBaseOperationProcessReasonError("invalid signing: %w", err), nil
	}

	policy, err := stostate.ExistsPolicy(fact.Contract(), getStateFunc)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError("sto policy not found, %s: %w", fact.Contract(), err), nil
	}

	controllers := policy.Controllers()
	if len(controllers) == 0 {
		return nil, base.NewBaseOperationProcessReasonError("empty controllers, %s", fact.Contract()), nil
	}

	for i, con := range controllers {
		if con.Equal(fact.Sender()) {
			break
		}

		if i == len(controllers)-1 {
			return nil, base.NewBaseOperationProcessReasonError("sender is not controller of sto, %q, %s", fact.Sender(), fact.Contract()), nil
		}
	}

	return ctx, nil, nil
}

func (opp *SetDocumentProcessor) Process(
	ctx context.Context, op base.Operation, getStateFunc base.GetStateFunc) (
	[]base.StateMergeValue, base.OperationProcessReasonError, error,
) {
	e := util.StringError("failed to process SetDocument")

	fact, ok := op.Fact().(SetDocumentFact)
	if !ok {
		return nil, nil, e.Wrap(errors.Errorf("expected SetDocumentFact, not %T", op.Fact()))
	}

	doc := stotypes.NewDocument(fact.Title(), fact.DocumentHash(), fact.URI())
	if err := doc.IsValid(nil); err != nil {
		return nil, base.NewBaseOperationProcessReasonError("invalid sto document, %q: %w", fact.DocumentHash(), err), nil
	}

	st, err := currencystate.ExistsState(stostate.StateKeyDesign(fact.Contract()), "key of sto design", getStateFunc)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError("sto design not found, %s: %w", fact.Contract(), err), nil
	}

	design, err := stostate.StateDesignValue(st)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError("sto design value not found, %s: %w", fact.Contract(), err), nil
	}
	Policy := design.Policy()

	Policy = stotypes.NewPolicy(Policy.Partitions(), Policy.Aggregate(), Policy.Controllers(), append(Policy.Documents(), doc))
	if err := Policy.IsValid(nil); err != nil {
		return nil, base.NewBaseOperationProcessReasonError("invalid sto policy, %s: %w", fact.Contract(), err), nil
	}

	design = stotypes.NewDesign(design.Granularity(), Policy)
	if err := design.IsValid(nil); err != nil {
		return nil, base.NewBaseOperationProcessReasonError("invalid sto design, %s: %w", fact.Contract(), err), nil
	}

	sts := make([]base.StateMergeValue, 2)

	sts[0] = currencystate.NewStateMergeValue(
		stostate.StateKeyDesign(fact.Contract()),
		stostate.NewDesignStateValue(design),
	)

	currencyPolicy, err := currencystate.ExistsCurrencyPolicy(fact.Currency(), getStateFunc)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError("currency not found, %q: %w", fact.Currency(), err), nil
	}

	fee, err := currencyPolicy.Feeer().Fee(common.ZeroBig)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError("failed to check fee of currency, %q: %w", fact.Currency(), err), nil
	}

	st, err = currencystate.ExistsState(currency.StateKeyBalance(fact.Sender(), fact.Currency()), "key of sender balance", getStateFunc)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError("sender balance not found, %q: %w", fact.Sender(), err), nil
	}
	sb := currencystate.NewStateMergeValue(st.Key(), st.Value())

	switch b, err := currency.StateBalanceValue(st); {
	case err != nil:
		return nil, base.NewBaseOperationProcessReasonError("failed to get balance value, %q: %w", currency.StateKeyBalance(fact.Sender(), fact.Currency()), err), nil
	case b.Big().Compare(fee) < 0:
		return nil, base.NewBaseOperationProcessReasonError("not enough balance of sender, %q", fact.Sender()), nil
	}

	v, ok := sb.Value().(currency.BalanceStateValue)
	if !ok {
		return nil, base.NewBaseOperationProcessReasonError("expected BalanceStateValue, not %T", sb.Value()), nil
	}
	sts[1] = currencystate.NewStateMergeValue(
		sb.Key(),
		currency.NewBalanceStateValue(v.Amount.WithBig(v.Amount.Big().Sub(fee))),
	)

	return sts, nil, nil
}

func (opp *SetDocumentProcessor) Close() error {
	setDocumentProcessorPool.Put(opp)

	return nil
}
