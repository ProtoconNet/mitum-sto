package sto

import (
	"encoding/json"

	"github.com/ProtoconNet/mitum-currency/v2/currency"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	jsonenc "github.com/ProtoconNet/mitum2/util/encoder/json"
)

type SetDocumentsFactJSONMarshaler struct {
	base.BaseFactJSONMarshaler
	Owner base.Address       `json:"sender"`
	Items []SetDocumentsItem `json:"items"`
}

func (fact SetDocumentsFact) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(SetDocumentsFactJSONMarshaler{
		BaseFactJSONMarshaler: fact.BaseFact.JSONMarshaler(),
		Owner:                 fact.sender,
		Items:                 fact.items,
	})
}

type SetDocumentsFactJSONUnMarshaler struct {
	base.BaseFactJSONUnmarshaler
	Owner string          `json:"sender"`
	Items json.RawMessage `json:"items"`
}

func (fact *SetDocumentsFact) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode json of SetDocumentsFact")

	var uf SetDocumentsFactJSONUnMarshaler
	if err := enc.Unmarshal(b, &uf); err != nil {
		return e(err, "")
	}

	fact.BaseFact.SetJSONUnmarshaler(uf.BaseFactJSONUnmarshaler)

	return fact.unpack(enc, uf.Owner, uf.Items)
}

type SetDocumentsMarshaler struct {
	currency.BaseOperationJSONMarshaler
}

func (op SetDocuments) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(SetDocumentsMarshaler{
		BaseOperationJSONMarshaler: op.BaseOperation.JSONMarshaler(),
	})
}

func (op *SetDocuments) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode json of SetDocuments")

	var ubo currency.BaseOperation
	if err := ubo.DecodeJSON(b, enc); err != nil {
		return e(err, "")
	}

	op.BaseOperation = ubo

	return nil
}
