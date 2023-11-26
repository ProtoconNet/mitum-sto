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

var transferByPartitionItemProcessorPool = sync.Pool{
	New: func() interface{} {
		return new(TransferByPartitionItemProcessor)
	},
}

type TransferByPartitionItemProcessor struct {
	h          util.Hash
	sender     base.Address
	item       TransferByPartitionItem
	partitions map[string][]typesto.Partition
	balances   map[string]common.Big
}

func (ipp *TransferByPartitionItemProcessor) PreProcess(
	ctx context.Context, op base.Operation, getStateFunc base.GetStateFunc,
) error {
	it := ipp.item

	if err := crcstate.CheckExistsState(
		stcurrency.StateKeyAccount(it.TokenHolder()), getStateFunc,
	); err != nil {
		return err
	}

	if err := crcstate.CheckNotExistsState(
		stextension.StateKeyContractAccount(it.TokenHolder()), getStateFunc,
	); err != nil {
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

	if err := crcstate.CheckExistsState(
		stcurrency.StateKeyAccount(it.Receiver()), getStateFunc,
	); err != nil {
		return err
	}

	if err := crcstate.CheckNotExistsState(
		stextension.StateKeyContractAccount(it.Receiver()), getStateFunc,
	); err != nil {
		return err
	}

	partitions := ipp.partitions[ststo.StateKeyTokenHolderPartitions(it.Contract(), it.TokenHolder())]
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
				it.Contract(), it.TokenHolder(), it.Partition())
		}
	}

	st, err = crcstate.ExistsState(
		ststo.StateKeyDesign(it.Contract()), "key of sto design", getStateFunc,
	)
	if err != nil {
		return err
	}

	design, err := ststo.StateDesignValue(st)
	if err != nil {
		return err
	}

	//policy := design.Policy()

	gn := new(big.Int)
	gn.SetUint64(design.Granularity())

	if mod := common.NewBigFromBigInt(new(big.Int)).Mod(it.Amount().Int, gn); common.NewBigFromBigInt(mod).OverZero() {
		return errors.Errorf("amount unit does not comply with sto granularity rule, %q, %q",
			it.Amount(), design.Granularity(),
		)
	}

	if err := crcstate.CheckExistsState(
		stcurrency.StateKeyCurrencyDesign(it.Currency()), getStateFunc,
	); err != nil {
		return err
	}

	return nil
}

func (ipp *TransferByPartitionItemProcessor) Process(
	ctx context.Context, op base.Operation, getStateFunc base.GetStateFunc,
) ([]base.StateMergeValue, error) {
	it := ipp.item

	partitionsKey := ststo.StateKeyTokenHolderPartitions(
		it.Contract(), it.TokenHolder(),
	)
	balanceKey := ststo.StateKeyTokenHolderPartitionBalance(
		it.Contract(), it.TokenHolder(), it.Partition(),
	)

	receiverPartitionsKey := ststo.StateKeyTokenHolderPartitions(
		it.Contract(), it.Receiver(),
	)
	receiverBalanceKey := ststo.StateKeyTokenHolderPartitionBalance(
		it.Contract(), it.Receiver(), it.Partition(),
	)

	balance := ipp.balances[balanceKey]
	partitions := ipp.partitions[partitionsKey]

	receiverBalance := ipp.balances[receiverBalanceKey]
	receiverPartitions := ipp.partitions[receiverPartitionsKey]

	balance = balance.Sub(it.Amount())
	receiverBalance = receiverBalance.Add(it.Amount())

	sts := []base.StateMergeValue{}

	if !balance.OverZero() {
		for i, p := range partitions {
			if p == it.Partition() {
				if i < len(partitions)-1 {
					copy(partitions[i:], partitions[i+1:])
				}
				partitions = partitions[:len(partitions)-1]
			}
		}

		opk := ststo.StateKeyTokenHolderPartitionOperators(
			it.Contract(), it.TokenHolder(), it.Partition(),
		)

		var operators []base.Address
		switch st, found, err := getStateFunc(opk); {
		case err != nil:
			return nil, err
		case found:
			operators, err = ststo.StateTokenHolderPartitionOperatorsValue(st)
			if err != nil {
				return nil, err
			}
		default:
			operators = []base.Address{}
		}

		sts = append(sts, crcstate.NewStateMergeValue(
			opk, ststo.NewTokenHolderPartitionOperatorsStateValue([]base.Address{}),
		))

		for _, op := range operators {
			thk := ststo.StateKeyOperatorTokenHolders(it.Contract(), op, it.Partition())

			st, err := crcstate.ExistsState(thk, "key of operator tokenholders", getStateFunc)
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

	if len(receiverPartitions) == 0 {
		receiverPartitions = append(receiverPartitions, it.Partition())
	} else {
		for i, p := range receiverPartitions {
			if p == it.Partition() {
				break
			}

			if i == len(receiverPartitions)-1 {
				receiverPartitions = append(receiverPartitions, it.Partition())
			}
		}
	}

	ipp.partitions[partitionsKey] = partitions
	ipp.partitions[receiverPartitionsKey] = receiverPartitions
	ipp.balances[balanceKey] = balance
	ipp.balances[receiverBalanceKey] = receiverBalance

	return sts, nil
}

func (ipp *TransferByPartitionItemProcessor) Close() error {
	ipp.h = nil
	ipp.sender = nil
	ipp.item = TransferByPartitionItem{}
	ipp.balances = nil
	ipp.partitions = nil

	transferByPartitionItemProcessorPool.Put(ipp)

	return nil
}
