package sto

import (
	"encoding/json"

	extensioncurrency "github.com/ProtoconNet/mitum-currency-extension/v2/currency"
	"github.com/ProtoconNet/mitum2/util"
	jsonenc "github.com/ProtoconNet/mitum2/util/encoder/json"
	"github.com/ProtoconNet/mitum2/util/hint"
)

type DesignJSONMarshaler struct {
	hint.BaseHinter
	STO         extensioncurrency.ContractID `json:"stoid"`
	Granularity uint64                       `json:"granularity"`
	Policy      Policy                       `json:"policy"`
}

func (de Design) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(DesignJSONMarshaler{
		BaseHinter:  de.BaseHinter,
		STO:         de.stoID,
		Granularity: de.granularity,
		Policy:      de.policy,
	})
}

type DesignJSONUnmarshaler struct {
	Hint        hint.Hint       `json:"_hint"`
	STO         string          `json:"stoid"`
	Granularity uint64          `json:"granularity"`
	Policy      json.RawMessage `json:"policy"`
}

func (de *Design) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode json of Design")

	var ud DesignJSONUnmarshaler
	if err := enc.Unmarshal(b, &ud); err != nil {
		return e(err, "")
	}

	return de.unpack(enc, ud.Hint, ud.STO, ud.Granularity, ud.Policy)
}
