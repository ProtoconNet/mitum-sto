package sto

import (
	currencyextension "github.com/ProtoconNet/mitum-currency-extension/v2/currency"
	"github.com/ProtoconNet/mitum-currency/v2/currency"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util/hint"
)

type BaseTransferSecurityTokensPartitionsItem struct {
	hint.BaseHinter
	stoID     currencyextension.ContractID // token id
	contract  base.Address                 // contract account
	receiver  base.Address                 // token holder
	partition Partition                    // partition
	amount    currency.Big                 // transfer amount
	currency  currency.CurrencyID          // fee
}

func NewBaseBaseTransferSecurityTokensPartitionsItem(ht hint.Hint, stoID currencyextension.ContractID, contract, receiver base.Address, partition Partition, amount currency.Big, currency currency.CurrencyID) BaseTransferSecurityTokensPartitionsItem {
	return BaseTransferSecurityTokensPartitionsItem{
		BaseHinter: hint.NewBaseHinter(ht),
		stoID:      stoID,
		contract:   contract,
		receiver:   receiver,
		partition:  partition,
		amount:     amount,
		currency:   currency,
	}
}
