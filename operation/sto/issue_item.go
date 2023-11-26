package sto

import (
	"github.com/ProtoconNet/mitum-currency/v3/common"
	currencytypes "github.com/ProtoconNet/mitum-currency/v3/types"
	stotypes "github.com/ProtoconNet/mitum-sto/types/sto"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
)

var IssueItemHint = hint.MustNewHint("mitum-sto-issue-item-v0.0.1")

type IssueItem struct {
	hint.BaseHinter
	contract  base.Address             // contract
	receiver  base.Address             // token holder
	amount    common.Big               // amount
	partition stotypes.Partition       // partition
	currency  currencytypes.CurrencyID // fee
}

func NewIssueItem(
	contract base.Address,
	receiver base.Address,
	amount common.Big,
	partition stotypes.Partition,
	currency currencytypes.CurrencyID,
) IssueItem {
	return IssueItem{
		BaseHinter: hint.NewBaseHinter(IssueItemHint),
		contract:   contract,
		receiver:   receiver,
		amount:     amount,
		partition:  partition,
		currency:   currency,
	}
}

func (it IssueItem) Bytes() []byte {
	return util.ConcatBytesSlice(
		it.contract.Bytes(),
		it.receiver.Bytes(),
		it.amount.Bytes(),
		it.partition.Bytes(),
		it.currency.Bytes(),
	)
}

func (it IssueItem) IsValid([]byte) error {
	if err := util.CheckIsValiders(nil, false,
		it.BaseHinter,
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

func (it IssueItem) Contract() base.Address {
	return it.contract
}

func (it IssueItem) Receiver() base.Address {
	return it.receiver
}

func (it IssueItem) Amount() common.Big {
	return it.amount
}

func (it IssueItem) Partition() stotypes.Partition {
	return it.partition
}

func (it IssueItem) Currency() currencytypes.CurrencyID {
	return it.currency
}

func (it IssueItem) Addresses() []base.Address {
	ad := make([]base.Address, 2)

	ad[0] = it.contract
	ad[1] = it.receiver

	return ad
}
