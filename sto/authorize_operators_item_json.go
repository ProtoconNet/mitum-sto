package sto

import (
	currencybase "github.com/ProtoconNet/mitum-currency/v3/base"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	jsonenc "github.com/ProtoconNet/mitum2/util/encoder/json"
	"github.com/ProtoconNet/mitum2/util/hint"
)

type AuthorizeOperatorsItemJSONMarshaler struct {
	hint.BaseHinter
	Contract  base.Address            `json:"contract"`
	STO       currencybase.ContractID `json:"stoid"`
	Operator  base.Address            `json:"operator"`
	Partition Partition               `json:"partition"`
	Currency  currencybase.CurrencyID `json:"currency"`
}

func (it AuthorizeOperatorsItem) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(AuthorizeOperatorsItemJSONMarshaler{
		BaseHinter: it.BaseHinter,
		Contract:   it.contract,
		STO:        it.stoID,
		Operator:   it.operator,
		Partition:  it.partition,
		Currency:   it.currency,
	})
}

type AuthorizeOperatorsItemJSONUnMarshaler struct {
	Hint      hint.Hint `json:"_hint"`
	Contract  string    `json:"contract"`
	STO       string    `json:"stoid"`
	Operator  string    `json:"operator"`
	Partition string    `json:"partition"`
	Currency  string    `json:"currency"`
}

func (it *AuthorizeOperatorsItem) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode json of AuthorizeOperatorsItem")

	var uit AuthorizeOperatorsItemJSONUnMarshaler
	if err := enc.Unmarshal(b, &uit); err != nil {
		return e(err, "")
	}

	return it.unpack(enc, uit.Hint, uit.Contract, uit.STO, uit.Operator, uit.Partition, uit.Currency)
}
