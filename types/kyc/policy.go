package kyc

import (
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
)

var (
	PolicyHint = hint.MustNewHint("mitum-kyc-policy-v0.0.1")
)

type Policy struct {
	hint.BaseHinter
	controllers []base.Address
}

func NewPolicy(controllers []base.Address) Policy {
	return Policy{
		BaseHinter:  hint.NewBaseHinter(PolicyHint),
		controllers: controllers,
	}
}

func (po Policy) Bytes() []byte {
	bs := make([][]byte, len(po.controllers))
	for i, p := range po.controllers {
		bs[i] = p.Bytes()
	}

	return util.ConcatBytesSlice(
		bs...,
	)
}

func (po Policy) IsValid([]byte) error {
	if err := util.CheckIsValiders(nil, false, po.BaseHinter); err != nil {
		return util.ErrInvalid.Errorf("invalid kyc policy: %v", err)
	}

	for _, p := range po.controllers {
		if err := p.IsValid(nil); err != nil {
			return util.ErrInvalid.Errorf("invalid Controller: %v", err)
		}
	}

	return nil
}

func (po Policy) Controllers() []base.Address {
	return po.controllers
}
