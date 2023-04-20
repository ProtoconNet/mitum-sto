package sto // nolint:dupl

import (
	bsonenc "github.com/ProtoconNet/mitum-currency/v2/digest/util/bson"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
	"go.mongodb.org/mongo-driver/bson"
)

func (it CreateSecurityTokensItem) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(
		bson.M{
			"_hint":             it.Hint().String(),
			"stoid":             it.stoID,
			"granularity":       it.granularity,
			"default_partition": it.defaultPartition,
			"controllers":       it.controllers,
			"currency":          it.currency,
		},
	)
}

type CreateSecurityTokensItemBSONUnmarshaler struct {
	Hint             string   `bson:"_hint"`
	STO              string   `bson:"stoid"`
	Granularity      uint64   `bson:"granularity"`
	DefaultPartition string   `bson:"default_partition"`
	Controllers      []string `bson:"controllers"`
	Currency         string   `bson:"currency"`
}

func (it *CreateSecurityTokensItem) DecodeBSON(b []byte, enc *bsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode bson of CreateSecurityTokensItem")

	var uit CreateSecurityTokensItemBSONUnmarshaler
	if err := bson.Unmarshal(b, &uit); err != nil {
		return e(err, "")
	}

	ht, err := hint.ParseHint(uit.Hint)
	if err != nil {
		return e(err, "")
	}

	return it.unpack(enc, ht, uit.STO, uit.Granularity, uit.DefaultPartition, uit.Controllers, uit.Currency)
}
