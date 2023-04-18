package sto

import (
	currencyextension "github.com/ProtoconNet/mitum-currency-extension/v2/currency"
	"github.com/ProtoconNet/mitum-currency/v2/currency"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util/hint"
)

type BaseRevokeOperatorsItem struct {
	hint.BaseHinter
	stoID       currencyextension.ContractID // token id
	contract    base.Address                 // contract account
	tokenHolder base.Address                 // token tokenHolder
	operator    base.Address                 // operator account
	partition   Partition                    // partition
	currency    currency.CurrencyID          // fee
}

func NewBaseRevokeOperatorsItem(ht hint.Hint, stoID currencyextension.ContractID, contract, tokenHolder, operator base.Address, partition Partition, currency currency.CurrencyID) BaseRevokeOperatorsItem {
	return BaseRevokeOperatorsItem{
		BaseHinter:  hint.NewBaseHinter(ht),
		stoID:       stoID,
		contract:    contract,
		tokenHolder: tokenHolder,
		operator:    operator,
		partition:   partition,
		currency:    currency,
	}
}
