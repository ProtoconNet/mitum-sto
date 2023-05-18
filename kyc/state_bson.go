package kyc

import (
	bsonenc "github.com/ProtoconNet/mitum-currency/v2/digest/util/bson"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
	"go.mongodb.org/mongo-driver/bson"
)

func (de KYCDesignStateValue) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(
		bson.M{
			"_hint": de.Hint().String(),
			"kyc":   de.Design,
		},
	)
}

type KYCDesignStateValueBSONUnmarshaler struct {
	Hint string   `bson:"_hint"`
	KYC  bson.Raw `bson:"kyc"`
}

func (de *KYCDesignStateValue) DecodeBSON(b []byte, enc *bsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode bson of KYCDesignStateValue")

	var u KYCDesignStateValueBSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return e(err, "")
	}

	ht, err := hint.ParseHint(u.Hint)
	if err != nil {
		return e(err, "")
	}

	de.BaseHinter = hint.NewBaseHinter(ht)

	var design KYCDesign
	if err := design.DecodeBSON(u.KYC, enc); err != nil {
		return e(err, "")
	}

	de.Design = design

	return nil
}

func (cm CustomerStateValue) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(
		bson.M{
			"_hint":  cm.Hint().String(),
			"status": cm.status,
		},
	)
}

type CustomerStateValueBSONUnmarshaler struct {
	Hint   string `bson:"_hint"`
	Status bool   `bson:"status"`
}

func (cm *CustomerStateValue) DecodeBSON(b []byte, enc *bsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode bson of CustomerStateValue")

	var u CustomerStateValueBSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return e(err, "")
	}

	ht, err := hint.ParseHint(u.Hint)
	if err != nil {
		return e(err, "")
	}

	cm.BaseHinter = hint.NewBaseHinter(ht)

	cm.status = Status(u.Status)

	return nil
}
