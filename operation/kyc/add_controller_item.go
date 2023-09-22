package kyc

import (
	currencytypes "github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
)

type KYCItem interface {
	util.Byter
	util.IsValider
	Currency() currencytypes.CurrencyID
}

var AddControllerItemHint = hint.MustNewHint("mitum-kyc-add-controller-item-v0.0.1")

type AddControllerItem struct {
	hint.BaseHinter
	contract   base.Address
	controller base.Address
	currency   currencytypes.CurrencyID
}

func NewAddControllersItem(
	contract base.Address,
	controller base.Address,
	currency currencytypes.CurrencyID,
) AddControllerItem {
	return AddControllerItem{
		BaseHinter: hint.NewBaseHinter(AddControllerItemHint),
		contract:   contract,
		controller: controller,
		currency:   currency,
	}
}

func (it AddControllerItem) Bytes() []byte {
	return util.ConcatBytesSlice(
		it.contract.Bytes(),
		it.controller.Bytes(),
		it.currency.Bytes(),
	)
}

func (it AddControllerItem) IsValid([]byte) error {
	if err := util.CheckIsValiders(nil, false, it.BaseHinter, it.contract, it.controller, it.currency); err != nil {
		return err
	}

	if it.contract.Equal(it.controller) {
		return util.ErrInvalid.Errorf("contract address is same with controller, %q", it.contract)
	}

	return nil
}

func (it AddControllerItem) Contract() base.Address {
	return it.contract
}

func (it AddControllerItem) Controller() base.Address {
	return it.controller
}

func (it AddControllerItem) Currency() currencytypes.CurrencyID {
	return it.currency
}

func (it AddControllerItem) Addresses() []base.Address {
	ad := make([]base.Address, 2)

	ad[0] = it.contract
	ad[1] = it.controller

	return ad
}
