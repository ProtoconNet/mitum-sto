package kyc

import (
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
)

var (
	DesignHint = hint.MustNewHint("mitum-kyc-design-v0.0.1")
)

type Design struct {
	hint.BaseHinter
	policy Policy
}

func NewDesign(policy Policy) Design {
	return Design{
		BaseHinter: hint.NewBaseHinter(DesignHint),
		policy:     policy,
	}
}

func (k Design) IsValid([]byte) error {
	if err := util.CheckIsValiders(nil, false,
		k.BaseHinter,
		k.policy,
	); err != nil {
		return util.ErrInvalid.Errorf("invalid KYCDesign: %v", err)
	}

	return k.policy.IsValid(nil)
}

func (k Design) Bytes() []byte {
	return util.ConcatBytesSlice(
		k.policy.Bytes(),
	)
}

func (k Design) Policy() Policy {
	return k.policy
}

func (k Design) SetPolicy(po Policy) Design {
	k.policy = po

	return k
}
