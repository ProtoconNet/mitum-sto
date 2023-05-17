package kyc // nolint:dupl

import (
	bsonenc "github.com/ProtoconNet/mitum-currency/v2/digest/util/bson"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
	"go.mongodb.org/mongo-driver/bson"
)

func (it AddCustomersItem) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(
		bson.M{
			"_hint":    it.Hint().String(),
			"contract": it.contract,
			"kycid":    it.kycID,
			"customer": it.customer,
			"status":   it.status,
			"currency": it.currency,
		},
	)
}

type AddCustomersItemBSONUnmarshaler struct {
	Hint     string `bson:"_hint"`
	Contract string `bson:"contract"`
	KYC      string `bson:"kycid"`
	Customer string `bson:"customer"`
	Status   bool   `bson:"status"`
	Currency string `bson:"currency"`
}

func (it *AddCustomersItem) DecodeBSON(b []byte, enc *bsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode bson of AddCustomersItem")

	var uit AddCustomersItemBSONUnmarshaler
	if err := bson.Unmarshal(b, &uit); err != nil {
		return e(err, "")
	}

	ht, err := hint.ParseHint(uit.Hint)
	if err != nil {
		return e(err, "")
	}

	return it.unpack(enc, ht, uit.Contract, uit.KYC, uit.Customer, uit.Status, uit.Currency)
}
