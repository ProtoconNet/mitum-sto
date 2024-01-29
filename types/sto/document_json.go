package sto

import (
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
	"github.com/ProtoconNet/mitum2/util/hint"
)

type DocumentJSONMarshaler struct {
	hint.BaseHinter
	Title string `json:"title"`
	Hash  string `json:"hash"`
	URI   URI    `json:"uri"`
}

func (doc Document) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(DocumentJSONMarshaler{
		BaseHinter: doc.BaseHinter,
		Title:      doc.title,
		Hash:       doc.hash,
		URI:        doc.uri,
	})
}

type DocumentJSONUnmarshaler struct {
	Hint  hint.Hint `json:"_hint"`
	Title string    `json:"title"`
	Hash  string    `json:"hash"`
	URI   string    `json:"uri"`
}

func (doc *Document) DecodeJSON(b []byte, enc encoder.Encoder) error {
	e := util.StringError("failed to decode json of Document")

	var ud DocumentJSONUnmarshaler
	if err := enc.Unmarshal(b, &ud); err != nil {
		return e.Wrap(err)
	}

	return doc.unpack(enc, ud.Hint, ud.Title, ud.Hash, ud.URI)
}
