package sto

import (
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
)

type BaseCreateContractAccountsItem struct {
	hint.BaseHinter
	design STODesign
}

func NewBaseCreateContractAccountsItem(ht hint.Hint, design STODesign) BaseCreateContractAccountsItem {
	return BaseCreateContractAccountsItem{
		BaseHinter: hint.NewBaseHinter(ht),
		design:     design,
	}
}

func (it BaseCreateContractAccountsItem) Bytes() []byte {
	return it.design.Bytes()
}

func (it BaseCreateContractAccountsItem) IsValid([]byte) error {
	if err := util.CheckIsValiders(nil, false, it.BaseHinter, it.design); err != nil {
		return err
	}

	return nil
}
