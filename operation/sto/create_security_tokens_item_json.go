package sto

import (
	currencytypes "github.com/ProtoconNet/mitum-currency/v3/types"
	stotypes "github.com/ProtoconNet/mitum-sto/types/sto"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	jsonenc "github.com/ProtoconNet/mitum2/util/encoder/json"
	"github.com/ProtoconNet/mitum2/util/hint"
)

type CreateSecurityTokensItemJSONMarshaler struct {
	hint.BaseHinter
	Contract         base.Address             `json:"contract"`
	Granularity      uint64                   `json:"granularity"`
	DefaultPartition stotypes.Partition       `json:"default_partition"`
	Controllers      []base.Address           `json:"controllers"`
	Currency         currencytypes.CurrencyID `json:"currency"`
}

func (it CreateSecurityTokensItem) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(CreateSecurityTokensItemJSONMarshaler{
		BaseHinter:       it.BaseHinter,
		Contract:         it.contract,
		Granularity:      it.granularity,
		DefaultPartition: it.defaultPartition,
		Controllers:      it.controllers,
		Currency:         it.currency,
	})
}

type CreateSecurityTokensItemJSONUnMarshaler struct {
	Hint             hint.Hint `json:"_hint"`
	Contract         string    `json:"contract"`
	Granularity      uint64    `json:"granularity"`
	DefaultPartition string    `json:"default_partition"`
	Controllers      []string  `json:"controllers"`
	Currency         string    `json:"currency"`
}

func (it *CreateSecurityTokensItem) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringError("failed to decode json of CreateSecurityTokensItem")

	var uit CreateSecurityTokensItemJSONUnMarshaler
	if err := enc.Unmarshal(b, &uit); err != nil {
		return e.Wrap(err)
	}

	return it.unpack(enc, uit.Hint, uit.Contract, uit.Granularity, uit.DefaultPartition, uit.Controllers, uit.Currency)
}
