package kyc

import (
	currencytypes "github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
)

var RemoveControllerItemHint = hint.MustNewHint("mitum-kyc-remove-controller-item-v0.0.1")

type RemoveControllerItem struct {
	hint.BaseHinter
	contract   base.Address
	controller base.Address
	currency   currencytypes.CurrencyID
}

func NewRemoveControllerItem(
	contract base.Address,
	controller base.Address,
	currency currencytypes.CurrencyID,
) RemoveControllerItem {
	return RemoveControllerItem{
		BaseHinter: hint.NewBaseHinter(RemoveControllerItemHint),
		contract:   contract,
		controller: controller,
		currency:   currency,
	}
}

func (it RemoveControllerItem) Bytes() []byte {
	return util.ConcatBytesSlice(
		it.contract.Bytes(),
		it.controller.Bytes(),
		it.currency.Bytes(),
	)
}

func (it RemoveControllerItem) IsValid([]byte) error {
	if err := util.CheckIsValiders(nil, false, it.BaseHinter, it.contract, it.controller, it.currency); err != nil {
		return err
	}

	if it.contract.Equal(it.controller) {
		return util.ErrInvalid.Errorf("contract address is same with controller, %q", it.contract)
	}

	return nil
}

func (it RemoveControllerItem) Contract() base.Address {
	return it.contract
}

func (it RemoveControllerItem) Controller() base.Address {
	return it.controller
}

func (it RemoveControllerItem) Currency() currencytypes.CurrencyID {
	return it.currency
}

func (it RemoveControllerItem) Addresses() []base.Address {
	ad := make([]base.Address, 2)

	ad[0] = it.contract
	ad[1] = it.controller

	return ad
}
