package sto

import (
	currencyextension "github.com/ProtoconNet/mitum-currency-extension/v2/currency"
	"github.com/ProtoconNet/mitum-currency/v2/currency"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util/hint"
)

type BaseIssueSecurityTokensItem struct {
	hint.BaseHinter
	stoID    currencyextension.ContractID // token id
	receiver base.Address                 // token holder
	currency currency.CurrencyID          // fee
}

func NewBaseIssueSecurityTokensItem(ht hint.Hint, stoID currencyextension.ContractID, receiver base.Address, currency currency.CurrencyID) BaseIssueSecurityTokensItem {
	return BaseIssueSecurityTokensItem{
		BaseHinter: hint.NewBaseHinter(ht),
		stoID:      stoID,
		receiver:   receiver,
		currency:   currency,
	}
}
