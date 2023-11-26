package sto

import (
	currencytypes "github.com/ProtoconNet/mitum-currency/v3/types"
	stotypes "github.com/ProtoconNet/mitum-sto/types/sto"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
	"github.com/ProtoconNet/mitum2/util/hint"
)

func (it *CreateSecurityTokenItem) unpack(enc encoder.Encoder, ht hint.Hint, ca string, granularity uint64, partition string, cid string) error {
	e := util.StringError("failed to unmarshal CreateSecurityTokenItem")

	it.BaseHinter = hint.NewBaseHinter(ht)
	it.granularity = granularity
	it.defaultPartition = stotypes.Partition(partition)
	it.currency = currencytypes.CurrencyID(cid)

	switch a, err := base.DecodeAddress(ca, enc); {
	case err != nil:
		return e.Wrap(err)
	default:
		it.contract = a
	}

	return nil
}
