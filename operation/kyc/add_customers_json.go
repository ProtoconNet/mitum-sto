package kyc

import (
	"encoding/json"
	"github.com/ProtoconNet/mitum-currency/v3/common"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
)

type AddCustomerFactJSONMarshaler struct {
	base.BaseFactJSONMarshaler
	Owner base.Address      `json:"sender"`
	Items []AddCustomerItem `json:"items"`
}

func (fact AddCustomerFact) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(AddCustomerFactJSONMarshaler{
		BaseFactJSONMarshaler: fact.BaseFact.JSONMarshaler(),
		Owner:                 fact.sender,
		Items:                 fact.items,
	})
}

type AddCustomerFactJSONUnMarshaler struct {
	base.BaseFactJSONUnmarshaler
	Owner string          `json:"sender"`
	Items json.RawMessage `json:"items"`
}

func (fact *AddCustomerFact) DecodeJSON(b []byte, enc encoder.Encoder) error {
	e := util.StringError("failed to decode json of AddCustomerFact")

	var uf AddCustomerFactJSONUnMarshaler
	if err := enc.Unmarshal(b, &uf); err != nil {
		return e.Wrap(err)
	}

	fact.BaseFact.SetJSONUnmarshaler(uf.BaseFactJSONUnmarshaler)

	return fact.unpack(enc, uf.Owner, uf.Items)
}

type AddCustomerMarshaler struct {
	common.BaseOperationJSONMarshaler
}

func (op AddCustomer) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(AddCustomerMarshaler{
		BaseOperationJSONMarshaler: op.BaseOperation.JSONMarshaler(),
	})
}

func (op *AddCustomer) DecodeJSON(b []byte, enc encoder.Encoder) error {
	e := util.StringError("failed to decode json of AddCustomer")

	var ubo common.BaseOperation
	if err := ubo.DecodeJSON(b, enc); err != nil {
		return e.Wrap(err)
	}

	op.BaseOperation = ubo

	return nil
}
