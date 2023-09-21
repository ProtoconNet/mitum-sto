package sto

import (
	currencytypes "github.com/ProtoconNet/mitum-currency/v3/types"
	stotypes "github.com/ProtoconNet/mitum-sto/types/sto"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
)

var CreateSecurityTokensItemHint = hint.MustNewHint("mitum-sto-create-security-tokens-item-v0.0.1")

type CreateSecurityTokensItem struct {
	hint.BaseHinter
	contract         base.Address             // contract account
	granularity      uint64                   // token granulariry
	defaultPartition stotypes.Partition       // default partitions
	controllers      []base.Address           // initial controllers
	currency         currencytypes.CurrencyID // fee
}

func NewCreateSecurityTokensItem(
	contract base.Address,
	granularity uint64,
	partition stotypes.Partition,
	controllers []base.Address,
	currency currencytypes.CurrencyID,
) CreateSecurityTokensItem {
	return CreateSecurityTokensItem{
		BaseHinter:       hint.NewBaseHinter(CreateSecurityTokensItemHint),
		contract:         contract,
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
		util.Uint64ToBytes(it.granularity),
		it.defaultPartition.Bytes(),
		util.ConcatBytesSlice(bc...),
		it.currency.Bytes(),
	)
}

func (it CreateSecurityTokensItem) IsValid([]byte) error {
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

	founds := map[string]struct{}{}
	for _, con := range it.controllers {
		if err := con.IsValid(nil); err != nil {
			return err
		}

		if con.Equal(it.contract) {
			return util.ErrInvalid.Errorf("controller address is same with contract, %q", con)
		}

		if _, found := founds[con.String()]; found {
			return util.ErrInvalid.Errorf("duplicated controller found, %s", con.String())
		}

		founds[con.String()] = struct{}{}
	}

	return nil
}

func (it CreateSecurityTokensItem) Contract() base.Address {
	return it.contract
}

func (it CreateSecurityTokensItem) Granularity() uint64 {
	return it.granularity
}

func (it CreateSecurityTokensItem) DefaultPartition() stotypes.Partition {
	return it.defaultPartition
}

func (it CreateSecurityTokensItem) Controllers() []base.Address {
	return it.controllers
}

func (it CreateSecurityTokensItem) Currency() currencytypes.CurrencyID {
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
