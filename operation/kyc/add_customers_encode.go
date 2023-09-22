package kyc

import (
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
	"github.com/pkg/errors"
)

func (fact *AddCustomerFact) unpack(enc encoder.Encoder, sa string, bit []byte) error {
	e := util.StringError("failed to unmarshal AddCustomerFact")

	switch a, err := base.DecodeAddress(sa, enc); {
	case err != nil:
		return e.Wrap(err)
	default:
		fact.sender = a
	}

	hit, err := enc.DecodeSlice(bit)
	if err != nil {
		return e.Wrap(err)
	}

	items := make([]AddCustomerItem, len(hit))
	for i := range hit {
		j, ok := hit[i].(AddCustomerItem)
		if !ok {
			return e.Wrap(errors.Errorf("expected AddCustomerItem, not %T", hit[i]))
		}

		items[i] = j
	}
	fact.items = items

	return nil
}
