package sto

import (
	currencyextension "github.com/ProtoconNet/mitum-currency-extension/v2/currency"
	"github.com/ProtoconNet/mitum-currency/v2/currency"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util/hint"
)

type BaseRedeemTokensItem struct {
	hint.BaseHinter
	stoID       currencyextension.ContractID // token id
	contract    base.Address                 // contract account
	tokenHolder base.Address                 // token tokenHolder
	amount      currency.Big                 // redeem amount
	partition   Partition                    // partition
	currency    currency.CurrencyID          // fee
}

func NewBaseRedeemTokensItem(ht hint.Hint, stoID currencyextension.ContractID, contract, tokenHolder base.Address, amount currency.Big, partition Partition, currency currency.CurrencyID) BaseRedeemTokensItem {
	return BaseRedeemTokensItem{
		BaseHinter:  hint.NewBaseHinter(ht),
		stoID:       stoID,
		contract:    contract,
		tokenHolder: tokenHolder,
		amount:      amount,
		partition:   partition,
		currency:    currency,
	}
}
