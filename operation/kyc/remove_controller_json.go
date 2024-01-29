package kyc

import (
	"encoding/json"

	"github.com/ProtoconNet/mitum-currency/v3/common"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
)

type RemoveControllerFactJSONMarshaler struct {
	base.BaseFactJSONMarshaler
	Owner base.Address           `json:"sender"`
	Items []RemoveControllerItem `json:"items"`
}

func (fact RemoveControllerFact) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(RemoveControllerFactJSONMarshaler{
		BaseFactJSONMarshaler: fact.BaseFact.JSONMarshaler(),
		Owner:                 fact.sender,
		Items:                 fact.items,
	})
}

type RemoveControllerFactJSONUnMarshaler struct {
	base.BaseFactJSONUnmarshaler
	Owner string          `json:"sender"`
	Items json.RawMessage `json:"items"`
}

func (fact *RemoveControllerFact) DecodeJSON(b []byte, enc encoder.Encoder) error {
	e := util.StringError("failed to decode json of RemoveControllerFact")

	var uf RemoveControllerFactJSONUnMarshaler
	if err := enc.Unmarshal(b, &uf); err != nil {
		return e.Wrap(err)
	}

	fact.BaseFact.SetJSONUnmarshaler(uf.BaseFactJSONUnmarshaler)

	return fact.unpack(enc, uf.Owner, uf.Items)
}

type RemoveControllerMarshaler struct {
	common.BaseOperationJSONMarshaler
}

func (op RemoveController) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(RemoveControllerMarshaler{
		BaseOperationJSONMarshaler: op.BaseOperation.JSONMarshaler(),
	})
}

func (op *RemoveController) DecodeJSON(b []byte, enc encoder.Encoder) error {
	e := util.StringError("failed to decode json of RemoveController")

	var ubo common.BaseOperation
	if err := ubo.DecodeJSON(b, enc); err != nil {
		return e.Wrap(err)
	}

	op.BaseOperation = ubo

	return nil
}
