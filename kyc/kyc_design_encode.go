package kyc

import (
	extensioncurrency "github.com/ProtoconNet/mitum-currency-extension/v2/currency"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
	"github.com/ProtoconNet/mitum2/util/hint"
)

func (de *KYCDesign) unpack(enc encoder.Encoder, ht hint.Hint, kyc string, bpo []byte) error {
	e := util.StringErrorFunc("failed to decode bson of KYCDesign")

	de.BaseHinter = hint.NewBaseHinter(ht)
	de.kycID = extensioncurrency.ContractID(kyc)

	if hinter, err := enc.Decode(bpo); err != nil {
		return e(err, "")
	} else if po, ok := hinter.(KYCPolicy); !ok {
		return e(util.ErrWrongType.Errorf("expected KYCPolicy, not %T", hinter), "")
	} else {
		de.policy = po
	}

	return nil
}
