package sto

import (
	"github.com/ProtoconNet/mitum-currency/v2/currency"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
	"github.com/ProtoconNet/mitum2/util/hint"
)

func (po *STOPolicy) unpack(enc encoder.Encoder, ht hint.Hint, ps []string, big string, bcs []string, bds []byte) error {
	e := util.StringErrorFunc("failed to decode bson of STOPolicy")

	po.BaseHinter = hint.NewBaseHinter(ht)

	partitions := make([]Partition, len(ps))
	for i, p := range ps {
		partitions[i] = Partition(p)
	}
	po.partitions = partitions

	if ag, err := currency.NewBigFromString(big); err != nil {
		return e(err, "")
	} else {
		po.aggregate = ag
	}

	controllers := make([]base.Address, len(bcs))
	for i := range bcs {
		ctr, err := base.DecodeAddress(bcs[i], enc)
		if err != nil {
			return e(err, "")
		}
		controllers[i] = ctr
	}
	po.controllers = controllers

	hds, err := enc.DecodeSlice(bds)
	if err != nil {
		return e(err, "")
	}

	documents := make([]Document, len(hds))
	for i := range hds {
		doc, ok := hds[i].(Document)
		if !ok {
			return e(util.ErrWrongType.Errorf("expected Document, not %T", hds[i]), "")
		}

		documents[i] = doc
	}
	po.documents = documents

	return nil
}
