package sto

import (
	currencybase "github.com/ProtoconNet/mitum-currency/v3/base"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
)

var RevokeOperatorsItemHint = hint.MustNewHint("mitum-sto-revoke-operators-item-v0.0.1")

type RevokeOperatorsItem struct {
	hint.BaseHinter
	contract  base.Address            // contract account
	stoID     currencybase.ContractID // token id
	operator  base.Address            // operator account
	partition Partition               // partition
	currency  currencybase.CurrencyID // fee
}

func NewRevokeOperatorsItem(
	contract base.Address,
	stoID currencybase.ContractID,
	operator base.Address,
	partition Partition,
	currency currencybase.CurrencyID,
) RevokeOperatorsItem {
	return RevokeOperatorsItem{
		BaseHinter: hint.NewBaseHinter(RevokeOperatorsItemHint),
		contract:   contract,
		stoID:      stoID,
		operator:   operator,
		partition:  partition,
		currency:   currency,
	}
}

func (it RevokeOperatorsItem) Bytes() []byte {
	return util.ConcatBytesSlice(
		it.contract.Bytes(),
		it.stoID.Bytes(),
		it.operator.Bytes(),
		it.partition.Bytes(),
		it.currency.Bytes(),
	)
}

func (it RevokeOperatorsItem) IsValid([]byte) error {
	if err := util.CheckIsValiders(nil, false,
		it.BaseHinter,
		it.contract,
		it.stoID,
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

func (it RevokeOperatorsItem) STO() currencybase.ContractID {
	return it.stoID
}

func (it RevokeOperatorsItem) Operator() base.Address {
	return it.operator
}

func (it RevokeOperatorsItem) Partition() Partition {
	return it.partition
}

func (it RevokeOperatorsItem) Currency() currencybase.CurrencyID {
	return it.currency
}

func (it RevokeOperatorsItem) Addresses() []base.Address {
	ad := make([]base.Address, 2)

	ad[0] = it.contract
	ad[2] = it.operator

	return ad
}
