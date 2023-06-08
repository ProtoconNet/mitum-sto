package kyc

import (
	"context"
	"sync"

	currencybase "github.com/ProtoconNet/mitum-currency/v3/base"
	currencyoperation "github.com/ProtoconNet/mitum-currency/v3/operation/currency"
	types "github.com/ProtoconNet/mitum-currency/v3/operation/type"
	currency "github.com/ProtoconNet/mitum-currency/v3/state/currency"
	extensioncurrency "github.com/ProtoconNet/mitum-currency/v3/state/extension"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/pkg/errors"
)

var removeControllersItemProcessorPool = sync.Pool{
	New: func() interface{} {
		return new(RemoveControllersItemProcessor)
	},
}

var removeControllersProcessorPool = sync.Pool{
	New: func() interface{} {
		return new(RemoveControllersProcessor)
	},
}

func (RemoveControllers) Process(
	ctx context.Context, getStateFunc base.GetStateFunc,
) ([]base.StateMergeValue, base.OperationProcessReasonError, error) {
	return nil, nil, nil
}

type RemoveControllersItemProcessor struct {
	h           util.Hash
	sender      base.Address
	item        RemoveControllersItem
	controllers *[]base.Address
}

func (ipp *RemoveControllersItemProcessor) PreProcess(
	ctx context.Context, op base.Operation, getStateFunc base.GetStateFunc,
) error {
	it := ipp.item

	st, err := existsState(extensioncurrency.StateKeyContractAccount(it.Contract()), "key of contract account", getStateFunc)
	if err != nil {
		return err
	}

	ca, err := extensioncurrency.StateContractAccountValue(st)
	if err != nil {
		return err
	}

	if !ca.Owner().Equal(ipp.sender) {
		return errors.Errorf("not contract account owner, %q", ca)
	}

	if len(*ipp.controllers) == 0 {
		return errors.Errorf("empty controllers, %s-%s", it.Contract(), it.KYC())
	}

	for i, ad := range *ipp.controllers {
		if ad.Equal(it.Controller()) {
			break
		}

		if i == len(*ipp.controllers)-1 {
			return errors.Errorf("controller not found in kyc policy controllers, %q", ad)
		}
	}

	if err := checkExistsState(currency.StateKeyCurrencyDesign(it.Currency()), getStateFunc); err != nil {
		return err
	}

	return nil
}

func (ipp *RemoveControllersItemProcessor) Process(
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
			return nil, errors.Errorf("controller not in kyc service controllers, %s-%s, %q", it.Contract(), it.KYC(), it.Controller())
		}
	}

	return nil, nil
}

func (ipp *RemoveControllersItemProcessor) Close() error {
	ipp.h = nil
	ipp.sender = nil
	ipp.item = RemoveControllersItem{}
	ipp.controllers = nil

	removeControllersItemProcessorPool.Put(ipp)

	return nil
}

type RemoveControllersProcessor struct {
	*base.BaseOperationProcessor
}

func NewRemoveControllersProcessor() types.GetNewProcessor {
	return func(
		height base.Height,
		getStateFunc base.GetStateFunc,
		newPreProcessConstraintFunc base.NewOperationProcessorProcessFunc,
		newProcessConstraintFunc base.NewOperationProcessorProcessFunc,
	) (base.OperationProcessor, error) {
		e := util.StringErrorFunc("failed to create new RemoveControllersProcessor")

		nopp := removeControllersProcessorPool.Get()
		opp, ok := nopp.(*RemoveControllersProcessor)
		if !ok {
			return nil, e(nil, "expected RemoveControllersProcessor, not %T", nopp)
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

func (opp *RemoveControllersProcessor) PreProcess(
	ctx context.Context, op base.Operation, getStateFunc base.GetStateFunc,
) (context.Context, base.OperationProcessReasonError, error) {
	e := util.StringErrorFunc("failed to preprocess RemoveControllers")

	fact, ok := op.Fact().(RemoveControllersFact)
	if !ok {
		return ctx, nil, e(nil, "expected RemoveControllersFact, not %T", op.Fact())
	}

	if err := fact.IsValid(nil); err != nil {
		return ctx, nil, e(err, "")
	}

	if err := checkExistsState(currency.StateKeyAccount(fact.Sender()), getStateFunc); err != nil {
		return ctx, base.NewBaseOperationProcessReasonError("sender not found, %q: %w", fact.Sender(), err), nil
	}

	if err := checkNotExistsState(extensioncurrency.StateKeyContractAccount(fact.Sender()), getStateFunc); err != nil {
		return ctx, base.NewBaseOperationProcessReasonError("contract account cannot remove its controllers, %q: %w", fact.Sender(), err), nil
	}

	if err := checkFactSignsByState(fact.sender, op.Signs(), getStateFunc); err != nil {
		return ctx, base.NewBaseOperationProcessReasonError("invalid signing: %w", err), nil
	}

	controllers := map[string]*[]base.Address{}

	for _, it := range fact.Items() {
		policy, err := existsPolicy(it.Contract(), it.KYC(), getStateFunc)
		if err != nil {
			return nil, base.NewBaseOperationProcessReasonError("failed to get kyc policy, %s-%s: %w", it.Contract(), it.KYC(), err), nil
		}
		cons := policy.Controllers()
		controllers[StateKeyDesign(it.Contract(), it.KYC())] = &cons
	}

	for _, it := range fact.Items() {
		ip := removeControllersItemProcessorPool.Get()
		ipc, ok := ip.(*RemoveControllersItemProcessor)
		if !ok {
			return nil, nil, e(nil, "expected RemoveControllersItemProcessor, not %T", ip)
		}

		ipc.h = op.Hash()
		ipc.sender = fact.Sender()
		ipc.item = it
		ipc.controllers = controllers[StateKeyDesign(it.Contract(), it.KYC())]

		if err := ipc.PreProcess(ctx, op, getStateFunc); err != nil {
			return nil, base.NewBaseOperationProcessReasonError("failed to preprocess RemoveControllersItem: %w", err), nil
		}

		ipc.Close()
	}

	return ctx, nil, nil
}

func (opp *RemoveControllersProcessor) Process( // nolint:dupl
	ctx context.Context, op base.Operation, getStateFunc base.GetStateFunc) (
	[]base.StateMergeValue, base.OperationProcessReasonError, error,
) {
	e := util.StringErrorFunc("failed to process RemoveControllers")

	fact, ok := op.Fact().(RemoveControllersFact)
	if !ok {
		return nil, nil, e(nil, "expected RemoveControllersFact, not %T", op.Fact())
	}

	var sts []base.StateMergeValue // nolint:prealloc

	controllers := map[string]map[string]*[]base.Address{}

	for _, it := range fact.Items() {
		policy, err := existsPolicy(it.Contract(), it.KYC(), getStateFunc)
		if err != nil {
			return nil, base.NewBaseOperationProcessReasonError("failed to get kyc policy, %s-%s: %w", it.Contract(), it.KYC(), err), nil
		}
		cons := policy.Controllers()
		controllers[StateKeyDesign(it.Contract(), it.KYC())][it.KYC().String()] = &cons
	}

	for _, it := range fact.Items() {
		ip := removeControllersItemProcessorPool.Get()
		ipc, ok := ip.(*RemoveControllersItemProcessor)
		if !ok {
			return nil, nil, e(nil, "expected RemoveControllersItemProcessor, not %T", ip)
		}

		ipc.h = op.Hash()
		ipc.sender = fact.Sender()
		ipc.item = it
		ipc.controllers = controllers[StateKeyDesign(it.Contract(), it.KYC())][it.KYC().String()]

		_, err := ipc.Process(ctx, op, getStateFunc)
		if err != nil {
			return nil, base.NewBaseOperationProcessReasonError("failed to process RemoveControllersItem: %w", err), nil
		}

		ipc.Close()
	}

	for k, m := range controllers {
		for id, cons := range m {
			policy := NewPolicy(*cons)
			design := NewDesign(currencybase.ContractID(id), policy)
			if err := design.IsValid(nil); err != nil {
				return nil, base.NewBaseOperationProcessReasonError("invalid design, %s: %w", k, err), nil
			}

			sts = append(sts, NewStateMergeValue(
				k,
				NewDesignStateValue(design),
			))
		}
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
			return nil, nil, e(nil, "expected BalanceStateValue, not %T", sb[i].Value())
		}
		stv := currency.NewBalanceStateValue(v.Amount.WithBig(v.Amount.Big().Sub(required[i][0])))
		sts = append(sts, NewStateMergeValue(sb[i].Key(), stv))
	}

	return sts, nil, nil
}

func (opp *RemoveControllersProcessor) Close() error {
	removeControllersProcessorPool.Put(opp)

	return nil
}
