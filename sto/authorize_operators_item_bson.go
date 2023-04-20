package sto // nolint:dupl

import (
	bsonenc "github.com/ProtoconNet/mitum-currency/v2/digest/util/bson"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
	"go.mongodb.org/mongo-driver/bson"
)

func (it AuthorizeOperatorsItem) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(
		bson.M{
			"_hint":    it.Hint().String(),
			"stoid":    it.stoID,
			"contract": it.contract,
			"operator": it.operator,
			"currency": it.currency,
		},
	)
}

type AuthorizeOperatorsItemBSONUnmarshaler struct {
	Hint     string `bson:"_hint"`
	STO      string `bson:"stoid"`
	Contract string `bson:"contract"`
	Operator string `bson:"operator"`
	Currency string `bson:"currency"`
}

func (it *AuthorizeOperatorsItem) DecodeBSON(b []byte, enc *bsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode bson of AuthorizeOperatorsItem")

	var uit AuthorizeOperatorsItemBSONUnmarshaler
	if err := bson.Unmarshal(b, &uit); err != nil {
		return e(err, "")
	}

	ht, err := hint.ParseHint(uit.Hint)
	if err != nil {
		return e(err, "")
	}

	return it.unpack(enc, ht, uit.STO, uit.Contract, uit.Operator, uit.Currency)
}
