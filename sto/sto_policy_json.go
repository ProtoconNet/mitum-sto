package sto

import (
	"encoding/json"

	"github.com/ProtoconNet/mitum-currency/v2/currency"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	jsonenc "github.com/ProtoconNet/mitum2/util/encoder/json"
	"github.com/ProtoconNet/mitum2/util/hint"
)

type STOPolicyJSONMarshaler struct {
	hint.BaseHinter
	Partitions  []Partition     `json:"partitions"`
	Aggregate   currency.Amount `json:"aggregate"`
	Controllers []base.Address  `json:"controllers"`
	Documents   []Document      `json:"documents"`
}

func (po STOPolicy) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(STOPolicyJSONMarshaler{
		BaseHinter:  po.BaseHinter,
		Partitions:  po.partitions,
		Aggregate:   po.aggregate,
		Controllers: po.controllers,
		Documents:   po.documents,
	})
}

type STOPolicyJSONUnmarshaler struct {
	Hint        hint.Hint       `json:"_hint"`
	Partitions  json.RawMessage `json:"partitions"`
	Aggregate   json.RawMessage `json:"aggregate"`
	Controllers []string        `json:"controllers"`
	Documents   json.RawMessage `json:"documents"`
}

func (po *STOPolicy) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode json of STOPolicy")

	var upo STOPolicyJSONUnmarshaler
	if err := enc.Unmarshal(b, &upo); err != nil {
		return e(err, "")
	}

	return po.unpack(enc, upo.Hint, upo.Partitions, upo.Aggregate, upo.Controllers, upo.Documents)
}
