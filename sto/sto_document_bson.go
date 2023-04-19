package sto

import (
	"go.mongodb.org/mongo-driver/bson"

	bsonenc "github.com/ProtoconNet/mitum-currency/v2/digest/util/bson"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
)

func (doc Document) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(
		bson.M{
			"_hint": doc.Hint().String(),
			"sto":   doc.stoID,
			"title": doc.title,
			"hash":  doc.hash,
			"uri":   doc.uri,
		},
	)
}

type DocumentBSONUnmarshaler struct {
	Hint  string `bson:"_hint"`
	STO   string `json:"sto"`
	Title string `json:"title"`
	Hash  string `json:"hash"`
	URI   string `json:"uri"`
}

func (doc *Document) DecodeBSON(b []byte, enc *bsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode bson of Document")

	var ud DocumentBSONUnmarshaler
	if err := enc.Unmarshal(b, &ud); err != nil {
		return e(err, "")
	}

	ht, err := hint.ParseHint(ud.Hint)
	if err != nil {
		return e(err, "")
	}

	return doc.unpack(enc, ht, ud.STO, ud.Title, ud.Hash, ud.URI)
}
