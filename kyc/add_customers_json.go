package kyc

import (
	"encoding/json"

	currencybase "github.com/ProtoconNet/mitum-currency/v3/base"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	jsonenc "github.com/ProtoconNet/mitum2/util/encoder/json"
)

type AddCustomersFactJSONMarshaler struct {
	base.BaseFactJSONMarshaler
	Owner base.Address       `json:"sender"`
	Items []AddCustomersItem `json:"items"`
}

func (fact AddCustomersFact) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(AddCustomersFactJSONMarshaler{
		BaseFactJSONMarshaler: fact.BaseFact.JSONMarshaler(),
		Owner:                 fact.sender,
		Items:                 fact.items,
	})
}

type AddCustomersFactJSONUnMarshaler struct {
	base.BaseFactJSONUnmarshaler
	Owner string          `json:"sender"`
	Items json.RawMessage `json:"items"`
}

func (fact *AddCustomersFact) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode json of AddCustomersFact")

	var uf AddCustomersFactJSONUnMarshaler
	if err := enc.Unmarshal(b, &uf); err != nil {
		return e(err, "")
	}

	fact.BaseFact.SetJSONUnmarshaler(uf.BaseFactJSONUnmarshaler)

	return fact.unpack(enc, uf.Owner, uf.Items)
}

type AddCustomersMarshaler struct {
	currencybase.BaseOperationJSONMarshaler
}

func (op AddCustomers) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(AddCustomersMarshaler{
		BaseOperationJSONMarshaler: op.BaseOperation.JSONMarshaler(),
	})
}

func (op *AddCustomers) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode json of AddCustomers")

	var ubo currencybase.BaseOperation
	if err := ubo.DecodeJSON(b, enc); err != nil {
		return e(err, "")
	}

	op.BaseOperation = ubo

	return nil
}
