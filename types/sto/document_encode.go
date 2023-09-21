package sto

import (
	"github.com/ProtoconNet/mitum2/util/encoder"
	"github.com/ProtoconNet/mitum2/util/hint"
)

func (doc *Document) unpack(_ encoder.Encoder, ht hint.Hint, title, hash, uri string) error {
	doc.BaseHinter = hint.NewBaseHinter(ht)
	doc.title = title
	doc.hash = hash
	doc.uri = URI(uri)

	return nil
}
