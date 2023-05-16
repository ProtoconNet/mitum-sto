package kyc

import (
	extensioncurrency "github.com/ProtoconNet/mitum-currency-extension/v2/currency"
	"github.com/ProtoconNet/mitum-currency/v2/currency"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
	"github.com/ProtoconNet/mitum2/util/hint"
)

func (it *UpdateCustomersItem) unpack(enc encoder.Encoder, ht hint.Hint, ca, kyc, ctm string, st bool, cid string) error {
	e := util.StringErrorFunc("failed to unmarshal UpdateCustomersItem")

	it.BaseHinter = hint.NewBaseHinter(ht)
	it.kycID = extensioncurrency.ContractID(kyc)
	it.status = st
	it.currency = currency.CurrencyID(cid)

	switch a, err := base.DecodeAddress(ca, enc); {
	case err != nil:
		return e(err, "")
	default:
		it.contract = a
	}

	switch a, err := base.DecodeAddress(ctm, enc); {
	case err != nil:
		return e(err, "")
	default:
		it.customer = a
	}

	return nil
}
