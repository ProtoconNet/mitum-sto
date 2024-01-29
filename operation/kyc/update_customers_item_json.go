package kyc

import (
	currencytypes "github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
	"github.com/ProtoconNet/mitum2/util/hint"
)

type UpdateCustomersItemJSONMarshaler struct {
	hint.BaseHinter
	Contract base.Address             `json:"contract"`
	Customer base.Address             `json:"customer"`
	Status   bool                     `json:"status"`
	Currency currencytypes.CurrencyID `json:"currency"`
}

func (it UpdateCustomersItem) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(UpdateCustomersItemJSONMarshaler{
		BaseHinter: it.BaseHinter,
		Contract:   it.contract,
		Customer:   it.customer,
		Status:     it.status,
		Currency:   it.currency,
	})
}

type UpdateCustomersItemJSONUnMarshaler struct {
	Hint     hint.Hint `json:"_hint"`
	Contract string    `json:"contract"`
	Customer string    `json:"customer"`
	Status   bool      `json:"status"`
	Currency string    `json:"currency"`
}

func (it *UpdateCustomersItem) DecodeJSON(b []byte, enc encoder.Encoder) error {
	e := util.StringError("failed to decode json of UpdateCustomersItem")

	var uit UpdateCustomersItemJSONUnMarshaler
	if err := enc.Unmarshal(b, &uit); err != nil {
		return e.Wrap(err)
	}

	return it.unpack(enc, uit.Hint, uit.Contract, uit.Customer, uit.Status, uit.Currency)
}
