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

var redeemProcessorPool = sync.Pool{
	New: func() interface{} {
		return new(RedeemProcessor)
	},
}

func (Redeem) Process(
	ctx context.Context, getStateFunc base.GetStateFunc,
) ([]base.StateMergeValue, base.OperationProcessReasonError, error) {
	return nil, nil, nil
}

type RedeemProcessor struct {
	*base.BaseOperationProcessor
}

func NewRedeemProcessor() crctypes.GetNewProcessor {
	return func(
		height base.Height,
		getStateFunc base.GetStateFunc,
		newPreProcessConstraintFunc base.NewOperationProcessorProcessFunc,
		newProcessConstraintFunc base.NewOperationProcessorProcessFunc,
	) (base.OperationProcessor, error) {
		e := util.StringError("failed to create new RedeemProcessor")

		nopp := redeemProcessorPool.Get()
		opp, ok := nopp.(*RedeemProcessor)
		if !ok {
			return nil, e.Wrap(errors.Errorf("expected RedeemProcessor, not %T", nopp))
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

func (opp *RedeemProcessor) PreProcess(
	ctx context.Context, op base.Operation, getStateFunc base.GetStateFunc,
) (context.Context, base.OperationProcessReasonError, error) {
	e := util.StringError("failed to preprocess Redeem")

	fact, ok := op.Fact().(RedeemFact)
	if !ok {
		return ctx, nil, e.Wrap(errors.Errorf("expected RedeemFact, not %T", op.Fact()))
	}

	if err := fact.IsValid(nil); err != nil {
		return ctx, nil, e.Wrap(err)
	}

	if err := crcstate.CheckExistsState(stcurrency.StateKeyAccount(fact.Sender()), getStateFunc); err != nil {
		return ctx, base.NewBaseOperationProcessReasonError("sender not found, %q: %w", fact.Sender(), err), nil
	}

	if err := crcstate.CheckNotExistsState(stextension.StateKeyContractAccount(fact.Sender()), getStateFunc); err != nil {
		return ctx, base.NewBaseOperationProcessReasonError("contract account cannot issue security tokens, %q: %w", fact.Sender(), err), nil
	}

	if err := crcstate.CheckFactSignsByState(fact.sender, op.Signs(), getStateFunc); err != nil {
		return ctx, base.NewBaseOperationProcessReasonError("invalid signing: %w", err), nil
	}

	stos := map[string]*typesto.Design{}

	for _, it := range fact.Items() {
		k := ststo.StateKeyDesign(it.Contract())

		if _, found := stos[k]; !found {
			st, err := crcstate.ExistsState(k, "key of sto design", getStateFunc)
			if err != nil {
				return nil, base.NewBaseOperationProcessReasonError("sto design doesn't exist, %q: %w", k, err), nil
			}

			design, err := ststo.StateDesignValue(st)
			if err != nil {
				return nil, base.NewBaseOperationProcessReasonError("failed to get sto design value, %q: %w", k, err), nil
			}

			stos[k] = &design
		}
	}

	_, err := checkEnoughPartitionBalance(getStateFunc, fact.Items())
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError("not enough partition balance: %w", err), nil
	}

	for _, it := range fact.Items() {
		ip := redeemItemProcessorPool.Get()
		ipc, ok := ip.(*RedeemItemProcessor)
		if !ok {
			return nil, nil, e.Wrap(errors.Errorf("expected RedeemItemProcessor, not %T", ip))
		}

		ipc.h = op.Hash()
		ipc.sender = fact.Sender()
		ipc.item = it
		ipc.sto = stos[ststo.StateKeyDesign(it.Contract())]
		ipc.partitionBalance = nil

		if err := ipc.PreProcess(ctx, op, getStateFunc); err != nil {
			return nil, base.NewBaseOperationProcessReasonError("fail to preprocess RedeemItem: %w", err), nil
		}

		ipc.Close()
	}

	return ctx, nil, nil
}

func (opp *RedeemProcessor) Process( // nolint:dupl
	ctx context.Context, op base.Operation, getStateFunc base.GetStateFunc) (
	[]base.StateMergeValue, base.OperationProcessReasonError, error,
) {
	e := util.StringError("failed to process Redeem")

	fact, ok := op.Fact().(RedeemFact)
	if !ok {
		return nil, nil, e.Wrap(errors.Errorf("expected RedeemFact, not %T", op.Fact()))
	}

	stos := map[string]*typesto.Design{}

	for _, it := range fact.Items() {
		k := ststo.StateKeyDesign(it.Contract())

		if _, found := stos[k]; !found {
			st, err := crcstate.ExistsState(k, "key of sto design", getStateFunc)
			if err != nil {
				return nil, base.NewBaseOperationProcessReasonError("sto design doesn't exist, %q: %w", k, err), nil
			}

			design, err := ststo.StateDesignValue(st)
			if err != nil {
				return nil, base.NewBaseOperationProcessReasonError("failed to get sto design value, %q: %w", k, err), nil
			}

			stos[k] = &design
		}
	}

	partitionBalances, err := checkEnoughPartitionBalance(getStateFunc, fact.Items())
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError("not enough partition balance: %w", err), nil
	}

	var sts []base.StateMergeValue // nolint:prealloc

	ipcs := make([]*RedeemItemProcessor, len(fact.Items()))
	for i, it := range fact.Items() {
		ip := redeemItemProcessorPool.Get()
		ipc, ok := ip.(*RedeemItemProcessor)
		if !ok {
			return nil, nil, e.Wrap(errors.Errorf("expected RedeemItemProcessor, not %T", ip))
		}

		ipc.h = op.Hash()
		ipc.sender = fact.Sender()
		ipc.item = it
		ipc.sto = stos[ststo.StateKeyDesign(it.Contract())]
		ipc.partitionBalance = partitionBalances[ststo.StateKeyPartitionBalance(it.Contract(), it.Partition())]

		s, err := ipc.Process(ctx, op, getStateFunc)
		if err != nil {
			return nil, base.NewBaseOperationProcessReasonError("failed to process RedeemItem: %w", err), nil
		}
		sts = append(sts, s...)

		ipcs[i] = ipc
	}

	for k, v := range stos {
		sts = append(sts, crcstate.NewStateMergeValue(k, ststo.NewDesignStateValue(*v)))
	}

	for k, v := range partitionBalances {
		sts = append(sts, crcstate.NewStateMergeValue(k, ststo.NewPartitionBalanceStateValue(*v)))
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

func (opp *RedeemProcessor) Close() error {
	redeemProcessorPool.Put(opp)

	return nil
}

func checkEnoughPartitionBalance(getStateFunc base.GetStateFunc, items []RedeemItem) (map[string]*common.Big, error) {
	balances := map[string]*common.Big{}
	amounts := map[string]common.Big{}

	for _, it := range items {
		k := ststo.StateKeyPartitionBalance(it.Contract(), it.Partition())

		if _, found := balances[k]; found {
			amounts[k] = amounts[k].Add(it.Amount())
			continue
		}

		st, err := crcstate.ExistsState(k, "key of partition balance", getStateFunc)
		if err != nil {
			return nil, err
		}

		balance, err := ststo.StatePartitionBalanceValue(st)
		if err != nil {
			return nil, err
		}

		balances[k] = &balance
		amounts[k] = it.Amount()
	}

	for k, balance := range balances {
		if balance.Compare(amounts[k]) < 0 {
			return nil, errors.Errorf("partition balance not over total amounts, %q, %q < %q", k, balance, amounts[k])
		}
	}

	return balances, nil
}
