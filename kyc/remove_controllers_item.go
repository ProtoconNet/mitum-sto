package kyc

import (
	currencybase "github.com/ProtoconNet/mitum-currency/v3/base"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
)

var RemoveControllersItemHint = hint.MustNewHint("mitum-kyc-remove-controllers-item-v0.0.1")

type RemoveControllersItem struct {
	hint.BaseHinter
	contract   base.Address
	kycID      currencybase.ContractID
	controller base.Address
	currency   currencybase.CurrencyID
}

func NewRemoveControllersItem(
	contract base.Address,
	kycID currencybase.ContractID,
	controller base.Address,
	currency currencybase.CurrencyID,
) RemoveControllersItem {
	return RemoveControllersItem{
		BaseHinter: hint.NewBaseHinter(RemoveControllersItemHint),
		contract:   contract,
		kycID:      kycID,
		controller: controller,
		currency:   currency,
	}
}

func (it RemoveControllersItem) Bytes() []byte {
	return util.ConcatBytesSlice(
		it.contract.Bytes(),
		it.kycID.Bytes(),
		it.controller.Bytes(),
		it.currency.Bytes(),
	)
}

func (it RemoveControllersItem) IsValid([]byte) error {
	if err := util.CheckIsValiders(nil, false, it.BaseHinter, it.kycID, it.contract, it.controller, it.currency); err != nil {
		return err
	}

	if it.contract.Equal(it.controller) {
		return util.ErrInvalid.Errorf("contract address is same with controller, %q", it.contract)
	}

	return nil
}

func (it RemoveControllersItem) KYC() currencybase.ContractID {
	return it.kycID
}

func (it RemoveControllersItem) Contract() base.Address {
	return it.contract
}

func (it RemoveControllersItem) Controller() base.Address {
	return it.controller
}

func (it RemoveControllersItem) Currency() currencybase.CurrencyID {
	return it.currency
}

func (it RemoveControllersItem) Addresses() []base.Address {
	ad := make([]base.Address, 2)

	ad[0] = it.contract
	ad[1] = it.controller

	return ad
}
