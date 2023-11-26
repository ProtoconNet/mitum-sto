package sto

import (
	"context"
	"sync"

	"github.com/ProtoconNet/mitum-currency/v3/operation/currency"
	crcstate "github.com/ProtoconNet/mitum-currency/v3/state"
	stcurrency "github.com/ProtoconNet/mitum-currency/v3/state/currency"
	stextension "github.com/ProtoconNet/mitum-currency/v3/state/extension"
	crctypes "github.com/ProtoconNet/mitum-currency/v3/types"
	ststo "github.com/ProtoconNet/mitum-sto/state/sto"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/pkg/errors"
)

var revokeOperatorProcessorPool = sync.Pool{
	New: func() interface{} {
		return new(RevokeOperatorProcessor)
	},
}

func (RevokeOperator) Process(
	ctx context.Context, getStateFunc base.GetStateFunc,
) ([]base.StateMergeValue, base.OperationProcessReasonError, error) {
	return nil, nil, nil
}

type RevokeOperatorProcessor struct {
	*base.BaseOperationProcessor
}

func NewRevokeOperatorProcessor() crctypes.GetNewProcessor {
	return func(
		height base.Height,
		getStateFunc base.GetStateFunc,
		newPreProcessConstraintFunc base.NewOperationProcessorProcessFunc,
		newProcessConstraintFunc base.NewOperationProcessorProcessFunc,
	) (base.OperationProcessor, error) {
		e := util.StringError("failed to create new RevokeOperatorProcessor")

		nopp := revokeOperatorProcessorPool.Get()
		opp, ok := nopp.(*RevokeOperatorProcessor)
		if !ok {
			return nil, e.Wrap(errors.Errorf("expected RevokeOperatorProcessor, not %T", nopp))
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

func (opp *RevokeOperatorProcessor) PreProcess(
	ctx context.Context, op base.Operation, getStateFunc base.GetStateFunc,
) (context.Context, base.OperationProcessReasonError, error) {
	e := util.StringError("failed to preprocess RevokeOperator")

	fact, ok := op.Fact().(RevokeOperatorFact)
	if !ok {
		return ctx, nil, e.Wrap(errors.Errorf("expected RevokeOperatorFact, not %T", op.Fact()))
	}

	if err := fact.IsValid(nil); err != nil {
		return ctx, nil, e.Wrap(err)
	}

	if err := crcstate.CheckExistsState(stcurrency.StateKeyAccount(fact.Sender()), getStateFunc); err != nil {
		return ctx, base.NewBaseOperationProcessReasonError("sender not found, %q: %w", fact.Sender(), err), nil
	}

	if err := crcstate.CheckNotExistsState(stextension.StateKeyContractAccount(fact.Sender()), getStateFunc); err != nil {
		return ctx, base.NewBaseOperationProcessReasonError("contract account cannot set its operators, %q: %w", fact.Sender(), err), nil
	}

	if err := crcstate.CheckFactSignsByState(fact.sender, op.Signs(), getStateFunc); err != nil {
		return ctx, base.NewBaseOperationProcessReasonError("invalid signing: %w", err), nil
	}

	operators := map[string]*[]base.Address{}
	holders := map[string]*[]base.Address{}

	for _, it := range fact.Items() {
		var ops, hds []base.Address

		k := ststo.StateKeyTokenHolderPartitionOperators(it.Contract(), fact.sender, it.Partition())
		if _, found := operators[k]; !found {
			switch st, found, err := getStateFunc(k); {
			case err != nil:
				return nil, base.NewBaseOperationProcessReasonError("failed to find token holder partition operators, %s: %w", k, err), nil
			case found:
				ops, err = ststo.StateTokenHolderPartitionOperatorsValue(st)
				if err != nil {
					return nil, base.NewBaseOperationProcessReasonError("failed to get token holder partition operators, %s: %w", k, err), nil
				}
			default:
				return nil, base.NewBaseOperationProcessReasonError("token holder partition operators not in state, %q", k), nil
			}
			operators[k] = &ops
		}

		k = ststo.StateKeyOperatorTokenHolders(it.Contract(), it.Operator(), it.Partition())
		if _, found := holders[k]; !found {
			switch st, found, err := getStateFunc(k); {
			case err != nil:
				return nil, base.NewBaseOperationProcessReasonError("failed to find operator tokenholders, %s: %w", k, err), nil
			case found:
				hds, err = ststo.StateOperatorTokenHoldersValue(st)
				if err != nil {
					return nil, base.NewBaseOperationProcessReasonError("failed to get operator tokenholders, %s: %w", k, err), nil
				}
			default:
				return nil, base.NewBaseOperationProcessReasonError("operator tokenholders not in state, %q", k), nil
			}
			holders[k] = &hds
		}
	}

	for _, it := range fact.Items() {
		ip := revokeOperatorItemProcessorPool.Get()
		ipc, ok := ip.(*RevokeOperatorItemProcessor)
		if !ok {
			return nil, nil, e.Wrap(errors.Errorf("expected RevokeOperatorItemProcessor, not %T", ip))
		}

		ipc.h = op.Hash()
		ipc.sender = fact.Sender()
		ipc.item = it
		ipc.operators = operators[ststo.StateKeyTokenHolderPartitionOperators(it.Contract(), fact.sender, it.Partition())]
		ipc.tokenHolders = holders[ststo.StateKeyOperatorTokenHolders(it.Contract(), it.Operator(), it.Partition())]

		if err := ipc.PreProcess(ctx, op, getStateFunc); err != nil {
			return nil, base.NewBaseOperationProcessReasonError("fail to preprocess RevokeOperatorItem: %w", err), nil
		}

		ipc.Close()
	}

	return ctx, nil, nil
}

func (opp *RevokeOperatorProcessor) Process( // nolint:dupl
	ctx context.Context, op base.Operation, getStateFunc base.GetStateFunc) (
	[]base.StateMergeValue, base.OperationProcessReasonError, error,
) {
	e := util.StringError("failed to process RevokeOperator")

	fact, ok := op.Fact().(RevokeOperatorFact)
	if !ok {
		return nil, nil, e.Wrap(errors.Errorf("expected RevokeOperatorFact, not %T", op.Fact()))
	}

	var sts []base.StateMergeValue // nolint:prealloc

	operators := map[string]*[]base.Address{}
	holders := map[string]*[]base.Address{}

	for _, it := range fact.Items() {
		var ops, hds []base.Address

		k := ststo.StateKeyTokenHolderPartitionOperators(it.Contract(), fact.sender, it.Partition())
		if _, found := operators[k]; !found {
			switch st, found, err := getStateFunc(k); {
			case err != nil:
				return nil, base.NewBaseOperationProcessReasonError("failed to find token holder partition operators, %s: %w", k, err), nil
			case found:
				ops, err = ststo.StateTokenHolderPartitionOperatorsValue(st)
				if err != nil {
					return nil, base.NewBaseOperationProcessReasonError("failed to get token holder partition operators, %s: %w", k, err), nil
				}
			default:
				return nil, base.NewBaseOperationProcessReasonError("token holder partition operators not in state, %q", k), nil
			}
			operators[k] = &ops
		}

		k = ststo.StateKeyOperatorTokenHolders(it.Contract(), it.Operator(), it.Partition())
		if _, found := holders[k]; !found {
			switch st, found, err := getStateFunc(k); {
			case err != nil:
				return nil, base.NewBaseOperationProcessReasonError("failed to find operator tokenholders, %s: %w", k, err), nil
			case found:
				hds, err = ststo.StateOperatorTokenHoldersValue(st)
				if err != nil {
					return nil, base.NewBaseOperationProcessReasonError("failed to get operator tokenholders, %s: %w", k, err), nil
				}
			default:
				return nil, base.NewBaseOperationProcessReasonError("operator tokenholders not in state, %q", k), nil
			}

			holders[k] = &hds
		}
	}

	ipcs := make([]*RevokeOperatorItemProcessor, len(fact.items))
	for i, it := range fact.Items() {
		ip := revokeOperatorItemProcessorPool.Get()
		ipc, ok := ip.(*RevokeOperatorItemProcessor)
		if !ok {
			return nil, nil, e.Wrap(errors.Errorf("expected RevokeOperatorItemProcessor, not %T", ip))
		}

		ipc.h = op.Hash()
		ipc.sender = fact.Sender()
		ipc.item = it
		ipc.operators = operators[ststo.StateKeyTokenHolderPartitionOperators(it.Contract(), fact.sender, it.Partition())]
		ipc.tokenHolders = holders[ststo.StateKeyOperatorTokenHolders(it.Contract(), it.Operator(), it.Partition())]

		s, err := ipc.Process(ctx, op, getStateFunc)
		if err != nil {
			return nil, base.NewBaseOperationProcessReasonError("failed to process RevokeOperatorItem: %w", err), nil
		}
		sts = append(sts, s...)

		ipcs[i] = ipc
	}

	for k, v := range operators {
		sts = append(sts, crcstate.NewStateMergeValue(
			k,
			ststo.NewTokenHolderPartitionOperatorsStateValue(*v),
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
		sts = append(sts, crcstate.NewStateMergeValue(sb[i].Key(), stv))
	}

	return sts, nil, nil
}

func (opp *RevokeOperatorProcessor) Close() error {
	revokeOperatorProcessorPool.Put(opp)

	return nil
}
