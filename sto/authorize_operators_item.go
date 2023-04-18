package sto

import (
	currencyextension "github.com/ProtoconNet/mitum-currency-extension/v2/currency"
	"github.com/ProtoconNet/mitum-currency/v2/currency"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util/hint"
)

type BaseAuthorizeOperatorsItem struct {
	hint.BaseHinter
	stoID    currencyextension.ContractID // token id
	contract base.Address
	operator base.Address        // initial controllers
	currency currency.CurrencyID // fee
}

func NewBaseAuthorizeOperatorsItem(ht hint.Hint, stoID currencyextension.ContractID, contract, operator base.Address, currency currency.CurrencyID) BaseAuthorizeOperatorsItem {
	return BaseAuthorizeOperatorsItem{
		BaseHinter: hint.NewBaseHinter(ht),
		stoID:      stoID,
		contract:   contract,
		operator:   operator,
		currency:   currency,
	}
}
