package sto

import (
	currencytypes "github.com/ProtoconNet/mitum-currency/v3/types"
	stotypes "github.com/ProtoconNet/mitum-sto/types/sto"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
	"github.com/ProtoconNet/mitum2/util/hint"
)

type CreateSecurityTokenItemJSONMarshaler struct {
	hint.BaseHinter
	Contract         base.Address             `json:"contract"`
	Granularity      uint64                   `json:"granularity"`
	DefaultPartition stotypes.Partition       `json:"default_partition"`
	Currency         currencytypes.CurrencyID `json:"currency"`
}

func (it CreateSecurityTokenItem) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(CreateSecurityTokenItemJSONMarshaler{
		BaseHinter:       it.BaseHinter,
		Contract:         it.contract,
		Granularity:      it.granularity,
		DefaultPartition: it.defaultPartition,
		Currency:         it.currency,
	})
}

type CreateSecurityTokenItemJSONUnMarshaler struct {
	Hint             hint.Hint `json:"_hint"`
	Contract         string    `json:"contract"`
	Granularity      uint64    `json:"granularity"`
	DefaultPartition string    `json:"default_partition"`
	Currency         string    `json:"currency"`
}

func (it *CreateSecurityTokenItem) DecodeJSON(b []byte, enc encoder.Encoder) error {
	e := util.StringError("failed to decode json of CreateSecurityTokenItem")

	var uit CreateSecurityTokenItemJSONUnMarshaler
	if err := enc.Unmarshal(b, &uit); err != nil {
		return e.Wrap(err)
	}

	return it.unpack(enc, uit.Hint, uit.Contract, uit.Granularity, uit.DefaultPartition, uit.Currency)
}
