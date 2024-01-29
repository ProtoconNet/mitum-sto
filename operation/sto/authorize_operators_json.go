package sto

import (
	"encoding/json"

	"github.com/ProtoconNet/mitum-currency/v3/common"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
)

type AuthorizeOperatorFactJSONMarshaler struct {
	base.BaseFactJSONMarshaler
	Owner base.Address            `json:"sender"`
	Items []AuthorizeOperatorItem `json:"items"`
}

func (fact AuthorizeOperatorFact) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(AuthorizeOperatorFactJSONMarshaler{
		BaseFactJSONMarshaler: fact.BaseFact.JSONMarshaler(),
		Owner:                 fact.sender,
		Items:                 fact.items,
	})
}

type AuthorizeOperatorFactJSONUnMarshaler struct {
	base.BaseFactJSONUnmarshaler
	Owner string          `json:"sender"`
	Items json.RawMessage `json:"items"`
}

func (fact *AuthorizeOperatorFact) DecodeJSON(b []byte, enc encoder.Encoder) error {
	e := util.StringError("failed to decode json of AuthorizeOperatorFact")

	var uf AuthorizeOperatorFactJSONUnMarshaler
	if err := enc.Unmarshal(b, &uf); err != nil {
		return e.Wrap(err)
	}

	fact.BaseFact.SetJSONUnmarshaler(uf.BaseFactJSONUnmarshaler)

	return fact.unpack(enc, uf.Owner, uf.Items)
}

type AuthorizeOperatorMarshaler struct {
	common.BaseOperationJSONMarshaler
}

func (op AuthorizeOperator) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(AuthorizeOperatorMarshaler{
		BaseOperationJSONMarshaler: op.BaseOperation.JSONMarshaler(),
	})
}

func (op *AuthorizeOperator) DecodeJSON(b []byte, enc encoder.Encoder) error {
	e := util.StringError("failed to decode json of AuthorizeOperator")

	var ubo common.BaseOperation
	if err := ubo.DecodeJSON(b, enc); err != nil {
		return e.Wrap(err)
	}

	op.BaseOperation = ubo

	return nil
}
