package kyc

import (
	"context"
	"sync"

	"github.com/ProtoconNet/mitum-currency/v3/common"
	"github.com/ProtoconNet/mitum-currency/v3/state"
	stcurrency "github.com/ProtoconNet/mitum-currency/v3/state/currency"
	stextension "github.com/ProtoconNet/mitum-currency/v3/state/extension"
	"github.com/ProtoconNet/mitum-currency/v3/types"
	stkyc "github.com/ProtoconNet/mitum-sto/state/kyc"
	typekyc "github.com/ProtoconNet/mitum-sto/types/kyc"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/pkg/errors"
)

var createServiceProcessorPool = sync.Pool{
	New: func() interface{} {
		return new(CreateServiceProcessor)
	},
}

func (CreateService) Process(
	ctx context.Context, getStateFunc base.GetStateFunc,
) ([]base.StateMergeValue, base.OperationProcessReasonError, error) {
	return nil, nil, nil
}

type CreateServiceProcessor struct {
	*base.BaseOperationProcessor
}

func NewCreateServiceProcessor() types.GetNewProcessor {
	return func(
		height base.Height,
		getStateFunc base.GetStateFunc,
		newPreProcessConstraintFunc base.NewOperationProcessorProcessFunc,
		newProcessConstraintFunc base.NewOperationProcessorProcessFunc,
	) (base.OperationProcessor, error) {
		e := util.StringError("failed to create new CreateServiceProcessor")

		nopp := createServiceProcessorPool.Get()
		opp, ok := nopp.(*CreateServiceProcessor)
		if !ok {
			return nil, errors.Errorf("expected CreateServiceProcessor, not %T", nopp)
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

func (opp *CreateServiceProcessor) PreProcess(
	ctx context.Context, op base.Operation, getStateFunc base.GetStateFunc,
) (context.Context, base.OperationProcessReasonError, error) {
	e := util.StringError("failed to preprocess CreateService")

	fact, ok := op.Fact().(CreateServiceFact)
	if !ok {
		return ctx, nil, e.Wrap(errors.Errorf("not CreateServiceFact, %T", op.Fact()))
	}

	if err := fact.IsValid(nil); err != nil {
		return ctx, nil, e.Wrap(err)
	}

	if err := state.CheckExistsState(stcurrency.StateKeyAccount(fact.Sender()), getStateFunc); err != nil {
		return nil, base.NewBaseOperationProcessReasonError("sender not found, %q: %w", fact.Sender(), err), nil
	}

	st, err := state.ExistsState(stextension.StateKeyContractAccount(fact.Sender()), "key of contract account", getStateFunc)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError("contract account cannot create kyc service, %q: %w", fact.Sender(), err), nil
	}

	ca, err := stextension.StateContractAccountValue(st)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError("contract account value not found, %q: %w", fact.Contract(), err), nil
	}

	if !ca.Owner().Equal(fact.sender) {
		return nil, base.NewBaseOperationProcessReasonError("not contract account owner, %q", fact.sender), nil
	}

	if err := state.CheckNotExistsState(stkyc.StateKeyDesign(fact.Contract()), getStateFunc); err != nil {
		return nil, base.NewBaseOperationProcessReasonError("kyc service already exists, %s: %w", fact.Contract(), err), nil
	}

	if err := state.CheckFactSignsByState(fact.Sender(), op.Signs(), getStateFunc); err != nil {
		return ctx, base.NewBaseOperationProcessReasonError("invalid signing: %w", err), nil
	}

	return ctx, nil, nil
}

func (opp *CreateServiceProcessor) Process(
	ctx context.Context, op base.Operation, getStateFunc base.GetStateFunc) (
	[]base.StateMergeValue, base.OperationProcessReasonError, error,
) {
	e := util.StringError("failed to process CreateService")

	fact, ok := op.Fact().(CreateServiceFact)
	if !ok {
		return nil, nil, e.Wrap(errors.Errorf("expected CreateServiceFact, not %T", op.Fact()))
	}

	policy := typekyc.NewPolicy(fact.Controllers())
	if err := policy.IsValid(nil); err != nil {
		return nil, base.NewBaseOperationProcessReasonError("invalid kyc policy, %s-%s: %w", fact.Contract(), err), nil
	}

	design := typekyc.NewDesign(policy)
	if err := design.IsValid(nil); err != nil {
		return nil, base.NewBaseOperationProcessReasonError("invalid kyc design, %s-%s: %w", fact.Contract(), err), nil
	}

	sts := make([]base.StateMergeValue, 2)

	sts[0] = state.NewStateMergeValue(
		stkyc.StateKeyDesign(fact.Contract()),
		stkyc.NewDesignStateValue(design),
	)

	currencyPolicy, err := state.ExistsCurrencyPolicy(fact.Currency(), getStateFunc)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError("currency not found, %q: %w", fact.Currency(), err), nil
	}

	fee, err := currencyPolicy.Feeer().Fee(common.ZeroBig)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError("failed to check fee of currency, %q: %w", fact.Currency(), err), nil
	}

	st, err := state.ExistsState(stcurrency.StateKeyBalance(fact.Sender(), fact.Currency()), "key of sender balance", getStateFunc)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError("sender balance not found, %q: %w", fact.Sender(), err), nil
	}
	sb := state.NewStateMergeValue(st.Key(), st.Value())

	switch b, err := stcurrency.StateBalanceValue(st); {
	case err != nil:
		return nil, base.NewBaseOperationProcessReasonError("failed to get balance value, %q: %w", stcurrency.StateKeyBalance(fact.Sender(), fact.Currency()), err), nil
	case b.Big().Compare(fee) < 0:
		return nil, base.NewBaseOperationProcessReasonError("not enough balance of sender, %q", fact.Sender()), nil
	}

	v, ok := sb.Value().(stcurrency.BalanceStateValue)
	if !ok {
		return nil, base.NewBaseOperationProcessReasonError("expected BalanceStateValue, not %T", sb.Value()), nil
	}
	sts[1] = state.NewStateMergeValue(
		sb.Key(),
		stcurrency.NewBalanceStateValue(v.Amount.WithBig(v.Amount.Big().Sub(fee))),
	)

	return sts, nil, nil
}

func (opp *CreateServiceProcessor) Close() error {
	createServiceProcessorPool.Put(opp)

	return nil
}
