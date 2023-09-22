package sto

import (
	"github.com/ProtoconNet/mitum-currency/v3/common"
	currencytypes "github.com/ProtoconNet/mitum-currency/v3/types"
	stotypes "github.com/ProtoconNet/mitum-sto/types/sto"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
)

var RedeemItemHint = hint.MustNewHint("mitum-sto-redeem-item-v0.0.1")

type RedeemItem struct {
	hint.BaseHinter
	contract    base.Address             // contract account
	tokenHolder base.Address             // token tokenHolder
	amount      common.Big               // redeem amount
	partition   stotypes.Partition       // partition
	currency    currencytypes.CurrencyID // fee
}

func NewRedeemItem(
	contract base.Address,
	tokenHolder base.Address,
	amount common.Big,
	partition stotypes.Partition,
	currency currencytypes.CurrencyID,
) RedeemItem {
	return RedeemItem{
		BaseHinter:  hint.NewBaseHinter(RedeemItemHint),
		contract:    contract,
		tokenHolder: tokenHolder,
		amount:      amount,
		partition:   partition,
		currency:    currency,
	}
}

func (it RedeemItem) Bytes() []byte {
	return util.ConcatBytesSlice(
		it.contract.Bytes(),
		it.tokenHolder.Bytes(),
		it.amount.Bytes(),
		it.partition.Bytes(),
		it.currency.Bytes(),
	)
}

func (it RedeemItem) IsValid([]byte) error {
	if err := util.CheckIsValiders(nil, false,
		it.BaseHinter,
		it.contract,
		it.tokenHolder,
		it.partition,
		it.currency,
	); err != nil {
		return err
	}

	if !it.amount.OverZero() {
		return util.ErrInvalid.Errorf("amount must be over zero")
	}

	if it.contract.Equal(it.tokenHolder) {
		return util.ErrInvalid.Errorf("contract address is same with tokenholder, %q", it.contract)
	}

	return nil
}

func (it RedeemItem) Contract() base.Address {
	return it.contract
}

func (it RedeemItem) TokenHolder() base.Address {
	return it.tokenHolder
}

func (it RedeemItem) Amount() common.Big {
	return it.amount
}

func (it RedeemItem) Partition() stotypes.Partition {
	return it.partition
}

func (it RedeemItem) Currency() currencytypes.CurrencyID {
	return it.currency
}

func (it RedeemItem) Addresses() []base.Address {
	ad := make([]base.Address, 2)

	ad[0] = it.contract
	ad[1] = it.tokenHolder

	return ad
}
