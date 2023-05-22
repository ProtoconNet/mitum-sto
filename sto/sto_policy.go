package sto

import (
	"github.com/ProtoconNet/mitum-currency/v2/currency"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
)

var (
	PolicyHint = hint.MustNewHint("mitum-sto-policy-v0.0.1")
)

type Policy struct {
	hint.BaseHinter
	partitions  []Partition
	aggregate   currency.Big
	controllers []base.Address
	documents   []Document
}

func NewPolicy(partitions []Partition, aggregate currency.Big, controllers []base.Address, documents []Document) Policy {
	return Policy{
		BaseHinter:  hint.NewBaseHinter(PolicyHint),
		partitions:  partitions,
		aggregate:   aggregate,
		controllers: controllers,
		documents:   documents,
	}
}

func (po Policy) Bytes() []byte {
	bs := make([][]byte, len(po.partitions)+len(po.controllers)+len(po.documents))
	for i, p := range po.partitions {
		bs[i] = p.Bytes()
	}
	for i, p := range po.controllers {
		bs[i+len(po.partitions)] = p.Bytes()
	}
	for i, p := range po.documents {
		bs[i+len(po.partitions)+len(po.controllers)] = p.Bytes()
	}

	return util.ConcatBytesSlice(
		util.ConcatBytesSlice(bs...),
		po.aggregate.Bytes(),
	)
}

func (po Policy) IsValid([]byte) error {
	if len(po.partitions) == 0 {
		return util.ErrInvalid.Errorf("empty partitions")
	}

	if err := util.CheckIsValiders(nil, false, po.BaseHinter); err != nil {
		return util.ErrInvalid.Errorf("invalid currency policy: %w", err)
	}

	for _, p := range po.partitions {
		if err := p.IsValid(nil); err != nil {
			return util.ErrInvalid.Errorf("invalid Partition: %w", err)
		}
	}
	for _, p := range po.controllers {
		if err := p.IsValid(nil); err != nil {
			return util.ErrInvalid.Errorf("invalid Controller: %w", err)
		}
	}
	for _, p := range po.documents {
		if err := p.IsValid(nil); err != nil {
			return util.ErrInvalid.Errorf("invalid Document: %w", err)
		}
	}

	return nil
}

func (po Policy) Partitions() []Partition {
	return po.partitions
}

func (po Policy) Aggregate() currency.Big {
	return po.aggregate
}

func (po Policy) Controllers() []base.Address {
	return po.controllers
}

func (po Policy) Documents() []Document {
	return po.documents
}
