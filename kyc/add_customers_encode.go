package kyc

import (
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
)

func (fact *AddCustomersFact) unpack(enc encoder.Encoder, sa string, bit []byte) error {
	e := util.StringErrorFunc("failed to unmarshal AddCustomersFact")

	switch a, err := base.DecodeAddress(sa, enc); {
	case err != nil:
		return e(err, "")
	default:
		fact.sender = a
	}

	hit, err := enc.DecodeSlice(bit)
	if err != nil {
		return e(err, "")
	}

	items := make([]AddCustomersItem, len(hit))
	for i := range hit {
		j, ok := hit[i].(AddCustomersItem)
		if !ok {
			return e(util.ErrWrongType.Errorf("expected AddCustomersItem, not %T", hit[i]), "")
		}

		items[i] = j
	}
	fact.items = items

	return nil
}
