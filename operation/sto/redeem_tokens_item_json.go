package sto

import (
	currencytypes "github.com/ProtoconNet/mitum-currency/v3/types"
	stotypes "github.com/ProtoconNet/mitum-sto/types/sto"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	jsonenc "github.com/ProtoconNet/mitum2/util/encoder/json"
	"github.com/ProtoconNet/mitum2/util/hint"
)

type RedeemTokensItemJSONMarshaler struct {
	hint.BaseHinter
	Contract    base.Address             `json:"contract"`
	STO         currencytypes.ContractID `json:"stoid"`
	TokenHolder base.Address             `json:"tokenholder"`
	Amount      string                   `json:"amount"`
	Partition   stotypes.Partition       `json:"partition"`
	Currency    currencytypes.CurrencyID `json:"currency"`
}

func (it RedeemTokensItem) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(RedeemTokensItemJSONMarshaler{
		BaseHinter:  it.BaseHinter,
		Contract:    it.contract,
		STO:         it.stoID,
		TokenHolder: it.tokenHolder,
		Amount:      it.amount.String(),
		Partition:   it.partition,
		Currency:    it.currency,
	})
}

type RedeemTokensItemJSONUnMarshaler struct {
	Hint        hint.Hint `json:"_hint"`
	Contract    string    `json:"contract"`
	STO         string    `json:"stoid"`
	TokenHolder string    `json:"tokenholder"`
	Amount      string    `json:"amount"`
	Partition   string    `json:"partition"`
	Currency    string    `json:"currency"`
}

func (it *RedeemTokensItem) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringError("failed to decode json of RedeemTokensItem")

	var uit RedeemTokensItemJSONUnMarshaler
	if err := enc.Unmarshal(b, &uit); err != nil {
		return e.Wrap(err)
	}

	return it.unpack(enc, uit.Hint, uit.Contract, uit.STO, uit.TokenHolder, uit.Amount, uit.Partition, uit.Currency)
}
