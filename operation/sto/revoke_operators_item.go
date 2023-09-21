package sto

import (
	currencytypes "github.com/ProtoconNet/mitum-currency/v3/types"
	stotypes "github.com/ProtoconNet/mitum-sto/types/sto"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
)

var RevokeOperatorsItemHint = hint.MustNewHint("mitum-sto-revoke-operators-item-v0.0.1")

type RevokeOperatorsItem struct {
	hint.BaseHinter
	contract  base.Address             // contract account
	operator  base.Address             // operator account
	partition stotypes.Partition       // partition
	currency  currencytypes.CurrencyID // fee
}

func NewRevokeOperatorsItem(
	contract base.Address,
	operator base.Address,
	partition stotypes.Partition,
	currency currencytypes.CurrencyID,
) RevokeOperatorsItem {
	return RevokeOperatorsItem{
		BaseHinter: hint.NewBaseHinter(RevokeOperatorsItemHint),
		contract:   contract,
		operator:   operator,
		partition:  partition,
		currency:   currency,
	}
}

func (it RevokeOperatorsItem) Bytes() []byte {
	return util.ConcatBytesSlice(
		it.contract.Bytes(),
		it.operator.Bytes(),
		it.partition.Bytes(),
		it.currency.Bytes(),
	)
}

func (it RevokeOperatorsItem) IsValid([]byte) error {
	if err := util.CheckIsValiders(nil, false,
		it.BaseHinter,
		it.contract,
		it.operator,
		it.partition,
		it.currency,
	); err != nil {
		return err
	}

	if it.contract.Equal(it.operator) {
		return util.ErrInvalid.Errorf("contract address is same with operator, %q", it.contract)
	}

	return nil
}

func (it RevokeOperatorsItem) Contract() base.Address {
	return it.contract
}

func (it RevokeOperatorsItem) Operator() base.Address {
	return it.operator
}

func (it RevokeOperatorsItem) Partition() stotypes.Partition {
	return it.partition
}

func (it RevokeOperatorsItem) Currency() currencytypes.CurrencyID {
	return it.currency
}

func (it RevokeOperatorsItem) Addresses() []base.Address {
	ad := make([]base.Address, 2)

	ad[0] = it.contract
	ad[2] = it.operator

	return ad
}
