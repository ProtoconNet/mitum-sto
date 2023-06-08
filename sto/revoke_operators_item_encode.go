package sto

import (
	currencybase "github.com/ProtoconNet/mitum-currency/v3/base"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
	"github.com/ProtoconNet/mitum2/util/hint"
)

func (it *RevokeOperatorsItem) unpack(enc encoder.Encoder, ht hint.Hint, ca, sto, oper, p, cid string) error {
	e := util.StringErrorFunc("failed to unmarshal RevokeOperatorsItem")

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

	switch a, err := base.DecodeAddress(oper, enc); {
	case err != nil:
		return e(err, "")
	default:
		it.operator = a
	}

	return nil
}
