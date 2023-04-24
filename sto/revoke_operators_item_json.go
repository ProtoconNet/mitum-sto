package sto

import (
	extensioncurrency "github.com/ProtoconNet/mitum-currency-extension/v2/currency"
	"github.com/ProtoconNet/mitum-currency/v2/currency"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	jsonenc "github.com/ProtoconNet/mitum2/util/encoder/json"
	"github.com/ProtoconNet/mitum2/util/hint"
)

type RevokeOperatorsItemJSONMarshaler struct {
	hint.BaseHinter
	Contract    base.Address                 `json:"contract"`
	STO         extensioncurrency.ContractID `json:"stoid"`
	TokenHolder base.Address                 `json:"token_holder"`
	Operator    base.Address                 `json:"operator"`
	Partition   Partition                    `json:"partition"`
	Currency    currency.CurrencyID          `json:"currency"`
}

func (it RevokeOperatorsItem) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(RevokeOperatorsItemJSONMarshaler{
		BaseHinter:  it.BaseHinter,
		Contract:    it.contract,
		STO:         it.stoID,
		TokenHolder: it.tokenHolder,
		Operator:    it.operator,
		Partition:   it.partition,
		Currency:    it.currency,
	})
}

type RevokeOperatorsItemJSONUnMarshaler struct {
	Hint        hint.Hint `json:"_hint"`
	Contract    string    `json:"contract"`
	STO         string    `json:"stoid"`
	TokenHolder string    `json:"token_holder"`
	Operator    string    `json:"operator"`
	Partition   string    `json:"partition"`
	Currency    string    `json:"currency"`
}

func (it *RevokeOperatorsItem) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode json of RevokeOperatorsItem")

	var uit RevokeOperatorsItemJSONUnMarshaler
	if err := enc.Unmarshal(b, &uit); err != nil {
		return e(err, "")
	}

	return it.unpack(enc, uit.Hint, uit.Contract, uit.STO, uit.TokenHolder, uit.Operator, uit.Partition, uit.Currency)
}
