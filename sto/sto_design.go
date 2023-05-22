package sto

import (
	extensioncurrency "github.com/ProtoconNet/mitum-currency-extension/v2/currency"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
)

var (
	DesignHint = hint.MustNewHint("mitum-sto-design-v0.0.1")
)

type Design struct {
	hint.BaseHinter
	stoID       extensioncurrency.ContractID
	granularity uint64
	policy      Policy
}

func NewDesign(stoID extensioncurrency.ContractID, granularity uint64, policy Policy) Design {
	return Design{
		BaseHinter:  hint.NewBaseHinter(DesignHint),
		stoID:       stoID,
		granularity: granularity,
		policy:      policy,
	}
}

func (s Design) IsValid([]byte) error {
	if err := util.CheckIsValiders(nil, false,
		s.BaseHinter,
		s.stoID,
		s.policy,
	); err != nil {
		return util.ErrInvalid.Errorf("invalid Design: %w", err)
	}

	if err := s.stoID.IsValid(nil); err != nil {
		return util.ErrInvalid.Errorf("invalid ContractID: %w", err)
	}

	return s.policy.IsValid(nil)
}

func (s Design) Bytes() []byte {
	return util.ConcatBytesSlice(
		s.stoID.Bytes(),
		util.Uint64ToBigBytes(s.granularity),
		s.policy.Bytes(),
	)
}

func (s Design) STO() extensioncurrency.ContractID {
	return s.stoID
}

func (s Design) Granularity() uint64 {
	return s.granularity
}

func (s Design) Policy() Policy {
	return s.policy
}

func (s Design) SetPolicy(po Policy) Design {
	s.policy = po

	return s
}
