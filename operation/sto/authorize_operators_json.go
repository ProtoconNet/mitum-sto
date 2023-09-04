package sto

import (
	"encoding/json"

	"github.com/ProtoconNet/mitum-currency/v3/common"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	jsonenc "github.com/ProtoconNet/mitum2/util/encoder/json"
)

type AuthorizeOperatorsFactJSONMarshaler struct {
	base.BaseFactJSONMarshaler
	Owner base.Address             `json:"sender"`
	Items []AuthorizeOperatorsItem `json:"items"`
}

func (fact AuthorizeOperatorsFact) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(AuthorizeOperatorsFactJSONMarshaler{
		BaseFactJSONMarshaler: fact.BaseFact.JSONMarshaler(),
		Owner:                 fact.sender,
		Items:                 fact.items,
	})
}

type AuthorizeOperatorsFactJSONUnMarshaler struct {
	base.BaseFactJSONUnmarshaler
	Owner string          `json:"sender"`
	Items json.RawMessage `json:"items"`
}

func (fact *AuthorizeOperatorsFact) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringError("failed to decode json of AuthorizeOperatorsFact")

	var uf AuthorizeOperatorsFactJSONUnMarshaler
	if err := enc.Unmarshal(b, &uf); err != nil {
		return e.Wrap(err)
	}

	fact.BaseFact.SetJSONUnmarshaler(uf.BaseFactJSONUnmarshaler)

	return fact.unpack(enc, uf.Owner, uf.Items)
}

type AuthorizeOperatorsMarshaler struct {
	common.BaseOperationJSONMarshaler
}

func (op AuthorizeOperators) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(AuthorizeOperatorsMarshaler{
		BaseOperationJSONMarshaler: op.BaseOperation.JSONMarshaler(),
	})
}

func (op *AuthorizeOperators) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringError("failed to decode json of AuthorizeOperators")

	var ubo common.BaseOperation
	if err := ubo.DecodeJSON(b, enc); err != nil {
		return e.Wrap(err)
	}

	op.BaseOperation = ubo

	return nil
}
