package sto

import (
	extensioncurrency "github.com/ProtoconNet/mitum-currency-extension/v2/currency"
	"github.com/ProtoconNet/mitum2/util"
	jsonenc "github.com/ProtoconNet/mitum2/util/encoder/json"
	"github.com/ProtoconNet/mitum2/util/hint"
)

type DocumentJSONMarshaler struct {
	hint.BaseHinter
	STO   extensioncurrency.ContractID `json:"sto"`
	Title string                       `json:"title"`
	Hash  string                       `json:"hash"`
	URI   URI                          `json:"uri"`
}

func (doc Document) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(DocumentJSONMarshaler{
		BaseHinter: doc.BaseHinter,
		STO:        doc.stoID,
		Title:      doc.title,
		Hash:       doc.hash,
		URI:        doc.uri,
	})
}

type DocumentJSONUnmarshaler struct {
	Hint  hint.Hint `json:"_hint"`
	STO   string    `json:"sto"`
	Title string    `json:"title"`
	Hash  string    `json:"hash"`
	URI   string    `json:"uri"`
}

func (doc *Document) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode json of Document")

	var ud DocumentJSONUnmarshaler
	if err := enc.Unmarshal(b, &ud); err != nil {
		return e(err, "")
	}

	return doc.unpack(enc, ud.Hint, ud.STO, ud.Title, ud.Hash, ud.URI)
}