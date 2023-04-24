package sto

import (
	extensioncurrency "github.com/ProtoconNet/mitum-currency-extension/v2/currency"
	"github.com/ProtoconNet/mitum-currency/v2/currency"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
)

var TransferSecurityTokensPartitionItemHint = hint.MustNewHint("mitum-sto-transfer-security-tokens-partition-item-v0.0.1")

type TransferSecurityTokensPartitionItem struct {
	hint.BaseHinter
	contract  base.Address                 // contract accounts
	stoID     extensioncurrency.ContractID // token id
	receiver  base.Address                 // token holder
	partition Partition                    // partition
	amount    currency.Big                 // transfer amount
	currency  currency.CurrencyID          // fee
}

func NewTransferSecurityTokensPartitionItem(
	contract base.Address,
	stoID extensioncurrency.ContractID,
	receiver base.Address,
	partition Partition,
	amount currency.Big,
	currency currency.CurrencyID,
) TransferSecurityTokensPartitionItem {
	return TransferSecurityTokensPartitionItem{
		BaseHinter: hint.NewBaseHinter(TransferSecurityTokensPartitionItemHint),
		contract:   contract,
		stoID:      stoID,
		receiver:   receiver,
		partition:  partition,
		amount:     amount,
		currency:   currency,
	}
}

func (it TransferSecurityTokensPartitionItem) Bytes() []byte {
	return util.ConcatBytesSlice(
		it.contract.Bytes(),
		it.stoID.Bytes(),
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

func (it TransferSecurityTokensPartitionItem) Contract() base.Address {
	return it.contract
}

func (it TransferSecurityTokensPartitionItem) STO() extensioncurrency.ContractID {
	return it.stoID
}

func (it TransferSecurityTokensPartitionItem) Receiver() base.Address {
	return it.receiver
}

func (it TransferSecurityTokensPartitionItem) Amount() currency.Big {
	return it.amount
}

func (it TransferSecurityTokensPartitionItem) Partition() Partition {
	return it.partition
}

func (it TransferSecurityTokensPartitionItem) Currency() currency.CurrencyID {
	return it.currency
}

func (it TransferSecurityTokensPartitionItem) Addresses() []base.Address {
	ad := make([]base.Address, 2)

	ad[0] = it.contract
	ad[1] = it.receiver

	return ad
}
