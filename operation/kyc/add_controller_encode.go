package kyc

import (
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
	"github.com/pkg/errors"
)

func (fact *AddControllerFact) unpack(enc encoder.Encoder, sa string, bit []byte) error {
	e := util.StringError("failed to unmarshal AddControllerFact")

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

	items := make([]AddControllerItem, len(hit))
	for i := range hit {
		j, ok := hit[i].(AddControllerItem)
		if !ok {
			return e.Wrap(errors.Errorf("expected AddControllerItem, not %T", hit[i]))
		}

		items[i] = j
	}
	fact.items = items

	return nil
}
