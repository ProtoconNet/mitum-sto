package kyc

import (
	currencytypes "github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
)

var (
	DesignHint = hint.MustNewHint("mitum-kyc-design-v0.0.1")
)

type Design struct {
	hint.BaseHinter
	kycID  currencytypes.ContractID
	policy Policy
}

func NewDesign(kycID currencytypes.ContractID, policy Policy) Design {
	return Design{
		BaseHinter: hint.NewBaseHinter(DesignHint),
		kycID:      kycID,
		policy:     policy,
	}
}

func (k Design) IsValid([]byte) error {
	if err := util.CheckIsValiders(nil, false,
		k.BaseHinter,
		k.kycID,
		k.policy,
	); err != nil {
		return util.ErrInvalid.Errorf("invalid KYCDesign: %v", err)
	}

	if err := k.kycID.IsValid(nil); err != nil {
		return util.ErrInvalid.Errorf("invalid ContractID: %v", err)
	}

	return k.policy.IsValid(nil)
}

func (k Design) Bytes() []byte {
	return util.ConcatBytesSlice(
		k.kycID.Bytes(),
		k.policy.Bytes(),
	)
}

func (k Design) KYC() currencytypes.ContractID {
	return k.kycID
}

func (k Design) Policy() Policy {
	return k.policy
}

func (k Design) SetPolicy(po Policy) Design {
	k.policy = po

	return k
}
