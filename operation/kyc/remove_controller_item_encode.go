package kyc

import (
	currencytypes "github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
	"github.com/ProtoconNet/mitum2/util/hint"
)

func (it *RemoveControllerItem) unpack(enc encoder.Encoder, ht hint.Hint, ca, con, cid string) error {
	e := util.StringError("failed to unmarshal RemoveControllerItem")

	it.BaseHinter = hint.NewBaseHinter(ht)
	it.currency = currencytypes.CurrencyID(cid)

	switch a, err := base.DecodeAddress(ca, enc); {
	case err != nil:
		return e.Wrap(err)
	default:
		it.contract = a
	}

	switch a, err := base.DecodeAddress(con, enc); {
	case err != nil:
		return e.Wrap(err)
	default:
		it.controller = a
	}

	return nil
}
