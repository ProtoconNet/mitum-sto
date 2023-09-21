package sto

import (
	"go.mongodb.org/mongo-driver/bson"

	bsonenc "github.com/ProtoconNet/mitum-currency/v3/digest/util/bson"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
)

func (doc Document) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(
		bson.M{
			"_hint": doc.Hint().String(),
			"title": doc.title,
			"hash":  doc.hash,
			"uri":   doc.uri,
		},
	)
}

type DocumentBSONUnmarshaler struct {
	Hint  string `bson:"_hint"`
	Title string `bson:"title"`
	Hash  string `bson:"hash"`
	URI   string `bson:"uri"`
}

func (doc *Document) DecodeBSON(b []byte, enc *bsonenc.Encoder) error {
	e := util.StringError("failed to decode bson of Document")

	var ud DocumentBSONUnmarshaler
	if err := enc.Unmarshal(b, &ud); err != nil {
		return e.Wrap(err)
	}

	ht, err := hint.ParseHint(ud.Hint)
	if err != nil {
		return e.Wrap(err)
	}

	return doc.unpack(enc, ht, ud.Title, ud.Hash, ud.URI)
}
