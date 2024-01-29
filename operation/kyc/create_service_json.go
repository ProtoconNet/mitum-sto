package kyc

import (
	"github.com/ProtoconNet/mitum-currency/v3/common"
	"github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
)

type CreateServiceFactJSONMarshaler struct {
	base.BaseFactJSONMarshaler
	Owner       base.Address     `json:"sender"`
	Contract    base.Address     `json:"contract"`
	Controllers []base.Address   `json:"controllers"`
	Currency    types.CurrencyID `json:"currency"`
}

func (fact CreateServiceFact) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(CreateServiceFactJSONMarshaler{
		BaseFactJSONMarshaler: fact.BaseFact.JSONMarshaler(),
		Owner:                 fact.sender,
		Contract:              fact.contract,
		Controllers:           fact.controllers,
		Currency:              fact.currency,
	})
}

type CreateServiceFactJSONUnMarshaler struct {
	base.BaseFactJSONUnmarshaler
	Owner       string   `json:"sender"`
	Contract    string   `json:"contract"`
	Controllers []string `json:"controllers"`
	Currency    string   `json:"currency"`
}

func (fact *CreateServiceFact) DecodeJSON(b []byte, enc encoder.Encoder) error {
	e := util.StringError("failed to decode json of CreateServiceFact")

	var uf CreateServiceFactJSONUnMarshaler
	if err := enc.Unmarshal(b, &uf); err != nil {
		return e.Wrap(err)
	}

	fact.BaseFact.SetJSONUnmarshaler(uf.BaseFactJSONUnmarshaler)

	return fact.unpack(enc, uf.Owner, uf.Contract, uf.Controllers, uf.Currency)
}

type CreateServiceMarshaler struct {
	common.BaseOperationJSONMarshaler
}

func (op CreateService) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(CreateServiceMarshaler{
		BaseOperationJSONMarshaler: op.BaseOperation.JSONMarshaler(),
	})
}

func (op *CreateService) DecodeJSON(b []byte, enc encoder.Encoder) error {
	e := util.StringError("failed to decode json of CreateService")

	var ubo common.BaseOperation
	if err := ubo.DecodeJSON(b, enc); err != nil {
		return e.Wrap(err)
	}

	op.BaseOperation = ubo

	return nil
}
