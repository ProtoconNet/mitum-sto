package kyc

import (
	currencytypes "github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
)

func (fact *CreateServiceFact) unpack(enc encoder.Encoder, sa, ca string, cons []string, cid string) error {
	e := util.StringError("failed to unmarshal CreateServiceFact")

	switch a, err := base.DecodeAddress(sa, enc); {
	case err != nil:
		return e.Wrap(err)
	default:
		fact.sender = a
	}

	switch a, err := base.DecodeAddress(ca, enc); {
	case err != nil:
		return e.Wrap(err)
	default:
		fact.contract = a
	}

	controllers := make([]base.Address, len(cons))
	for i, con := range cons {
		if a, err := base.DecodeAddress(con, enc); err != nil {
			return e.Wrap(err)
		} else {
			controllers[i] = a
		}
	}
	fact.controllers = controllers

	fact.currency = currencytypes.CurrencyID(cid)

	return nil
}
