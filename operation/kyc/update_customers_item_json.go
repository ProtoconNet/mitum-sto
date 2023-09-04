package kyc

import (
	currencytypes "github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	jsonenc "github.com/ProtoconNet/mitum2/util/encoder/json"
	"github.com/ProtoconNet/mitum2/util/hint"
)

type UpdateCustomersItemJSONMarshaler struct {
	hint.BaseHinter
	Contract base.Address             `json:"contract"`
	KYC      currencytypes.ContractID `json:"kycid"`
	Customer base.Address             `json:"customer"`
	Status   bool                     `json:"status"`
	Currency currencytypes.CurrencyID `json:"currency"`
}

func (it UpdateCustomersItem) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(UpdateCustomersItemJSONMarshaler{
		BaseHinter: it.BaseHinter,
		Contract:   it.contract,
		KYC:        it.kycID,
		Customer:   it.customer,
		Status:     it.status,
		Currency:   it.currency,
	})
}

type UpdateCustomersItemJSONUnMarshaler struct {
	Hint     hint.Hint `json:"_hint"`
	Contract string    `json:"contract"`
	KYC      string    `json:"kycid"`
	Customer string    `json:"customer"`
	Status   bool      `json:"status"`
	Currency string    `json:"currency"`
}

func (it *UpdateCustomersItem) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringError("failed to decode json of UpdateCustomersItem")

	var uit UpdateCustomersItemJSONUnMarshaler
	if err := enc.Unmarshal(b, &uit); err != nil {
		return e.Wrap(err)
	}

	return it.unpack(enc, uit.Hint, uit.Contract, uit.KYC, uit.Customer, uit.Status, uit.Currency)
}
