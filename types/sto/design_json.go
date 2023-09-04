package sto

import (
	"encoding/json"

	currencytypes "github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum2/util"
	jsonenc "github.com/ProtoconNet/mitum2/util/encoder/json"
	"github.com/ProtoconNet/mitum2/util/hint"
)

type DesignJSONMarshaler struct {
	hint.BaseHinter
	STO         currencytypes.ContractID `json:"stoid"`
	Granularity uint64                   `json:"granularity"`
	Policy      Policy                   `json:"policy"`
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
	e := util.StringError("failed to decode json of Design")

	var ud DesignJSONUnmarshaler
	if err := enc.Unmarshal(b, &ud); err != nil {
		return e.Wrap(err)
	}

	return de.unpack(enc, ud.Hint, ud.STO, ud.Granularity, ud.Policy)
}
