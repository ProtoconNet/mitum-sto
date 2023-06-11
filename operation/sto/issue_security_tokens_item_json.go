package sto

import (
	currencytypes "github.com/ProtoconNet/mitum-currency/v3/types"
	stotypes "github.com/ProtoconNet/mitum-sto/types/sto"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	jsonenc "github.com/ProtoconNet/mitum2/util/encoder/json"
	"github.com/ProtoconNet/mitum2/util/hint"
)

type IssueSecurityTokensItemJSONMarshaler struct {
	hint.BaseHinter
	Contract  base.Address             `json:"contract"`
	STO       currencytypes.ContractID `json:"stoid"`
	Receiver  base.Address             `json:"receiver"`
	Amount    string                   `json:"amount"`
	Partition stotypes.Partition       `json:"partition"`
	Currency  currencytypes.CurrencyID `json:"currency"`
}

func (it IssueSecurityTokensItem) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(IssueSecurityTokensItemJSONMarshaler{
		BaseHinter: it.BaseHinter,
		Contract:   it.contract,
		STO:        it.stoID,
		Receiver:   it.receiver,
		Amount:     it.amount.String(),
		Partition:  it.partition,
		Currency:   it.currency,
	})
}

type IssueSecurityTokensItemJSONUnMarshaler struct {
	Hint      hint.Hint `json:"_hint"`
	Contract  string    `json:"contract"`
	STO       string    `json:"stoid"`
	Receiver  string    `json:"receiver"`
	Amount    string    `json:"amount"`
	Partition string    `json:"partition"`
	Currency  string    `json:"currency"`
}

func (it *IssueSecurityTokensItem) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode json of IssueSecurityTokensItem")

	var uit IssueSecurityTokensItemJSONUnMarshaler
	if err := enc.Unmarshal(b, &uit); err != nil {
		return e(err, "")
	}

	return it.unpack(enc, uit.Hint, uit.Contract, uit.STO, uit.Receiver, uit.Amount, uit.Partition, uit.Currency)
}
