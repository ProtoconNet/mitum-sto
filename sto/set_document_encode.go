package sto

import (
	extensioncurrency "github.com/ProtoconNet/mitum-currency-extension/v2/currency"
	"github.com/ProtoconNet/mitum-currency/v2/currency"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
)

func (fact *SetDocumentFact) unpack(enc encoder.Encoder, sa, ca, stoid, title, uri, dochash, cid string) error {
	e := util.StringErrorFunc("failed to unmarshal SetDocumentFact")

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

	fact.stoID = extensioncurrency.ContractID(stoid)
	fact.title = title
	fact.uri = URI(uri)
	fact.documentHash = dochash
	fact.currency = currency.CurrencyID(cid)

	return nil
}
