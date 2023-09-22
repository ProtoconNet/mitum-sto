package kyc

import (
	"context"
	"sync"

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

var removeControllerItemProcessorPool = sync.Pool{
	New: func() interface{} {
		return new(RemoveControllerItemProcessor)
	},
}

var removeControllerProcessorPool = sync.Pool{
	New: func() interface{} {
		return new(RemoveControllerProcessor)
	},
}

func (RemoveController) Process(
	ctx context.Context, getStateFunc base.GetStateFunc,
) ([]base.StateMergeValue, base.OperationProcessReasonError, error) {
	return nil, nil, nil
}

type RemoveControllerItemProcessor struct {
	h           util.Hash
	sender      base.Address
	item        RemoveControllerItem
	controllers *[]base.Address
}

func (ipp *RemoveControllerItemProcessor) PreProcess(
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
		return errors.Errorf("not contract account owner, %q", ipp.sender)
	}

	if len(*ipp.controllers) == 0 {
		return errors.Errorf("empty controllers, %s", it.Contract())
	}

	for i, ad := range *ipp.controllers {
		if ad.Equal(it.Controller()) {
			break
		}

		if i == len(*ipp.controllers)-1 {
			return errors.Errorf("controller not found in kyc policy controllers, %q", ad)
		}
	}

	if err := currencystate.CheckExistsState(currency.StateKeyCurrencyDesign(it.Currency()), getStateFunc); err != nil {
		return err
	}

	return nil
}

func (ipp *RemoveControllerItemProcessor) Process(
	ctx context.Context, op base.Operation, getStateFunc base.GetStateFunc,
) ([]base.StateMergeValue, error) {
	it := ipp.item

	for i, ad := range *ipp.controllers {
		if ad.Equal(it.Controller()) {
			if i < len(*ipp.controllers)-1 {
				copy((*ipp.controllers)[i:], (*ipp.controllers)[i+1:])
			}
			*ipp.controllers = (*ipp.controllers)[:len(*ipp.controllers)-1]
			break
		}

		if i == len(*ipp.controllers)-1 {
			return nil, errors.Errorf("controller not in kyc service controllers, %s, %q", it.Contract(), it.Controller())
		}
	}

	return nil, nil
}

func (ipp *RemoveControllerItemProcessor) Close() error {
	ipp.h = nil
	ipp.sender = nil
	ipp.item = RemoveControllerItem{}
	ipp.controllers = nil

	removeControllerItemProcessorPool.Put(ipp)

	return nil
}

type RemoveControllerProcessor struct {
	*base.BaseOperationProcessor
}

func NewRemoveControllerProcessor() currencytypes.GetNewProcessor {
	return func(
		height base.Height,
		getStateFunc base.GetStateFunc,
		newPreProcessConstraintFunc base.NewOperationProcessorProcessFunc,
		newProcessConstraintFunc base.NewOperationProcessorProcessFunc,
	) (base.OperationProcessor, error) {
		e := util.StringError("failed to create new RemoveControllerProcessor")

		nopp := removeControllerProcessorPool.Get()
		opp, ok := nopp.(*RemoveControllerProcessor)
		if !ok {
			return nil, e.Wrap(errors.Errorf("expected RemoveControllerProcessor, not %T", nopp))
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

func (opp *RemoveControllerProcessor) PreProcess(
	ctx context.Context, op base.Operation, getStateFunc base.GetStateFunc,
) (context.Context, base.OperationProcessReasonError, error) {
	e := util.StringError("failed to preprocess RemoveController")

	fact, ok := op.Fact().(RemoveControllerFact)
	if !ok {
		return ctx, nil, e.Wrap(errors.Errorf("expected RemoveControllerFact, not %T", op.Fact()))
	}

	if err := fact.IsValid(nil); err != nil {
		return ctx, nil, e.Wrap(err)
	}

	if err := currencystate.CheckExistsState(currency.StateKeyAccount(fact.Sender()), getStateFunc); err != nil {
		return ctx, base.NewBaseOperationProcessReasonError("sender not found, %q: %w", fact.Sender(), err), nil
	}

	if err := currencystate.CheckNotExistsState(extensioncurrency.StateKeyContractAccount(fact.Sender()), getStateFunc); err != nil {
		return ctx, base.NewBaseOperationProcessReasonError("contract account cannot remove its controllers, %q: %w", fact.Sender(), err), nil
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
		ip := removeControllerItemProcessorPool.Get()
		ipc, ok := ip.(*RemoveControllerItemProcessor)
		if !ok {
			return nil, nil, e.Wrap(errors.Errorf("expected RemoveControllerItemProcessor, not %T", ip))
		}

		ipc.h = op.Hash()
		ipc.sender = fact.Sender()
		ipc.item = it
		ipc.controllers = controllers[kycstate.StateKeyDesign(it.Contract())]

		if err := ipc.PreProcess(ctx, op, getStateFunc); err != nil {
			return nil, base.NewBaseOperationProcessReasonError("failed to preprocess RemoveControllerItem: %w", err), nil
		}

		ipc.Close()
	}

	return ctx, nil, nil
}

func (opp *RemoveControllerProcessor) Process( // nolint:dupl
	ctx context.Context, op base.Operation, getStateFunc base.GetStateFunc) (
	[]base.StateMergeValue, base.OperationProcessReasonError, error,
) {
	e := util.StringError("failed to process RemoveController")

	fact, ok := op.Fact().(RemoveControllerFact)
	if !ok {
		return nil, nil, e.Wrap(errors.Errorf("expected RemoveControllerFact, not %T", op.Fact()))
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
		ip := removeControllerItemProcessorPool.Get()
		ipc, ok := ip.(*RemoveControllerItemProcessor)
		if !ok {
			return nil, nil, e.Wrap(errors.Errorf("expected RemoveControllerItemProcessor, not %T", ip))
		}

		ipc.h = op.Hash()
		ipc.sender = fact.Sender()
		ipc.item = it
		ipc.controllers = controllers[kycstate.StateKeyDesign(it.Contract())]

		_, err := ipc.Process(ctx, op, getStateFunc)
		if err != nil {
			return nil, base.NewBaseOperationProcessReasonError("failed to process RemoveControllerItem: %w", err), nil
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

func (opp *RemoveControllerProcessor) Close() error {
	removeControllerProcessorPool.Put(opp)

	return nil
}
