package sto

import (
	currencytypes "github.com/ProtoconNet/mitum-currency/v3/types"
	stotypes "github.com/ProtoconNet/mitum-sto/types/sto"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	jsonenc "github.com/ProtoconNet/mitum2/util/encoder/json"
	"github.com/ProtoconNet/mitum2/util/hint"
)

type IssueItemJSONMarshaler struct {
	hint.BaseHinter
	Contract  base.Address             `json:"contract"`
	Receiver  base.Address             `json:"receiver"`
	Amount    string                   `json:"amount"`
	Partition stotypes.Partition       `json:"partition"`
	Currency  currencytypes.CurrencyID `json:"currency"`
}

func (it IssueItem) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(IssueItemJSONMarshaler{
		BaseHinter: it.BaseHinter,
		Contract:   it.contract,
		Receiver:   it.receiver,
		Amount:     it.amount.String(),
		Partition:  it.partition,
		Currency:   it.currency,
	})
}

type IssueItemJSONUnMarshaler struct {
	Hint      hint.Hint `json:"_hint"`
	Contract  string    `json:"contract"`
	Receiver  string    `json:"receiver"`
	Amount    string    `json:"amount"`
	Partition string    `json:"partition"`
	Currency  string    `json:"currency"`
}

func (it *IssueItem) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringError("failed to decode json of IssueItem")

	var uit IssueItemJSONUnMarshaler
	if err := enc.Unmarshal(b, &uit); err != nil {
		return e.Wrap(err)
	}

	return it.unpack(enc, uit.Hint, uit.Contract, uit.Receiver, uit.Amount, uit.Partition, uit.Currency)
}
