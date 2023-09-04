package sto

import (
	currencytypes "github.com/ProtoconNet/mitum-currency/v3/types"
	stotypes "github.com/ProtoconNet/mitum-sto/types/sto"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
	"github.com/ProtoconNet/mitum2/util/hint"
)

func (it *RevokeOperatorsItem) unpack(enc encoder.Encoder, ht hint.Hint, ca, sto, oper, p, cid string) error {
	e := util.StringError("failed to unmarshal RevokeOperatorsItem")

	it.BaseHinter = hint.NewBaseHinter(ht)
	it.stoID = currencytypes.ContractID(sto)
	it.partition = stotypes.Partition(p)
	it.currency = currencytypes.CurrencyID(cid)

	switch a, err := base.DecodeAddress(ca, enc); {
	case err != nil:
		return e.Wrap(err)
	default:
		it.contract = a
	}

	switch a, err := base.DecodeAddress(oper, enc); {
	case err != nil:
		return e.Wrap(err)
	default:
		it.operator = a
	}

	return nil
}
