package kyc // nolint: dupl

import (
	"go.mongodb.org/mongo-driver/bson"

	"github.com/ProtoconNet/mitum-currency/v2/currency"
	bsonenc "github.com/ProtoconNet/mitum-currency/v2/digest/util/bson"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
	"github.com/ProtoconNet/mitum2/util/valuehash"
)

func (fact AddCustomersFact) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(
		bson.M{
			"_hint":  fact.Hint().String(),
			"sender": fact.sender,
			"items":  fact.items,
			"hash":   fact.BaseFact.Hash().String(),
			"token":  fact.BaseFact.Token(),
		},
	)
}

type AddCustomersFactBSONUnmarshaler struct {
	Hint   string   `bson:"_hint"`
	Sender string   `bson:"sender"`
	Items  bson.Raw `bson:"items"`
}

func (fact *AddCustomersFact) DecodeBSON(b []byte, enc *bsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode bson of AddCustomersFact")

	var ubf currency.BaseFactBSONUnmarshaler

	if err := enc.Unmarshal(b, &ubf); err != nil {
		return e(err, "")
	}

	fact.BaseFact.SetHash(valuehash.NewBytesFromString(ubf.Hash))
	fact.BaseFact.SetToken(ubf.Token)

	var uf AddCustomersFactBSONUnmarshaler
	if err := bson.Unmarshal(b, &uf); err != nil {
		return e(err, "")
	}

	ht, err := hint.ParseHint(uf.Hint)
	if err != nil {
		return e(err, "")
	}
	fact.BaseHinter = hint.NewBaseHinter(ht)

	return fact.unpack(enc, uf.Sender, uf.Items)
}

func (op AddCustomers) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(
		bson.M{
			"_hint": op.Hint().String(),
			"hash":  op.Hash().String(),
			"fact":  op.Fact(),
			"signs": op.Signs(),
		})
}

func (op *AddCustomers) DecodeBSON(b []byte, enc *bsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode bson of AddCustomers")

	var ubo currency.BaseOperation
	if err := ubo.DecodeBSON(b, enc); err != nil {
		return e(err, "")
	}

	op.BaseOperation = ubo

	return nil
}
