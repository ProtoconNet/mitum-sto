package kyc

import (
	currencytypes "github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
)

var AddCustomerItemHint = hint.MustNewHint("mitum-kyc-add-customer-item-v0.0.1")

type AddCustomerItem struct {
	hint.BaseHinter
	contract base.Address
	customer base.Address
	status   bool
	currency currencytypes.CurrencyID
}

func NewAddCustomerItem(
	contract base.Address,
	customer base.Address,
	status bool,
	currency currencytypes.CurrencyID,
) AddCustomerItem {
	return AddCustomerItem{
		BaseHinter: hint.NewBaseHinter(AddCustomerItemHint),
		contract:   contract,
		customer:   customer,
		status:     status,
		currency:   currency,
	}
}

func (it AddCustomerItem) Bytes() []byte {
	b := []byte{0}
	if it.status {
		b[0] = 1
	}

	return util.ConcatBytesSlice(
		it.contract.Bytes(),
		it.customer.Bytes(),
		b,
		it.currency.Bytes(),
	)
}

func (it AddCustomerItem) IsValid([]byte) error {
	if err := util.CheckIsValiders(nil, false, it.BaseHinter, it.contract, it.customer, it.currency); err != nil {
		return err
	}

	if it.contract.Equal(it.customer) {
		return util.ErrInvalid.Errorf("contract address is same with customer, %q", it.contract)
	}

	return nil
}

func (it AddCustomerItem) Contract() base.Address {
	return it.contract
}

func (it AddCustomerItem) Customer() base.Address {
	return it.customer
}

func (it AddCustomerItem) Status() bool {
	return it.status
}

func (it AddCustomerItem) Currency() currencytypes.CurrencyID {
	return it.currency
}

func (it AddCustomerItem) Addresses() []base.Address {
	ad := make([]base.Address, 2)

	ad[0] = it.contract
	ad[1] = it.customer

	return ad
}
