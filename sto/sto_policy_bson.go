package sto

import (
	"go.mongodb.org/mongo-driver/bson"

	bsonenc "github.com/ProtoconNet/mitum-currency/v2/digest/util/bson"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
)

func (po STOPolicy) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(
		bson.M{
			"_hint":       po.Hint().String(),
			"partitions":  po.partitions,
			"aggregate":   po.aggregate,
			"controllers": po.controllers,
			"documents":   po.documents,
		},
	)
}

type STOPolicyBSONUnmarshaler struct {
	Hint        string   `bson:"_hint"`
	Partitions  bson.Raw `json:"partitions"`
	Aggregate   bson.Raw `json:"aggregate"`
	Controllers []string `json:"controllers"`
	Documents   bson.Raw `json:"documents"`
}

func (po *STOPolicy) DecodeBSON(b []byte, enc *bsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode bson of STOPolicy")

	var upo STOPolicyBSONUnmarshaler
	if err := enc.Unmarshal(b, &upo); err != nil {
		return e(err, "")
	}

	ht, err := hint.ParseHint(upo.Hint)
	if err != nil {
		return e(err, "")
	}

	return po.unpack(enc, ht, upo.Partitions, upo.Aggregate, upo.Controllers, upo.Documents)
}
