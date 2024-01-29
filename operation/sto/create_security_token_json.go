package sto

import (
	"encoding/json"

	"github.com/ProtoconNet/mitum-currency/v3/common"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
)

type CreateSecurityTokenFactJSONMarshaler struct {
	base.BaseFactJSONMarshaler
	Owner base.Address              `json:"sender"`
	Items []CreateSecurityTokenItem `json:"items"`
}

func (fact CreateSecurityTokenFact) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(CreateSecurityTokenFactJSONMarshaler{
		BaseFactJSONMarshaler: fact.BaseFact.JSONMarshaler(),
		Owner:                 fact.sender,
		Items:                 fact.items,
	})
}

type CreateSecurityTokenFactJSONUnMarshaler struct {
	base.BaseFactJSONUnmarshaler
	Owner string          `json:"sender"`
	Items json.RawMessage `json:"items"`
}

func (fact *CreateSecurityTokenFact) DecodeJSON(b []byte, enc encoder.Encoder) error {
	e := util.StringError("failed to decode json of CreateSecurityTokenFact")

	var uf CreateSecurityTokenFactJSONUnMarshaler
	if err := enc.Unmarshal(b, &uf); err != nil {
		return e.Wrap(err)
	}

	fact.BaseFact.SetJSONUnmarshaler(uf.BaseFactJSONUnmarshaler)

	return fact.unpack(enc, uf.Owner, uf.Items)
}

type CreateSecurityTokenMarshaler struct {
	common.BaseOperationJSONMarshaler
}

func (op CreateSecurityToken) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(CreateSecurityTokenMarshaler{
		BaseOperationJSONMarshaler: op.BaseOperation.JSONMarshaler(),
	})
}

func (op *CreateSecurityToken) DecodeJSON(b []byte, enc encoder.Encoder) error {
	e := util.StringError("failed to decode json of CreateSecurityToken")

	var ubo common.BaseOperation
	if err := ubo.DecodeJSON(b, enc); err != nil {
		return e.Wrap(err)
	}

	op.BaseOperation = ubo

	return nil
}
