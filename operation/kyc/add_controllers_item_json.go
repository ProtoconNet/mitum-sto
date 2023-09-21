package kyc

import (
	currencytypes "github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	jsonenc "github.com/ProtoconNet/mitum2/util/encoder/json"
	"github.com/ProtoconNet/mitum2/util/hint"
)

type AddControllersItemJSONMarshaler struct {
	hint.BaseHinter
	Contract   base.Address             `json:"contract"`
	Controller base.Address             `json:"controller"`
	Currency   currencytypes.CurrencyID `json:"currency"`
}

func (it AddControllersItem) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(AddControllersItemJSONMarshaler{
		BaseHinter: it.BaseHinter,
		Contract:   it.contract,
		Controller: it.controller,
		Currency:   it.currency,
	})
}

type AddControllersItemJSONUnMarshaler struct {
	Hint       hint.Hint `json:"_hint"`
	Contract   string    `json:"contract"`
	Controller string    `json:"controller"`
	Currency   string    `json:"currency"`
}

func (it *AddControllersItem) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringError("failed to decode json of AddControllersItem")

	var uit AddControllersItemJSONUnMarshaler
	if err := enc.Unmarshal(b, &uit); err != nil {
		return e.Wrap(err)
	}

	return it.unpack(enc, uit.Hint, uit.Contract, uit.Controller, uit.Currency)
}
