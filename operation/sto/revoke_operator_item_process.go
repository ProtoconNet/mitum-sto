package sto

import (
	"context"
	"sync"

	crcstate "github.com/ProtoconNet/mitum-currency/v3/state"
	"github.com/ProtoconNet/mitum-currency/v3/state/currency"
	stextension "github.com/ProtoconNet/mitum-currency/v3/state/extension"
	ststo "github.com/ProtoconNet/mitum-sto/state/sto"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/pkg/errors"
)

var revokeOperatorItemProcessorPool = sync.Pool{
	New: func() interface{} {
		return new(RevokeOperatorItemProcessor)
	},
}

type RevokeOperatorItemProcessor struct {
	h            util.Hash
	sender       base.Address
	item         RevokeOperatorItem
	operators    *[]base.Address
	tokenHolders *[]base.Address
}

func (ipp *RevokeOperatorItemProcessor) PreProcess(
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

	if len(*ipp.operators) == 0 {
		return errors.Errorf("empty token holder operators, %s-%s-%s", it.Contract(), it.Partition(), ipp.sender)
	}

	for i, ad := range *ipp.operators {
		if ad.Equal(it.Operator()) {
			break
		}

		if i == len(*ipp.operators)-1 {
			return errors.Errorf("operator not in token holder operators, %s-%s-%s, %q", it.Contract(), it.Partition(), ipp.sender, it.Operator())
		}
	}

	if len(*ipp.tokenHolders) == 0 {
		return errors.Errorf("empty operator tokenholders, %s-%s-%s", it.Contract(), it.Partition(), it.Operator())
	}

	for i, ad := range *ipp.tokenHolders {
		if ad.Equal(ipp.sender) {
			break
		}

		if i == len(*ipp.tokenHolders)-1 {
			return errors.Errorf("sender not in operator tokenholders, %s-%s-%s, %q", it.Contract(), it.Partition(), it.Operator(), ipp.sender)
		}
	}

	if err := crcstate.CheckExistsState(currency.StateKeyCurrencyDesign(it.Currency()), getStateFunc); err != nil {
		return err
	}

	return nil
}

func (ipp *RevokeOperatorItemProcessor) Process(
	ctx context.Context, op base.Operation, getStateFunc base.GetStateFunc,
) ([]base.StateMergeValue, error) {
	sts := make([]base.StateMergeValue, 1)

	it := ipp.item

	if len(*ipp.operators) == 0 {
		return nil, errors.Errorf("empty token holder operators, %s-%s-%s", it.Contract(), it.Partition(), ipp.sender)
	}

	for i, ad := range *ipp.operators {
		if ad.Equal(it.Operator()) {
			if i < len(*ipp.operators)-1 {
				copy((*ipp.operators)[i:], (*ipp.operators)[i+1:])
			}
			*ipp.operators = (*ipp.operators)[:len(*ipp.operators)-1]
			break
		}

		if i == len(*ipp.operators)-1 {
			return nil, errors.Errorf("operator not in token holder operators, %s-%s-%s, %q", it.Contract(), it.Partition(), ipp.sender, it.Operator())
		}
	}

	holders := *ipp.tokenHolders
	if len(holders) == 0 {
		return nil, errors.Errorf("empty operator tokenholders, %s-%s-%s", it.Contract(), it.Partition(), it.Operator())
	}

	for i, ad := range holders {
		if ad.Equal(ipp.sender) {
			if i < len(holders)-1 {
				copy((holders)[i:], (holders)[i+1:])
			}
			holders = (holders)[:len(holders)-1]
			break
		}

		if i == len(holders)-1 {
			return nil, errors.Errorf("sender not in operator tokenholders, %s-%s-%s, %q", it.Contract(), it.Partition(), it.Operator(), ipp.sender)
		}
	}

	sts[0] = crcstate.NewStateMergeValue(
		ststo.StateKeyOperatorTokenHolders(it.Contract(), it.Operator(), it.Partition()),
		ststo.NewOperatorTokenHoldersStateValue(holders),
	)

	return sts, nil
}

func (ipp *RevokeOperatorItemProcessor) Close() error {
	ipp.h = nil
	ipp.sender = nil
	ipp.item = RevokeOperatorItem{}
	ipp.operators = nil
	ipp.tokenHolders = nil

	revokeOperatorItemProcessorPool.Put(ipp)

	return nil
}
