package kyc // nolint:dupl

import (
	bsonenc "github.com/ProtoconNet/mitum-currency/v3/digest/util/bson"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
	"go.mongodb.org/mongo-driver/bson"
)

func (it AddControllersItem) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(
		bson.M{
			"_hint":      it.Hint().String(),
			"contract":   it.contract,
			"kycid":      it.kycID,
			"controller": it.controller,
			"currency":   it.currency,
		},
	)
}

type AddControllersItemBSONUnmarshaler struct {
	Hint       string `bson:"_hint"`
	Contract   string `bson:"contract"`
	KYC        string `bson:"kycid"`
	Controller string `bson:"controller"`
	Currency   string `bson:"currency"`
}

func (it *AddControllersItem) DecodeBSON(b []byte, enc *bsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode bson of AddControllersItem")

	var uit AddControllersItemBSONUnmarshaler
	if err := bson.Unmarshal(b, &uit); err != nil {
		return e(err, "")
	}

	ht, err := hint.ParseHint(uit.Hint)
	if err != nil {
		return e(err, "")
	}

	return it.unpack(enc, ht, uit.Contract, uit.KYC, uit.Controller, uit.Currency)
}
