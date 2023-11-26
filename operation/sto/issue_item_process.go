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

var issueItemProcessorPool = sync.Pool{
	New: func() interface{} {
		return new(IssueItemProcessor)
	},
}

type IssueItemProcessor struct {
	h      util.Hash
	sender base.Address
	item   IssueItem
}

func (ipp *IssueItemProcessor) PreProcess(
	ctx context.Context, op base.Operation, getStateFunc base.GetStateFunc,
) error {
	it := ipp.item

	st, err := crcstate.ExistsState(stextension.StateKeyContractAccount(it.Contract()), "key of contract account", getStateFunc)
	if err != nil {
		return err
	}

	ca, err := stextension.StateContractAccountValue(st)
	if err != nil {
		return err
	}

	if !(ca.Owner().Equal(ipp.sender) || ca.IsOperator(ipp.sender)) {
		return errors.Errorf("sender is neither the owner nor the operator of the target contract account, %q", ipp.sender)
	}

	if err := crcstate.CheckExistsState(stcurrency.StateKeyAccount(it.Receiver()), getStateFunc); err != nil {
		return err
	}

	if err := crcstate.CheckNotExistsState(stextension.StateKeyContractAccount(it.Receiver()), getStateFunc); err != nil {
		return err
	}

	st, err = crcstate.ExistsState(ststo.StateKeyDesign(it.Contract()), "key of sto design", getStateFunc)
	if err != nil {
		return err
	}

	design, err := ststo.StateDesignValue(st)
	if err != nil {
		return err
	}

	gn := new(big.Int)
	gn.SetUint64(design.Granularity())

	if mod := common.NewBigFromBigInt(new(big.Int)).Mod(it.Amount().Int, gn); common.NewBigFromBigInt(mod).OverZero() {
		return errors.Errorf("amount unit does not comply with sto granularity rule, %q, %q", it.Amount(), design.Granularity())
	}

	if err := crcstate.CheckExistsState(stcurrency.StateKeyCurrencyDesign(it.Currency()), getStateFunc); err != nil {
		return err
	}

	return nil
}

func (ipp *IssueItemProcessor) Process(
	ctx context.Context, op base.Operation, getStateFunc base.GetStateFunc,
) ([]base.StateMergeValue, error) {
	sts := make([]base.StateMergeValue, 4)

	it := ipp.item

	st, err := crcstate.ExistsState(ststo.StateKeyDesign(it.Contract()), "key of sto design", getStateFunc)
	if err != nil {
		return nil, err
	}

	design, err := ststo.StateDesignValue(st)
	if err != nil {
		return nil, err
	}
	p := design.Policy()
	dps := p.Partitions()

	var pb common.Big
	switch st, found, err := getStateFunc(ststo.StateKeyPartitionBalance(it.Contract(), it.Partition())); {
	case err != nil:
		return nil, err
	case found:
		pb, err = ststo.StatePartitionBalanceValue(st)
		if err != nil {
			return nil, err
		}
		pb = pb.Add(it.Amount())
	default:
		pb = it.Amount()
		dps = append(dps, it.Partition())
	}

	policy := typesto.NewPolicy(dps, it.Amount().Add(p.Aggregate()), p.Documents())
	if err := policy.IsValid(nil); err != nil {
		return nil, err
	}

	design = typesto.NewDesign(design.Granularity(), policy)
	if err := design.IsValid(nil); err != nil {
		return nil, err
	}

	sts[0] = crcstate.NewStateMergeValue(
		ststo.StateKeyDesign(it.Contract()),
		ststo.NewDesignStateValue(design),
	)

	sts[1] = crcstate.NewStateMergeValue(
		ststo.StateKeyPartitionBalance(it.Contract(), it.Partition()),
		ststo.NewPartitionBalanceStateValue(pb),
	)

	var ps []typesto.Partition
	switch st, found, err := getStateFunc(ststo.StateKeyTokenHolderPartitions(it.Contract(), it.Receiver())); {
	case err != nil:
		return nil, err
	case found:
		ps, err = ststo.StateTokenHolderPartitionsValue(st)
		if err != nil {
			return nil, err
		}
	default:
		ps = []typesto.Partition{}
	}

	if len(ps) == 0 {
		ps = append(ps, it.Partition())
	} else {
		for i, pt := range ps {
			if pt == it.Partition() {
				break
			}

			if i == len(ps)-1 {
				ps = append(ps, it.Partition())
			}
		}
	}

	sts[2] = crcstate.NewStateMergeValue(
		ststo.StateKeyTokenHolderPartitions(it.Contract(), it.Receiver()),
		ststo.NewTokenHolderPartitionsStateValue(ps),
	)

	var am common.Big
	switch st, found, err := getStateFunc(ststo.StateKeyTokenHolderPartitionBalance(it.Contract(), it.Receiver(), it.Partition())); {
	case err != nil:
		return nil, err
	case found:
		am, err = ststo.StateTokenHolderPartitionBalanceValue(st)
		if err != nil {
			return nil, err
		}
	default:
		am = common.ZeroBig
	}

	am = am.Add(it.Amount())

	sts[3] = crcstate.NewStateMergeValue(
		ststo.StateKeyTokenHolderPartitionBalance(it.Contract(), it.Receiver(), it.Partition()),
		ststo.NewTokenHolderPartitionBalanceStateValue(am, it.Partition()),
	)

	return sts, nil
}

func (ipp *IssueItemProcessor) Close() error {
	ipp.h = nil
	ipp.sender = nil
	ipp.item = IssueItem{}

	issueItemProcessorPool.Put(ipp)

	return nil
}
