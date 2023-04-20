package sto

import (
	extensioncurrency "github.com/ProtoconNet/mitum-currency-extension/v2/currency"
	"github.com/ProtoconNet/mitum-currency/v2/currency"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
)

var AuthorizeOperatorsItemHint = hint.MustNewHint("mitum-sto-authorize-operators-item-v0.0.1")

type AuthorizeOperatorsItem struct {
	hint.BaseHinter
	stoID    extensioncurrency.ContractID // token id
	contract base.Address
	operator base.Address        // initial controllers
	currency currency.CurrencyID // fee
}

func NewBaseAuthorizeOperatorsItem(stoID extensioncurrency.ContractID, contract, operator base.Address, currency currency.CurrencyID) AuthorizeOperatorsItem {
	return AuthorizeOperatorsItem{
		BaseHinter: hint.NewBaseHinter(AuthorizeOperatorsItemHint),
		stoID:      stoID,
		contract:   contract,
		operator:   operator,
		currency:   currency,
	}
}

func (it AuthorizeOperatorsItem) Bytes() []byte {
	return util.ConcatBytesSlice(
		it.stoID.Bytes(),
		it.contract.Bytes(),
		it.operator.Bytes(),
		it.currency.Bytes(),
	)
}

func (it AuthorizeOperatorsItem) IsValid([]byte) error {
	if err := util.CheckIsValiders(nil, false, it.BaseHinter, it.stoID, it.contract, it.operator, it.currency); err != nil {
		return err
	}

	if it.contract.Equal(it.operator) {
		return util.ErrInvalid.Errorf("contract and operation address are same, %q == %q", it.contract, it.operator)
	}

	return nil
}

func (it AuthorizeOperatorsItem) STO() extensioncurrency.ContractID {
	return it.stoID
}

func (it AuthorizeOperatorsItem) Contract() base.Address {
	return it.contract
}

func (it AuthorizeOperatorsItem) Operator() base.Address {
	return it.operator
}

func (it AuthorizeOperatorsItem) Addresses() []base.Address {
	ad := make([]base.Address, 2)

	ad[0] = it.contract
	ad[1] = it.operator

	return ad
}

func (it AuthorizeOperatorsItem) Rebuild() AuthorizeOperatorsItem {
	return it
}
