package kyc

import (
	extensioncurrency "github.com/ProtoconNet/mitum-currency-extension/v2/currency"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
)

var (
	KYCDesignHint = hint.MustNewHint("mitum-kyc-design-v0.0.1")
)

type KYCDesign struct {
	hint.BaseHinter
	kycID  extensioncurrency.ContractID
	policy KYCPolicy
}

func NewKYCDesign(kycID extensioncurrency.ContractID, policy KYCPolicy) KYCDesign {
	return KYCDesign{
		BaseHinter: hint.NewBaseHinter(KYCDesignHint),
		kycID:      kycID,
		policy:     policy,
	}
}

func (k KYCDesign) IsValid([]byte) error {
	if err := util.CheckIsValiders(nil, false,
		k.BaseHinter,
		k.kycID,
		k.policy,
	); err != nil {
		return util.ErrInvalid.Errorf("invalid KYCKYCDesign: %w", err)
	}

	if err := k.kycID.IsValid(nil); err != nil {
		return util.ErrInvalid.Errorf("invalid ContractID: %w", err)
	}

	return k.policy.IsValid(nil)
}

func (k KYCDesign) Bytes() []byte {
	return util.ConcatBytesSlice(
		k.kycID.Bytes(),
		k.policy.Bytes(),
	)
}

func (k KYCDesign) KYC() extensioncurrency.ContractID {
	return k.kycID
}

func (k KYCDesign) Policy() KYCPolicy {
	return k.policy
}

func (k KYCDesign) SetPolicy(po KYCPolicy) KYCDesign {
	k.policy = po

	return k
}
