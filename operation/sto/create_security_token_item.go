package sto

import (
	currencytypes "github.com/ProtoconNet/mitum-currency/v3/types"
	stotypes "github.com/ProtoconNet/mitum-sto/types/sto"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
)

var CreateSecurityTokenItemHint = hint.MustNewHint("mitum-sto-create-security-token-item-v0.0.1")

type CreateSecurityTokenItem struct {
	hint.BaseHinter
	contract         base.Address             // contract account
	granularity      uint64                   // token granulariry
	defaultPartition stotypes.Partition       // default partitions
	currency         currencytypes.CurrencyID // fee
}

func NewCreateSecurityTokenItem(
	contract base.Address,
	granularity uint64,
	partition stotypes.Partition,
	currency currencytypes.CurrencyID,
) CreateSecurityTokenItem {
	return CreateSecurityTokenItem{
		BaseHinter:       hint.NewBaseHinter(CreateSecurityTokenItemHint),
		contract:         contract,
		granularity:      granularity,
		defaultPartition: partition,
		currency:         currency,
	}
}

func (it CreateSecurityTokenItem) Bytes() []byte {
	return util.ConcatBytesSlice(
		it.contract.Bytes(),
		util.Uint64ToBytes(it.granularity),
		it.defaultPartition.Bytes(),
		it.currency.Bytes(),
	)
}

func (it CreateSecurityTokenItem) IsValid([]byte) error {
	if err := util.CheckIsValiders(nil, false,
		it.BaseHinter,
		it.contract,
		it.defaultPartition,
		it.currency,
	); err != nil {
		return err
	}

	if it.granularity < uint64(1) {
		return util.ErrInvalid.Errorf("zero granularity")
	}

	return nil
}

func (it CreateSecurityTokenItem) Contract() base.Address {
	return it.contract
}

func (it CreateSecurityTokenItem) Granularity() uint64 {
	return it.granularity
}

func (it CreateSecurityTokenItem) DefaultPartition() stotypes.Partition {
	return it.defaultPartition
}

func (it CreateSecurityTokenItem) Currency() currencytypes.CurrencyID {
	return it.currency
}

func (it CreateSecurityTokenItem) Addresses() []base.Address {
	return []base.Address{it.contract}
}
