package sto

import (
	currencyextension "github.com/ProtoconNet/mitum-currency-extension/v2/currency"
	"github.com/ProtoconNet/mitum-currency/v2/currency"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util/hint"
)

type BaseCreateSecurityTokensItem struct {
	hint.BaseHinter
	stoID            currencyextension.ContractID // token id
	granularity      uint64                       // token granulariry
	defaultPartition Partition                    // default partitions
	controllers      []base.Address               // initial controllers
	currency         currency.CurrencyID          // fee
}

func NewBaseCreateSecurityTokensItem(ht hint.Hint, stoID currencyextension.ContractID, granularity uint64, partition Partition, controllers []base.Address, currency currency.CurrencyID) BaseCreateSecurityTokensItem {
	return BaseCreateSecurityTokensItem{
		BaseHinter:       hint.NewBaseHinter(ht),
		stoID:            stoID,
		granularity:      granularity,
		defaultPartition: partition,
		controllers:      controllers,
		currency:         currency,
	}
}
