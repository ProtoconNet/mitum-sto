package sto

import (
	extensioncurrency "github.com/ProtoconNet/mitum-currency-extension/v2/currency"
	"github.com/ProtoconNet/mitum-currency/v2/currency"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
	"github.com/ProtoconNet/mitum2/util/hint"
)

func (it *BaseAuthorizeOperatorsItem) unpack(enc encoder.Encoder, ht hint.Hint, sto, ca, op, cid string) error {
	e := util.StringErrorFunc("failed to unmarshal BaseAuthorizeOperatorsItem")

	it.BaseHinter = hint.NewBaseHinter(ht)
	it.stoID = extensioncurrency.ContractID(sto)
	it.currency = currency.CurrencyID(cid)

	switch a, err := base.DecodeAddress(ca, enc); {
	case err != nil:
		return e(err, "")
	default:
		it.contract = a
	}

	switch a, err := base.DecodeAddress(op, enc); {
	case err != nil:
		return e(err, "")
	default:
		it.operator = a
	}

	return nil
}
