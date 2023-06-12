package kyc

import (
	currencytypes "github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
)

func (fact *CreateKYCServiceFact) unpack(enc encoder.Encoder, sa, ca, kycid string, cons []string, cid string) error {
	e := util.StringErrorFunc("failed to unmarshal CreateKYCServiceFact")

	switch a, err := base.DecodeAddress(sa, enc); {
	case err != nil:
		return e(err, "")
	default:
		fact.sender = a
	}

	switch a, err := base.DecodeAddress(ca, enc); {
	case err != nil:
		return e(err, "")
	default:
		fact.contract = a
	}

	controllers := make([]base.Address, len(cons))
	for i, con := range cons {
		if a, err := base.DecodeAddress(con, enc); err != nil {
			return e(err, "")
		} else {
			controllers[i] = a
		}
	}
	fact.controllers = controllers

	fact.kycID = currencytypes.ContractID(kycid)
	fact.currency = currencytypes.CurrencyID(cid)

	return nil
}
