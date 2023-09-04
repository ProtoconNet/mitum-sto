package sto // nolint: dupl

import (
	"go.mongodb.org/mongo-driver/bson"

	"github.com/ProtoconNet/mitum-currency/v3/common"
	bsonenc "github.com/ProtoconNet/mitum-currency/v3/digest/util/bson"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
	"github.com/ProtoconNet/mitum2/util/valuehash"
)

func (fact SetDocumentFact) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(
		bson.M{
			"_hint":        fact.Hint().String(),
			"sender":       fact.sender,
			"contract":     fact.contract,
			"stoid":        fact.stoID,
			"title":        fact.title,
			"uri":          fact.uri,
			"documenthash": fact.documentHash,
			"currency":     fact.currency,
			"hash":         fact.BaseFact.Hash().String(),
			"token":        fact.BaseFact.Token(),
		},
	)
}

type SetDocumentFactBSONUnmarshaler struct {
	Hint         string `bson:"_hint"`
	Sender       string `bson:"sender"`
	Contract     string `bson:"contract"`
	STOID        string `bson:"stoid"`
	Title        string `bson:"title"`
	Uri          string `bson:"uri"`
	DocumentHash string `bson:"documenthash"`
	Currency     string `bson:"currency"`
}

func (fact *SetDocumentFact) DecodeBSON(b []byte, enc *bsonenc.Encoder) error {
	e := util.StringError("failed to decode bson of SetDocumentFact")

	var ubf common.BaseFactBSONUnmarshaler

	if err := enc.Unmarshal(b, &ubf); err != nil {
		return e.Wrap(err)
	}

	fact.BaseFact.SetHash(valuehash.NewBytesFromString(ubf.Hash))
	fact.BaseFact.SetToken(ubf.Token)

	var uf SetDocumentFactBSONUnmarshaler
	if err := bson.Unmarshal(b, &uf); err != nil {
		return e.Wrap(err)
	}

	ht, err := hint.ParseHint(uf.Hint)
	if err != nil {
		return e.Wrap(err)
	}
	fact.BaseHinter = hint.NewBaseHinter(ht)

	return fact.unpack(enc, uf.Sender, uf.Contract, uf.STOID, uf.Title, uf.Uri, uf.DocumentHash, uf.Currency)
}

func (op SetDocument) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(
		bson.M{
			"_hint": op.Hint().String(),
			"hash":  op.Hash().String(),
			"fact":  op.Fact(),
			"signs": op.Signs(),
		})
}

func (op *SetDocument) DecodeBSON(b []byte, enc *bsonenc.Encoder) error {
	e := util.StringError("failed to decode bson of SetDocument")

	var ubo common.BaseOperation
	if err := ubo.DecodeBSON(b, enc); err != nil {
		return e.Wrap(err)
	}

	op.BaseOperation = ubo

	return nil
}
