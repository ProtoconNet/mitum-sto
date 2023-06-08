package kyc

import (
	"context"
	"sync"

	extensioncurrency "github.com/ProtoconNet/mitum-currency-extension/v2/currency"
	"github.com/ProtoconNet/mitum-currency/v2/currency"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/pkg/errors"
)

var updateCustomersItemProcessorPool = sync.Pool{
	New: func() interface{} {
		return new(UpdateCustomersItemProcessor)
	},
}

var updateCustomersProcessorPool = sync.Pool{
	New: func() interface{} {
		return new(UpdateCustomersProcessor)
	},
}

func (UpdateCustomers) Process(
	ctx context.Context, getStateFunc base.GetStateFunc,
) ([]base.StateMergeValue, base.OperationProcessReasonError, error) {
	return nil, nil, nil
}

type UpdateCustomersItemProcessor struct {
	h      util.Hash
	sender base.Address
	item   UpdateCustomersItem
}

func (ipp *UpdateCustomersItemProcessor) PreProcess(
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
		policy, err := existsPolicy(it.Contract(), it.KYC(), getStateFunc)
		if err != nil {
			return err
		}

		controllers := policy.Controllers()
		if len(controllers) == 0 {
			return errors.Errorf("not contract account owner neither its controller, %s-%s", it.Contract(), it.KYC())
		}

		for i, con := range controllers {
			if con.Equal(ipp.sender) {
				break
			}

			if i == len(controllers)-1 {
				return errors.Errorf("not contract account owner neither its controller, %s-%s", it.Contract(), it.KYC())
			}
		}
	}

	st, err = existsState(StateKeyCustomer(it.Contract(), it.KYC(), it.Customer()), "key of customer status", getStateFunc)
	if err != nil {
		return err
	}

	status, err := StateCustomerValue(st)
	if err != nil {
		return err
	}

	if bool(*status) == it.Status() {
		return errors.Errorf("customer status already reflected, %s-%s-%s", it.Contract(), it.KYC(), it.Customer())
	}

	if err := checkExistsState(currency.StateKeyCurrencyDesign(it.Currency()), getStateFunc); err != nil {
		return err
	}

	return nil
}

func (ipp *UpdateCustomersItemProcessor) Process(
	ctx context.Context, op base.Operation, getStateFunc base.GetStateFunc,
) ([]base.StateMergeValue, error) {
	it := ipp.item

	v := NewStateMergeValue(
		StateKeyCustomer(it.Contract(), it.KYC(), it.Customer()),
		NewCustomerStateValue(Status(it.Status())),
	)

	sts := []base.StateMergeValue{v}

	return sts, nil
}

func (ipp *UpdateCustomersItemProcessor) Close() error {
	ipp.h = nil
	ipp.sender = nil
	ipp.item = UpdateCustomersItem{}

	updateCustomersItemProcessorPool.Put(ipp)

	return nil
}

type UpdateCustomersProcessor struct {
	*base.BaseOperationProcessor
}

func NewUpdateCustomersProcessor() extensioncurrency.GetNewProcessor {
	return func(
		height base.Height,
		getStateFunc base.GetStateFunc,
		newPreProcessConstraintFunc base.NewOperationProcessorProcessFunc,
		newProcessConstraintFunc base.NewOperationProcessorProcessFunc,
	) (base.OperationProcessor, error) {
		e := util.StringErrorFunc("failed to create new UpdateCustomersProcessor")

		nopp := updateCustomersProcessorPool.Get()
		opp, ok := nopp.(*UpdateCustomersProcessor)
		if !ok {
			return nil, e(nil, "expected UpdateCustomersProcessor, not %T", nopp)
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

func (opp *UpdateCustomersProcessor) PreProcess(
	ctx context.Context, op base.Operation, getStateFunc base.GetStateFunc,
) (context.Context, base.OperationProcessReasonError, error) {
	e := util.StringErrorFunc("failed to preprocess UpdateCustomers")

	fact, ok := op.Fact().(UpdateCustomersFact)
	if !ok {
		return ctx, nil, e(nil, "expected UpdateCustomersFact, not %T", op.Fact())
	}

	if err := fact.IsValid(nil); err != nil {
		return ctx, nil, e(err, "")
	}

	if err := checkExistsState(currency.StateKeyAccount(fact.Sender()), getStateFunc); err != nil {
		return ctx, base.NewBaseOperationProcessReasonError("sender not found, %q: %w", fact.Sender(), err), nil
	}

	if err := checkNotExistsState(extensioncurrency.StateKeyContractAccount(fact.Sender()), getStateFunc); err != nil {
		return ctx, base.NewBaseOperationProcessReasonError("contract account cannot update customer status, %q: %w", fact.Sender(), err), nil
	}

	if err := checkFactSignsByState(fact.sender, op.Signs(), getStateFunc); err != nil {
		return ctx, base.NewBaseOperationProcessReasonError("invalid signing: %w", err), nil
	}

	for _, it := range fact.Items() {
		ip := updateCustomersItemProcessorPool.Get()
		ipc, ok := ip.(*UpdateCustomersItemProcessor)
		if !ok {
			return nil, nil, e(nil, "expected UpdateCustomersItemProcessor, not %T", ip)
		}

		ipc.h = op.Hash()
		ipc.sender = fact.Sender()
		ipc.item = it

		if err := ipc.PreProcess(ctx, op, getStateFunc); err != nil {
			return nil, base.NewBaseOperationProcessReasonError("failed to preprocess UpdateCustomersItem: %w", err), nil
		}

		ipc.Close()
	}

	return ctx, nil, nil
}

func (opp *UpdateCustomersProcessor) Process( // nolint:dupl
	ctx context.Context, op base.Operation, getStateFunc base.GetStateFunc) (
	[]base.StateMergeValue, base.OperationProcessReasonError, error,
) {
	e := util.StringErrorFunc("failed to process UpdateCustomers")

	fact, ok := op.Fact().(UpdateCustomersFact)
	if !ok {
		return nil, nil, e(nil, "expected UpdateCustomersFact, not %T", op.Fact())
	}

	var sts []base.StateMergeValue // nolint:prealloc

	for _, it := range fact.Items() {
		ip := updateCustomersItemProcessorPool.Get()
		ipc, ok := ip.(*UpdateCustomersItemProcessor)
		if !ok {
			return nil, nil, e(nil, "expected UpdateCustomersItemProcessor, not %T", ip)
		}

		ipc.h = op.Hash()
		ipc.sender = fact.Sender()
		ipc.item = it

		st, err := ipc.Process(ctx, op, getStateFunc)
		if err != nil {
			return nil, base.NewBaseOperationProcessReasonError("failed to process UpdateCustomersItem: %w", err), nil
		}

		sts = append(sts, st...)

		ipc.Close()
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
		v, ok := sb[i].Value().(currency.BalanceStateValue)
		if !ok {
			return nil, nil, e(nil, "expected BalanceStateValue, not %T", sb[i].Value())
		}
		stv := currency.NewBalanceStateValue(v.Amount.WithBig(v.Amount.Big().Sub(required[i][0])))
		sts = append(sts, currency.NewBalanceStateMergeValue(sb[i].Key(), stv))
	}

	return sts, nil, nil
}

func (opp *UpdateCustomersProcessor) Close() error {
	updateCustomersProcessorPool.Put(opp)

	return nil
}
