package kyc

import (
	extensioncurrency "github.com/ProtoconNet/mitum-currency-extension/v2/currency"
	"github.com/ProtoconNet/mitum-currency/v2/currency"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	jsonenc "github.com/ProtoconNet/mitum2/util/encoder/json"
	"github.com/ProtoconNet/mitum2/util/hint"
)

type RemoveControllersItemJSONMarshaler struct {
	hint.BaseHinter
	Contract   base.Address                 `json:"contract"`
	KYC        extensioncurrency.ContractID `json:"kycid"`
	Controller base.Address                 `json:"controller"`
	Currency   currency.CurrencyID          `json:"currency"`
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
	e := util.StringErrorFunc("failed to decode json of RemoveControllersItem")

	var uit RemoveControllersItemJSONUnMarshaler
	if err := enc.Unmarshal(b, &uit); err != nil {
		return e(err, "")
	}

	return it.unpack(enc, uit.Hint, uit.Contract, uit.KYC, uit.Controller, uit.Currency)
}
