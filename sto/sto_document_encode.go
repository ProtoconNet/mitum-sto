package sto

import (
	extensioncurrency "github.com/ProtoconNet/mitum-currency-extension/v2/currency"
	"github.com/ProtoconNet/mitum2/util/encoder"
	"github.com/ProtoconNet/mitum2/util/hint"
)

func (doc *Document) unpack(enc encoder.Encoder, ht hint.Hint, sto, title, hash, uri string) error {
	doc.BaseHinter = hint.NewBaseHinter(ht)
	doc.stoID = extensioncurrency.ContractID(sto)
	doc.title = title
	doc.hash = hash
	doc.uri = URI(uri)

	return nil
}
