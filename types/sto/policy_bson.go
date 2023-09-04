package sto

import (
	"go.mongodb.org/mongo-driver/bson"

	bsonenc "github.com/ProtoconNet/mitum-currency/v3/digest/util/bson"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
)

func (po Policy) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(
		bson.M{
			"_hint":       po.Hint().String(),
			"partitions":  po.partitions,
			"aggregate":   po.aggregate.String(),
			"controllers": po.controllers,
			"documents":   po.documents,
		},
	)
}

type PolicyBSONUnmarshaler struct {
	Hint        string   `bson:"_hint"`
	Partitions  []string `bson:"partitions"`
	Aggregate   string   `bson:"aggregate"`
	Controllers []string `bson:"controllers"`
	Documents   bson.Raw `bson:"documents"`
}

func (po *Policy) DecodeBSON(b []byte, enc *bsonenc.Encoder) error {
	e := util.StringError("failed to decode bson of Policy")

	var upo PolicyBSONUnmarshaler
	if err := enc.Unmarshal(b, &upo); err != nil {
		return e.Wrap(err)
	}

	ht, err := hint.ParseHint(upo.Hint)
	if err != nil {
		return e.Wrap(err)
	}

	return po.unpack(enc, ht, upo.Partitions, upo.Aggregate, upo.Controllers, upo.Documents)
}
