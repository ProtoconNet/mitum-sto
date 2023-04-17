package sto

import (
	"encoding/json"

	"github.com/ProtoconNet/mitum-currency/v2/currency"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	jsonenc "github.com/ProtoconNet/mitum2/util/encoder/json"
)

type MintSecurityTokensFactJSONMarshaler struct {
	base.BaseFactJSONMarshaler
	Owner base.Address             `json:"sender"`
	Items []MintSecurityTokensItem `json:"items"`
}

func (fact MintSecurityTokensFact) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(MintSecurityTokensFactJSONMarshaler{
		BaseFactJSONMarshaler: fact.BaseFact.JSONMarshaler(),
		Owner:                 fact.sender,
		Items:                 fact.items,
	})
}

type MintSecurityTokensFactJSONUnMarshaler struct {
	base.BaseFactJSONUnmarshaler
	Owner string          `json:"sender"`
	Items json.RawMessage `json:"items"`
}

func (fact *MintSecurityTokensFact) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode json of MintSecurityTokensFact")

	var uf MintSecurityTokensFactJSONUnMarshaler
	if err := enc.Unmarshal(b, &uf); err != nil {
		return e(err, "")
	}

	fact.BaseFact.SetJSONUnmarshaler(uf.BaseFactJSONUnmarshaler)

	return fact.unpack(enc, uf.Owner, uf.Items)
}

type MintSecurityTokensMarshaler struct {
	currency.BaseOperationJSONMarshaler
}

func (op MintSecurityTokens) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(MintSecurityTokensMarshaler{
		BaseOperationJSONMarshaler: op.BaseOperation.JSONMarshaler(),
	})
}

func (op *MintSecurityTokens) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode json of MintSecurityTokens")

	var ubo currency.BaseOperation
	if err := ubo.DecodeJSON(b, enc); err != nil {
		return e(err, "")
	}

	op.BaseOperation = ubo

	return nil
}
