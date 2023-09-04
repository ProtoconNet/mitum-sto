package sto

import (
	"encoding/json"

	"github.com/ProtoconNet/mitum-currency/v3/common"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	jsonenc "github.com/ProtoconNet/mitum2/util/encoder/json"
)

type CreateSecurityTokensFactJSONMarshaler struct {
	base.BaseFactJSONMarshaler
	Owner base.Address               `json:"sender"`
	Items []CreateSecurityTokensItem `json:"items"`
}

func (fact CreateSecurityTokensFact) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(CreateSecurityTokensFactJSONMarshaler{
		BaseFactJSONMarshaler: fact.BaseFact.JSONMarshaler(),
		Owner:                 fact.sender,
		Items:                 fact.items,
	})
}

type CreateSecurityTokensFactJSONUnMarshaler struct {
	base.BaseFactJSONUnmarshaler
	Owner string          `json:"sender"`
	Items json.RawMessage `json:"items"`
}

func (fact *CreateSecurityTokensFact) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringError("failed to decode json of CreateSecurityTokensFact")

	var uf CreateSecurityTokensFactJSONUnMarshaler
	if err := enc.Unmarshal(b, &uf); err != nil {
		return e.Wrap(err)
	}

	fact.BaseFact.SetJSONUnmarshaler(uf.BaseFactJSONUnmarshaler)

	return fact.unpack(enc, uf.Owner, uf.Items)
}

type CreateSecurityTokensMarshaler struct {
	common.BaseOperationJSONMarshaler
}

func (op CreateSecurityTokens) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(CreateSecurityTokensMarshaler{
		BaseOperationJSONMarshaler: op.BaseOperation.JSONMarshaler(),
	})
}

func (op *CreateSecurityTokens) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringError("failed to decode json of CreateSecurityTokens")

	var ubo common.BaseOperation
	if err := ubo.DecodeJSON(b, enc); err != nil {
		return e.Wrap(err)
	}

	op.BaseOperation = ubo

	return nil
}
