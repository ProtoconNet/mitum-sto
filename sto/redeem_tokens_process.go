package sto

import (
	"context"
	"fmt"
	"math/big"
	"sync"

	extensioncurrency "github.com/ProtoconNet/mitum-currency-extension/v2/currency"
	"github.com/ProtoconNet/mitum-currency/v2/currency"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/pkg/errors"
)

var redeemTokensItemProcessorPool = sync.Pool{
	New: func() interface{} {
		return new(RedeemTokensItemProcessor)
	},
}

var redeemTokensProcessorPool = sync.Pool{
	New: func() interface{} {
		return new(RedeemTokensProcessor)
	},
}

func (RedeemTokens) Process(
	ctx context.Context, getStateFunc base.GetStateFunc,
) ([]base.StateMergeValue, base.OperationProcessReasonError, error) {
	return nil, nil, nil
}

type RedeemTokensItemProcessor struct {
	h                util.Hash
	sender           base.Address
	item             RedeemTokensItem
	sto              *STODesign
	partitionBalance *currency.Big
}

func (ipp *RedeemTokensItemProcessor) PreProcess(
	ctx context.Context, op base.Operation, getStateFunc base.GetStateFunc,
) error {
	it := ipp.item

	if err := checkExistsState(extensioncurrency.StateKeyContractAccount(it.Contract()), getStateFunc); err != nil {
		return err
	}

	if err := checkExistsState(currency.StateKeyAccount(it.TokenHolder()), getStateFunc); err != nil {
		return err
	}

	if err := checkNotExistsState(extensioncurrency.StateKeyContractAccount(it.TokenHolder()), getStateFunc); err != nil {
		return err
	}

	design := ipp.sto

	if !it.TokenHolder().Equal(ipp.sender) {
		policy := ipp.sto.Policy()

		controllers := policy.Controllers()
		isController, isOperator := false, false

		for _, con := range controllers {
			if con.Equal(ipp.sender) {
				isController = true
				break
			}
		}

		if !isController {
			st, err := existsState(StateKeyTokenHolderPartitionOperators(it.Contract(), it.STO(), it.TokenHolder(), it.Partition()), "key of tokenholder partition operators", getStateFunc)
			if err != nil {
				return err
			}

			operators, err := StateTokenHolderPartitionOperatorsValue(st)
			if err != nil {
				return err
			}

			for _, op := range operators {
				if op.Equal(ipp.sender) {
					isOperator = true
					break
				}
			}
		}

		if !(isController || isOperator) {
			return errors.Errorf("sender is neither controller nor operator, %s, %q", it.Partition(), ipp.sender)
		}
	}

	partitions, err := existsTokenHolderPartitions(it.Contract(), it.STO(), it.TokenHolder(), getStateFunc)
	if err != nil {
		return err
	}

	for i, p := range partitions {
		if p == it.Partition() {
			break
		}

		if i == len(partitions)-1 {
			return errors.Errorf("partition not in tokenholder partitions, %s-%s-%s, %q", it.Contract(), it.STO(), it.TokenHolder(), it.Partition())
		}
	}

	balance, err := existsTokenHolderPartitionBalance(it.Contract(), it.STO(), it.TokenHolder(), it.Partition(), getStateFunc)
	if err != nil {
		return err
	}

	if balance.Compare(it.Amount()) < 0 {
		k := fmt.Sprintf("%s-%s-%s-%s", it.Contract(), it.STO(), it.TokenHolder(), it.Partition())
		return errors.Errorf("tokenholder partition balance not over item amount, %q, %q < %q", k, balance, it.Amount())
	}

	gn := new(big.Int)
	gn.SetUint64(design.Granularity())

	if mod := currency.NewBigFromBigInt(new(big.Int)).Mod(it.Amount().Int, gn); currency.NewBigFromBigInt(mod).OverZero() {
		return errors.Errorf("amount unit does not comply with sto granularity rule, %q, %q", it.Amount(), design.Granularity())
	}

	if err := checkExistsState(currency.StateKeyCurrencyDesign(it.Currency()), getStateFunc); err != nil {
		return err
	}

	return nil
}

func (ipp *RedeemTokensItemProcessor) Process(
	ctx context.Context, op base.Operation, getStateFunc base.GetStateFunc,
) ([]base.StateMergeValue, error) {
	it := ipp.item

	*ipp.partitionBalance = (*ipp.partitionBalance).Sub(it.Amount())

	design := *ipp.sto
	policy := design.Policy()

	aggr := policy.Aggregate().Sub(it.Amount())

	if (*ipp.partitionBalance).OverZero() {
		policy = NewSTOPolicy(policy.Partitions(), aggr, policy.Controllers(), policy.Documents())
		if err := policy.IsValid(nil); err != nil {
			return nil, err
		}
	} else {
		partitions := policy.Partitions()
		for i, p := range partitions {
			if p == it.Partition() {
				if i < len(partitions)-1 {
					copy(partitions[i:], partitions[i+1:])
				}
				partitions = partitions[:len(partitions)-1]
			}
		}

		policy = NewSTOPolicy(partitions, aggr, policy.Controllers(), policy.Documents())
		if err := policy.IsValid(nil); err != nil {
			return nil, err
		}
	}

	design = NewSTODesign(it.STO(), design.Granularity(), policy)
	if err := design.IsValid(nil); err != nil {
		return nil, err
	}

	*ipp.sto = design

	balance, err := existsTokenHolderPartitionBalance(it.Contract(), it.STO(), it.TokenHolder(), it.Partition(), getStateFunc)
	if err != nil {
		return nil, err
	}

	tokenholderPartitions, err := existsTokenHolderPartitions(it.Contract(), it.STO(), it.TokenHolder(), getStateFunc)
	if err != nil {
		return nil, err
	}

	sts := []base.StateMergeValue{}

	balance = balance.Sub(it.Amount())
	if !balance.OverZero() {
		for i, p := range tokenholderPartitions {
			if p == it.Partition() {
				if i < len(tokenholderPartitions)-1 {
					copy(tokenholderPartitions[i:], tokenholderPartitions[i+1:])
				}
				tokenholderPartitions = tokenholderPartitions[:len(tokenholderPartitions)-1]
			}
		}

		opk := StateKeyTokenHolderPartitionOperators(it.Contract(), it.STO(), it.TokenHolder(), it.Partition())

		st, err := existsState(opk, "key of tokenholder partition operators", getStateFunc)
		if err != nil {
			return nil, err
		}

		operators, err := StateTokenHolderPartitionOperatorsValue(st)
		if err != nil {
			return nil, err
		}

		sts = append(sts, NewStateMergeValue(
			opk, NewTokenHolderPartitionOperatorsStateValue([]base.Address{}),
		))

		for _, op := range operators {
			thk := StateKeyOperatorTokenHolders(it.Contract(), it.STO(), op, it.Partition())

			st, err := existsState(thk, "key of operator tokenholders", getStateFunc)
			if err != nil {
				return nil, err
			}

			holders, err := StateOperatorTokenHoldersValue(st)
			if err != nil {
				return nil, err
			}

			for i, th := range holders {
				if th.Equal(it.TokenHolder()) {
					if i < len(holders)-1 {
						copy(holders[i:], holders[i+1:])
					}
					holders = holders[:len(holders)-1]
				}
			}

			sts = append(sts, NewStateMergeValue(
				thk, NewOperatorTokenHoldersStateValue(holders),
			))
		}
	}

	sts = append(sts, NewStateMergeValue(
		StateKeyTokenHolderPartitionBalance(it.Contract(), it.STO(), it.TokenHolder(), it.Partition()),
		NewTokenHolderPartitionBalanceStateValue(balance, it.Partition()),
	))
	sts = append(sts, NewStateMergeValue(
		StateKeyTokenHolderPartitions(it.Contract(), it.STO(), it.TokenHolder()),
		NewTokenHolderPartitionsStateValue(tokenholderPartitions),
	))

	return sts, nil
}

func (ipp *RedeemTokensItemProcessor) Close() error {
	ipp.h = nil
	ipp.sender = nil
	ipp.item = RedeemTokensItem{}
	ipp.sto = nil
	ipp.partitionBalance = nil

	redeemTokensItemProcessorPool.Put(ipp)

	return nil
}

type RedeemTokensProcessor struct {
	*base.BaseOperationProcessor
}

func NewRedeemTokensProcessor() extensioncurrency.GetNewProcessor {
	return func(
		height base.Height,
		getStateFunc base.GetStateFunc,
		newPreProcessConstraintFunc base.NewOperationProcessorProcessFunc,
		newProcessConstraintFunc base.NewOperationProcessorProcessFunc,
	) (base.OperationProcessor, error) {
		e := util.StringErrorFunc("failed to create new RedeemTokensProcessor")

		nopp := redeemTokensProcessorPool.Get()
		opp, ok := nopp.(*RedeemTokensProcessor)
		if !ok {
			return nil, e(nil, "expected RedeemTokensProcessor, not %T", nopp)
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

func (opp *RedeemTokensProcessor) PreProcess(
	ctx context.Context, op base.Operation, getStateFunc base.GetStateFunc,
) (context.Context, base.OperationProcessReasonError, error) {
	e := util.StringErrorFunc("failed to preprocess RedeemTokens")

	fact, ok := op.Fact().(RedeemTokensFact)
	if !ok {
		return ctx, nil, e(nil, "expected RedeemTokensFact, not %T", op.Fact())
	}

	if err := fact.IsValid(nil); err != nil {
		return ctx, nil, e(err, "")
	}

	if err := checkExistsState(currency.StateKeyAccount(fact.Sender()), getStateFunc); err != nil {
		return ctx, base.NewBaseOperationProcessReasonError("sender not found, %q: %w", fact.Sender(), err), nil
	}

	if err := checkNotExistsState(extensioncurrency.StateKeyContractAccount(fact.Sender()), getStateFunc); err != nil {
		return ctx, base.NewBaseOperationProcessReasonError("contract account cannot issue security tokens, %q: %w", fact.Sender(), err), nil
	}

	if err := checkFactSignsByState(fact.sender, op.Signs(), getStateFunc); err != nil {
		return ctx, base.NewBaseOperationProcessReasonError("invalid signing: %w", err), nil
	}

	stos := map[string]*STODesign{}

	for _, it := range fact.Items() {
		k := StateKeySTODesign(it.Contract(), it.STO())

		if _, found := stos[k]; !found {
			st, err := existsState(k, "key of sto design", getStateFunc)
			if err != nil {
				return nil, base.NewBaseOperationProcessReasonError("sto design doesn't exist, %q: %w", k, err), nil
			}

			design, err := StateSTODesignValue(st)
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
		ip := redeemTokensItemProcessorPool.Get()
		ipc, ok := ip.(*RedeemTokensItemProcessor)
		if !ok {
			return nil, nil, e(nil, "expected RedeemTokensItemProcessor, not %T", ip)
		}

		ipc.h = op.Hash()
		ipc.sender = fact.Sender()
		ipc.item = it
		ipc.sto = stos[StateKeySTODesign(it.Contract(), it.STO())]
		ipc.partitionBalance = nil

		if err := ipc.PreProcess(ctx, op, getStateFunc); err != nil {
			return nil, base.NewBaseOperationProcessReasonError("fail to preprocess RedeemTokensItem: %w", err), nil
		}

		ipc.Close()
	}

	return ctx, nil, nil
}

func (opp *RedeemTokensProcessor) Process( // nolint:dupl
	ctx context.Context, op base.Operation, getStateFunc base.GetStateFunc) (
	[]base.StateMergeValue, base.OperationProcessReasonError, error,
) {
	e := util.StringErrorFunc("failed to process RedeemTokens")

	fact, ok := op.Fact().(RedeemTokensFact)
	if !ok {
		return nil, nil, e(nil, "expected RedeemTokensFact, not %T", op.Fact())
	}

	stos := map[string]*STODesign{}

	for _, it := range fact.Items() {
		k := StateKeySTODesign(it.Contract(), it.STO())

		if _, found := stos[k]; !found {
			st, err := existsState(k, "key of sto design", getStateFunc)
			if err != nil {
				return nil, base.NewBaseOperationProcessReasonError("sto design doesn't exist, %q: %w", k, err), nil
			}

			design, err := StateSTODesignValue(st)
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

	ipcs := make([]*RedeemTokensItemProcessor, len(fact.Items()))
	for i, it := range fact.Items() {
		ip := redeemTokensItemProcessorPool.Get()
		ipc, ok := ip.(*RedeemTokensItemProcessor)
		if !ok {
			return nil, nil, e(nil, "expected RedeemTokensItemProcessor, not %T", ip)
		}

		ipc.h = op.Hash()
		ipc.sender = fact.Sender()
		ipc.item = it
		ipc.sto = stos[StateKeySTODesign(it.Contract(), it.STO())]
		ipc.partitionBalance = partitionBalances[StateKeyPartitionBalance(it.Contract(), it.STO(), it.Partition())]

		s, err := ipc.Process(ctx, op, getStateFunc)
		if err != nil {
			return nil, base.NewBaseOperationProcessReasonError("failed to process RedeemTokensItem: %w", err), nil
		}
		sts = append(sts, s...)

		ipcs[i] = ipc
	}

	for k, v := range stos {
		sts = append(sts, NewStateMergeValue(k, NewSTODesignStateValue(*v)))
	}

	for k, v := range partitionBalances {
		sts = append(sts, NewStateMergeValue(k, NewPartitionBalanceStateValue(*v)))
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
		v, ok := sb[i].Value().(currency.BalanceStateValue)
		if !ok {
			return nil, nil, e(nil, "expected BalanceStateValue, not %T", sb[i].Value())
		}
		stv := currency.NewBalanceStateValue(v.Amount.WithBig(v.Amount.Big().Sub(required[i][0])))
		sts = append(sts, currency.NewBalanceStateMergeValue(sb[i].Key(), stv))
	}

	return sts, nil, nil
}

func (opp *RedeemTokensProcessor) Close() error {
	redeemTokensProcessorPool.Put(opp)

	return nil
}

func checkEnoughPartitionBalance(getStateFunc base.GetStateFunc, items []RedeemTokensItem) (map[string]*currency.Big, error) {
	balances := map[string]*currency.Big{}
	amounts := map[string]currency.Big{}

	for _, it := range items {
		k := StateKeyPartitionBalance(it.Contract(), it.STO(), it.Partition())

		if _, found := balances[k]; found {
			amounts[k] = amounts[k].Add(it.Amount())
			continue
		}

		st, err := existsState(k, "key of partition balance", getStateFunc)
		if err != nil {
			return nil, err
		}

		balance, err := StatePartitionBalanceValue(st)
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
