package sto

import (
	"context"
	"math/big"
	"sync"

	"github.com/ProtoconNet/mitum-currency/v3/common"
	crcstate "github.com/ProtoconNet/mitum-currency/v3/state"
	stcurrency "github.com/ProtoconNet/mitum-currency/v3/state/currency"
	stextension "github.com/ProtoconNet/mitum-currency/v3/state/extension"
	ststo "github.com/ProtoconNet/mitum-sto/state/sto"
	typesto "github.com/ProtoconNet/mitum-sto/types/sto"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/pkg/errors"
)

var redeemItemProcessorPool = sync.Pool{
	New: func() interface{} {
		return new(RedeemItemProcessor)
	},
}

type RedeemItemProcessor struct {
	h                util.Hash
	sender           base.Address
	item             RedeemItem
	sto              *typesto.Design
	partitionBalance *common.Big
}

func (ipp *RedeemItemProcessor) PreProcess(
	ctx context.Context, op base.Operation, getStateFunc base.GetStateFunc,
) error {
	it := ipp.item

	if err := crcstate.CheckExistsState(stcurrency.StateKeyAccount(it.TokenHolder()), getStateFunc); err != nil {
		return err
	}

	if err := crcstate.CheckNotExistsState(stextension.StateKeyContractAccount(it.TokenHolder()), getStateFunc); err != nil {
		return err
	}

	st, err := crcstate.ExistsState(stextension.StateKeyContractAccount(
		it.Contract()),
		"key of contract account",
		getStateFunc,
	)
	if err != nil {
		return err
	}

	ca, err := stextension.StateContractAccountValue(st)
	if err != nil {
		return err
	}

	if !it.TokenHolder().Equal(ipp.sender) {
		isOperator := false

		if !(ca.Owner().Equal(ipp.sender) || ca.IsOperator(ipp.sender)) {
			st, err := crcstate.ExistsState(
				ststo.StateKeyTokenHolderPartitionOperators(it.Contract(), it.TokenHolder(), it.Partition()),
				"key of token holder partition operators",
				getStateFunc,
			)
			if err != nil {
				return err
			}

			operators, err := ststo.StateTokenHolderPartitionOperatorsValue(st)
			if err != nil {
				return err
			}

			for _, op := range operators {
				if op.Equal(ipp.sender) {
					isOperator = true
					break
				}
			}
			if !isOperator {
				return errors.Errorf(
					"sender is neither contract owner nor contract operator and sto operator, %s, %q",
					it.Partition(),
					ipp.sender,
				)
			}
		}
	}

	design := ipp.sto

	partitions, err := ststo.ExistsTokenHolderPartitions(it.Contract(), it.TokenHolder(), getStateFunc)
	if err != nil {
		return err
	}

	if len(partitions) == 0 {
		return errors.Errorf("empty token holder partitions, %s-%s", it.Contract(), it.TokenHolder())
	}

	for i, p := range partitions {
		if p == it.Partition() {
			break
		}

		if i == len(partitions)-1 {
			return errors.Errorf(
				"partition not in token holder partitions, %s-%s, %q",
				it.Contract(), it.TokenHolder(), it.Partition(),
			)
		}
	}

	balance, err := ststo.ExistsTokenHolderPartitionBalance(it.Contract(), it.TokenHolder(), it.Partition(), getStateFunc)
	if err != nil {
		return err
	}

	if balance.Compare(it.Amount()) < 0 {
		return errors.Errorf(
			"token holder partition, %s-%s-%s balance not over item amount, %q < %q",
			it.Contract(), it.TokenHolder(), it.Partition(), balance, it.Amount(),
		)
	}

	gn := new(big.Int)
	gn.SetUint64(design.Granularity())

	if mod := common.NewBigFromBigInt(new(big.Int)).Mod(it.Amount().Int, gn); common.NewBigFromBigInt(mod).OverZero() {
		return errors.Errorf(
			"amount unit does not comply with sto granularity rule, %q, %q",
			it.Amount(), design.Granularity(),
		)
	}

	if err := crcstate.CheckExistsState(stcurrency.StateKeyCurrencyDesign(it.Currency()), getStateFunc); err != nil {
		return err
	}

	return nil
}

func (ipp *RedeemItemProcessor) Process(
	ctx context.Context, op base.Operation, getStateFunc base.GetStateFunc,
) ([]base.StateMergeValue, error) {
	it := ipp.item

	*ipp.partitionBalance = (*ipp.partitionBalance).Sub(it.Amount())

	design := *ipp.sto
	policy := design.Policy()

	aggr := policy.Aggregate().Sub(it.Amount())

	if (*ipp.partitionBalance).OverZero() {
		policy = typesto.NewPolicy(policy.Partitions(), aggr, policy.Documents())
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

		policy = typesto.NewPolicy(partitions, aggr, policy.Documents())
		if err := policy.IsValid(nil); err != nil {
			return nil, err
		}
	}

	design = typesto.NewDesign(design.Granularity(), policy)
	if err := design.IsValid(nil); err != nil {
		return nil, err
	}

	*ipp.sto = design

	balance, err := ststo.ExistsTokenHolderPartitionBalance(it.Contract(), it.TokenHolder(), it.Partition(), getStateFunc)
	if err != nil {
		return nil, err
	}

	tokenHolderPartitions, err := ststo.ExistsTokenHolderPartitions(it.Contract(), it.TokenHolder(), getStateFunc)
	if err != nil {
		return nil, err
	}

	var sts []base.StateMergeValue

	balance = balance.Sub(it.Amount())
	if !balance.OverZero() {
		for i, p := range tokenHolderPartitions {
			if p == it.Partition() {
				if i < len(tokenHolderPartitions)-1 {
					copy(tokenHolderPartitions[i:], tokenHolderPartitions[i+1:])
				}
				tokenHolderPartitions = tokenHolderPartitions[:len(tokenHolderPartitions)-1]
			}
		}

		opk := ststo.StateKeyTokenHolderPartitionOperators(it.Contract(), it.TokenHolder(), it.Partition())

		st, err := crcstate.ExistsState(opk, "key of token holder partition operators", getStateFunc)
		if err != nil {
			return nil, err
		}

		operators, err := ststo.StateTokenHolderPartitionOperatorsValue(st)
		if err != nil {
			return nil, err
		}

		sts = append(sts, crcstate.NewStateMergeValue(
			opk, ststo.NewTokenHolderPartitionOperatorsStateValue([]base.Address{}),
		))

		for _, op := range operators {
			thk := ststo.StateKeyOperatorTokenHolders(it.Contract(), op, it.Partition())

			st, err := crcstate.ExistsState(thk, "key of operator token holders", getStateFunc)
			if err != nil {
				return nil, err
			}

			holders, err := ststo.StateOperatorTokenHoldersValue(st)
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

			sts = append(sts, crcstate.NewStateMergeValue(
				thk, ststo.NewOperatorTokenHoldersStateValue(holders),
			))
		}
	}

	sts = append(sts, crcstate.NewStateMergeValue(
		ststo.StateKeyTokenHolderPartitionBalance(it.Contract(), it.TokenHolder(), it.Partition()),
		ststo.NewTokenHolderPartitionBalanceStateValue(balance, it.Partition()),
	))
	sts = append(sts, crcstate.NewStateMergeValue(
		ststo.StateKeyTokenHolderPartitions(it.Contract(), it.TokenHolder()),
		ststo.NewTokenHolderPartitionsStateValue(tokenHolderPartitions),
	))

	return sts, nil
}

func (ipp *RedeemItemProcessor) Close() error {
	ipp.h = nil
	ipp.sender = nil
	ipp.item = RedeemItem{}
	ipp.sto = nil
	ipp.partitionBalance = nil

	redeemItemProcessorPool.Put(ipp)

	return nil
}
