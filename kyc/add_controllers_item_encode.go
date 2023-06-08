package kyc

import (
	currencybase "github.com/ProtoconNet/mitum-currency/v3/base"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
	"github.com/ProtoconNet/mitum2/util/hint"
)

func (it *AddControllersItem) unpack(enc encoder.Encoder, ht hint.Hint, ca, kyc, con, cid string) error {
	e := util.StringErrorFunc("failed to unmarshal AddControllersItem")

	it.BaseHinter = hint.NewBaseHinter(ht)
	it.kycID = currencybase.ContractID(kyc)
	it.currency = currencybase.CurrencyID(cid)

	switch a, err := base.DecodeAddress(ca, enc); {
	case err != nil:
		return e(err, "")
	default:
		it.contract = a
	}

	switch a, err := base.DecodeAddress(con, enc); {
	case err != nil:
		return e(err, "")
	default:
		it.controller = a
	}

	return nil
}
