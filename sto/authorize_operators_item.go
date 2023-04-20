package sto

import (
	extensioncurrency "github.com/ProtoconNet/mitum-currency-extension/v2/currency"
	"github.com/ProtoconNet/mitum-currency/v2/currency"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
)

type BaseAuthorizeOperatorsItem struct {
	hint.BaseHinter
	stoID    extensioncurrency.ContractID // token id
	contract base.Address
	operator base.Address        // initial controllers
	currency currency.CurrencyID // fee
}

func NewBaseAuthorizeOperatorsItem(ht hint.Hint, stoID extensioncurrency.ContractID, contract, operator base.Address, currency currency.CurrencyID) BaseAuthorizeOperatorsItem {
	return BaseAuthorizeOperatorsItem{
		BaseHinter: hint.NewBaseHinter(ht),
		stoID:      stoID,
		contract:   contract,
		operator:   operator,
		currency:   currency,
	}
}

func (it BaseAuthorizeOperatorsItem) Bytes() []byte {
	return util.ConcatBytesSlice(
		it.stoID.Bytes(),
		it.contract.Bytes(),
		it.operator.Bytes(),
		it.currency.Bytes(),
	)
}

func (it BaseAuthorizeOperatorsItem) IsValid([]byte) error {
	if err := util.CheckIsValiders(nil, false, it.BaseHinter, it.stoID, it.contract, it.operator, it.currency); err != nil {
		return err
	}

	if it.contract.Equal(it.operator) {
		return util.ErrInvalid.Errorf("contract and operation address are same, %q == %q", it.contract, it.operator)
	}

	return nil
}

func (it BaseAuthorizeOperatorsItem) STO() extensioncurrency.ContractID {
	return it.stoID
}

func (it BaseAuthorizeOperatorsItem) Contract() base.Address {
	return it.contract
}

func (it BaseAuthorizeOperatorsItem) Operator() base.Address {
	return it.operator
}

func (it BaseAuthorizeOperatorsItem) Addresses() []base.Address {
	ad := make([]base.Address, 2)

	ad[0] = it.contract
	ad[1] = it.operator

	return ad
}

func (it BaseAuthorizeOperatorsItem) Rebuild() AuthorizeOperatorsItem {
	return it
}
