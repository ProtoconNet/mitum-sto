package sto

import (
	"encoding/json"

	"github.com/ProtoconNet/mitum-currency/v3/common"
	stotypes "github.com/ProtoconNet/mitum-sto/types/sto"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
	"github.com/ProtoconNet/mitum2/util/hint"
)

type DesignStateValueJSONMarshaler struct {
	hint.BaseHinter
	STO stotypes.Design `json:"sto"`
}

func (de DesignStateValue) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(DesignStateValueJSONMarshaler{
		BaseHinter: de.BaseHinter,
		STO:        de.Design,
	})
}

type DesignStateValueJSONUnmarshaler struct {
	STO json.RawMessage `json:"sto"`
}

func (de *DesignStateValue) DecodeJSON(b []byte, enc encoder.Encoder) error {
	e := util.StringError("failed to decode json of DesignStateValue")

	var u DesignStateValueJSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return e.Wrap(err)
	}

	var design stotypes.Design

	if err := design.DecodeJSON(u.STO, enc); err != nil {
		return e.Wrap(err)
	} else if err := design.IsValid(nil); err != nil {
		return e.Wrap(err)
	} else {
		de.Design = design
	}

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

func (p *PartitionBalanceStateValue) DecodeJSON(b []byte, enc encoder.Encoder) error {
	e := util.StringError("failed to decode json of PartitionBalanceStateValue")

	var u PartitionBalanceStateValueJSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return e.Wrap(err)
	}

	big, err := common.NewBigFromString(u.Amount)
	if err != nil {
		return e.Wrap(err)
	}

	p.Amount = big

	return nil
}

type TokenHolderPartitionsStateValueJSONMarshaler struct {
	hint.BaseHinter
	Partitions []stotypes.Partition `json:"partitions"`
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

func (p *TokenHolderPartitionsStateValue) DecodeJSON(b []byte, enc encoder.Encoder) error {
	e := util.StringError("failed to decode json of TokenHolderPartitionsStateValue")

	var u TokenHolderPartitionsStateValueJSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return e.Wrap(err)
	}

	partitions := make([]stotypes.Partition, len(u.Partitions))
	for i, s := range u.Partitions {
		partitions[i] = stotypes.Partition(s)
	}

	p.Partitions = partitions

	return nil
}

type TokenHolderPartitionBalanceStateValueJSONMarshaler struct {
	hint.BaseHinter
	Amount    string             `json:"amount"`
	Partition stotypes.Partition `json:"partition"`
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

func (p *TokenHolderPartitionBalanceStateValue) DecodeJSON(b []byte, enc encoder.Encoder) error {
	e := util.StringError("failed to decode json of TokenHolderPartitionBalanceStateValue")

	var u TokenHolderPartitionBalanceStateValueJSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return e.Wrap(err)
	}

	big, err := common.NewBigFromString(u.Amount)
	if err != nil {
		return e.Wrap(err)
	}
	p.Amount = big

	p.Partition = stotypes.Partition(u.Partition)

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

func (ops *TokenHolderPartitionOperatorsStateValue) DecodeJSON(b []byte, enc encoder.Encoder) error {
	e := util.StringError("failed to decode json of TokenHolderPartitionOperatorsStateValue")

	var u TokenHolderPartitionOperatorsStateValueJSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return e.Wrap(err)
	}

	operators := make([]base.Address, len(u.Operators))
	for i := range u.Operators {
		a, err := base.DecodeAddress(u.Operators[i], enc)
		if err != nil {
			return e.Wrap(err)
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

func (oth *OperatorTokenHoldersStateValue) DecodeJSON(b []byte, enc encoder.Encoder) error {
	e := util.StringError("failed to decode json of OperatorTokenHoldersStateValue")

	var u OperatorTokenHoldersStateValueJSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return e.Wrap(err)
	}

	holders := make([]base.Address, len(u.TokenHolders))
	for i := range u.TokenHolders {
		a, err := base.DecodeAddress(u.TokenHolders[i], enc)
		if err != nil {
			return e.Wrap(err)
		}
		holders[i] = a
	}
	oth.TokenHolders = holders

	return nil
}
