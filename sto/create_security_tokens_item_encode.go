package sto

import (
	extensioncurrency "github.com/ProtoconNet/mitum-currency-extension/v2/currency"
	"github.com/ProtoconNet/mitum-currency/v2/currency"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
	"github.com/ProtoconNet/mitum2/util/hint"
)

func (it *CreateSecurityTokensItem) unpack(enc encoder.Encoder, ht hint.Hint, ca, sto string, granularity uint64, partition string, bcs []string, cid string) error {
	e := util.StringErrorFunc("failed to unmarshal CreateSecurityTokensItem")

	it.BaseHinter = hint.NewBaseHinter(ht)
	it.stoID = extensioncurrency.ContractID(sto)
	it.granularity = granularity
	it.defaultPartition = Partition(partition)
	it.currency = currency.CurrencyID(cid)

	switch a, err := base.DecodeAddress(ca, enc); {
	case err != nil:
		return e(err, "")
	default:
		it.contract = a
	}

	controllers := make([]base.Address, len(bcs))
	for i := range bcs {
		ctr, err := base.DecodeAddress(bcs[i], enc)
		if err != nil {
			return e(err, "")
		}
		controllers[i] = ctr
	}
	it.controllers = controllers

	return nil
}
