package kyc

import (
	currencytypes "github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
	"github.com/ProtoconNet/mitum2/util/hint"
)

type RemoveControllerItemJSONMarshaler struct {
	hint.BaseHinter
	Contract   base.Address             `json:"contract"`
	Controller base.Address             `json:"controller"`
	Currency   currencytypes.CurrencyID `json:"currency"`
}

func (it RemoveControllerItem) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(RemoveControllerItemJSONMarshaler{
		BaseHinter: it.BaseHinter,
		Contract:   it.contract,
		Controller: it.controller,
		Currency:   it.currency,
	})
}

type RemoveControllerItemJSONUnMarshaler struct {
	Hint       hint.Hint `json:"_hint"`
	Contract   string    `json:"contract"`
	Controller string    `json:"controller"`
	Currency   string    `json:"currency"`
}

func (it *RemoveControllerItem) DecodeJSON(b []byte, enc encoder.Encoder) error {
	e := util.StringError("failed to decode json of RemoveControllerItem")

	var uit RemoveControllerItemJSONUnMarshaler
	if err := enc.Unmarshal(b, &uit); err != nil {
		return e.Wrap(err)
	}

	return it.unpack(enc, uit.Hint, uit.Contract, uit.Controller, uit.Currency)
}
