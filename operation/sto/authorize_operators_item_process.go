package sto

import (
	"context"
	"sync"

	crcstate "github.com/ProtoconNet/mitum-currency/v3/state"
	stcurrency "github.com/ProtoconNet/mitum-currency/v3/state/currency"
	stextension "github.com/ProtoconNet/mitum-currency/v3/state/extension"
	ststo "github.com/ProtoconNet/mitum-sto/state/sto"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/pkg/errors"
)

var authorizeOperatorItemProcessorPool = sync.Pool{
	New: func() interface{} {
		return new(AuthorizeOperatorItemProcessor)
	},
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

	if err := crcstate.CheckExistsState(stextension.StateKeyContractAccount(it.Contract()), getStateFunc); err != nil {
		return err
	}

	if err := crcstate.CheckExistsState(ststo.StateKeyDesign(it.Contract()), getStateFunc); err != nil {
		return err
	}

	if err := crcstate.CheckExistsState(ststo.StateKeyPartitionBalance(it.Contract(), it.Partition()), getStateFunc); err != nil {
		return err
	}

	partitions, err := ststo.ExistsTokenHolderPartitions(it.Contract(), ipp.sender, getStateFunc)
	if err != nil {
		return err
	}

	if len(partitions) == 0 {
		return errors.Errorf("empty token holder partitions, %s-%s", it.Contract(), ipp.sender)
	}

	for i, p := range partitions {
		if p == it.Partition() {
			break
		}

		if i == len(partitions)-1 {
			return errors.Errorf("partition not in token holder partitions, %s-%s, %s", it.Contract(), ipp.sender, it.Partition())
		}
	}

	for _, ad := range *ipp.operators {
		if ad.Equal(it.Operator()) {
			return errors.Errorf("operator is already in token holder operators, %q", ad)
		}
	}

	for _, ad := range *ipp.tokenHolders {
		if ad.Equal(ipp.sender) {
			return errors.Errorf("sender is already in operator tokenholders, %q", ad)
		}
	}

	if err := crcstate.CheckExistsState(stcurrency.StateKeyCurrencyDesign(it.Currency()), getStateFunc); err != nil {
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

	sts[0] = crcstate.NewStateMergeValue(
		ststo.StateKeyOperatorTokenHolders(it.Contract(), it.Operator(), it.Partition()),
		ststo.NewOperatorTokenHoldersStateValue(holders),
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
