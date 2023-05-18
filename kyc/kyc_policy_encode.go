package kyc

import (
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
	"github.com/ProtoconNet/mitum2/util/hint"
)

func (po *KYCPolicy) unpack(enc encoder.Encoder, ht hint.Hint, bcs []string) error {
	e := util.StringErrorFunc("failed to decode bson of KYCPolicy")

	po.BaseHinter = hint.NewBaseHinter(ht)

	controllers := make([]base.Address, len(bcs))
	for i := range bcs {
		ctr, err := base.DecodeAddress(bcs[i], enc)
		if err != nil {
			return e(err, "")
		}
		controllers[i] = ctr
	}
	po.controllers = controllers

	return nil
}
