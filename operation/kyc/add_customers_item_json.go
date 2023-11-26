package kyc

import (
	"github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	jsonenc "github.com/ProtoconNet/mitum2/util/encoder/json"
	"github.com/ProtoconNet/mitum2/util/hint"
)

type AddCustomerItemJSONMarshaler struct {
	hint.BaseHinter
	Contract base.Address     `json:"contract"`
	Customer base.Address     `json:"customer"`
	Status   bool             `json:"status"`
	Currency types.CurrencyID `json:"currency"`
}

func (it AddCustomerItem) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(AddCustomerItemJSONMarshaler{
		BaseHinter: it.BaseHinter,
		Contract:   it.contract,
		Customer:   it.customer,
		Status:     it.status,
		Currency:   it.currency,
	})
}

type AddCustomerItemJSONUnMarshaler struct {
	Hint     hint.Hint `json:"_hint"`
	Contract string    `json:"contract"`
	Customer string    `json:"customer"`
	Status   bool      `json:"status"`
	Currency string    `json:"currency"`
}

func (it *AddCustomerItem) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringError("failed to decode json of AddCustomerItem")

	var uit AddCustomerItemJSONUnMarshaler
	if err := enc.Unmarshal(b, &uit); err != nil {
		return e.Wrap(err)
	}

	return it.unpack(enc, uit.Hint, uit.Contract, uit.Customer, uit.Status, uit.Currency)
}
