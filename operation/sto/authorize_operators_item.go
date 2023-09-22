package sto

import (
	currencytypes "github.com/ProtoconNet/mitum-currency/v3/types"
	stotypes "github.com/ProtoconNet/mitum-sto/types/sto"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
)

var AuthorizeOperatorItemHint = hint.MustNewHint("mitum-sto-authorize-operator-item-v0.0.1")

type AuthorizeOperatorItem struct {
	hint.BaseHinter
	contract  base.Address             // contract address
	operator  base.Address             // initial controllers
	partition stotypes.Partition       // partition
	currency  currencytypes.CurrencyID // fee
}

func NewAuthorizeOperatorItem(
	contract base.Address,
	operator base.Address,
	partition stotypes.Partition,
	currency currencytypes.CurrencyID,
) AuthorizeOperatorItem {
	return AuthorizeOperatorItem{
		BaseHinter: hint.NewBaseHinter(AuthorizeOperatorItemHint),
		contract:   contract,
		operator:   operator,
		partition:  partition,
		currency:   currency,
	}
}

func (it AuthorizeOperatorItem) Bytes() []byte {
	return util.ConcatBytesSlice(
		it.contract.Bytes(),
		it.operator.Bytes(),
		it.partition.Bytes(),
		it.currency.Bytes(),
	)
}

func (it AuthorizeOperatorItem) IsValid([]byte) error {
	if err := util.CheckIsValiders(nil, false, it.BaseHinter, it.contract, it.operator, it.partition, it.currency); err != nil {
		return err
	}

	if it.contract.Equal(it.operator) {
		return util.ErrInvalid.Errorf("contract address is same with operator, %q", it.contract)
	}

	return nil
}

func (it AuthorizeOperatorItem) Contract() base.Address {
	return it.contract
}

func (it AuthorizeOperatorItem) Operator() base.Address {
	return it.operator
}

func (it AuthorizeOperatorItem) Partition() stotypes.Partition {
	return it.partition
}

func (it AuthorizeOperatorItem) Currency() currencytypes.CurrencyID {
	return it.currency
}

func (it AuthorizeOperatorItem) Addresses() []base.Address {
	ad := make([]base.Address, 2)

	ad[0] = it.contract
	ad[1] = it.operator

	return ad
}
