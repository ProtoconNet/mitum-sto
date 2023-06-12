package sto

import (
	"github.com/ProtoconNet/mitum-currency/v3/common"
	currencytypes "github.com/ProtoconNet/mitum-currency/v3/types"
	stotypes "github.com/ProtoconNet/mitum-sto/types/sto"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
)

var IssueSecurityTokensItemHint = hint.MustNewHint("mitum-sto-issue-security-tokens-item-v0.0.1")

type IssueSecurityTokensItem struct {
	hint.BaseHinter
	contract  base.Address             // contract
	stoID     currencytypes.ContractID // token id
	receiver  base.Address             // tokenholder
	amount    common.Big               // amount
	partition stotypes.Partition       // partition
	currency  currencytypes.CurrencyID // fee
}

func NewIssueSecurityTokensItem(
	contract base.Address,
	stoID currencytypes.ContractID,
	receiver base.Address,
	amount common.Big,
	partition stotypes.Partition,
	currency currencytypes.CurrencyID,
) IssueSecurityTokensItem {
	return IssueSecurityTokensItem{
		BaseHinter: hint.NewBaseHinter(IssueSecurityTokensItemHint),
		contract:   contract,
		stoID:      stoID,
		receiver:   receiver,
		amount:     amount,
		partition:  partition,
		currency:   currency,
	}
}

func (it IssueSecurityTokensItem) Bytes() []byte {
	return util.ConcatBytesSlice(
		it.contract.Bytes(),
		it.stoID.Bytes(),
		it.receiver.Bytes(),
		it.amount.Bytes(),
		it.partition.Bytes(),
		it.currency.Bytes(),
	)
}

func (it IssueSecurityTokensItem) IsValid([]byte) error {
	if err := util.CheckIsValiders(nil, false,
		it.BaseHinter,
		it.contract,
		it.stoID,
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

func (it IssueSecurityTokensItem) Contract() base.Address {
	return it.contract
}

func (it IssueSecurityTokensItem) STO() currencytypes.ContractID {
	return it.stoID
}

func (it IssueSecurityTokensItem) Receiver() base.Address {
	return it.receiver
}

func (it IssueSecurityTokensItem) Amount() common.Big {
	return it.amount
}

func (it IssueSecurityTokensItem) Partition() stotypes.Partition {
	return it.partition
}

func (it IssueSecurityTokensItem) Currency() currencytypes.CurrencyID {
	return it.currency
}

func (it IssueSecurityTokensItem) Addresses() []base.Address {
	ad := make([]base.Address, 2)

	ad[0] = it.contract
	ad[1] = it.receiver

	return ad
}
