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

	st, err := currencystate.ExistsState(extensioncurrency.StateKeyContractAccount(fact.Contract()), "key of contract account", getStateFunc)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError("target contract account not found, %q; %w", fact.Contract(), err), nil
	}

	ca, err := extensioncurrency.StateContractAccountValue(st)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError("failed to get state value of contract account, %q; %w", fact.Contract(), err), nil
	}

	if !(ca.Owner().Equal(fact.Sender()) || ca.IsOperator(fact.Sender())) {
		return nil, base.NewBaseOperationProcessReasonError("sender is neither the owner nor the operator of the target contract account, %q", fact.sender), nil
	}

	if err := currencystate.CheckFactSignsByState(fact.Sender(), op.Signs(), getStateFunc); err != nil {
		return ctx, base.NewBaseOperationProcessReasonError("invalid signing: %w", err), nil
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

	Policy = stotypes.NewPolicy(Policy.Partitions(), Policy.Aggregate(), append(Policy.Documents(), doc))
	if err := Policy.IsValid(nil); err != nil {
		return nil, base.NewBaseOperationProcessReasonError("invalid sto policy, %s: %w", fact.Contract(), err), nil
	}

	design = stotypes.NewDesign(design.Granularity(), Policy)
	if err := design.IsValid(nil); err != nil {
		return nil, base.NewBaseOperationProcessReasonError("invalid sto design, %s: %w", fact.Contract(), err), nil
	}

	var sts []base.StateMergeValue

	sts = append(sts, currencystate.NewStateMergeValue(
		stostate.StateKeyDesign(fact.Contract()),
		stostate.NewDesignStateValue(design),
	))

	currencyPolicy, err := currencystate.ExistsCurrencyPolicy(fact.Currency(), getStateFunc)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError("currency not found, %q; %w", fact.Currency(), err), nil
	}

	if currencyPolicy.Feeer().Receiver() == nil {
		return sts, nil, nil
	}

	fee, err := currencyPolicy.Feeer().Fee(common.ZeroBig)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError(
			"failed to check fee of currency, %q; %w",
			fact.Currency(),
			err,
		), nil
	}

	senderBalSt, err := currencystate.ExistsState(
		currency.StateKeyBalance(fact.Sender(), fact.Currency()),
		"key of sender balance",
		getStateFunc,
	)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError(
			"sender balance not found, %q; %w",
			fact.Sender(),
			err,
		), nil
	}

	switch senderBal, err := currency.StateBalanceValue(senderBalSt); {
	case err != nil:
		return nil, base.NewBaseOperationProcessReasonError(
			"failed to get balance value, %q; %w",
			currency.StateKeyBalance(fact.Sender(), fact.Currency()),
			err,
		), nil
	case senderBal.Big().Compare(fee) < 0:
		return nil, base.NewBaseOperationProcessReasonError(
			"not enough balance of sender, %q",
			fact.Sender(),
		), nil
	}

	v, ok := senderBalSt.Value().(currency.BalanceStateValue)
	if !ok {
		return nil, base.NewBaseOperationProcessReasonError("expected BalanceStateValue, not %T", senderBalSt.Value()), nil
	}

	if err := currencystate.CheckExistsState(currency.StateKeyAccount(currencyPolicy.Feeer().Receiver()), getStateFunc); err != nil {
		return nil, nil, err
	} else if feeRcvrSt, found, err := getStateFunc(currency.StateKeyBalance(currencyPolicy.Feeer().Receiver(), fact.currency)); err != nil {
		return nil, nil, err
	} else if !found {
		return nil, nil, errors.Errorf("feeer receiver %s not found", currencyPolicy.Feeer().Receiver())
	} else if feeRcvrSt.Key() != senderBalSt.Key() {
		r, ok := feeRcvrSt.Value().(currency.BalanceStateValue)
		if !ok {
			return nil, nil, errors.Errorf("expected %T, not %T", currency.BalanceStateValue{}, feeRcvrSt.Value())
		}
		sts = append(sts, common.NewBaseStateMergeValue(
			feeRcvrSt.Key(),
			currency.NewAddBalanceStateValue(r.Amount.WithBig(fee)),
			func(height base.Height, st base.State) base.StateValueMerger {
				return currency.NewBalanceStateValueMerger(height, feeRcvrSt.Key(), fact.currency, st)
			},
		))

		sts = append(sts, common.NewBaseStateMergeValue(
			senderBalSt.Key(),
			currency.NewDeductBalanceStateValue(v.Amount.WithBig(fee)),
			func(height base.Height, st base.State) base.StateValueMerger {
				return currency.NewBalanceStateValueMerger(height, senderBalSt.Key(), fact.currency, st)
			},
		))
	}

	return sts, nil, nil
}

func (opp *SetDocumentProcessor) Close() error {
	setDocumentProcessorPool.Put(opp)

	return nil
}
