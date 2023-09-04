package sto

import (
	currencytypes "github.com/ProtoconNet/mitum-currency/v3/types"
	stotypes "github.com/ProtoconNet/mitum-sto/types/sto"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	jsonenc "github.com/ProtoconNet/mitum2/util/encoder/json"
	"github.com/ProtoconNet/mitum2/util/hint"
)

type TransferSecurityTokensPartitionItemJSONMarshaler struct {
	hint.BaseHinter
	Contract    base.Address             `json:"contract"`
	STO         currencytypes.ContractID `json:"stoid"`
	TokenHolder base.Address             `json:"tokenholder"`
	Receiver    base.Address             `json:"receiver"`
	Partition   stotypes.Partition       `json:"partition"`
	Amount      string                   `json:"amount"`
	Currency    currencytypes.CurrencyID `json:"currency"`
}

func (it TransferSecurityTokensPartitionItem) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(TransferSecurityTokensPartitionItemJSONMarshaler{
		BaseHinter:  it.BaseHinter,
		Contract:    it.contract,
		STO:         it.stoID,
		TokenHolder: it.tokenholder,
		Receiver:    it.receiver,
		Partition:   it.partition,
		Amount:      it.amount.String(),
		Currency:    it.currency,
	})
}

type TransferSecurityTokensPartitionItemJSONUnMarshaler struct {
	Hint        hint.Hint `json:"_hint"`
	Contract    string    `json:"contract"`
	STO         string    `json:"stoid"`
	TokenHolder string    `json:"tokenholder"`
	Receiver    string    `json:"receiver"`
	Partition   string    `json:"partition"`
	Amount      string    `json:"amount"`
	Currency    string    `json:"currency"`
}

func (it *TransferSecurityTokensPartitionItem) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringError("failed to decode json of TransferSecurityTokensPartitionItem")

	var uit TransferSecurityTokensPartitionItemJSONUnMarshaler
	if err := enc.Unmarshal(b, &uit); err != nil {
		return e.Wrap(err)
	}

	return it.unpack(enc, uit.Hint, uit.Contract, uit.STO, uit.TokenHolder, uit.Receiver, uit.Partition, uit.Amount, uit.Currency)
}
