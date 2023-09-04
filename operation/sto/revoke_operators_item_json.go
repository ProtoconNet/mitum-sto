package sto

import (
	currencytypes "github.com/ProtoconNet/mitum-currency/v3/types"
	stotypes "github.com/ProtoconNet/mitum-sto/types/sto"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	jsonenc "github.com/ProtoconNet/mitum2/util/encoder/json"
	"github.com/ProtoconNet/mitum2/util/hint"
)

type RevokeOperatorsItemJSONMarshaler struct {
	hint.BaseHinter
	Contract  base.Address             `json:"contract"`
	STO       currencytypes.ContractID `json:"stoid"`
	Operator  base.Address             `json:"operator"`
	Partition stotypes.Partition       `json:"partition"`
	Currency  currencytypes.CurrencyID `json:"currency"`
}

func (it RevokeOperatorsItem) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(RevokeOperatorsItemJSONMarshaler{
		BaseHinter: it.BaseHinter,
		Contract:   it.contract,
		STO:        it.stoID,
		Operator:   it.operator,
		Partition:  it.partition,
		Currency:   it.currency,
	})
}

type RevokeOperatorsItemJSONUnMarshaler struct {
	Hint      hint.Hint `json:"_hint"`
	Contract  string    `json:"contract"`
	STO       string    `json:"stoid"`
	Operator  string    `json:"operator"`
	Partition string    `json:"partition"`
	Currency  string    `json:"currency"`
}

func (it *RevokeOperatorsItem) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringError("failed to decode json of RevokeOperatorsItem")

	var uit RevokeOperatorsItemJSONUnMarshaler
	if err := enc.Unmarshal(b, &uit); err != nil {
		return e.Wrap(err)
	}

	return it.unpack(enc, uit.Hint, uit.Contract, uit.STO, uit.Operator, uit.Partition, uit.Currency)
}
