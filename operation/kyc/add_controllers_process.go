package kyc

import (
	"context"
	"sync"

	"github.com/ProtoconNet/mitum-currency/v3/common"
	currencyoperation "github.com/ProtoconNet/mitum-currency/v3/operation/currency"
	currencystate "github.com/ProtoconNet/mitum-currency/v3/state"
	currency "github.com/ProtoconNet/mitum-currency/v3/state/currency"
	extensioncurrency "github.com/ProtoconNet/mitum-currency/v3/state/extension"
	currencytypes "github.com/ProtoconNet/mitum-currency/v3/types"
	kycstate "github.com/ProtoconNet/mitum-sto/state/kyc"
	kyctypes "github.com/ProtoconNet/mitum-sto/types/kyc"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/pkg/errors"
)

var addControllersItemProcessorPool = sync.Pool{
	New: func() interface{} {
		return new(AddControllersItemProcessor)
	},
}

var addControllersProcessorPool = sync.Pool{
	New: func() interface{} {
		return new(AddControllersProcessor)
	},
}

func (AddControllers) Process(
	ctx context.Context, getStateFunc base.GetStateFunc,
) ([]base.StateMergeValue, base.OperationProcessReasonError, error) {
	return nil, nil, nil
}

type AddControllersItemProcessor struct {
	h           util.Hash
	sender      base.Address
	item        AddControllersItem
	controllers *[]base.Address
}

func (ipp *AddControllersItemProcessor) PreProcess(
	ctx context.Context, op base.Operation, getStateFunc base.GetStateFunc,
) error {
	it := ipp.item

	st, err := currencystate.ExistsState(extensioncurrency.StateKeyContractAccount(it.Contract()), "key of contract account", getStateFunc)
	if err != nil {
		return err
	}

	ca, err := extensioncurrency.StateContractAccountValue(st)
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

	if err := currencystate.CheckExistsState(currency.StateKeyCurrencyDesign(it.Currency()), getStateFunc); err != nil {
		return err
	}

	return nil
}

func (ipp *AddControllersItemProcessor) Process(
	ctx context.Context, op base.Operation, getStateFunc base.GetStateFunc,
) ([]base.StateMergeValue, error) {
	*ipp.controllers = append(*ipp.controllers, ipp.item.Controller())
	return nil, nil
}

func (ipp *AddControllersItemProcessor) Close() error {
	ipp.h = nil
	ipp.sender = nil
	ipp.item = AddControllersItem{}
	ipp.controllers = nil

	addControllersItemProcessorPool.Put(ipp)

	return nil
}

type AddControllersProcessor struct {
	*base.BaseOperationProcessor
}

func NewAddControllersProcessor() currencytypes.GetNewProcessor {
	return func(
		height base.Height,
		getStateFunc base.GetStateFunc,
		newPreProcessConstraintFunc base.NewOperationProcessorProcessFunc,
		newProcessConstraintFunc base.NewOperationProcessorProcessFunc,
	) (base.OperationProcessor, error) {
		e := util.StringError("failed to create new AddControllersProcessor")

		nopp := addControllersProcessorPool.Get()
		opp, ok := nopp.(*AddControllersProcessor)
		if !ok {
			return nil, e.Wrap(errors.Errorf("expected AddControllersProcessor, not %T", nopp))
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

func (opp *AddControllersProcessor) PreProcess(
	ctx context.Context, op base.Operation, getStateFunc base.GetStateFunc,
) (context.Context, base.OperationProcessReasonError, error) {
	e := util.StringError("failed to preprocess AddControllers")

	fact, ok := op.Fact().(AddControllersFact)
	if !ok {
		return ctx, nil, e.Wrap(errors.Errorf("expected AddControllersFact, not %T", op.Fact()))
	}

	if err := fact.IsValid(nil); err != nil {
		return ctx, nil, e.Wrap(err)
	}

	if err := currencystate.CheckExistsState(currency.StateKeyAccount(fact.Sender()), getStateFunc); err != nil {
		return ctx, base.NewBaseOperationProcessReasonError("sender not found, %q: %w", fact.Sender(), err), nil
	}

	if err := currencystate.CheckNotExistsState(extensioncurrency.StateKeyContractAccount(fact.Sender()), getStateFunc); err != nil {
		return ctx, base.NewBaseOperationProcessReasonError("contract account cannot set its controllers, %q: %w", fact.Sender(), err), nil
	}

	if err := currencystate.CheckFactSignsByState(fact.sender, op.Signs(), getStateFunc); err != nil {
		return ctx, base.NewBaseOperationProcessReasonError("invalid signing: %w", err), nil
	}

	controllers := map[string]*[]base.Address{}

	for _, it := range fact.Items() {
		policy, err := kycstate.ExistsPolicy(it.Contract(), getStateFunc)
		if err != nil {
			return nil, base.NewBaseOperationProcessReasonError("failed to get kyc policy, %s: %w", it.Contract(), err), nil
		}
		cons := policy.Controllers()
		controllers[kycstate.StateKeyDesign(it.Contract())] = &cons
	}

	for _, it := range fact.Items() {
		ip := addControllersItemProcessorPool.Get()
		ipc, ok := ip.(*AddControllersItemProcessor)
		if !ok {
			return nil, nil, e.Wrap(errors.Errorf("expected AddControllersItemProcessor, not %T", ip))
		}

		ipc.h = op.Hash()
		ipc.sender = fact.Sender()
		ipc.item = it
		ipc.controllers = controllers[kycstate.StateKeyDesign(it.Contract())]

		if err := ipc.PreProcess(ctx, op, getStateFunc); err != nil {
			return nil, base.NewBaseOperationProcessReasonError("failed to preprocess AddControllersItem: %w", err), nil
		}

		ipc.Close()
	}

	return ctx, nil, nil
}

func (opp *AddControllersProcessor) Process( // nolint:dupl
	ctx context.Context, op base.Operation, getStateFunc base.GetStateFunc) (
	[]base.StateMergeValue, base.OperationProcessReasonError, error,
) {
	e := util.StringError("failed to process AddControllers")

	fact, ok := op.Fact().(AddControllersFact)
	if !ok {
		return nil, nil, e.Wrap(errors.Errorf("expected AddControllersFact, not %T", op.Fact()))
	}

	var sts []base.StateMergeValue // nolint:prealloc

	controllers := map[string]*[]base.Address{}

	for _, it := range fact.Items() {
		policy, err := kycstate.ExistsPolicy(it.Contract(), getStateFunc)
		if err != nil {
			return nil, base.NewBaseOperationProcessReasonError("failed to get kyc policy, %s: %w", it.Contract(), err), nil
		}
		cons := policy.Controllers()
		controllers[kycstate.StateKeyDesign(it.Contract())] = &cons
	}

	for _, it := range fact.Items() {
		ip := addControllersItemProcessorPool.Get()
		ipc, ok := ip.(*AddControllersItemProcessor)
		if !ok {
			return nil, nil, e.Wrap(errors.Errorf("expected AddControllersItemProcessor, not %T", ip))
		}

		ipc.h = op.Hash()
		ipc.sender = fact.Sender()
		ipc.item = it
		ipc.controllers = controllers[kycstate.StateKeyDesign(it.Contract())]

		_, err := ipc.Process(ctx, op, getStateFunc)
		if err != nil {
			return nil, base.NewBaseOperationProcessReasonError("failed to process AddControllersItem: %w", err), nil
		}

		ipc.Close()
	}

	for k, m := range controllers {
		policy := kyctypes.NewPolicy(*m)
		design := kyctypes.NewDesign(policy)
		if err := design.IsValid(nil); err != nil {
			return nil, base.NewBaseOperationProcessReasonError("invalid design, %s: %w", k, err), nil
		}

		sts = append(sts, currencystate.NewStateMergeValue(
			k,
			kycstate.NewDesignStateValue(design),
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
	sb, err := currencyoperation.CheckEnoughBalance(fact.sender, required, getStateFunc)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError("failed to check enough balance: %w", err), nil
	}

	for i := range sb {
		v, ok := sb[i].Value().(currency.BalanceStateValue)
		if !ok {
			return nil, nil, e.Wrap(errors.Errorf("expected BalanceStateValue, not %T", sb[i].Value()))
		}
		stv := currency.NewBalanceStateValue(v.Amount.WithBig(v.Amount.Big().Sub(required[i][0])))
		sts = append(sts, currencystate.NewStateMergeValue(sb[i].Key(), stv))
	}

	return sts, nil, nil
}

func (opp *AddControllersProcessor) Close() error {
	addControllersProcessorPool.Put(opp)

	return nil
}

func calculateKYCItemsFee(getStateFunc base.GetStateFunc, items []KYCItem) (map[currencytypes.CurrencyID][2]common.Big, error) {
	required := map[currencytypes.CurrencyID][2]common.Big{}

	for _, item := range items {
		rq := [2]common.Big{common.ZeroBig, common.ZeroBig}

		if k, found := required[item.Currency()]; found {
			rq = k
		}

		policy, err := currencystate.ExistsCurrencyPolicy(item.Currency(), getStateFunc)
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
