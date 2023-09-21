package sto // nolint:dupl

import (
	bsonenc "github.com/ProtoconNet/mitum-currency/v3/digest/util/bson"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
	"go.mongodb.org/mongo-driver/bson"
)

func (it RedeemTokensItem) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(
		bson.M{
			"_hint":       it.Hint().String(),
			"contract":    it.contract,
			"tokenholder": it.tokenHolder,
			"amount":      it.amount.String(),
			"partition":   it.partition,
			"currency":    it.currency,
		},
	)
}

type RedeemTokensItemBSONUnmarshaler struct {
	Hint        string `bson:"_hint"`
	Contract    string `bson:"contract"`
	TokenHolder string `bson:"tokenholder"`
	Amount      string `bson:"amount"`
	Partition   string `bson:"partition"`
	Currency    string `bson:"currency"`
}

func (it *RedeemTokensItem) DecodeBSON(b []byte, enc *bsonenc.Encoder) error {
	e := util.StringError("failed to decode bson of RedeemTokensItem")

	var uit RedeemTokensItemBSONUnmarshaler
	if err := bson.Unmarshal(b, &uit); err != nil {
		return e.Wrap(err)
	}

	ht, err := hint.ParseHint(uit.Hint)
	if err != nil {
		return e.Wrap(err)
	}

	return it.unpack(enc, ht, uit.Contract, uit.TokenHolder, uit.Amount, uit.Partition, uit.Currency)
}
