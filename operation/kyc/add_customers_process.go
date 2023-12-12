package kyc

import (
	"context"
	"github.com/ProtoconNet/mitum-currency/v3/common"
	"sync"

	"github.com/ProtoconNet/mitum-currency/v3/operation/currency"
	"github.com/ProtoconNet/mitum-currency/v3/state"
	stcurrency "github.com/ProtoconNet/mitum-currency/v3/state/currency"
	stextension "github.com/ProtoconNet/mitum-currency/v3/state/extension"
	"github.com/ProtoconNet/mitum-currency/v3/types"
	stkyc "github.com/ProtoconNet/mitum-sto/state/kyc"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/pkg/errors"
)

var addCustomerItemProcessorPool = sync.Pool{
	New: func() interface{} {
		return new(AddCustomerItemProcessor)
	},
}

var addCustomerProcessorPool = sync.Pool{
	New: func() interface{} {
		return new(AddCustomerProcessor)
	},
}

func (AddCustomer) Process(
	ctx context.Context, getStateFunc base.GetStateFunc,
) ([]base.StateMergeValue, base.OperationProcessReasonError, error) {
	return nil, nil, nil
}

type AddCustomerItemProcessor struct {
	h      util.Hash
	sender base.Address
	item   AddCustomerItem
}

func (ipp *AddCustomerItemProcessor) PreProcess(
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
		policy, err := stkyc.ExistsPolicy(it.Contract(), getStateFunc)
		if err != nil {
			return err
		}

		controllers := policy.Controllers()
		if len(controllers) == 0 {
			return errors.Errorf("not contract account owner neither its controller, %s", it.Contract())
		}

		for i, con := range controllers {
			if con.Equal(ipp.sender) {
				break
			}

			if i == len(controllers)-1 {
				return errors.Errorf("not contract account owner neither its controller, %s-%s", it.Contract())
			}
		}
	}

	if err := state.CheckNotExistsState(stkyc.StateKeyCustomer(it.Contract(), it.Customer()), getStateFunc); err != nil {
		return err
	}

	if err := state.CheckExistsState(stcurrency.StateKeyCurrencyDesign(it.Currency()), getStateFunc); err != nil {
		return err
	}

	return nil
}

func (ipp *AddCustomerItemProcessor) Process(
	ctx context.Context, op base.Operation, getStateFunc base.GetStateFunc,
) ([]base.StateMergeValue, error) {
	it := ipp.item

	v := state.NewStateMergeValue(
		stkyc.StateKeyCustomer(it.Contract(), it.Customer()),
		stkyc.NewCustomerStateValue(stkyc.Status(it.Status())),
	)

	sts := []base.StateMergeValue{v}

	return sts, nil
}

func (ipp *AddCustomerItemProcessor) Close() {
	ipp.h = nil
	ipp.sender = nil
	ipp.item = AddCustomerItem{}

	addCustomerItemProcessorPool.Put(ipp)
}

type AddCustomerProcessor struct {
	*base.BaseOperationProcessor
}

func NewAddCustomerProcessor() types.GetNewProcessor {
	return func(
		height base.Height,
		getStateFunc base.GetStateFunc,
		newPreProcessConstraintFunc base.NewOperationProcessorProcessFunc,
		newProcessConstraintFunc base.NewOperationProcessorProcessFunc,
	) (base.OperationProcessor, error) {
		e := util.StringError("failed to create new AddCustomerProcessor")

		nopp := addCustomerProcessorPool.Get()
		opp, ok := nopp.(*AddCustomerProcessor)
		if !ok {
			return nil, e.Wrap(errors.Errorf("expected AddCustomerProcessor, not %T", nopp))
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

func (opp *AddCustomerProcessor) PreProcess(
	ctx context.Context, op base.Operation, getStateFunc base.GetStateFunc,
) (context.Context, base.OperationProcessReasonError, error) {
	e := util.StringError("failed to preprocess AddCustomer")

	fact, ok := op.Fact().(AddCustomerFact)
	if !ok {
		return ctx, nil, e.Wrap(errors.Errorf("expected AddCustomerFact, not %T", op.Fact()))
	}

	if err := fact.IsValid(nil); err != nil {
		return ctx, nil, e.Wrap(err)
	}

	if err := state.CheckExistsState(stcurrency.StateKeyAccount(fact.Sender()), getStateFunc); err != nil {
		return ctx, base.NewBaseOperationProcessReasonError("sender not found, %q: %w", fact.Sender(), err), nil
	}

	if err := state.CheckNotExistsState(stextension.StateKeyContractAccount(fact.Sender()), getStateFunc); err != nil {
		return ctx, base.NewBaseOperationProcessReasonError("contract account cannot add customer status, %q: %w", fact.Sender(), err), nil
	}

	if err := state.CheckFactSignsByState(fact.sender, op.Signs(), getStateFunc); err != nil {
		return ctx, base.NewBaseOperationProcessReasonError("invalid signing: %w", err), nil
	}

	for _, it := range fact.Items() {
		ip := addCustomerItemProcessorPool.Get()
		ipc, ok := ip.(*AddCustomerItemProcessor)
		if !ok {
			return nil, nil, e.Wrap(errors.Errorf("expected AddCustomerItemProcessor, not %T", ip))
		}

		ipc.h = op.Hash()
		ipc.sender = fact.Sender()
		ipc.item = it

		if err := ipc.PreProcess(ctx, op, getStateFunc); err != nil {
			return nil, base.NewBaseOperationProcessReasonError("failed to preprocess AddCustomerItem: %w", err), nil
		}

		ipc.Close()
	}

	return ctx, nil, nil
}

func (opp *AddCustomerProcessor) Process( // nolint:dupl
	ctx context.Context, op base.Operation, getStateFunc base.GetStateFunc) (
	[]base.StateMergeValue, base.OperationProcessReasonError, error,
) {
	e := util.StringError("failed to process AddCustomer")

	fact, ok := op.Fact().(AddCustomerFact)
	if !ok {
		return nil, nil, e.Wrap(errors.Errorf("expected AddCustomerFact, not %T", op.Fact()))
	}

	var sts []base.StateMergeValue // nolint:prealloc

	for _, it := range fact.Items() {
		ip := addCustomerItemProcessorPool.Get()
		ipc, ok := ip.(*AddCustomerItemProcessor)
		if !ok {
			return nil, nil, e.Wrap(errors.Errorf("expected AddCustomerItemProcessor, not %T", ip))
		}

		ipc.h = op.Hash()
		ipc.sender = fact.Sender()
		ipc.item = it

		st, err := ipc.Process(ctx, op, getStateFunc)
		if err != nil {
			return nil, base.NewBaseOperationProcessReasonError("failed to process AddCustomerItem: %w", err), nil
		}

		sts = append(sts, st...)

		ipc.Close()
	}

	items := make([]KYCItem, len(fact.Items()))
	for i := range fact.Items() {
		items[i] = fact.Items()[i]
	}

	feeReceiveBalSts, required, err := calculateKYCItemsFee(getStateFunc, items)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError("failed to calculate fee; %w", err), nil
	}
	sb, err := currency.CheckEnoughBalance(fact.sender, required, getStateFunc)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError("failed to check enough balance; %w", err), nil
	}

	for cid := range sb {
		v, ok := sb[cid].Value().(stcurrency.BalanceStateValue)
		if !ok {
			return nil, nil, e.Errorf("expected BalanceStateValue, not %T", sb[cid].Value())
		}

		if sb[cid].Key() != feeReceiveBalSts[cid].Key() {
			stmv := common.NewBaseStateMergeValue(
				sb[cid].Key(),
				stcurrency.NewDeductBalanceStateValue(v.Amount.WithBig(required[cid][1])),
				func(height base.Height, st base.State) base.StateValueMerger {
					return stcurrency.NewBalanceStateValueMerger(height, sb[cid].Key(), cid, st)
				},
			)

			r, ok := feeReceiveBalSts[cid].Value().(stcurrency.BalanceStateValue)
			if !ok {
				return nil, base.NewBaseOperationProcessReasonError("expected %T, not %T", stcurrency.BalanceStateValue{}, feeReceiveBalSts[cid].Value()), nil
			}
			sts = append(
				sts,
				common.NewBaseStateMergeValue(
					feeReceiveBalSts[cid].Key(),
					stcurrency.NewAddBalanceStateValue(r.Amount.WithBig(required[cid][1])),
					func(height base.Height, st base.State) base.StateValueMerger {
						return stcurrency.NewBalanceStateValueMerger(height, feeReceiveBalSts[cid].Key(), cid, st)
					},
				),
			)

			sts = append(sts, stmv)
		}
	}

	return sts, nil, nil
}

func (opp *AddCustomerProcessor) Close() error {
	addCustomerProcessorPool.Put(opp)

	return nil
}
