package kyc

import (
	extensioncurrency "github.com/ProtoconNet/mitum-currency-extension/v2/currency"
	"github.com/ProtoconNet/mitum-currency/v2/currency"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
)

type KYCItem interface {
	util.Byter
	util.IsValider
	Currency() currency.CurrencyID
}

var AddControllersItemHint = hint.MustNewHint("mitum-kyc-add-controllers-item-v0.0.1")

type AddControllersItem struct {
	hint.BaseHinter
	contract   base.Address
	kycID      extensioncurrency.ContractID
	controller base.Address
	currency   currency.CurrencyID
}

func NewAddControllersItem(
	contract base.Address,
	kycID extensioncurrency.ContractID,
	controller base.Address,
	currency currency.CurrencyID,
) AddControllersItem {
	return AddControllersItem{
		BaseHinter: hint.NewBaseHinter(AddControllersItemHint),
		contract:   contract,
		kycID:      kycID,
		controller: controller,
		currency:   currency,
	}
}

func (it AddControllersItem) Bytes() []byte {
	return util.ConcatBytesSlice(
		it.contract.Bytes(),
		it.kycID.Bytes(),
		it.controller.Bytes(),
		it.currency.Bytes(),
	)
}

func (it AddControllersItem) IsValid([]byte) error {
	if err := util.CheckIsValiders(nil, false, it.BaseHinter, it.kycID, it.contract, it.controller, it.currency); err != nil {
		return err
	}

	if it.contract.Equal(it.controller) {
		return util.ErrInvalid.Errorf("contract address is same with controller, %q", it.contract)
	}

	return nil
}

func (it AddControllersItem) KYC() extensioncurrency.ContractID {
	return it.kycID
}

func (it AddControllersItem) Contract() base.Address {
	return it.contract
}

func (it AddControllersItem) Controller() base.Address {
	return it.controller
}

func (it AddControllersItem) Currency() currency.CurrencyID {
	return it.currency
}

func (it AddControllersItem) Addresses() []base.Address {
	ad := make([]base.Address, 2)

	ad[0] = it.contract
	ad[1] = it.controller

	return ad
}
