package sto

import (
	currencyextension "github.com/ProtoconNet/mitum-currency-extension/v2/currency"
	"github.com/ProtoconNet/mitum-currency/v2/currency"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
)

type IssueSecurityTokensItem struct {
	hint.BaseHinter
	stoID     currencyextension.ContractID // token id
	contract  base.Address                 // contract
	receiver  base.Address                 // token holder
	amount    currency.Big                 // amount
	partition Partition                    // partition
	currency  currency.CurrencyID          // fee
}

func NewIssueSecurityTokensItem(
	ht hint.Hint,
	stoID currencyextension.ContractID,
	contract, receiver base.Address,
	amount currency.Big,
	partition Partition,
	currency currency.CurrencyID,
) IssueSecurityTokensItem {
	return IssueSecurityTokensItem{
		BaseHinter: hint.NewBaseHinter(ht),
		stoID:      stoID,
		contract:   contract,
		receiver:   receiver,
		amount:     amount,
		partition:  partition,
		currency:   currency,
	}
}

func (it IssueSecurityTokensItem) Bytes() []byte {
	return util.ConcatBytesSlice(
		it.stoID.Bytes(),
		it.contract.Bytes(),
		it.receiver.Bytes(),
		it.amount.Bytes(),
		it.partition.Bytes(),
		it.currency.Bytes(),
	)
}

func (it IssueSecurityTokensItem) IsValid([]byte) error {
	if err := util.CheckIsValiders(nil, false,
		it.BaseHinter,
		it.stoID,
		it.contract,
		it.receiver,
		it.partition,
		it.currency,
	); err != nil {
		return err
	}

	if !it.amount.OverZero() {
		return util.ErrInvalid.Errorf("amount must be over zero")
	}

	if it.contract.Equal(it.receiver) {
		return util.ErrInvalid.Errorf("contract address is same with receiver, %q", it.contract)
	}

	return nil
}

func (it IssueSecurityTokensItem) STO() currencyextension.ContractID {
	return it.stoID
}

func (it IssueSecurityTokensItem) Contract() base.Address {
	return it.contract
}

func (it IssueSecurityTokensItem) Receiver() (base.Address, error) {
	return it.receiver, nil
}

func (it IssueSecurityTokensItem) Amount() currency.Big {
	return it.amount
}

func (it IssueSecurityTokensItem) Partition() Partition {
	return it.partition
}

func (it IssueSecurityTokensItem) Currency() currency.CurrencyID {
	return it.currency
}

func (it IssueSecurityTokensItem) Addresses() []base.Address {
	ad := make([]base.Address, 2)

	ad[0] = it.contract
	ad[1] = it.receiver

	return ad
}
