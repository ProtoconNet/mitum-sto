package kyc

import (
	currencytypes "github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
	"github.com/ProtoconNet/mitum2/util/hint"
)

func (it *UpdateCustomersItem) unpack(enc encoder.Encoder, ht hint.Hint, ca, ctm string, st bool, cid string) error {
	e := util.StringError("failed to unmarshal UpdateCustomersItem")

	it.BaseHinter = hint.NewBaseHinter(ht)
	it.status = st
	it.currency = currencytypes.CurrencyID(cid)

	switch a, err := base.DecodeAddress(ca, enc); {
	case err != nil:
		return e.Wrap(err)
	default:
		it.contract = a
	}

	switch a, err := base.DecodeAddress(ctm, enc); {
	case err != nil:
		return e.Wrap(err)
	default:
		it.customer = a
	}

	return nil
}
