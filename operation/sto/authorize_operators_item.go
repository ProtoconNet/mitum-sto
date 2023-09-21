package sto

import (
	currencytypes "github.com/ProtoconNet/mitum-currency/v3/types"
	stotypes "github.com/ProtoconNet/mitum-sto/types/sto"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
)

var AuthorizeOperatorsItemHint = hint.MustNewHint("mitum-sto-authorize-operators-item-v0.0.1")

type AuthorizeOperatorsItem struct {
	hint.BaseHinter
	contract  base.Address             // contract address
	operator  base.Address             // initial controllers
	partition stotypes.Partition       // partition
	currency  currencytypes.CurrencyID // fee
}

func NewAuthorizeOperatorsItem(
	contract base.Address,
	operator base.Address,
	partition stotypes.Partition,
	currency currencytypes.CurrencyID,
) AuthorizeOperatorsItem {
	return AuthorizeOperatorsItem{
		BaseHinter: hint.NewBaseHinter(AuthorizeOperatorsItemHint),
		contract:   contract,
		operator:   operator,
		partition:  partition,
		currency:   currency,
	}
}

func (it AuthorizeOperatorsItem) Bytes() []byte {
	return util.ConcatBytesSlice(
		it.contract.Bytes(),
		it.operator.Bytes(),
		it.partition.Bytes(),
		it.currency.Bytes(),
	)
}

func (it AuthorizeOperatorsItem) IsValid([]byte) error {
	if err := util.CheckIsValiders(nil, false, it.BaseHinter, it.contract, it.operator, it.partition, it.currency); err != nil {
		return err
	}

	if it.contract.Equal(it.operator) {
		return util.ErrInvalid.Errorf("contract address is same with operator, %q", it.contract)
	}

	return nil
}

func (it AuthorizeOperatorsItem) Contract() base.Address {
	return it.contract
}

func (it AuthorizeOperatorsItem) Operator() base.Address {
	return it.operator
}

func (it AuthorizeOperatorsItem) Partition() stotypes.Partition {
	return it.partition
}

func (it AuthorizeOperatorsItem) Currency() currencytypes.CurrencyID {
	return it.currency
}

func (it AuthorizeOperatorsItem) Addresses() []base.Address {
	ad := make([]base.Address, 2)

	ad[0] = it.contract
	ad[1] = it.operator

	return ad
}
