package sto

import (
	"github.com/ProtoconNet/mitum-currency/v3/common"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
	"github.com/ProtoconNet/mitum2/util/hint"
	"github.com/pkg/errors"
)

func (po *Policy) unpack(enc encoder.Encoder, ht hint.Hint, ps []string, big string, bcs []string, bds []byte) error {
	e := util.StringError("failed to decode bson of Policy")

	po.BaseHinter = hint.NewBaseHinter(ht)

	partitions := make([]Partition, len(ps))
	for i, p := range ps {
		partitions[i] = Partition(p)
	}
	po.partitions = partitions

	if ag, err := common.NewBigFromString(big); err != nil {
		return e.Wrap(err)
	} else {
		po.aggregate = ag
	}

	controllers := make([]base.Address, len(bcs))
	for i := range bcs {
		ctr, err := base.DecodeAddress(bcs[i], enc)
		if err != nil {
			return e.Wrap(err)
		}
		controllers[i] = ctr
	}
	po.controllers = controllers

	hds, err := enc.DecodeSlice(bds)
	if err != nil {
		return e.Wrap(err)
	}

	documents := make([]Document, len(hds))
	for i := range hds {
		doc, ok := hds[i].(Document)
		if !ok {
			return e.Wrap(errors.Errorf("expected Document, not %T", hds[i]))
		}

		documents[i] = doc
	}
	po.documents = documents

	return nil
}
