package sto

import (
	"encoding/json"

	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
	"github.com/ProtoconNet/mitum2/util/hint"
)

type PolicyJSONMarshaler struct {
	hint.BaseHinter
	Partitions []Partition `json:"partitions"`
	Aggregate  string      `json:"aggregate"`
	//Controllers []base.Address `json:"controllers"`
	Documents []Document `json:"documents"`
}

func (po Policy) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(PolicyJSONMarshaler{
		BaseHinter: po.BaseHinter,
		Partitions: po.partitions,
		Aggregate:  po.aggregate.String(),
		//Controllers: po.controllers,
		Documents: po.documents,
	})
}

type PolicyJSONUnmarshaler struct {
	Hint       hint.Hint `json:"_hint"`
	Partitions []string  `json:"partitions"`
	Aggregate  string    `json:"aggregate"`
	//Controllers []string        `json:"controllers"`
	Documents json.RawMessage `json:"documents"`
}

func (po *Policy) DecodeJSON(b []byte, enc encoder.Encoder) error {
	e := util.StringError("failed to decode json of Policy")

	var upo PolicyJSONUnmarshaler
	if err := enc.Unmarshal(b, &upo); err != nil {
		return e.Wrap(err)
	}

	return po.unpack(enc, upo.Hint, upo.Partitions, upo.Aggregate, upo.Documents)
}
