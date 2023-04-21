package sto

import (
	extensioncurrency "github.com/ProtoconNet/mitum-currency-extension/v2/currency"
	"github.com/ProtoconNet/mitum-currency/v2/currency"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	jsonenc "github.com/ProtoconNet/mitum2/util/encoder/json"
	"github.com/ProtoconNet/mitum2/util/hint"
)

type CreateSecurityTokensItemJSONMarshaler struct {
	hint.BaseHinter
	Contract         base.Address                 `json:"contract"`
	STO              extensioncurrency.ContractID `json:"stoid"`
	Granularity      uint64                       `json:"granularity"`
	DefaultPartition Partition                    `json:"default_partition"`
	Controllers      []base.Address               `json:"controllers"`
	Currency         currency.CurrencyID          `json:"currency"`
}

func (it CreateSecurityTokensItem) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(CreateSecurityTokensItemJSONMarshaler{
		BaseHinter:       it.BaseHinter,
		Contract:         it.contract,
		STO:              it.stoID,
		Granularity:      it.granularity,
		DefaultPartition: it.defaultPartition,
		Controllers:      it.controllers,
		Currency:         it.currency,
	})
}

type CreateSecurityTokensItemJSONUnMarshaler struct {
	Hint             hint.Hint `json:"_hint"`
	Contract         string    `json:"contract"`
	STO              string    `json:"stoid"`
	Granularity      uint64    `json:"granularity"`
	DefaultPartition string    `json:"default_partition"`
	Controllers      []string  `json:"controllers"`
	Currency         string    `json:"currency"`
}

func (it *CreateSecurityTokensItem) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode json of CreateSecurityTokensItem")

	var uit CreateSecurityTokensItemJSONUnMarshaler
	if err := enc.Unmarshal(b, &uit); err != nil {
		return e(err, "")
	}

	return it.unpack(enc, uit.Hint, uit.Contract, uit.STO, uit.Granularity, uit.DefaultPartition, uit.Controllers, uit.Currency)
}
