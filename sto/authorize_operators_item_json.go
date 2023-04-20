package sto

import (
	extensioncurrency "github.com/ProtoconNet/mitum-currency-extension/v2/currency"
	"github.com/ProtoconNet/mitum-currency/v2/currency"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	jsonenc "github.com/ProtoconNet/mitum2/util/encoder/json"
	"github.com/ProtoconNet/mitum2/util/hint"
)

type AuthorizeOperatorsItemJSONMarshaler struct {
	hint.BaseHinter
	STO      extensioncurrency.ContractID `json:"stoid"`
	Contract base.Address                 `json:"contract"`
	Operator base.Address                 `json:"operator"`
	Currency currency.CurrencyID          `json:"currency"`
}

func (it AuthorizeOperatorsItem) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(AuthorizeOperatorsItemJSONMarshaler{
		BaseHinter: it.BaseHinter,
		STO:        it.stoID,
		Contract:   it.contract,
		Operator:   it.operator,
		Currency:   it.currency,
	})
}

type AuthorizeOperatorsItemJSONUnMarshaler struct {
	Hint     hint.Hint `json:"_hint"`
	STO      string    `json:"stoid"`
	Contract string    `json:"contract"`
	Operator string    `json:"operator"`
	Currency string    `json:"currency"`
}

func (it *AuthorizeOperatorsItem) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode json of AuthorizeOperatorsItem")

	var uit AuthorizeOperatorsItemJSONUnMarshaler
	if err := enc.Unmarshal(b, &uit); err != nil {
		return e(err, "")
	}

	return it.unpack(enc, uit.Hint, uit.STO, uit.Contract, uit.Operator, uit.Currency)
}
