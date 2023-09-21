package sto // nolint:dupl

import (
	bsonenc "github.com/ProtoconNet/mitum-currency/v3/digest/util/bson"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
	"go.mongodb.org/mongo-driver/bson"
)

func (it TransferSecurityTokensPartitionItem) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(
		bson.M{
			"_hint":       it.Hint().String(),
			"contract":    it.contract,
			"tokenholder": it.tokenholder,
			"receiver":    it.receiver,
			"partition":   it.partition,
			"amount":      it.amount.String(),
			"currency":    it.currency,
		},
	)
}

type TransferSecurityTokensPartitionItemBSONUnmarshaler struct {
	Hint        string `bson:"_hint"`
	Contract    string `bson:"contract"`
	TokenHolder string `bson:"tokenholder"`
	Receiver    string `bson:"receiver"`
	Partition   string `bson:"partition"`
	Amount      string `bson:"amount"`
	Currency    string `bson:"currency"`
}

func (it *TransferSecurityTokensPartitionItem) DecodeBSON(b []byte, enc *bsonenc.Encoder) error {
	e := util.StringError("failed to decode bson of TransferSecurityTokensPartitionItem")

	var uit TransferSecurityTokensPartitionItemBSONUnmarshaler
	if err := bson.Unmarshal(b, &uit); err != nil {
		return e.Wrap(err)
	}

	ht, err := hint.ParseHint(uit.Hint)
	if err != nil {
		return e.Wrap(err)
	}

	return it.unpack(enc, ht, uit.Contract, uit.TokenHolder, uit.Receiver, uit.Partition, uit.Amount, uit.Currency)
}
