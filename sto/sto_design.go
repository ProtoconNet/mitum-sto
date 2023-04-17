package sto

import (
	extensioncurrency "github.com/ProtoconNet/mitum-currency-extension/v2/currency"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
)

var (
	STODesignHint = hint.MustNewHint("mitum-sto-design-v0.0.1")
)

type STODesign struct {
	hint.BaseHinter
	stoID  extensioncurrency.ContractID
	policy STOPolicy
}

func NewSTODesign(stoID extensioncurrency.ContractID, policy STOPolicy) STODesign {
	return STODesign{
		BaseHinter: hint.NewBaseHinter(STODesignHint),
		stoID:      stoID,
		policy:     policy,
	}
}

func (s STODesign) IsValid([]byte) error {
	if err := util.CheckIsValiders(nil, false,
		s.BaseHinter,
		s.stoID,
		s.policy,
	); err != nil {
		return util.ErrInvalid.Errorf("invalid STODesign: %w", err)
	}

	if err := s.stoID.IsValid(nil); err != nil {
		return util.ErrInvalid.Errorf("invalid ContractID: %w", err)
	}

	return s.policy.IsValid(nil)
}

func (s STODesign) Bytes() []byte {
	return util.ConcatBytesSlice(
		s.stoID.Bytes(),
		s.policy.Bytes(),
	)
}

func (s STODesign) STO() extensioncurrency.ContractID {
	return s.stoID
}

func (s STODesign) Policy() STOPolicy {
	return s.policy
}

func (s STODesign) SetPolicy(po STOPolicy) STODesign {
	s.policy = po

	return s
}
