package sto

import (
	"encoding/json"

	"github.com/ProtoconNet/mitum-currency/v2/currency"
	"github.com/ProtoconNet/mitum2/util"
	jsonenc "github.com/ProtoconNet/mitum2/util/encoder/json"
	"github.com/ProtoconNet/mitum2/util/hint"
)

type STODesignStateValueJSONMarshaler struct {
	hint.BaseHinter
	STO STODesign `json:"sto"`
}

func (de STODesignStateValue) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(STODesignStateValueJSONMarshaler{
		BaseHinter: de.BaseHinter,
		STO:        de.Design,
	})
}

type STODesignStateValueJSONUnmarshaler struct {
	STO json.RawMessage `json:"sto"`
}

func (de *STODesignStateValue) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode STODesignStateValue")

	var u STODesignStateValueJSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return e(err, "")
	}

	var design STODesign

	if err := de.DecodeJSON(u.STO, enc); err != nil {
		return e(err, "")
	}

	de.Design = design

	return nil
}

type PartitionBalanceStateValueJSONMarshaler struct {
	hint.BaseHinter
	Amount string `json:"amount"`
}

func (p PartitionBalanceStateValue) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(PartitionBalanceStateValueJSONMarshaler{
		BaseHinter: p.BaseHinter,
		Amount:     p.Amount.String(),
	})
}

type PartitionBalanceStateValueJSONUnmarshaler struct {
	Amount string `json:"amount"`
}

func (p *PartitionBalanceStateValue) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode PartitionBalancetateValue")

	var u PartitionBalanceStateValueJSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return e(err, "")
	}

	big, err := currency.NewBigFromString(u.Amount)
	if err != nil {
		return e(err, "")
	}

	p.Amount = big

	return nil
}
