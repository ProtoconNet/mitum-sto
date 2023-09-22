package sto

import (
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
	"github.com/pkg/errors"
)

func (fact *RevokeOperatorFact) unpack(enc encoder.Encoder, sa string, bit []byte) error {
	e := util.StringError("failed to unmarshal RevokeOperatorFact")

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

	items := make([]RevokeOperatorItem, len(hit))
	for i := range hit {
		j, ok := hit[i].(RevokeOperatorItem)
		if !ok {
			return e.Wrap(errors.Errorf("expected RevokeOperatorItem, not %T", hit[i]))
		}

		items[i] = j
	}
	fact.items = items

	return nil
}
