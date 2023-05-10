package kyc

import (
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
)

var (
	KYCPolicyHint = hint.MustNewHint("mitum-kyc-policy-v0.0.1")
)

type KYCPolicy struct {
	hint.BaseHinter
	controllers []base.Address
}

func NewKYCPolicy(controllers []base.Address) KYCPolicy {
	return KYCPolicy{
		BaseHinter:  hint.NewBaseHinter(KYCPolicyHint),
		controllers: controllers,
	}
}

func (po KYCPolicy) Bytes() []byte {
	bs := make([][]byte, len(po.controllers))
	for i, p := range po.controllers {
		bs[i] = p.Bytes()
	}

	return util.ConcatBytesSlice(
		bs...,
	)
}

func (po KYCPolicy) IsValid([]byte) error {
	if err := util.CheckIsValiders(nil, false, po.BaseHinter); err != nil {
		return util.ErrInvalid.Errorf("invalid kyc policy: %w", err)
	}

	for _, p := range po.controllers {
		if err := p.IsValid(nil); err != nil {
			return util.ErrInvalid.Errorf("invalid Controller: %w", err)
		}
	}

	return nil
}

func (po KYCPolicy) Controllers() []base.Address {
	return po.controllers
}
