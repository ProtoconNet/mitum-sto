package sto

import (
	currencytypes "github.com/ProtoconNet/mitum-currency/v3/types"
	stotypes "github.com/ProtoconNet/mitum-sto/types/sto"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	jsonenc "github.com/ProtoconNet/mitum2/util/encoder/json"
	"github.com/ProtoconNet/mitum2/util/hint"
)

type RevokeOperatorItemJSONMarshaler struct {
	hint.BaseHinter
	Contract  base.Address             `json:"contract"`
	Operator  base.Address             `json:"operator"`
	Partition stotypes.Partition       `json:"partition"`
	Currency  currencytypes.CurrencyID `json:"currency"`
}

func (it RevokeOperatorItem) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(RevokeOperatorItemJSONMarshaler{
		BaseHinter: it.BaseHinter,
		Contract:   it.contract,
		Operator:   it.operator,
		Partition:  it.partition,
		Currency:   it.currency,
	})
}

type RevokeOperatorItemJSONUnMarshaler struct {
	Hint      hint.Hint `json:"_hint"`
	Contract  string    `json:"contract"`
	Operator  string    `json:"operator"`
	Partition string    `json:"partition"`
	Currency  string    `json:"currency"`
}

func (it *RevokeOperatorItem) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringError("failed to decode json of RevokeOperatorItem")

	var uit RevokeOperatorItemJSONUnMarshaler
	if err := enc.Unmarshal(b, &uit); err != nil {
		return e.Wrap(err)
	}

	return it.unpack(enc, uit.Hint, uit.Contract, uit.Operator, uit.Partition, uit.Currency)
}
