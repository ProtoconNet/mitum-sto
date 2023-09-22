package sto

import (
	currencytypes "github.com/ProtoconNet/mitum-currency/v3/types"
	stotypes "github.com/ProtoconNet/mitum-sto/types/sto"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
	"github.com/ProtoconNet/mitum2/util/hint"
)

func (it *AuthorizeOperatorItem) unpack(enc encoder.Encoder, ht hint.Hint, ca, op, pt, cid string) error {
	e := util.StringError("failed to unmarshal AuthorizeOperatorItem")

	it.BaseHinter = hint.NewBaseHinter(ht)
	it.partition = stotypes.Partition(pt)
	it.currency = currencytypes.CurrencyID(cid)

	switch a, err := base.DecodeAddress(ca, enc); {
	case err != nil:
		return e.Wrap(err)
	default:
		it.contract = a
	}

	switch a, err := base.DecodeAddress(op, enc); {
	case err != nil:
		return e.Wrap(err)
	default:
		it.operator = a
	}

	return nil
}
