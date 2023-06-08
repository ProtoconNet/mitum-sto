package sto

import (
	"encoding/json"

	currencybase "github.com/ProtoconNet/mitum-currency/v3/base"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	jsonenc "github.com/ProtoconNet/mitum2/util/encoder/json"
)

type RevokeOperatorsFactJSONMarshaler struct {
	base.BaseFactJSONMarshaler
	Owner base.Address          `json:"sender"`
	Items []RevokeOperatorsItem `json:"items"`
}

func (fact RevokeOperatorsFact) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(RevokeOperatorsFactJSONMarshaler{
		BaseFactJSONMarshaler: fact.BaseFact.JSONMarshaler(),
		Owner:                 fact.sender,
		Items:                 fact.items,
	})
}

type RevokeOperatorsFactJSONUnMarshaler struct {
	base.BaseFactJSONUnmarshaler
	Owner string          `json:"sender"`
	Items json.RawMessage `json:"items"`
}

func (fact *RevokeOperatorsFact) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode json of RevokeOperatorsFact")

	var uf RevokeOperatorsFactJSONUnMarshaler
	if err := enc.Unmarshal(b, &uf); err != nil {
		return e(err, "")
	}

	fact.BaseFact.SetJSONUnmarshaler(uf.BaseFactJSONUnmarshaler)

	return fact.unpack(enc, uf.Owner, uf.Items)
}

type RevokeOperatorsMarshaler struct {
	currencybase.BaseOperationJSONMarshaler
}

func (op RevokeOperators) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(RevokeOperatorsMarshaler{
		BaseOperationJSONMarshaler: op.BaseOperation.JSONMarshaler(),
	})
}

func (op *RevokeOperators) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode json of RevokeOperators")

	var ubo currencybase.BaseOperation
	if err := ubo.DecodeJSON(b, enc); err != nil {
		return e(err, "")
	}

	op.BaseOperation = ubo

	return nil
}
