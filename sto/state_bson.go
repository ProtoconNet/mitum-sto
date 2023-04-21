package sto

import (
	"github.com/ProtoconNet/mitum-currency/v2/currency"
	bsonenc "github.com/ProtoconNet/mitum-currency/v2/digest/util/bson"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
	"go.mongodb.org/mongo-driver/bson"
)

func (de STODesignStateValue) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(
		bson.M{
			"_hint": de.Hint().String(),
			"sto":   de.Design,
		},
	)
}

type STODesignStateValueBSONUnmarshaler struct {
	Hint string   `bson:"_hint"`
	STO  bson.Raw `bson:"sto"`
}

func (de *STODesignStateValue) DecodeBSON(b []byte, enc *bsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode STODesignStateValue")

	var u STODesignStateValueBSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return e(err, "")
	}

	ht, err := hint.ParseHint(u.Hint)
	if err != nil {
		return e(err, "")
	}

	de.BaseHinter = hint.NewBaseHinter(ht)

	var design STODesign
	if err := design.DecodeBSON(u.STO, enc); err != nil {
		return e(err, "")
	}

	de.Design = design

	return nil
}

func (p PartitionBalanceStateValue) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(
		bson.M{
			"_hint":  p.Hint().String(),
			"amount": p.Amount.String(),
		},
	)
}

type PartitionBalanceStateValueBSONUnmarshaler struct {
	Hint   string `bson:"_hint"`
	Amount string `bson:"amount"`
}

func (de *PartitionBalanceStateValue) DecodeBSON(b []byte, enc *bsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode PartitionBalanceStateValue")

	var u PartitionBalanceStateValueBSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return e(err, "")
	}

	ht, err := hint.ParseHint(u.Hint)
	if err != nil {
		return e(err, "")
	}

	de.BaseHinter = hint.NewBaseHinter(ht)

	big, err := currency.NewBigFromString(u.Amount)
	if err != nil {
		return err
	}
	de.Amount = big

	return nil
}
