package sto

import (
	"context"
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

var createSecurityTokenItemProcessorPool = sync.Pool{
	New: func() interface{} {
		return new(CreateSecurityTokenItemProcessor)
	},
}

type CreateSecurityTokenItemProcessor struct {
	h      util.Hash
	sender base.Address
	item   CreateSecurityTokenItem
}

func (ipp *CreateSecurityTokenItemProcessor) PreProcess(
	ctx context.Context, op base.Operation, getStateFunc base.GetStateFunc,
) error {
	it := ipp.item

	if err := crcstate.CheckExistsState(stextension.StateKeyContractAccount(it.Contract()), getStateFunc); err != nil {
		return err
	}

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

	if ca.IsActive() {
		return errors.Errorf("a contract account is already used, %q", it.Contract().String())
	}

	if err := crcstate.CheckNotExistsState(ststo.StateKeyPartitionBalance(it.Contract(), it.DefaultPartition()), getStateFunc); err != nil {
		return err
	}

	if err := crcstate.CheckExistsState(stcurrency.StateKeyCurrencyDesign(it.Currency()), getStateFunc); err != nil {
		return err
	}

	return nil
}

func (ipp *CreateSecurityTokenItemProcessor) Process(
	ctx context.Context, op base.Operation, getStateFunc base.GetStateFunc,
) ([]base.StateMergeValue, error) {
	sts := make([]base.StateMergeValue, 2)

	it := ipp.item

	partition := it.DefaultPartition()
	partitions := []typesto.Partition{partition}
	var documents []typesto.Document

	policy := typesto.NewPolicy(partitions, common.NewBig(0), documents)
	design := typesto.NewDesign(it.Granularity(), policy)

	if err := design.IsValid(nil); err != nil {
		return nil, err
	}

	sts[0] = crcstate.NewStateMergeValue(
		ststo.StateKeyDesign(it.Contract()),
		ststo.NewDesignStateValue(design),
	)
	sts[1] = crcstate.NewStateMergeValue(
		ststo.StateKeyPartitionBalance(it.Contract(), it.DefaultPartition()),
		ststo.NewPartitionBalanceStateValue(common.ZeroBig),
	)

	return sts, nil
}

func (ipp *CreateSecurityTokenItemProcessor) Close() error {
	ipp.h = nil
	ipp.sender = nil
	ipp.item = CreateSecurityTokenItem{}

	createSecurityTokenItemProcessorPool.Put(ipp)

	return nil
}
