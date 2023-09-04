package kyc

import (
	currencytypes "github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	jsonenc "github.com/ProtoconNet/mitum2/util/encoder/json"
	"github.com/ProtoconNet/mitum2/util/hint"
)

type RemoveControllersItemJSONMarshaler struct {
	hint.BaseHinter
	Contract   base.Address             `json:"contract"`
	KYC        currencytypes.ContractID `json:"kycid"`
	Controller base.Address             `json:"controller"`
	Currency   currencytypes.CurrencyID `json:"currency"`
}

func (it RemoveControllersItem) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(RemoveControllersItemJSONMarshaler{
		BaseHinter: it.BaseHinter,
		Contract:   it.contract,
		KYC:        it.kycID,
		Controller: it.controller,
		Currency:   it.currency,
	})
}

type RemoveControllersItemJSONUnMarshaler struct {
	Hint       hint.Hint `json:"_hint"`
	Contract   string    `json:"contract"`
	KYC        string    `json:"kycid"`
	Controller string    `json:"controller"`
	Currency   string    `json:"currency"`
}

func (it *RemoveControllersItem) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringError("failed to decode json of RemoveControllersItem")

	var uit RemoveControllersItemJSONUnMarshaler
	if err := enc.Unmarshal(b, &uit); err != nil {
		return e.Wrap(err)
	}

	return it.unpack(enc, uit.Hint, uit.Contract, uit.KYC, uit.Controller, uit.Currency)
}
