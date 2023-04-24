package sto // nolint:dupl

import (
	bsonenc "github.com/ProtoconNet/mitum-currency/v2/digest/util/bson"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
	"go.mongodb.org/mongo-driver/bson"
)

func (it RedeemTokensItem) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(
		bson.M{
			"_hint":        it.Hint().String(),
			"contract":     it.contract,
			"stoid":        it.stoID,
			"token_holder": it.tokenHolder,
			"amount":       it.amount.String(),
			"partition":    it.partition,
			"currency":     it.currency,
		},
	)
}

type RedeemTokensItemBSONUnmarshaler struct {
	Hint        string `bson:"_hint"`
	Contract    string `bson:"contract"`
	STO         string `bson:"stoid"`
	TokenHolder string `bson:"token_holder"`
	Amount      string `bson:"amount"`
	Partition   string `bson:"partition"`
	Currency    string `bson:"currency"`
}

func (it *RedeemTokensItem) DecodeBSON(b []byte, enc *bsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode bson of RedeemTokensItem")

	var uit RedeemTokensItemBSONUnmarshaler
	if err := bson.Unmarshal(b, &uit); err != nil {
		return e(err, "")
	}

	ht, err := hint.ParseHint(uit.Hint)
	if err != nil {
		return e(err, "")
	}

	return it.unpack(enc, ht, uit.Contract, uit.STO, uit.TokenHolder, uit.Amount, uit.Partition, uit.Currency)
}
