package sto

import (
	"encoding/json"

	"github.com/ProtoconNet/mitum-currency/v3/common"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	jsonenc "github.com/ProtoconNet/mitum2/util/encoder/json"
)

type TransferByPartitionFactJSONMarshaler struct {
	base.BaseFactJSONMarshaler
	Owner base.Address              `json:"sender"`
	Items []TransferByPartitionItem `json:"items"`
}

func (fact TransferByPartitionFact) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(TransferByPartitionFactJSONMarshaler{
		BaseFactJSONMarshaler: fact.BaseFact.JSONMarshaler(),
		Owner:                 fact.sender,
		Items:                 fact.items,
	})
}

type TransferByPartitionFactJSONUnMarshaler struct {
	base.BaseFactJSONUnmarshaler
	Owner string          `json:"sender"`
	Items json.RawMessage `json:"items"`
}

func (fact *TransferByPartitionFact) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringError("failed to decode json of TransferByPartitionFact")

	var uf TransferByPartitionFactJSONUnMarshaler
	if err := enc.Unmarshal(b, &uf); err != nil {
		return e.Wrap(err)
	}

	fact.BaseFact.SetJSONUnmarshaler(uf.BaseFactJSONUnmarshaler)

	return fact.unpack(enc, uf.Owner, uf.Items)
}

type TransferByPartitionMarshaler struct {
	common.BaseOperationJSONMarshaler
}

func (op TransferByPartition) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(TransferByPartitionMarshaler{
		BaseOperationJSONMarshaler: op.BaseOperation.JSONMarshaler(),
	})
}

func (op *TransferByPartition) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringError("failed to decode json of TransferByPartition")

	var ubo common.BaseOperation
	if err := ubo.DecodeJSON(b, enc); err != nil {
		return e.Wrap(err)
	}

	op.BaseOperation = ubo

	return nil
}
