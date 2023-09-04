package kyc

import (
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
	"github.com/pkg/errors"
)

func (fact *AddCustomersFact) unpack(enc encoder.Encoder, sa string, bit []byte) error {
	e := util.StringError("failed to unmarshal AddCustomersFact")

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

	items := make([]AddCustomersItem, len(hit))
	for i := range hit {
		j, ok := hit[i].(AddCustomersItem)
		if !ok {
			return e.Wrap(errors.Errorf("expected AddCustomersItem, not %T", hit[i]))
		}

		items[i] = j
	}
	fact.items = items

	return nil
}
