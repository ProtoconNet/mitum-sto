package sto

import (
	"github.com/ProtoconNet/mitum-currency/v2/currency"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
	"github.com/ProtoconNet/mitum2/util/hint"
)

func (po *STOPolicy) unpack(enc encoder.Encoder, ht hint.Hint, bps []byte, big string, bcs []string, bds []byte) error {
	e := util.StringErrorFunc("failed to decode bson of STOPolicy")

	po.BaseHinter = hint.NewBaseHinter(ht)

	hps, err := enc.DecodeSlice(bps)
	if err != nil {
		return e(err, "")
	}

	partitions := make([]Partition, len(hps))
	for i := range hps {
		p, ok := hps[i].(Partition)
		if !ok {
			return util.ErrWrongType.Errorf("expected Partition, not %T", hps[i])
		}

		partitions[i] = p
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
	for i := range hps {
		doc, ok := hds[i].(Document)
		if !ok {
			return util.ErrWrongType.Errorf("expected Document, not %T", hds[i])
		}

		documents[i] = doc
	}
	po.documents = documents

	return nil
}
