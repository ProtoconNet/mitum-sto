package sto

import (
	"encoding/json"

	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
	"github.com/ProtoconNet/mitum2/util/hint"
)

type DesignJSONMarshaler struct {
	hint.BaseHinter
	Granularity uint64 `json:"granularity"`
	Policy      Policy `json:"policy"`
}

func (de Design) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(DesignJSONMarshaler{
		BaseHinter:  de.BaseHinter,
		Granularity: de.granularity,
		Policy:      de.policy,
	})
}

type DesignJSONUnmarshaler struct {
	Hint        hint.Hint       `json:"_hint"`
	Granularity uint64          `json:"granularity"`
	Policy      json.RawMessage `json:"policy"`
}

func (de *Design) DecodeJSON(b []byte, enc encoder.Encoder) error {
	e := util.StringError("failed to decode json of Design")

	var ud DesignJSONUnmarshaler
	if err := enc.Unmarshal(b, &ud); err != nil {
		return e.Wrap(err)
	}

	return de.unpack(enc, ud.Hint, ud.Granularity, ud.Policy)
}
