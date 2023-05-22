package kyc

import (
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	jsonenc "github.com/ProtoconNet/mitum2/util/encoder/json"
	"github.com/ProtoconNet/mitum2/util/hint"
)

type KYCPolicyJSONMarshaler struct {
	hint.BaseHinter
	Controllers []base.Address `json:"controllers"`
}

func (po KYCPolicy) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(KYCPolicyJSONMarshaler{
		BaseHinter:  po.BaseHinter,
		Controllers: po.controllers,
	})
}

type KYCPolicyJSONUnmarshaler struct {
	Hint        hint.Hint `json:"_hint"`
	Controllers []string  `json:"controllers"`
}

func (po *KYCPolicy) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode json of KYCPolicy")

	var upo KYCPolicyJSONUnmarshaler
	if err := enc.Unmarshal(b, &upo); err != nil {
		return e(err, "")
	}

	return po.unpack(enc, upo.Hint, upo.Controllers)
}
