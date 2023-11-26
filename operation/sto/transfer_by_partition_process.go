package sto

import (
	"context"
	"sync"

	"github.com/ProtoconNet/mitum-currency/v3/common"
	"github.com/ProtoconNet/mitum-currency/v3/operation/currency"
	crcstate "github.com/ProtoconNet/mitum-currency/v3/state"
	stcurrency "github.com/ProtoconNet/mitum-currency/v3/state/currency"
	stextension "github.com/ProtoconNet/mitum-currency/v3/state/extension"
	crctypes "github.com/ProtoconNet/mitum-currency/v3/types"
	ststo "github.com/ProtoconNet/mitum-sto/state/sto"
	typesto "github.com/ProtoconNet/mitum-sto/types/sto"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/pkg/errors"
)

var transferByPartitionProcessorPool = sync.Pool{
	New: func() interface{} {
		return new(TransferByPartitionProcessor)
	},
}

func (TransferByPartition) Process(
	ctx context.Context, getStateFunc base.GetStateFunc,
) ([]base.StateMergeValue, base.OperationProcessReasonError, error) {
	return nil, nil, nil
}

type TransferByPartitionProcessor struct {
	*base.BaseOperationProcessor
}

func NewTransferByPartitionProcessor() crctypes.GetNewProcessor {
	return func(
		height base.Height,
		getStateFunc base.GetStateFunc,
		newPreProcessConstraintFunc base.NewOperationProcessorProcessFunc,
		newProcessConstraintFunc base.NewOperationProcessorProcessFunc,
	) (base.OperationProcessor, error) {
		e := util.StringError("failed to create new TransferByPartitionProcessor")

		nopp := transferByPartitionProcessorPool.Get()
		opp, ok := nopp.(*TransferByPartitionProcessor)
		if !ok {
			return nil, e.Wrap(errors.Errorf("expected TransferByPartitionProcessor, not %T", nopp))
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

func (opp *TransferByPartitionProcessor) PreProcess(
	ctx context.Context, op base.Operation, getStateFunc base.GetStateFunc,
) (context.Context, base.OperationProcessReasonError, error) {
	e := util.StringError("failed to preprocess TransferByPartition")

	fact, ok := op.Fact().(TransferByPartitionFact)
	if !ok {
		return ctx, nil, e.Wrap(errors.Errorf("expected TransferByPartitionFact, not %T", op.Fact()))
	}

	if err := fact.IsValid(nil); err != nil {
		return ctx, nil, e.Wrap(err)
	}

	if err := crcstate.CheckExistsState(stcurrency.StateKeyAccount(fact.Sender()), getStateFunc); err != nil {
		return ctx, base.NewBaseOperationProcessReasonError("sender not found, %q: %w", fact.Sender(), err), nil
	}

	if err := crcstate.CheckNotExistsState(stextension.StateKeyContractAccount(fact.Sender()), getStateFunc); err != nil {
		return ctx, base.NewBaseOperationProcessReasonError(
			"contract account cannot issue security tokens, %q: %w",
			fact.Sender(), err,
		), nil
	}

	if err := crcstate.CheckFactSignsByState(fact.sender, op.Signs(), getStateFunc); err != nil {
		return ctx, base.NewBaseOperationProcessReasonError("invalid signing: %w", err), nil
	}

	partitions := map[string][]typesto.Partition{}

	for _, it := range fact.Items() {
		k := ststo.StateKeyTokenHolderPartitions(it.Contract(), it.TokenHolder())

		if _, found := partitions[k]; !found {
			pts, err := ststo.ExistsTokenHolderPartitions(it.Contract(), it.TokenHolder(), getStateFunc)
			if err != nil {
				return nil, base.NewBaseOperationProcessReasonError(
					"failed to get token holder partitions value, %q: %w",
					k, err,
				), nil
			}

			partitions[k] = pts
		}
	}

	for _, it := range fact.Items() {
		ip := transferByPartitionItemProcessorPool.Get()
		ipc, ok := ip.(*TransferByPartitionItemProcessor)
		if !ok {
			return nil, nil, e.Wrap(errors.Errorf("expected TransferByPartitionItemProcessor, not %T", ip))
		}

		ipc.h = op.Hash()
		ipc.sender = fact.Sender()
		ipc.item = it
		ipc.partitions = partitions
		ipc.balances = nil

		if err := ipc.PreProcess(ctx, op, getStateFunc); err != nil {
			return nil, base.NewBaseOperationProcessReasonError(
				"fail to preprocess TransferByPartitionItem: %w", err,
			), nil
		}

		ipc.Close()
	}

	if err := checkEnoughTokenHolderBalance(getStateFunc, fact.Items()); err != nil {
		return nil, base.NewBaseOperationProcessReasonError(
			"not enough token holder partition balance: %w", err,
		), nil
	}

	return ctx, nil, nil
}

func (opp *TransferByPartitionProcessor) Process( // nolint:dupl
	ctx context.Context, op base.Operation, getStateFunc base.GetStateFunc) (
	[]base.StateMergeValue, base.OperationProcessReasonError, error,
) {
	e := util.StringError("failed to process TransferByPartition")

	fact, ok := op.Fact().(TransferByPartitionFact)
	if !ok {
		return nil, nil, e.Wrap(errors.Errorf("expected TransferByPartitionFact, not %T", op.Fact()))
	}

	partitions := map[string][]typesto.Partition{}
	balances := map[string]common.Big{}

	for _, it := range fact.Items() {
		k := ststo.StateKeyTokenHolderPartitions(it.Contract(), it.TokenHolder())

		if _, found := partitions[k]; !found {
			pts, err := ststo.ExistsTokenHolderPartitions(it.Contract(), it.TokenHolder(), getStateFunc)
			if err != nil {
				return nil, base.NewBaseOperationProcessReasonError(
					"failed to get token holder partitions value, %q: %w",
					k, err,
				), nil
			}

			partitions[k] = pts
		}

		k = ststo.StateKeyTokenHolderPartitionBalance(it.Contract(), it.TokenHolder(), it.Partition())

		if _, found := balances[k]; !found {
			balance, err := ststo.ExistsTokenHolderPartitionBalance(it.Contract(), it.TokenHolder(), it.Partition(), getStateFunc)
			if err != nil {
				return nil, base.NewBaseOperationProcessReasonError(
					"failed to get token holder partition balance value, %q: %w",
					k, err,
				), nil
			}

			balances[k] = balance
		}

		k = ststo.StateKeyTokenHolderPartitions(it.Contract(), it.Receiver())

		if _, found := partitions[k]; !found {
			var pts []typesto.Partition

			switch st, found, err := getStateFunc(k); {
			case err != nil:
				return nil, base.NewBaseOperationProcessReasonError(
					"failed to get token holder partitions, %q: %w", k, err,
				), nil
			case found:
				pts, err = ststo.StateTokenHolderPartitionsValue(st)
				if err != nil {
					return nil, base.NewBaseOperationProcessReasonError(
						"failed to get token holder partitions value, %q: %w",
						k, err,
					), nil
				}
			default:
				pts = []typesto.Partition{}
			}

			partitions[k] = pts
		}

		k = ststo.StateKeyTokenHolderPartitionBalance(it.Contract(), it.Receiver(), it.Partition())

		if _, found := balances[k]; !found {
			var am common.Big

			switch st, found, err := getStateFunc(k); {
			case err != nil:
				return nil, base.NewBaseOperationProcessReasonError(
					"failed to get token holder partition balance, %q: %w", k, err,
				), nil
			case found:
				am, err = ststo.StateTokenHolderPartitionBalanceValue(st)
				if err != nil {
					return nil, base.NewBaseOperationProcessReasonError(
						"failed to get token holder partition balance value, %q: %w", k, err,
					), nil
				}
			default:
				am = common.ZeroBig
			}

			balances[k] = am
		}
	}

	var sts []base.StateMergeValue // nolint:prealloc

	ipcs := make([]*TransferByPartitionItemProcessor, len(fact.Items()))
	for i, it := range fact.Items() {
		ip := transferByPartitionItemProcessorPool.Get()
		ipc, ok := ip.(*TransferByPartitionItemProcessor)
		if !ok {
			return nil, nil, e.Wrap(errors.Errorf("expected TransferByPartitionItemProcessor, not %T", ip))
		}

		ipc.h = op.Hash()
		ipc.sender = fact.Sender()
		ipc.item = it
		ipc.partitions = partitions
		ipc.balances = balances

		s, err := ipc.Process(ctx, op, getStateFunc)
		if err != nil {
			return nil, base.NewBaseOperationProcessReasonError(
				"failed to process TransferByPartitionItem: %w", err,
			), nil
		}
		sts = append(sts, s...)

		ipcs[i] = ipc
	}

	for k, v := range partitions {
		sts = append(sts, crcstate.NewStateMergeValue(k, ststo.NewTokenHolderPartitionsStateValue(v)))
	}

	for _, it := range fact.Items() {
		k := ststo.StateKeyTokenHolderPartitionBalance(it.Contract(), it.TokenHolder(), it.Partition())
		sts = append(
			sts,
			crcstate.NewStateMergeValue(k, ststo.NewTokenHolderPartitionBalanceStateValue(balances[k], it.Partition())),
		)

		k = ststo.StateKeyTokenHolderPartitionBalance(it.Contract(), it.Receiver(), it.Partition())
		sts = append(
			sts,
			crcstate.NewStateMergeValue(k, ststo.NewTokenHolderPartitionBalanceStateValue(balances[k], it.Partition())),
		)
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

func (opp *TransferByPartitionProcessor) Close() error {
	transferByPartitionProcessorPool.Put(opp)

	return nil
}

func checkEnoughTokenHolderBalance(getStateFunc base.GetStateFunc, items []TransferByPartitionItem) error {
	balances := map[string]common.Big{}
	amounts := map[string]common.Big{}

	for _, it := range items {
		k := ststo.StateKeyTokenHolderPartitionBalance(it.Contract(), it.TokenHolder(), it.Partition())

		if _, found := balances[k]; found {
			amounts[k] = amounts[k].Add(it.Amount())
			continue
		}

		balance, err := ststo.ExistsTokenHolderPartitionBalance(it.Contract(), it.TokenHolder(), it.Partition(), getStateFunc)
		if err != nil {
			return err
		}

		balances[k] = balance
		amounts[k] = it.Amount()
	}

	for k, balance := range balances {
		if balance.Compare(amounts[k]) < 0 {
			return errors.Errorf("token holder partition balance not over total amounts, %q, %q < %q", k, balance, amounts[k])
		}
	}

	return nil
}
