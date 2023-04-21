package sto

import (
	extensioncurrency "github.com/ProtoconNet/mitum-currency-extension/v2/currency"
	"github.com/ProtoconNet/mitum-currency/v2/currency"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
)

var CreateSecurityTokensItemHint = hint.MustNewHint("mitum-sto-create-security-tokens-item-v0.0.1")

type CreateSecurityTokensItem struct {
	hint.BaseHinter
	contract         base.Address                 // contract account
	stoID            extensioncurrency.ContractID // token id
	granularity      uint64                       // token granulariry
	defaultPartition Partition                    // default partitions
	controllers      []base.Address               // initial controllers
	currency         currency.CurrencyID          // fee
}

func NewCreateSecurityTokensItem(contract base.Address, stoID extensioncurrency.ContractID, granularity uint64, partition Partition, controllers []base.Address, currency currency.CurrencyID) CreateSecurityTokensItem {
	return CreateSecurityTokensItem{
		BaseHinter:       hint.NewBaseHinter(CreateSecurityTokensItemHint),
		contract:         contract,
		stoID:            stoID,
		granularity:      granularity,
		defaultPartition: partition,
		controllers:      controllers,
		currency:         currency,
	}
}

func (it CreateSecurityTokensItem) Bytes() []byte {
	bc := make([][]byte, len(it.controllers))

	for i, con := range it.controllers {
		bc[i] = con.Bytes()
	}

	return util.ConcatBytesSlice(
		it.contract.Bytes(),
		it.stoID.Bytes(),
		util.Uint64ToBytes(it.granularity),
		it.defaultPartition.Bytes(),
		util.ConcatBytesSlice(bc...),
		it.currency.Bytes(),
	)
}

func (it CreateSecurityTokensItem) IsValid([]byte) error {
	if err := util.CheckIsValiders(nil, false, it.BaseHinter, it.contract, it.stoID, it.defaultPartition, it.currency); err != nil {
		return err
	}

	if it.granularity < uint64(1) {
		return util.ErrInvalid.Errorf("zero granularity")
	}

	founds := map[string]struct{}{}
	for i := range it.controllers {
		if err := it.controllers[i].IsValid(nil); err != nil {
			return err
		}

		if _, found := founds[it.controllers[i].String()]; found {
			return util.ErrInvalid.Errorf("duplicated controller found, %s", it.controllers[i].String())
		}

		founds[it.controllers[i].String()] = struct{}{}
	}

	return nil
}

func (it CreateSecurityTokensItem) Contract() base.Address {
	return it.contract
}

func (it CreateSecurityTokensItem) STO() extensioncurrency.ContractID {
	return it.stoID
}

func (it CreateSecurityTokensItem) Granularity() uint64 {
	return it.granularity
}

func (it CreateSecurityTokensItem) DefaultPartitions() Partition {
	return it.defaultPartition
}

func (it CreateSecurityTokensItem) Controllers() []base.Address {
	return it.controllers
}

func (it CreateSecurityTokensItem) Currency() currency.CurrencyID {
	return it.currency
}

func (it CreateSecurityTokensItem) Addresses() []base.Address {
	ad := make([]base.Address, len(it.controllers)+1)

	ad[0] = it.contract
	for i, con := range it.controllers {
		ad[i+1] = con
	}

	return ad
}
