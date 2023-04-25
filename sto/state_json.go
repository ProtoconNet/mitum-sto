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
	e := util.StringErrorFunc("failed to decode json of STODesignStateValue")

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
	e := util.StringErrorFunc("failed to decode json of PartitionBalanceStateValue")

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

type TokenHolderPartitionsStateValueJSONMarshaler struct {
	hint.BaseHinter
	Partitions []Partition `json:"partitions"`
}

func (p TokenHolderPartitionsStateValue) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(TokenHolderPartitionsStateValueJSONMarshaler{
		BaseHinter: p.BaseHinter,
		Partitions: p.Partitions,
	})
}

type TokenHolderPartitionsStateValueJSONUnmarshaler struct {
	Partitions []string `json:"partitions"`
}

func (p *TokenHolderPartitionsStateValue) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode json of TokenHolderPartitionsStateValue")

	var u TokenHolderPartitionsStateValueJSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return e(err, "")
	}

	partitions := make([]Partition, len(u.Partitions))
	for i, s := range u.Partitions {
		partitions[i] = Partition(s)
	}

	p.Partitions = partitions

	return nil
}

type TokenHolderPartitionBalanceStateValueJSONMarshaler struct {
	hint.BaseHinter
	Amount    string    `json:"amount"`
	Partition Partition `json:"partition"`
}

func (p TokenHolderPartitionBalanceStateValue) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(TokenHolderPartitionBalanceStateValueJSONMarshaler{
		BaseHinter: p.BaseHinter,
		Amount:     p.Amount.String(),
		Partition:  p.Partition,
	})
}

type TokenHolderPartitionBalanceStateValueJSONUnmarshaler struct {
	Amount    string `json:"amount"`
	Partition string `json:"partition"`
}

func (p *TokenHolderPartitionBalanceStateValue) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode json of TokenHolderPartitionBalanceStateValue")

	var u TokenHolderPartitionBalanceStateValueJSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return e(err, "")
	}

	big, err := currency.NewBigFromString(u.Amount)
	if err != nil {
		return e(err, "")
	}
	p.Amount = big

	p.Partition = Partition(u.Partition)

	return nil
}
