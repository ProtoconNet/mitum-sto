package kyc // nolint: dupl

import (
	"go.mongodb.org/mongo-driver/bson"

	"github.com/ProtoconNet/mitum-currency/v2/currency"
	bsonenc "github.com/ProtoconNet/mitum-currency/v2/digest/util/bson"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
	"github.com/ProtoconNet/mitum2/util/valuehash"
)

func (fact CreateKYCServiceFact) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(
		bson.M{
			"_hint":       fact.Hint().String(),
			"sender":      fact.sender,
			"contract":    fact.contract,
			"kycid":       fact.kycID,
			"controllers": fact.controllers,
			"currency":    fact.currency,
			"hash":        fact.BaseFact.Hash().String(),
			"token":       fact.BaseFact.Token(),
		},
	)
}

type CreateKYCServiceFactBSONUnmarshaler struct {
	Hint         string   `bson:"_hint"`
	Sender       string   `bson:"sender"`
	Contract     string   `bson:"contract"`
	KYCID        string   `bson:"kycid"`
	Controllers  []string `bson:"controllers"`
	DocumentHash string   `bson:"documenthash"`
	Currency     string   `bson:"currency"`
}

func (fact *CreateKYCServiceFact) DecodeBSON(b []byte, enc *bsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode bson of CreateKYCServiceFact")

	var ubf currency.BaseFactBSONUnmarshaler

	if err := enc.Unmarshal(b, &ubf); err != nil {
		return e(err, "")
	}

	fact.BaseFact.SetHash(valuehash.NewBytesFromString(ubf.Hash))
	fact.BaseFact.SetToken(ubf.Token)

	var uf CreateKYCServiceFactBSONUnmarshaler
	if err := bson.Unmarshal(b, &uf); err != nil {
		return e(err, "")
	}

	ht, err := hint.ParseHint(uf.Hint)
	if err != nil {
		return e(err, "")
	}
	fact.BaseHinter = hint.NewBaseHinter(ht)

	return fact.unpack(enc, uf.Sender, uf.Contract, uf.KYCID, uf.Controllers, uf.Currency)
}

func (op CreateKYCService) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(
		bson.M{
			"_hint": op.Hint().String(),
			"hash":  op.Hash().String(),
			"fact":  op.Fact(),
			"signs": op.Signs(),
		})
}

func (op *CreateKYCService) DecodeBSON(b []byte, enc *bsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode bson of CreateKYCService")

	var ubo currency.BaseOperation
	if err := ubo.DecodeBSON(b, enc); err != nil {
		return e(err, "")
	}

	op.BaseOperation = ubo

	return nil
}
