package sto

import (
	currencytypes "github.com/ProtoconNet/mitum-currency/v3/types"
	stotypes "github.com/ProtoconNet/mitum-sto/types/sto"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	jsonenc "github.com/ProtoconNet/mitum2/util/encoder/json"
	"github.com/ProtoconNet/mitum2/util/hint"
)

type RedeemItemJSONMarshaler struct {
	hint.BaseHinter
	Contract    base.Address             `json:"contract"`
	TokenHolder base.Address             `json:"tokenholder"`
	Amount      string                   `json:"amount"`
	Partition   stotypes.Partition       `json:"partition"`
	Currency    currencytypes.CurrencyID `json:"currency"`
}

func (it RedeemItem) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(RedeemItemJSONMarshaler{
		BaseHinter:  it.BaseHinter,
		Contract:    it.contract,
		TokenHolder: it.tokenHolder,
		Amount:      it.amount.String(),
		Partition:   it.partition,
		Currency:    it.currency,
	})
}

type RedeemItemJSONUnMarshaler struct {
	Hint        hint.Hint `json:"_hint"`
	Contract    string    `json:"contract"`
	TokenHolder string    `json:"tokenholder"`
	Amount      string    `json:"amount"`
	Partition   string    `json:"partition"`
	Currency    string    `json:"currency"`
}

func (it *RedeemItem) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringError("failed to decode json of RedeemItem")

	var uit RedeemItemJSONUnMarshaler
	if err := enc.Unmarshal(b, &uit); err != nil {
		return e.Wrap(err)
	}

	return it.unpack(enc, uit.Hint, uit.Contract, uit.TokenHolder, uit.Amount, uit.Partition, uit.Currency)
}
