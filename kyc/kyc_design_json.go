package kyc

import (
	"encoding/json"

	extensioncurrency "github.com/ProtoconNet/mitum-currency-extension/v2/currency"
	"github.com/ProtoconNet/mitum2/util"
	jsonenc "github.com/ProtoconNet/mitum2/util/encoder/json"
	"github.com/ProtoconNet/mitum2/util/hint"
)

type KYCDesignJSONMarshaler struct {
	hint.BaseHinter
	KYC    extensioncurrency.ContractID `json:"kycid"`
	Policy KYCPolicy                    `json:"policy"`
}

func (de KYCDesign) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(KYCDesignJSONMarshaler{
		BaseHinter: de.BaseHinter,
		KYC:        de.kycID,
		Policy:     de.policy,
	})
}

type KYCDesignJSONUnmarshaler struct {
	Hint   hint.Hint       `json:"_hint"`
	KYC    string          `json:"kycid"`
	Policy json.RawMessage `json:"policy"`
}

func (de *KYCDesign) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode json of KYCDesign")

	var ud KYCDesignJSONUnmarshaler
	if err := enc.Unmarshal(b, &ud); err != nil {
		return e(err, "")
	}

	return de.unpack(enc, ud.Hint, ud.KYC, ud.Policy)
}
