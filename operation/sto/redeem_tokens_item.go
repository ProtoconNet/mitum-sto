package sto

import (
	"github.com/ProtoconNet/mitum-currency/v3/common"
	currencytypes "github.com/ProtoconNet/mitum-currency/v3/types"
	stotypes "github.com/ProtoconNet/mitum-sto/types/sto"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
)

var RedeemTokensItemHint = hint.MustNewHint("mitum-sto-redeem-tokens-item-v0.0.1")

type RedeemTokensItem struct {
	hint.BaseHinter
	contract    base.Address             // contract account
	tokenHolder base.Address             // token tokenHolder
	amount      common.Big               // redeem amount
	partition   stotypes.Partition       // partition
	currency    currencytypes.CurrencyID // fee
}

func NewRedeemTokensItem(
	contract base.Address,
	tokenHolder base.Address,
	amount common.Big,
	partition stotypes.Partition,
	currency currencytypes.CurrencyID,
) RedeemTokensItem {
	return RedeemTokensItem{
		BaseHinter:  hint.NewBaseHinter(RedeemTokensItemHint),
		contract:    contract,
		tokenHolder: tokenHolder,
		amount:      amount,
		partition:   partition,
		currency:    currency,
	}
}

func (it RedeemTokensItem) Bytes() []byte {
	return util.ConcatBytesSlice(
		it.contract.Bytes(),
		it.tokenHolder.Bytes(),
		it.amount.Bytes(),
		it.partition.Bytes(),
		it.currency.Bytes(),
	)
}

func (it RedeemTokensItem) IsValid([]byte) error {
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

func (it RedeemTokensItem) Contract() base.Address {
	return it.contract
}

func (it RedeemTokensItem) TokenHolder() base.Address {
	return it.tokenHolder
}

func (it RedeemTokensItem) Amount() common.Big {
	return it.amount
}

func (it RedeemTokensItem) Partition() stotypes.Partition {
	return it.partition
}

func (it RedeemTokensItem) Currency() currencytypes.CurrencyID {
	return it.currency
}

func (it RedeemTokensItem) Addresses() []base.Address {
	ad := make([]base.Address, 2)

	ad[0] = it.contract
	ad[1] = it.tokenHolder

	return ad
}
