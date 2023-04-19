package sto

import (
	"encoding/json"

	extensioncurrency "github.com/ProtoconNet/mitum-currency-extension/v2/currency"
	"github.com/ProtoconNet/mitum2/util"
	jsonenc "github.com/ProtoconNet/mitum2/util/encoder/json"
	"github.com/ProtoconNet/mitum2/util/hint"
)

type STODesignJSONMarshaler struct {
	hint.BaseHinter
	STO    extensioncurrency.ContractID `json:"sto"`
	Policy STOPolicy                    `json:"policy"`
}

func (de STODesign) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(STODesignJSONMarshaler{
		BaseHinter: de.BaseHinter,
		STO:        de.stoID,
		Policy:     de.policy,
	})
}

type STODesignJSONUnmarshaler struct {
	Hint   hint.Hint       `json:"_hint"`
	STO    string          `json:"sto"`
	Policy json.RawMessage `json:"policy"`
}

func (de *STODesign) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode json of STODesign")

	var ud STODesignJSONUnmarshaler
	if err := enc.Unmarshal(b, &ud); err != nil {
		return e(err, "")
	}

	return de.unpack(enc, ud.Hint, ud.STO, ud.Policy)
}
