package sto

import (
	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/util"
	"github.com/spikeekips/mitum/util/hint"
)

var (
	STOPolicyHint = hint.MustNewHint("mitum-sto-policy-v0.0.1")
)

type STOPolicy struct {
	hint.BaseHinter
	partitions  []Partition
	aggregate   currency.Amount
	controllers []currency.Account
}

func NewSTOPolicy(partitions []Partition, aggregate currency.Amount, controllers []currency.Account) STOPolicy {
	return STOPolicy{
		BaseHinter:  hint.NewBaseHinter(STOPolicyHint),
		partitions:  partitions,
		aggregate:   aggregate,
		controllers: controllers,
	}
}

func (po STOPolicy) Bytes() []byte {
	bs := make([][]byte, len(po.partitions))
	for i, p := range po.partitions {
		bs[i] = p.Bytes()
	}
	return util.ConcatBytesSlice(
		util.ConcatBytesSlice(bs...),
		po.aggregate.Bytes(),
	)
}

func (po STOPolicy) IsValid([]byte) error {
	if len(po.partitions) == 0 {
		return util.ErrInvalid.Errorf("empty partitions")
	}

	if !po.aggregate.Big().OverZero() {
		return util.ErrInvalid.Errorf("aggregate not over zero")
	}

	if err := util.CheckIsValiders(nil, false, po.BaseHinter, po.controllers); err != nil {
		return util.ErrInvalid.Errorf("invalid currency policy: %w", err)
	}

	for _, p := range po.partitions {
		if err := p.IsValid(nil); err != nil {
			return util.ErrInvalid.Errorf("invalid Partition: %w", err)
		}
	}

	return nil
}

func (po STOPolicy) Partitions() []Partition {
	return po.partitions
}

func (po STOPolicy) Aggregate() currency.Amount {
	return po.aggregate
}

func (po STOPolicy) Controllers() []currency.Account {
	return po.controllers
}
