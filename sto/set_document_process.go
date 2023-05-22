package sto

import (
	"context"
	"sync"

	extensioncurrency "github.com/ProtoconNet/mitum-currency-extension/v2/currency"
	"github.com/ProtoconNet/mitum-currency/v2/currency"
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

func NewSetDocumentProcessor() extensioncurrency.GetNewProcessor {
	return func(
		height base.Height,
		getStateFunc base.GetStateFunc,
		newPreProcessConstraintFunc base.NewOperationProcessorProcessFunc,
		newProcessConstraintFunc base.NewOperationProcessorProcessFunc,
	) (base.OperationProcessor, error) {
		e := util.StringErrorFunc("failed to create new SetDocumentProcessor")

		nopp := setDocumentProcessorPool.Get()
		opp, ok := nopp.(*SetDocumentProcessor)
		if !ok {
			return nil, errors.Errorf("expected SetDocumentProcessor, not %T", nopp)
		}

		b, err := base.NewBaseOperationProcessor(
			height, getStateFunc, newPreProcessConstraintFunc, newProcessConstraintFunc)
		if err != nil {
			return nil, e(err, "")
		}

		opp.BaseOperationProcessor = b

		return opp, nil
	}
}

func (opp *SetDocumentProcessor) PreProcess(
	ctx context.Context, op base.Operation, getStateFunc base.GetStateFunc,
) (context.Context, base.OperationProcessReasonError, error) {
	e := util.StringErrorFunc("failed to preprocess SetDocument")

	fact, ok := op.Fact().(SetDocumentFact)
	if !ok {
		return ctx, nil, e(nil, "not SetDocumentFact, %T", op.Fact())
	}

	if err := fact.IsValid(nil); err != nil {
		return ctx, nil, e(err, "")
	}

	if err := checkExistsState(currency.StateKeyAccount(fact.Sender()), getStateFunc); err != nil {
		return nil, base.NewBaseOperationProcessReasonError("sender not found, %q: %w", fact.Sender(), err), nil
	}

	if err := checkNotExistsState(extensioncurrency.StateKeyContractAccount(fact.Sender()), getStateFunc); err != nil {
		return nil, base.NewBaseOperationProcessReasonError("contract account cannot update sto documents, %q: %w", fact.Sender(), err), nil
	}

	if err := checkFactSignsByState(fact.Sender(), op.Signs(), getStateFunc); err != nil {
		return ctx, base.NewBaseOperationProcessReasonError("invalid signing: %w", err), nil
	}

	policy, err := existsSTOPolicy(fact.Contract(), fact.STO(), getStateFunc)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError("sto policy not found, %s-%s: %w", fact.Contract(), fact.STO(), err), nil
	}

	controllers := policy.Controllers()
	if len(controllers) == 0 {
		return nil, base.NewBaseOperationProcessReasonError("empty controllers, %s-%s", fact.Contract(), fact.STO()), nil
	}

	for i, con := range controllers {
		if con.Equal(fact.Sender()) {
			break
		}

		if i == len(controllers)-1 {
			return nil, base.NewBaseOperationProcessReasonError("sender is not controller of sto, %q, %s-%s", fact.Sender(), fact.Contract(), fact.STO()), nil
		}
	}

	return ctx, nil, nil
}

func (opp *SetDocumentProcessor) Process(
	ctx context.Context, op base.Operation, getStateFunc base.GetStateFunc) (
	[]base.StateMergeValue, base.OperationProcessReasonError, error,
) {
	e := util.StringErrorFunc("failed to process SetDocument")

	fact, ok := op.Fact().(SetDocumentFact)
	if !ok {
		return nil, nil, e(nil, "expected SetDocumentFact, not %T", op.Fact())
	}

	doc := NewDocument(fact.STO(), fact.Title(), fact.DocumentHash(), fact.URI())
	if err := doc.IsValid(nil); err != nil {
		return nil, base.NewBaseOperationProcessReasonError("invalid sto document, %q: %w", fact.DocumentHash(), err), nil
	}

	st, err := existsState(StateKeyDesign(fact.Contract(), fact.STO()), "key of sto design", getStateFunc)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError("sto design not found, %s-%s: %w", fact.Contract(), fact.STO(), err), nil
	}

	design, err := StateDesignValue(st)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError("sto design value not found, %s-%s: %w", fact.Contract(), fact.STO(), err), nil
	}
	stoPolicy := design.Policy()

	stoPolicy = NewSTOPolicy(stoPolicy.Partitions(), stoPolicy.Aggregate(), stoPolicy.Controllers(), append(stoPolicy.Documents(), doc))
	if err := stoPolicy.IsValid(nil); err != nil {
		return nil, base.NewBaseOperationProcessReasonError("invalid sto policy, %s-%s: %w", fact.Contract(), fact.STO(), err), nil
	}

	design = NewDesign(design.STO(), design.Granularity(), stoPolicy)
	if err := design.IsValid(nil); err != nil {
		return nil, base.NewBaseOperationProcessReasonError("invalid sto design, %s-%s: %w", fact.Contract(), fact.STO(), err), nil
	}

	sts := make([]base.StateMergeValue, 2)

	sts[0] = NewStateMergeValue(
		StateKeyDesign(fact.Contract(), fact.STO()),
		NewDesignStateValue(design),
	)

	currencyPolicy, err := existsCurrencyPolicy(fact.Currency(), getStateFunc)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError("currency not found, %q: %w", fact.Currency(), err), nil
	}

	fee, err := currencyPolicy.Feeer().Fee(currency.ZeroBig)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError("failed to check fee of currency, %q: %w", fact.Currency(), err), nil
	}

	st, err = existsState(currency.StateKeyBalance(fact.Sender(), fact.Currency()), "key of sender balance", getStateFunc)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError("sender balance not found, %q: %w", fact.Sender(), err), nil
	}
	sb := currency.NewBalanceStateMergeValue(st.Key(), st.Value())

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
	sts[1] = currency.NewBalanceStateMergeValue(
		sb.Key(),
		currency.NewBalanceStateValue(v.Amount.WithBig(v.Amount.Big().Sub(fee))),
	)

	return sts, nil, nil
}

func (opp *SetDocumentProcessor) Close() error {
	setDocumentProcessorPool.Put(opp)

	return nil
}
