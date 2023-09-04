package kyc

import (
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	jsonenc "github.com/ProtoconNet/mitum2/util/encoder/json"
	"github.com/ProtoconNet/mitum2/util/hint"
)

type PolicyJSONMarshaler struct {
	hint.BaseHinter
	Controllers []base.Address `json:"controllers"`
}

func (po Policy) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(PolicyJSONMarshaler{
		BaseHinter:  po.BaseHinter,
		Controllers: po.controllers,
	})
}

type PolicyJSONUnmarshaler struct {
	Hint        hint.Hint `json:"_hint"`
	Controllers []string  `json:"controllers"`
}

func (po *Policy) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringError("failed to decode json of Policy")

	var upo PolicyJSONUnmarshaler
	if err := enc.Unmarshal(b, &upo); err != nil {
		return e.Wrap(err)
	}

	return po.unpack(enc, upo.Hint, upo.Controllers)
}
