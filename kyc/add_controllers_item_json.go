package kyc

import (
	currencybase "github.com/ProtoconNet/mitum-currency/v3/base"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	jsonenc "github.com/ProtoconNet/mitum2/util/encoder/json"
	"github.com/ProtoconNet/mitum2/util/hint"
)

type AddControllersItemJSONMarshaler struct {
	hint.BaseHinter
	Contract   base.Address            `json:"contract"`
	KYC        currencybase.ContractID `json:"kycid"`
	Controller base.Address            `json:"controller"`
	Currency   currencybase.CurrencyID `json:"currency"`
}

func (it AddControllersItem) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(AddControllersItemJSONMarshaler{
		BaseHinter: it.BaseHinter,
		Contract:   it.contract,
		KYC:        it.kycID,
		Controller: it.controller,
		Currency:   it.currency,
	})
}

type AddControllersItemJSONUnMarshaler struct {
	Hint       hint.Hint `json:"_hint"`
	Contract   string    `json:"contract"`
	KYC        string    `json:"kycid"`
	Controller string    `json:"controller"`
	Currency   string    `json:"currency"`
}

func (it *AddControllersItem) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode json of AddControllersItem")

	var uit AddControllersItemJSONUnMarshaler
	if err := enc.Unmarshal(b, &uit); err != nil {
		return e(err, "")
	}

	return it.unpack(enc, uit.Hint, uit.Contract, uit.KYC, uit.Controller, uit.Currency)
}
