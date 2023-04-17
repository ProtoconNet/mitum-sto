package sto

import (
	"github.com/ProtoconNet/mitum-currency/v2/currency"
	"github.com/ProtoconNet/mitum2/util/hint"
)

var (
	CreateContractAccountsItemSingleAmountHint = hint.MustNewHint("mitum-currency-create-contract-accounts-single-amount-v0.0.1")
)

type CreateContractAccountsItemSingleAmount struct {
	BaseCreateContractAccountsItem
}

func NewCreateContractAccountsItemSingleAmount(keys currency.AccountKeys, design STODesign) CreateContractAccountsItemSingleAmount {
	return CreateContractAccountsItemSingleAmount{
		BaseCreateContractAccountsItem: NewBaseCreateContractAccountsItem(CreateContractAccountsItemSingleAmountHint, design),
	}
}

func (it CreateContractAccountsItemSingleAmount) IsValid([]byte) error {
	if err := it.BaseCreateContractAccountsItem.IsValid(nil); err != nil {
		return err
	}

	return nil
}
