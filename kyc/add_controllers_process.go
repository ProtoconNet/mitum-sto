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

	for _, ad := range *ipp.controllers {
		if ad.Equal(it.Controller()) {
			return errors.Errorf("controller is already in kyc policy controllers, %q", ad)
		}
	}

	if err := checkExistsState(currency.StateKeyCurrencyDesign(it.Currency()), getStateFunc); err != nil {
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

func NewAddControllersProcessor() types.GetNewProcessor {
	return func(
		height base.Height,
		getStateFunc base.GetStateFunc,
		newPreProcessConstraintFunc base.NewOperationProcessorProcessFunc,
		newProcessConstraintFunc base.NewOperationProcessorProcessFunc,
	) (base.OperationProcessor, error) {
		e := util.StringErrorFunc("failed to create new AddControllersProcessor")

		nopp := addControllersProcessorPool.Get()
		opp, ok := nopp.(*AddControllersProcessor)
		if !ok {
			return nil, e(nil, "expected AddControllersProcessor, not %T", nopp)
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

func (opp *AddControllersProcessor) PreProcess(
	ctx context.Context, op base.Operation, getStateFunc base.GetStateFunc,
) (context.Context, base.OperationProcessReasonError, error) {
	e := util.StringErrorFunc("failed to preprocess AddControllers")

	fact, ok := op.Fact().(AddControllersFact)
	if !ok {
		return ctx, nil, e(nil, "expected AddControllersFact, not %T", op.Fact())
	}

	if err := fact.IsValid(nil); err != nil {
		return ctx, nil, e(err, "")
	}

	if err := checkExistsState(currency.StateKeyAccount(fact.Sender()), getStateFunc); err != nil {
		return ctx, base.NewBaseOperationProcessReasonError("sender not found, %q: %w", fact.Sender(), err), nil
	}

	if err := checkNotExistsState(extensioncurrency.StateKeyContractAccount(fact.Sender()), getStateFunc); err != nil {
		return ctx, base.NewBaseOperationProcessReasonError("contract account cannot set its controllers, %q: %w", fact.Sender(), err), nil
	}

	if err := checkFactSignsByState(fact.sender, op.Signs(), getStateFunc); err != nil {
		return ctx, base.NewBaseOperationProcessReasonError("invalid signing: %w", err), nil
	}

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
		ip := addControllersItemProcessorPool.Get()
		ipc, ok := ip.(*AddControllersItemProcessor)
		if !ok {
			return nil, nil, e(nil, "expected AddControllersItemProcessor, not %T", ip)
		}

		ipc.h = op.Hash()
		ipc.sender = fact.Sender()
		ipc.item = it
		ipc.controllers = controllers[StateKeyDesign(it.Contract(), it.KYC())][it.KYC().String()]

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
	e := util.StringErrorFunc("failed to process AddControllers")

	fact, ok := op.Fact().(AddControllersFact)
	if !ok {
		return nil, nil, e(nil, "expected AddControllersFact, not %T", op.Fact())
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
		ip := addControllersItemProcessorPool.Get()
		ipc, ok := ip.(*AddControllersItemProcessor)
		if !ok {
			return nil, nil, e(nil, "expected AddControllersItemProcessor, not %T", ip)
		}

		ipc.h = op.Hash()
		ipc.sender = fact.Sender()
		ipc.item = it
		ipc.controllers = controllers[StateKeyDesign(it.Contract(), it.KYC())][it.KYC().String()]

		_, err := ipc.Process(ctx, op, getStateFunc)
		if err != nil {
			return nil, base.NewBaseOperationProcessReasonError("failed to process AddControllersItem: %w", err), nil
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

func (opp *AddControllersProcessor) Close() error {
	addControllersProcessorPool.Put(opp)

	return nil
}

func calculateKYCItemsFee(getStateFunc base.GetStateFunc, items []KYCItem) (map[currencybase.CurrencyID][2]currencybase.Big, error) {
	required := map[currencybase.CurrencyID][2]currencybase.Big{}

	for _, item := range items {
		rq := [2]currencybase.Big{currencybase.ZeroBig, currencybase.ZeroBig}

		if k, found := required[item.Currency()]; found {
			rq = k
		}

		policy, err := existsCurrencyPolicy(item.Currency(), getStateFunc)
		if err != nil {
			return nil, err
		}

		switch k, err := policy.Feeer().Fee(currencybase.ZeroBig); {
		case err != nil:
			return nil, err
		case !k.OverZero():
			required[item.Currency()] = [2]currencybase.Big{rq[0], rq[1]}
		default:
			required[item.Currency()] = [2]currencybase.Big{rq[0].Add(k), rq[1].Add(k)}
		}

	}

	return required, nil

}
