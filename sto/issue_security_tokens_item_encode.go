package sto

import (
	extensioncurrency "github.com/ProtoconNet/mitum-currency-extension/v2/currency"
	"github.com/ProtoconNet/mitum-currency/v2/currency"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
	"github.com/ProtoconNet/mitum2/util/hint"
)

func (it *IssueSecurityTokensItem) unpack(enc encoder.Encoder, ht hint.Hint, ca, sto, rc, am, p, cid string) error {
	e := util.StringErrorFunc("failed to unmarshal IssueSecurityTokensItem")

	it.BaseHinter = hint.NewBaseHinter(ht)
	it.stoID = extensioncurrency.ContractID(sto)
	it.partition = Partition(p)
	it.currency = currency.CurrencyID(cid)

	switch a, err := base.DecodeAddress(ca, enc); {
	case err != nil:
		return e(err, "")
	default:
		it.contract = a
	}

	switch a, err := base.DecodeAddress(rc, enc); {
	case err != nil:
		return e(err, "")
	default:
		it.receiver = a
	}

	amount, err := currency.NewBigFromString(am)
	if err != nil {
		return e(err, "")
	}
	it.amount = amount

	return nil
}
