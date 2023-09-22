package sto

import (
	currencytypes "github.com/ProtoconNet/mitum-currency/v3/types"
	stotypes "github.com/ProtoconNet/mitum-sto/types/sto"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	jsonenc "github.com/ProtoconNet/mitum2/util/encoder/json"
	"github.com/ProtoconNet/mitum2/util/hint"
)

type TransferByPartitionItemJSONMarshaler struct {
	hint.BaseHinter
	Contract    base.Address             `json:"contract"`
	TokenHolder base.Address             `json:"tokenholder"`
	Receiver    base.Address             `json:"receiver"`
	Partition   stotypes.Partition       `json:"partition"`
	Amount      string                   `json:"amount"`
	Currency    currencytypes.CurrencyID `json:"currency"`
}

func (it TransferByPartitionItem) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(TransferByPartitionItemJSONMarshaler{
		BaseHinter:  it.BaseHinter,
		Contract:    it.contract,
		TokenHolder: it.tokenholder,
		Receiver:    it.receiver,
		Partition:   it.partition,
		Amount:      it.amount.String(),
		Currency:    it.currency,
	})
}

type TransferByPartitionItemJSONUnMarshaler struct {
	Hint        hint.Hint `json:"_hint"`
	Contract    string    `json:"contract"`
	TokenHolder string    `json:"tokenholder"`
	Receiver    string    `json:"receiver"`
	Partition   string    `json:"partition"`
	Amount      string    `json:"amount"`
	Currency    string    `json:"currency"`
}

func (it *TransferByPartitionItem) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringError("failed to decode json of TransferByPartitionItem")

	var uit TransferByPartitionItemJSONUnMarshaler
	if err := enc.Unmarshal(b, &uit); err != nil {
		return e.Wrap(err)
	}

	return it.unpack(enc, uit.Hint, uit.Contract, uit.TokenHolder, uit.Receiver, uit.Partition, uit.Amount, uit.Currency)
}
