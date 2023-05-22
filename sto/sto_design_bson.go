package sto

import (
	"go.mongodb.org/mongo-driver/bson"

	bsonenc "github.com/ProtoconNet/mitum-currency/v2/digest/util/bson"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
)

func (de Design) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(
		bson.M{
			"_hint":       de.Hint().String(),
			"stoid":       de.stoID,
			"granularity": de.granularity,
			"policy":      de.policy,
		},
	)
}

type DesignBSONUnmarshaler struct {
	Hint        string   `bson:"_hint"`
	STO         string   `bson:"stoid"`
	Granularity uint64   `bson:"granularity"`
	Policy      bson.Raw `bson:"policy"`
}

func (de *Design) DecodeBSON(b []byte, enc *bsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode bson of Design")

	var ud DesignBSONUnmarshaler
	if err := enc.Unmarshal(b, &ud); err != nil {
		return e(err, "")
	}

	ht, err := hint.ParseHint(ud.Hint)
	if err != nil {
		return e(err, "")
	}

	return de.unpack(enc, ht, ud.STO, ud.Granularity, ud.Policy)
}
