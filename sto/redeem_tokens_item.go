package sto

import (
	extensioncurrency "github.com/ProtoconNet/mitum-currency-extension/v2/currency"
	"github.com/ProtoconNet/mitum-currency/v2/currency"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
)

var RedeemTokensItemHint = hint.MustNewHint("mitum-sto-redeem-tokens-item-v0.0.1")

type RedeemTokensItem struct {
	hint.BaseHinter
	contract    base.Address                 // contract account
	stoID       extensioncurrency.ContractID // token id
	tokenHolder base.Address                 // token tokenHolder
	amount      currency.Big                 // redeem amount
	partition   Partition                    // partition
	currency    currency.CurrencyID          // fee
}

func NewRedeemTokensItem(
	contract base.Address,
	stoID extensioncurrency.ContractID,
	tokenHolder base.Address,
	amount currency.Big,
	partition Partition,
	currency currency.CurrencyID,
) RedeemTokensItem {
	return RedeemTokensItem{
		BaseHinter:  hint.NewBaseHinter(RedeemTokensItemHint),
		contract:    contract,
		stoID:       stoID,
		tokenHolder: tokenHolder,
		amount:      amount,
		partition:   partition,
		currency:    currency,
	}
}

func (it RedeemTokensItem) Bytes() []byte {
	return util.ConcatBytesSlice(
		it.contract.Bytes(),
		it.stoID.Bytes(),
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
		it.stoID,
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

func (it RedeemTokensItem) STO() extensioncurrency.ContractID {
	return it.stoID
}

func (it RedeemTokensItem) TokenHolder() base.Address {
	return it.tokenHolder
}

func (it RedeemTokensItem) Amount() currency.Big {
	return it.amount
}

func (it RedeemTokensItem) Partition() Partition {
	return it.partition
}

func (it RedeemTokensItem) Currency() currency.CurrencyID {
	return it.currency
}

func (it RedeemTokensItem) Addresses() []base.Address {
	ad := make([]base.Address, 2)

	ad[0] = it.contract
	ad[1] = it.tokenHolder

	return ad
}
