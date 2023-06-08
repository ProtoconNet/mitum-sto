package sto

import (
	currencybase "github.com/ProtoconNet/mitum-currency/v3/base"
	"github.com/ProtoconNet/mitum2/util/encoder"
	"github.com/ProtoconNet/mitum2/util/hint"
)

func (doc *Document) unpack(enc encoder.Encoder, ht hint.Hint, sto, title, hash, uri string) error {
	doc.BaseHinter = hint.NewBaseHinter(ht)
	doc.stoID = currencybase.ContractID(sto)
	doc.title = title
	doc.hash = hash
	doc.uri = URI(uri)

	return nil
}
