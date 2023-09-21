package kyc

import (
	"github.com/ProtoconNet/mitum-currency/v3/common"
	currencytypes "github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	jsonenc "github.com/ProtoconNet/mitum2/util/encoder/json"
)

type CreateKYCServiceFactJSONMarshaler struct {
	base.BaseFactJSONMarshaler
	Owner       base.Address             `json:"sender"`
	Contract    base.Address             `json:"contract"`
	Controllers []base.Address           `json:"controllers"`
	Currency    currencytypes.CurrencyID `json:"currency"`
}

func (fact CreateKYCServiceFact) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(CreateKYCServiceFactJSONMarshaler{
		BaseFactJSONMarshaler: fact.BaseFact.JSONMarshaler(),
		Owner:                 fact.sender,
		Contract:              fact.contract,
		Controllers:           fact.controllers,
		Currency:              fact.currency,
	})
}

type CreateKYCServiceFactJSONUnMarshaler struct {
	base.BaseFactJSONUnmarshaler
	Owner       string   `json:"sender"`
	Contract    string   `json:"contract"`
	Controllers []string `json:"controllers"`
	Currency    string   `json:"currency"`
}

func (fact *CreateKYCServiceFact) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringError("failed to decode json of CreateKYCServiceFact")

	var uf CreateKYCServiceFactJSONUnMarshaler
	if err := enc.Unmarshal(b, &uf); err != nil {
		return e.Wrap(err)
	}

	fact.BaseFact.SetJSONUnmarshaler(uf.BaseFactJSONUnmarshaler)

	return fact.unpack(enc, uf.Owner, uf.Contract, uf.Controllers, uf.Currency)
}

type CreateKYCServiceMarshaler struct {
	common.BaseOperationJSONMarshaler
}

func (op CreateKYCService) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(CreateKYCServiceMarshaler{
		BaseOperationJSONMarshaler: op.BaseOperation.JSONMarshaler(),
	})
}

func (op *CreateKYCService) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringError("failed to decode json of CreateKYCService")

	var ubo common.BaseOperation
	if err := ubo.DecodeJSON(b, enc); err != nil {
		return e.Wrap(err)
	}

	op.BaseOperation = ubo

	return nil
}
