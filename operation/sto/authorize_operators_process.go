package sto

import (
	"context"
	"sync"

	currencyoperation "github.com/ProtoconNet/mitum-currency/v3/operation/currency"
	currencystate "github.com/ProtoconNet/mitum-currency/v3/state"
	currency "github.com/ProtoconNet/mitum-currency/v3/state/currency"
	extensioncurrency "github.com/ProtoconNet/mitum-currency/v3/state/extension"
	currencytypes "github.com/ProtoconNet/mitum-currency/v3/types"
	stostate "github.com/ProtoconNet/mitum-sto/state/sto"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/pkg/errors"
)

var authorizeOperatorItemProcessorPool = sync.Pool{
	New: func() interface{} {
		return new(AuthorizeOperatorItemProcessor)
	},
}

var authorizeOperatorProcessorPool = sync.Pool{
	New: func() interface{} {
		return new(AuthorizeOperatorProcessor)
	},
}

func (AuthorizeOperator) Process(
	ctx context.Context, getStateFunc base.GetStateFunc,
) ([]base.StateMergeValue, base.OperationProcessReasonError, error) {
	return nil, nil, nil
}

type AuthorizeOperatorItemProcessor struct {
	h            util.Hash
	sender       base.Address
	item         AuthorizeOperatorItem
	operators    *[]base.Address
	tokenHolders *[]base.Address
}

func (ipp *AuthorizeOperatorItemProcessor) PreProcess(
	ctx context.Context, op base.Operation, getStateFunc base.GetStateFunc,
) error {
	it := ipp.item

	if err := currencystate.CheckExistsState(extensioncurrency.StateKeyContractAccount(it.Contract()), getStateFunc); err != nil {
		return err
	}

	if err := currencystate.CheckExistsState(stostate.StateKeyDesign(it.Contract()), getStateFunc); err != nil {
		return err
	}

	if err := currencystate.CheckExistsState(stostate.StateKeyPartitionBalance(it.Contract(), it.Partition()), getStateFunc); err != nil {
		return err
	}

	partitions, err := stostate.ExistsTokenHolderPartitions(it.Contract(), ipp.sender, getStateFunc)
	if err != nil {
		return err
	}

	if len(partitions) == 0 {
		return errors.Errorf("empty tokenholder partitions, %s-%s", it.Contract(), ipp.sender)
	}

	for i, p := range partitions {
		if p == it.Partition() {
			break
		}

		if i == len(partitions)-1 {
			return errors.Errorf("partition not in tokenholder partitions, %s-%s, %s", it.Contract(), ipp.sender, it.Partition())
		}
	}

	for _, ad := range *ipp.operators {
		if ad.Equal(it.Operator()) {
			return errors.Errorf("operator is already in tokenholder operators, %q", ad)
		}
	}

	for _, ad := range *ipp.tokenHolders {
		if ad.Equal(ipp.sender) {
			return errors.Errorf("sender is already in operator tokenholders, %q", ad)
		}
	}

	if err := currencystate.CheckExistsState(currency.StateKeyCurrencyDesign(it.Currency()), getStateFunc); err != nil {
		return err
	}

	return nil
}

func (ipp *AuthorizeOperatorItemProcessor) Process(
	ctx context.Context, op base.Operation, getStateFunc base.GetStateFunc,
) ([]base.StateMergeValue, error) {
	sts := make([]base.StateMergeValue, 1)

	it := ipp.item

	*ipp.operators = append(*ipp.operators, it.Operator())
	holders := append(*ipp.tokenHolders, ipp.sender)

	sts[0] = currencystate.NewStateMergeValue(
		stostate.StateKeyOperatorTokenHolders(it.Contract(), it.Operator(), it.Partition()),
		stostate.NewOperatorTokenHoldersStateValue(holders),
	)

	return sts, nil
}

func (ipp *AuthorizeOperatorItemProcessor) Close() error {
	ipp.h = nil
	ipp.sender = nil
	ipp.item = AuthorizeOperatorItem{}
	ipp.operators = nil
	ipp.tokenHolders = nil

	authorizeOperatorItemProcessorPool.Put(ipp)

	return nil
}

type AuthorizeOperatorProcessor struct {
	*base.BaseOperationProcessor
}

func NewAuthorizeOperatorsProcessor() currencytypes.GetNewProcessor {
	return func(
		height base.Height,
		getStateFunc base.GetStateFunc,
		newPreProcessConstraintFunc base.NewOperationProcessorProcessFunc,
		newProcessConstraintFunc base.NewOperationProcessorProcessFunc,
	) (base.OperationProcessor, error) {
		e := util.StringError("failed to create new AuthorizeOperatorProcessor")

		nopp := authorizeOperatorProcessorPool.Get()
		opp, ok := nopp.(*AuthorizeOperatorProcessor)
		if !ok {
			return nil, e.Wrap(errors.Errorf("expected AuthorizeOperatorProcessor, not %T", nopp))
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

func (opp *AuthorizeOperatorProcessor) PreProcess(
	ctx context.Context, op base.Operation, getStateFunc base.GetStateFunc,
) (context.Context, base.OperationProcessReasonError, error) {
	e := util.StringError("failed to preprocess AuthorizeOperator")

	fact, ok := op.Fact().(AuthorizeOperatorFact)
	if !ok {
		return ctx, nil, e.Wrap(errors.Errorf("expected AuthorizeOperatorFact, not %T", op.Fact()))
	}

	if err := fact.IsValid(nil); err != nil {
		return ctx, nil, e.Wrap(err)
	}

	if err := currencystate.CheckExistsState(currency.StateKeyAccount(fact.Sender()), getStateFunc); err != nil {
		return ctx, base.NewBaseOperationProcessReasonError("sender not found, %q: %w", fact.Sender(), err), nil
	}

	if err := currencystate.CheckNotExistsState(extensioncurrency.StateKeyContractAccount(fact.Sender()), getStateFunc); err != nil {
		return ctx, base.NewBaseOperationProcessReasonError("contract account cannot set its operators, %q: %w", fact.Sender(), err), nil
	}

	if err := currencystate.CheckFactSignsByState(fact.sender, op.Signs(), getStateFunc); err != nil {
		return ctx, base.NewBaseOperationProcessReasonError("invalid signing: %w", err), nil
	}

	operators := map[string]*[]base.Address{}
	holders := map[string]*[]base.Address{}

	for _, it := range fact.Items() {
		var ops, hds []base.Address

		k := stostate.StateKeyTokenHolderPartitionOperators(it.Contract(), fact.sender, it.Partition())
		if _, found := operators[k]; !found {
			switch st, found, err := getStateFunc(k); {
			case err != nil:
				return nil, base.NewBaseOperationProcessReasonError("failed to find tokenholder partition operators, %s: %w", k, err), nil
			case found:
				ops, err = stostate.StateTokenHolderPartitionOperatorsValue(st)
				if err != nil {
					return nil, base.NewBaseOperationProcessReasonError("failed to get tokenholder partition operators, %s: %w", k, err), nil
				}
			default:
				ops = []base.Address{}
			}
			operators[k] = &ops
		}

		k = stostate.StateKeyOperatorTokenHolders(it.Contract(), it.Operator(), it.Partition())
		if _, found := holders[k]; !found {
			switch st, found, err := getStateFunc(k); {
			case err != nil:
				return nil, base.NewBaseOperationProcessReasonError("failed to find operator tokenholders, %s: %w", k, err), nil
			case found:
				hds, err = stostate.StateOperatorTokenHoldersValue(st)
				if err != nil {
					return nil, base.NewBaseOperationProcessReasonError("failed to get operator tokenholders, %s: %w", k, err), nil
				}
			default:
				hds = []base.Address{}
			}

			holders[k] = &hds
		}
	}

	for _, it := range fact.Items() {
		ip := authorizeOperatorItemProcessorPool.Get()
		ipc, ok := ip.(*AuthorizeOperatorItemProcessor)
		if !ok {
			return nil, nil, e.Wrap(errors.Errorf("expected AuthorizeOperatorItemProcessor, not %T", ip))
		}

		ipc.h = op.Hash()
		ipc.sender = fact.Sender()
		ipc.item = it
		ipc.operators = operators[stostate.StateKeyTokenHolderPartitionOperators(it.Contract(), fact.sender, it.Partition())]
		ipc.tokenHolders = holders[stostate.StateKeyOperatorTokenHolders(it.Contract(), it.Operator(), it.Partition())]

		if err := ipc.PreProcess(ctx, op, getStateFunc); err != nil {
			return nil, base.NewBaseOperationProcessReasonError("fail to preprocess AuthorizeOperatorItem: %w", err), nil
		}

		ipc.Close()
	}

	return ctx, nil, nil
}

func (opp *AuthorizeOperatorProcessor) Process( // nolint:dupl
	ctx context.Context, op base.Operation, getStateFunc base.GetStateFunc) (
	[]base.StateMergeValue, base.OperationProcessReasonError, error,
) {
	e := util.StringError("failed to process AuthorizeOperator")

	fact, ok := op.Fact().(AuthorizeOperatorFact)
	if !ok {
		return nil, nil, e.Wrap(errors.Errorf("expected AuthorizeOperatorFact, not %T", op.Fact()))
	}

	var sts []base.StateMergeValue // nolint:prealloc

	operators := map[string]*[]base.Address{}
	holders := map[string]*[]base.Address{}

	for _, it := range fact.Items() {
		var ops, hds []base.Address

		k := stostate.StateKeyTokenHolderPartitionOperators(it.Contract(), fact.sender, it.Partition())
		if _, found := operators[k]; !found {
			switch st, found, err := getStateFunc(k); {
			case err != nil:
				return nil, base.NewBaseOperationProcessReasonError("failed to find tokenholder partition operators, %s: %w", k, err), nil
			case found:
				ops, err = stostate.StateTokenHolderPartitionOperatorsValue(st)
				if err != nil {
					return nil, base.NewBaseOperationProcessReasonError("failed to get tokenholder partition operators, %s: %w", k, err), nil
				}
			default:
				ops = []base.Address{}
			}
			operators[k] = &ops
		}

		k = stostate.StateKeyOperatorTokenHolders(it.Contract(), it.Operator(), it.Partition())
		if _, found := holders[k]; !found {
			switch st, found, err := getStateFunc(k); {
			case err != nil:
				return nil, base.NewBaseOperationProcessReasonError("failed to find operator tokenholders, %s: %w", k, err), nil
			case found:
				hds, err = stostate.StateOperatorTokenHoldersValue(st)
				if err != nil {
					return nil, base.NewBaseOperationProcessReasonError("failed to get operator tokenholders, %s: %w", k, err), nil
				}
			default:
				hds = []base.Address{}
			}

			holders[k] = &hds
		}
	}

	ipcs := make([]*AuthorizeOperatorItemProcessor, len(fact.items))
	for i, it := range fact.Items() {
		ip := authorizeOperatorItemProcessorPool.Get()
		ipc, ok := ip.(*AuthorizeOperatorItemProcessor)
		if !ok {
			return nil, nil, e.Wrap(errors.Errorf("expected AuthorizeOperatorItemProcessor, not %T", ip))
		}

		ipc.h = op.Hash()
		ipc.sender = fact.Sender()
		ipc.item = it
		ipc.operators = operators[stostate.StateKeyTokenHolderPartitionOperators(it.Contract(), fact.sender, it.Partition())]
		ipc.tokenHolders = holders[stostate.StateKeyOperatorTokenHolders(it.Contract(), it.Operator(), it.Partition())]

		s, err := ipc.Process(ctx, op, getStateFunc)
		if err != nil {
			return nil, base.NewBaseOperationProcessReasonError("failed to process AuthorizeOperatorItem: %w", err), nil
		}
		sts = append(sts, s...)

		ipcs[i] = ipc
	}

	for k, v := range operators {
		sts = append(sts, currencystate.NewStateMergeValue(
			k,
			stostate.NewTokenHolderPartitionOperatorsStateValue(*v),
		))
	}

	for _, ipc := range ipcs {
		ipc.Close()
	}

	fitems := fact.Items()
	items := make([]STOItem, len(fitems))
	for i := range fact.Items() {
		items[i] = fitems[i]
	}

	required, err := calculateSTOItemsFee(getStateFunc, items)
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

func (opp *AuthorizeOperatorProcessor) Close() error {
	authorizeOperatorProcessorPool.Put(opp)

	return nil
}
