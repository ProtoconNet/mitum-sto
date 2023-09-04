package sto

import (
	"encoding/json"

	"github.com/ProtoconNet/mitum-currency/v3/common"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	jsonenc "github.com/ProtoconNet/mitum2/util/encoder/json"
)

type IssueSecurityTokensFactJSONMarshaler struct {
	base.BaseFactJSONMarshaler
	Owner base.Address              `json:"sender"`
	Items []IssueSecurityTokensItem `json:"items"`
}

func (fact IssueSecurityTokensFact) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(IssueSecurityTokensFactJSONMarshaler{
		BaseFactJSONMarshaler: fact.BaseFact.JSONMarshaler(),
		Owner:                 fact.sender,
		Items:                 fact.items,
	})
}

type IssueSecurityTokensFactJSONUnMarshaler struct {
	base.BaseFactJSONUnmarshaler
	Owner string          `json:"sender"`
	Items json.RawMessage `json:"items"`
}

func (fact *IssueSecurityTokensFact) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringError("failed to decode json of IssueSecurityTokensFact")

	var uf IssueSecurityTokensFactJSONUnMarshaler
	if err := enc.Unmarshal(b, &uf); err != nil {
		return e.Wrap(err)
	}

	fact.BaseFact.SetJSONUnmarshaler(uf.BaseFactJSONUnmarshaler)

	return fact.unpack(enc, uf.Owner, uf.Items)
}

type IssueSecurityTokensMarshaler struct {
	common.BaseOperationJSONMarshaler
}

func (op IssueSecurityTokens) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(IssueSecurityTokensMarshaler{
		BaseOperationJSONMarshaler: op.BaseOperation.JSONMarshaler(),
	})
}

func (op *IssueSecurityTokens) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringError("failed to decode json of MintSecurityTokens")

	var ubo common.BaseOperation
	if err := ubo.DecodeJSON(b, enc); err != nil {
		return e.Wrap(err)
	}

	op.BaseOperation = ubo

	return nil
}
