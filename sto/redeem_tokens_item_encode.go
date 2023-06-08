package sto

import (
	currencybase "github.com/ProtoconNet/mitum-currency/v3/base"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
	"github.com/ProtoconNet/mitum2/util/hint"
)

func (it *RedeemTokensItem) unpack(enc encoder.Encoder, ht hint.Hint, ca, sto, th, am, p, cid string) error {
	e := util.StringErrorFunc("failed to unmarshal RedeemTokensItem")

	it.BaseHinter = hint.NewBaseHinter(ht)
	it.stoID = currencybase.ContractID(sto)
	it.partition = Partition(p)
	it.currency = currencybase.CurrencyID(cid)

	switch a, err := base.DecodeAddress(ca, enc); {
	case err != nil:
		return e(err, "")
	default:
		it.contract = a
	}

	switch a, err := base.DecodeAddress(th, enc); {
	case err != nil:
		return e(err, "")
	default:
		it.tokenHolder = a
	}

	amount, err := currencybase.NewBigFromString(am)
	if err != nil {
		return e(err, "")
	}
	it.amount = amount

	return nil
}
