package kyc

import (
	currencybase "github.com/ProtoconNet/mitum-currency/v3/base"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	jsonenc "github.com/ProtoconNet/mitum2/util/encoder/json"
	"github.com/ProtoconNet/mitum2/util/hint"
)

type AddCustomersItemJSONMarshaler struct {
	hint.BaseHinter
	Contract base.Address            `json:"contract"`
	KYC      currencybase.ContractID `json:"kycid"`
	Customer base.Address            `json:"customer"`
	Status   bool                    `json:"status"`
	Currency currencybase.CurrencyID `json:"currency"`
}

func (it AddCustomersItem) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(AddCustomersItemJSONMarshaler{
		BaseHinter: it.BaseHinter,
		Contract:   it.contract,
		KYC:        it.kycID,
		Customer:   it.customer,
		Status:     it.status,
		Currency:   it.currency,
	})
}

type AddCustomersItemJSONUnMarshaler struct {
	Hint     hint.Hint `json:"_hint"`
	Contract string    `json:"contract"`
	KYC      string    `json:"kycid"`
	Customer string    `json:"customer"`
	Status   bool      `json:"status"`
	Currency string    `json:"currency"`
}

func (it *AddCustomersItem) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode json of AddCustomersItem")

	var uit AddCustomersItemJSONUnMarshaler
	if err := enc.Unmarshal(b, &uit); err != nil {
		return e(err, "")
	}

	return it.unpack(enc, uit.Hint, uit.Contract, uit.KYC, uit.Customer, uit.Status, uit.Currency)
}
