package sto

import (
	"encoding/json"

	"github.com/ProtoconNet/mitum-currency/v2/currency"
	"github.com/ProtoconNet/mitum2/base"
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

	if err := design.DecodeJSON(u.STO, enc); err != nil {
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

type TokenHolderPartitionOperatorsStateValueJSONMarshaler struct {
	hint.BaseHinter
	Operators []base.Address `json:"operators"`
}

func (ops TokenHolderPartitionOperatorsStateValue) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(TokenHolderPartitionOperatorsStateValueJSONMarshaler{
		BaseHinter: ops.BaseHinter,
		Operators:  ops.Operators,
	})
}

type TokenHolderPartitionOperatorsStateValueJSONUnmarshaler struct {
	Operators []string `json:"operators"`
}

func (ops *TokenHolderPartitionOperatorsStateValue) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode json of TokenHolderPartitionOperatorsStateValue")

	var u TokenHolderPartitionOperatorsStateValueJSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return e(err, "")
	}

	operators := make([]base.Address, len(u.Operators))
	for i := range u.Operators {
		a, err := base.DecodeAddress(u.Operators[i], enc)
		if err != nil {
			return e(err, "")
		}
		operators[i] = a
	}
	ops.Operators = operators

	return nil
}

type OperatorTokenHoldersStateValueJSONMarshaler struct {
	hint.BaseHinter
	TokenHolders []base.Address `json:"tokenholders"`
}

func (oth OperatorTokenHoldersStateValue) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(OperatorTokenHoldersStateValueJSONMarshaler{
		BaseHinter:   oth.BaseHinter,
		TokenHolders: oth.TokenHolders,
	})
}

type OperatorTokenHoldersStateValueJSONUnmarshaler struct {
	TokenHolders []string `json:"tokenholders"`
}

func (oth *OperatorTokenHoldersStateValue) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode json of OperatorTokenHoldersStateValue")

	var u OperatorTokenHoldersStateValueJSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return e(err, "")
	}

	holders := make([]base.Address, len(u.TokenHolders))
	for i := range u.TokenHolders {
		a, err := base.DecodeAddress(u.TokenHolders[i], enc)
		if err != nil {
			return e(err, "")
		}
		holders[i] = a
	}
	oth.TokenHolders = holders

	return nil
}
