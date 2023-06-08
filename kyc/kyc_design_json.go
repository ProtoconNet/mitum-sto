package kyc

import (
	"encoding/json"

	currencybase "github.com/ProtoconNet/mitum-currency/v3/base"
	"github.com/ProtoconNet/mitum2/util"
	jsonenc "github.com/ProtoconNet/mitum2/util/encoder/json"
	"github.com/ProtoconNet/mitum2/util/hint"
)

type DesignJSONMarshaler struct {
	hint.BaseHinter
	KYC    currencybase.ContractID `json:"kycid"`
	Policy Policy                  `json:"policy"`
}

func (de Design) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(DesignJSONMarshaler{
		BaseHinter: de.BaseHinter,
		KYC:        de.kycID,
		Policy:     de.policy,
	})
}

type DesignJSONUnmarshaler struct {
	Hint   hint.Hint       `json:"_hint"`
	KYC    string          `json:"kycid"`
	Policy json.RawMessage `json:"policy"`
}

func (de *Design) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode json of Design")

	var ud DesignJSONUnmarshaler
	if err := enc.Unmarshal(b, &ud); err != nil {
		return e(err, "")
	}

	return de.unpack(enc, ud.Hint, ud.KYC, ud.Policy)
}
