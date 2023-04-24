package sto

import (
	extensioncurrency "github.com/ProtoconNet/mitum-currency-extension/v2/currency"
	"github.com/ProtoconNet/mitum-currency/v2/currency"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
)

var RevokeOperatorsItemHint = hint.MustNewHint("mitum-sto-revoke-operators-item-v0.0.1")

type RevokeOperatorsItem struct {
	hint.BaseHinter
	contract    base.Address                 // contract account
	stoID       extensioncurrency.ContractID // token id
	tokenHolder base.Address                 // token tokenHolder
	operator    base.Address                 // operator account
	partition   Partition                    // partition
	currency    currency.CurrencyID          // fee
}

func NewRevokeOperatorsItem(
	contract base.Address,
	stoID extensioncurrency.ContractID,
	tokenHolder, operator base.Address,
	partition Partition,
	currency currency.CurrencyID,
) RevokeOperatorsItem {
	return RevokeOperatorsItem{
		BaseHinter:  hint.NewBaseHinter(RevokeOperatorsItemHint),
		contract:    contract,
		stoID:       stoID,
		tokenHolder: tokenHolder,
		operator:    operator,
		partition:   partition,
		currency:    currency,
	}
}

func (it RevokeOperatorsItem) Bytes() []byte {
	return util.ConcatBytesSlice(
		it.contract.Bytes(),
		it.stoID.Bytes(),
		it.tokenHolder.Bytes(),
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
		it.tokenHolder,
		it.operator,
		it.partition,
		it.currency,
	); err != nil {
		return err
	}

	if it.contract.Equal(it.tokenHolder) {
		return util.ErrInvalid.Errorf("contract address is same with token holder, %q", it.contract)
	}

	if it.contract.Equal(it.operator) {
		return util.ErrInvalid.Errorf("contract address is same with operator, %q", it.contract)
	}

	if it.tokenHolder.Equal(it.operator) {
		return util.ErrInvalid.Errorf("token holder address is same with operator, %q", it.operator)
	}

	return nil
}

func (it RevokeOperatorsItem) Contract() base.Address {
	return it.contract
}

func (it RevokeOperatorsItem) STO() extensioncurrency.ContractID {
	return it.stoID
}

func (it RevokeOperatorsItem) TokenHolder() base.Address {
	return it.tokenHolder
}

func (it RevokeOperatorsItem) Operator() base.Address {
	return it.operator
}

func (it RevokeOperatorsItem) Partition() Partition {
	return it.partition
}

func (it RevokeOperatorsItem) Currency() currency.CurrencyID {
	return it.currency
}

func (it RevokeOperatorsItem) Addresses() []base.Address {
	ad := make([]base.Address, 3)

	ad[0] = it.contract
	ad[1] = it.tokenHolder
	ad[2] = it.operator

	return ad
}
