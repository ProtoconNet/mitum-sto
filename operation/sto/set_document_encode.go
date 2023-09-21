package sto

import (
	currencytypes "github.com/ProtoconNet/mitum-currency/v3/types"
	stotypes "github.com/ProtoconNet/mitum-sto/types/sto"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
)

func (fact *SetDocumentFact) unpack(enc encoder.Encoder, sa, ca, title, uri, dochash, cid string) error {
	e := util.StringError("failed to unmarshal SetDocumentFact")

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

	fact.title = title
	fact.uri = stotypes.URI(uri)
	fact.documentHash = dochash
	fact.currency = currencytypes.CurrencyID(cid)

	return nil
}
