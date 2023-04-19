package sto

import (
	currencyextension "github.com/ProtoconNet/mitum-currency-extension/v2/currency"
	"github.com/ProtoconNet/mitum-currency/v2/currency"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util/hint"
)

type BaseIssueSecurityTokensItem struct {
	hint.BaseHinter
	stoID     currencyextension.ContractID // token id
	contract  base.Address                 // contract
	receiver  base.Address                 // token holder
	amount    currency.Big                 // amount
	partition Partition                    // partition
	currency  currency.CurrencyID          // fee
}

func NewBaseIssueSecurityTokensItem(ht hint.Hint, stoID currencyextension.ContractID, contract, receiver base.Address, amount currency.Big, partition Partition, currency currency.CurrencyID) BaseIssueSecurityTokensItem {
	return BaseIssueSecurityTokensItem{
		BaseHinter: hint.NewBaseHinter(ht),
		stoID:      stoID,
		contract:   contract,
		receiver:   receiver,
		amount:     amount,
		partition:  partition,
		currency:   currency,
	}
}
