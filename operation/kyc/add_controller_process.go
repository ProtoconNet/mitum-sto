package kyc

import (
	"context"
	"sync"

	"github.com/ProtoconNet/mitum-currency/v3/common"
	"github.com/ProtoconNet/mitum-currency/v3/operation/currency"
	state "github.com/ProtoconNet/mitum-currency/v3/state"
	stcurrency "github.com/ProtoconNet/mitum-currency/v3/state/currency"
	stextension "github.com/ProtoconNet/mitum-currency/v3/state/extension"
	"github.com/ProtoconNet/mitum-currency/v3/types"
	stkyc "github.com/ProtoconNet/mitum-sto/state/kyc"
	typekyc "github.com/ProtoconNet/mitum-sto/types/kyc"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/pkg/errors"
)

var addControllerItemProcessorPool = sync.Pool{
	New: func() interface{} {
		return new(AddControllerItemProcessor)
	},
}

var addControllerProcessorPool = sync.Pool{
	New: func() interface{} {
		return new(AddControllerProcessor)
	},
}

func (AddController) Process(
	ctx context.Context, getStateFunc base.GetStateFunc,
) ([]base.StateMergeValue, base.OperationProcessReasonError, error) {
	return nil, nil, nil
}

type AddControllerItemProcessor struct {
	h           util.Hash
	sender      base.Address
	item        AddControllerItem
	controllers *[]base.Address
}

func (ipp *AddControllerItemProcessor) PreProcess(
	ctx context.Context, op base.Operation, getStateFunc base.GetStateFunc,
) error {
	it := ipp.item

	st, err := state.ExistsState(stextension.StateKeyContractAccount(it.Contract()), "key of contract account", getStateFunc)
	if err != nil {
		return err
	}

	ca, err := stextension.StateContractAccountValue(st)
	if err != nil {
		return err
	}

	if !ca.Owner().Equal(ipp.sender) {
		return errors.Errorf("not contract account owner, %q", it.Contract())
	}

	for _, ad := range *ipp.controllers {
		if ad.Equal(it.Controller()) {
			return errors.Errorf("controller is already in kyc policy controllers, %q", ad)
		}
	}

	if err := state.CheckExistsState(stcurrency.StateKeyCurrencyDesign(it.Currency()), getStateFunc); err != nil {
		return err
	}

	return nil
}

func (ipp *AddControllerItemProcessor) Process(
	ctx context.Context, op base.Operation, getStateFunc base.GetStateFunc,
) ([]base.StateMergeValue, error) {
	*ipp.controllers = append(*ipp.controllers, ipp.item.Controller())
	return nil, nil
}

func (ipp *AddControllerItemProcessor) Close() error {
	ipp.h = nil
	ipp.sender = nil
	ipp.item = AddControllerItem{}
	ipp.controllers = nil

	addControllerItemProcessorPool.Put(ipp)

	return nil
}

type AddControllerProcessor struct {
	*base.BaseOperationProcessor
}

func NewAddControllerProcessor() types.GetNewProcessor {
	return func(
		height base.Height,
		getStateFunc base.GetStateFunc,
		newPreProcessConstraintFunc base.NewOperationProcessorProcessFunc,
		newProcessConstraintFunc base.NewOperationProcessorProcessFunc,
	) (base.OperationProcessor, error) {
		e := util.StringError("failed to create new AddControllerProcessor")

		nopp := addControllerProcessorPool.Get()
		opp, ok := nopp.(*AddControllerProcessor)
		if !ok {
			return nil, e.Wrap(errors.Errorf("expected AddControllerProcessor, not %T", nopp))
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

func (opp *AddControllerProcessor) PreProcess(
	ctx context.Context, op base.Operation, getStateFunc base.GetStateFunc,
) (context.Context, base.OperationProcessReasonError, error) {
	e := util.StringError("failed to preprocess AddController")

	fact, ok := op.Fact().(AddControllerFact)
	if !ok {
		return ctx, nil, e.Wrap(errors.Errorf("expected AddControllerFact, not %T", op.Fact()))
	}

	if err := fact.IsValid(nil); err != nil {
		return ctx, nil, e.Wrap(err)
	}

	if err := state.CheckExistsState(stcurrency.StateKeyAccount(fact.Sender()), getStateFunc); err != nil {
		return ctx, base.NewBaseOperationProcessReasonError("sender not found, %q: %w", fact.Sender(), err), nil
	}

	if err := state.CheckNotExistsState(stextension.StateKeyContractAccount(fact.Sender()), getStateFunc); err != nil {
		return ctx, base.NewBaseOperationProcessReasonError("contract account cannot set its controllers, %q: %w", fact.Sender(), err), nil
	}

	if err := state.CheckFactSignsByState(fact.sender, op.Signs(), getStateFunc); err != nil {
		return ctx, base.NewBaseOperationProcessReasonError("invalid signing: %w", err), nil
	}

	controllers := map[string]*[]base.Address{}

	for _, it := range fact.Items() {
		policy, err := stkyc.ExistsPolicy(it.Contract(), getStateFunc)
		if err != nil {
			return nil, base.NewBaseOperationProcessReasonError("failed to get kyc policy, %s: %w", it.Contract(), err), nil
		}
		cons := policy.Controllers()
		controllers[stkyc.StateKeyDesign(it.Contract())] = &cons
	}

	for _, it := range fact.Items() {
		ip := addControllerItemProcessorPool.Get()
		ipc, ok := ip.(*AddControllerItemProcessor)
		if !ok {
			return nil, nil, e.Wrap(errors.Errorf("expected AddControllerItemProcessor, not %T", ip))
		}

		ipc.h = op.Hash()
		ipc.sender = fact.Sender()
		ipc.item = it
		ipc.controllers = controllers[stkyc.StateKeyDesign(it.Contract())]

		if err := ipc.PreProcess(ctx, op, getStateFunc); err != nil {
			return nil, base.NewBaseOperationProcessReasonError("failed to preprocess AddControllerItem: %w", err), nil
		}

		ipc.Close()
	}

	return ctx, nil, nil
}

func (opp *AddControllerProcessor) Process( // nolint:dupl
	ctx context.Context, op base.Operation, getStateFunc base.GetStateFunc) (
	[]base.StateMergeValue, base.OperationProcessReasonError, error,
) {
	e := util.StringError("failed to process AddController")

	fact, ok := op.Fact().(AddControllerFact)
	if !ok {
		return nil, nil, e.Wrap(errors.Errorf("expected AddControllerFact, not %T", op.Fact()))
	}

	var sts []base.StateMergeValue // nolint:prealloc

	controllers := map[string]*[]base.Address{}

	for _, it := range fact.Items() {
		policy, err := stkyc.ExistsPolicy(it.Contract(), getStateFunc)
		if err != nil {
			return nil, base.NewBaseOperationProcessReasonError("failed to get kyc policy, %s: %w", it.Contract(), err), nil
		}
		cons := policy.Controllers()
		controllers[stkyc.StateKeyDesign(it.Contract())] = &cons
	}

	for _, it := range fact.Items() {
		ip := addControllerItemProcessorPool.Get()
		ipc, ok := ip.(*AddControllerItemProcessor)
		if !ok {
			return nil, nil, e.Wrap(errors.Errorf("expected AddControllerItemProcessor, not %T", ip))
		}

		ipc.h = op.Hash()
		ipc.sender = fact.Sender()
		ipc.item = it
		ipc.controllers = controllers[stkyc.StateKeyDesign(it.Contract())]

		_, err := ipc.Process(ctx, op, getStateFunc)
		if err != nil {
			return nil, base.NewBaseOperationProcessReasonError("failed to process AddControllerItem: %w", err), nil
		}

		ipc.Close()
	}

	for k, m := range controllers {
		policy := typekyc.NewPolicy(*m)
		design := typekyc.NewDesign(policy)
		if err := design.IsValid(nil); err != nil {
			return nil, base.NewBaseOperationProcessReasonError("invalid design, %s: %w", k, err), nil
		}

		sts = append(sts, state.NewStateMergeValue(
			k,
			stkyc.NewDesignStateValue(design),
		))
	}

	fitems := fact.Items()
	items := make([]KYCItem, len(fitems))
	for i := range fact.Items() {
		items[i] = fitems[i]
	}

	required, err := calculateKYCItemsFee(getStateFunc, items)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError("failed to calculate fee: %w", err), nil
	}
	sb, err := currency.CheckEnoughBalance(fact.sender, required, getStateFunc)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError("failed to check enough balance: %w", err), nil
	}

	for i := range sb {
		v, ok := sb[i].Value().(stcurrency.BalanceStateValue)
		if !ok {
			return nil, nil, e.Wrap(errors.Errorf("expected BalanceStateValue, not %T", sb[i].Value()))
		}
		stv := stcurrency.NewBalanceStateValue(v.Amount.WithBig(v.Amount.Big().Sub(required[i][0])))
		sts = append(sts, state.NewStateMergeValue(sb[i].Key(), stv))
	}

	return sts, nil, nil
}

func (opp *AddControllerProcessor) Close() error {
	addControllerProcessorPool.Put(opp)

	return nil
}

func calculateKYCItemsFee(getStateFunc base.GetStateFunc, items []KYCItem) (map[types.CurrencyID][2]common.Big, error) {
	required := map[types.CurrencyID][2]common.Big{}

	for _, item := range items {
		rq := [2]common.Big{common.ZeroBig, common.ZeroBig}

		if k, found := required[item.Currency()]; found {
			rq = k
		}

		policy, err := state.ExistsCurrencyPolicy(item.Currency(), getStateFunc)
		if err != nil {
			return nil, err
		}

		switch k, err := policy.Feeer().Fee(common.ZeroBig); {
		case err != nil:
			return nil, err
		case !k.OverZero():
			required[item.Currency()] = [2]common.Big{rq[0], rq[1]}
		default:
			required[item.Currency()] = [2]common.Big{rq[0].Add(k), rq[1].Add(k)}
		}

	}

	return required, nil

}
