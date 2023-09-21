package sto

import (
	"github.com/ProtoconNet/mitum-currency/v3/common"
	currencytypes "github.com/ProtoconNet/mitum-currency/v3/types"
	stotypes "github.com/ProtoconNet/mitum-sto/types/sto"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
)

var TransferSecurityTokensPartitionItemHint = hint.MustNewHint("mitum-sto-transfer-security-tokens-partition-item-v0.0.1")

type TransferSecurityTokensPartitionItem struct {
	hint.BaseHinter
	contract    base.Address // contract accounts
	tokenholder base.Address
	receiver    base.Address             // token holder
	partition   stotypes.Partition       // partition
	amount      common.Big               // transfer amount
	currency    currencytypes.CurrencyID // fee
}

func NewTransferSecurityTokensPartitionItem(
	contract base.Address,
	tokenholder, receiver base.Address,
	partition stotypes.Partition,
	amount common.Big,
	currency currencytypes.CurrencyID,
) TransferSecurityTokensPartitionItem {
	return TransferSecurityTokensPartitionItem{
		BaseHinter:  hint.NewBaseHinter(TransferSecurityTokensPartitionItemHint),
		contract:    contract,
		tokenholder: tokenholder,
		receiver:    receiver,
		partition:   partition,
		amount:      amount,
		currency:    currency,
	}
}

func (it TransferSecurityTokensPartitionItem) Bytes() []byte {
	return util.ConcatBytesSlice(
		it.contract.Bytes(),
		it.tokenholder.Bytes(),
		it.receiver.Bytes(),
		it.partition.Bytes(),
		it.amount.Bytes(),
		it.currency.Bytes(),
	)
}

func (it TransferSecurityTokensPartitionItem) IsValid([]byte) error {
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
		return util.ErrInvalid.Errorf("contract address is same with tokenholder, %q", it.contract)
	}

	if it.contract.Equal(it.receiver) {
		return util.ErrInvalid.Errorf("contract address is same with receiver, %q", it.contract)
	}

	if it.receiver.Equal(it.tokenholder) {
		return util.ErrInvalid.Errorf("tokenholder is same with receiver, %q", it.receiver)
	}

	return nil
}

func (it TransferSecurityTokensPartitionItem) Contract() base.Address {
	return it.contract
}

func (it TransferSecurityTokensPartitionItem) TokenHolder() base.Address {
	return it.tokenholder
}

func (it TransferSecurityTokensPartitionItem) Receiver() base.Address {
	return it.receiver
}

func (it TransferSecurityTokensPartitionItem) Amount() common.Big {
	return it.amount
}

func (it TransferSecurityTokensPartitionItem) Partition() stotypes.Partition {
	return it.partition
}

func (it TransferSecurityTokensPartitionItem) Currency() currencytypes.CurrencyID {
	return it.currency
}

func (it TransferSecurityTokensPartitionItem) Addresses() []base.Address {
	ad := make([]base.Address, 3)

	ad[0] = it.contract
	ad[1] = it.receiver
	ad[2] = it.tokenholder

	return ad
}
