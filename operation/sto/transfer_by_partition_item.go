package sto

import (
	"github.com/ProtoconNet/mitum-currency/v3/common"
	currencytypes "github.com/ProtoconNet/mitum-currency/v3/types"
	stotypes "github.com/ProtoconNet/mitum-sto/types/sto"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
)

var TransferByPartitionItemHint = hint.MustNewHint("mitum-sto-transfer-by-partition-item-v0.0.1")

type TransferByPartitionItem struct {
	hint.BaseHinter
	contract    base.Address // contract accounts
	tokenholder base.Address
	receiver    base.Address             // token holder
	partition   stotypes.Partition       // partition
	amount      common.Big               // transfer amount
	currency    currencytypes.CurrencyID // fee
}

func NewTransferByPartitionItem(
	contract base.Address,
	tokenHolder, receiver base.Address,
	partition stotypes.Partition,
	amount common.Big,
	currency currencytypes.CurrencyID,
) TransferByPartitionItem {
	return TransferByPartitionItem{
		BaseHinter:  hint.NewBaseHinter(TransferByPartitionItemHint),
		contract:    contract,
		tokenholder: tokenHolder,
		receiver:    receiver,
		partition:   partition,
		amount:      amount,
		currency:    currency,
	}
}

func (it TransferByPartitionItem) Bytes() []byte {
	return util.ConcatBytesSlice(
		it.contract.Bytes(),
		it.tokenholder.Bytes(),
		it.receiver.Bytes(),
		it.partition.Bytes(),
		it.amount.Bytes(),
		it.currency.Bytes(),
	)
}

func (it TransferByPartitionItem) IsValid([]byte) error {
	if err := util.CheckIsValiders(nil, false,
		it.BaseHinter,
		it.contract,
		it.tokenholder,
		it.receiver,
		it.partition,
		it.currency,
	); err != nil {
		return err
	}

	if !it.amount.OverZero() {
		return util.ErrInvalid.Errorf("amount must be over zero")
	}

	if it.contract.Equal(it.tokenholder) {
		return util.ErrInvalid.Errorf("contract address is same with token holder, %q", it.contract)
	}

	if it.contract.Equal(it.receiver) {
		return util.ErrInvalid.Errorf("contract address is same with receiver, %q", it.contract)
	}

	if it.receiver.Equal(it.tokenholder) {
		return util.ErrInvalid.Errorf("token holder is same with receiver, %q", it.receiver)
	}

	return nil
}

func (it TransferByPartitionItem) Contract() base.Address {
	return it.contract
}

func (it TransferByPartitionItem) TokenHolder() base.Address {
	return it.tokenholder
}

func (it TransferByPartitionItem) Receiver() base.Address {
	return it.receiver
}

func (it TransferByPartitionItem) Amount() common.Big {
	return it.amount
}

func (it TransferByPartitionItem) Partition() stotypes.Partition {
	return it.partition
}

func (it TransferByPartitionItem) Currency() currencytypes.CurrencyID {
	return it.currency
}

func (it TransferByPartitionItem) Addresses() []base.Address {
	ad := make([]base.Address, 3)

	ad[0] = it.contract
	ad[1] = it.receiver
	ad[2] = it.tokenholder

	return ad
}
