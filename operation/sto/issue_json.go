package sto

import (
	"encoding/json"

	"github.com/ProtoconNet/mitum-currency/v3/common"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
)

type IssueFactJSONMarshaler struct {
	base.BaseFactJSONMarshaler
	Owner base.Address `json:"sender"`
	Items []IssueItem  `json:"items"`
}

func (fact IssueFact) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(IssueFactJSONMarshaler{
		BaseFactJSONMarshaler: fact.BaseFact.JSONMarshaler(),
		Owner:                 fact.sender,
		Items:                 fact.items,
	})
}

type IssueFactJSONUnMarshaler struct {
	base.BaseFactJSONUnmarshaler
	Owner string          `json:"sender"`
	Items json.RawMessage `json:"items"`
}

func (fact *IssueFact) DecodeJSON(b []byte, enc encoder.Encoder) error {
	e := util.StringError("failed to decode json of IssueFact")

	var uf IssueFactJSONUnMarshaler
	if err := enc.Unmarshal(b, &uf); err != nil {
		return e.Wrap(err)
	}

	fact.BaseFact.SetJSONUnmarshaler(uf.BaseFactJSONUnmarshaler)

	return fact.unpack(enc, uf.Owner, uf.Items)
}

type IssueMarshaler struct {
	common.BaseOperationJSONMarshaler
}

func (op Issue) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(IssueMarshaler{
		BaseOperationJSONMarshaler: op.BaseOperation.JSONMarshaler(),
	})
}

func (op *Issue) DecodeJSON(b []byte, enc encoder.Encoder) error {
	e := util.StringError("failed to decode json of MintSecurityTokens")

	var ubo common.BaseOperation
	if err := ubo.DecodeJSON(b, enc); err != nil {
		return e.Wrap(err)
	}

	op.BaseOperation = ubo

	return nil
}
