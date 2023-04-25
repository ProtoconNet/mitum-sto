package sto

import (
	extensioncurrency "github.com/ProtoconNet/mitum-currency-extension/v2/currency"
	"github.com/ProtoconNet/mitum-currency/v2/currency"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	jsonenc "github.com/ProtoconNet/mitum2/util/encoder/json"
)

type SetDocumentsFactJSONMarshaler struct {
	base.BaseFactJSONMarshaler
	Owner        base.Address                 `json:"sender"`
	Contract     base.Address                 `json:"contract"`
	STOID        extensioncurrency.ContractID `json:"stoid"`
	Title        string                       `json:"title"`
	Uri          URI                          `json:"uri"`
	DocumentHash string                       `json:"documenthash"`
	Currency     currency.CurrencyID          `json:"currencyid"`
}

func (fact SetDocumentsFact) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(SetDocumentsFactJSONMarshaler{
		BaseFactJSONMarshaler: fact.BaseFact.JSONMarshaler(),
		Owner:                 fact.sender,
		Contract:              fact.contract,
		STOID:                 fact.stoID,
		Title:                 fact.title,
		Uri:                   fact.uri,
		DocumentHash:          fact.documentHash,
		Currency:              fact.currency,
	})
}

type SetDocumentsFactJSONUnMarshaler struct {
	base.BaseFactJSONUnmarshaler
	Owner        string `json:"sender"`
	Contract     string `json:"contract"`
	STOID        string `json:"stoid"`
	Title        string `json:"title"`
	Uri          string `json:"uri"`
	DocumentHash string `json:"documenthash"`
	Currency     string `json:"currencyid"`
}

func (fact *SetDocumentsFact) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode json of SetDocumentsFact")

	var uf SetDocumentsFactJSONUnMarshaler
	if err := enc.Unmarshal(b, &uf); err != nil {
		return e(err, "")
	}

	fact.BaseFact.SetJSONUnmarshaler(uf.BaseFactJSONUnmarshaler)

	return fact.unpack(enc, uf.Owner, uf.Contract, uf.STOID, uf.Title, uf.Uri, uf.DocumentHash, uf.Currency)
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
