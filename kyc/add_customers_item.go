package kyc

import (
	currencybase "github.com/ProtoconNet/mitum-currency/v3/base"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
)

var AddCustomersItemHint = hint.MustNewHint("mitum-kyc-add-customers-item-v0.0.1")

type AddCustomersItem struct {
	hint.BaseHinter
	contract base.Address
	kycID    currencybase.ContractID
	customer base.Address
	status   bool
	currency currencybase.CurrencyID
}

func NewAddCustomersItem(
	contract base.Address,
	kycID currencybase.ContractID,
	customer base.Address,
	status bool,
	currency currencybase.CurrencyID,
) AddCustomersItem {
	return AddCustomersItem{
		BaseHinter: hint.NewBaseHinter(AddCustomersItemHint),
		contract:   contract,
		kycID:      kycID,
		customer:   customer,
		status:     status,
		currency:   currency,
	}
}

func (it AddCustomersItem) Bytes() []byte {
	b := []byte{0}
	if it.status {
		b[0] = 1
	}

	return util.ConcatBytesSlice(
		it.contract.Bytes(),
		it.kycID.Bytes(),
		it.customer.Bytes(),
		b,
		it.currency.Bytes(),
	)
}

func (it AddCustomersItem) IsValid([]byte) error {
	if err := util.CheckIsValiders(nil, false, it.BaseHinter, it.kycID, it.contract, it.customer, it.currency); err != nil {
		return err
	}

	if it.contract.Equal(it.customer) {
		return util.ErrInvalid.Errorf("contract address is same with customer, %q", it.contract)
	}

	return nil
}

func (it AddCustomersItem) KYC() currencybase.ContractID {
	return it.kycID
}

func (it AddCustomersItem) Contract() base.Address {
	return it.contract
}

func (it AddCustomersItem) Customer() base.Address {
	return it.customer
}

func (it AddCustomersItem) Status() bool {
	return it.status
}

func (it AddCustomersItem) Currency() currencybase.CurrencyID {
	return it.currency
}

func (it AddCustomersItem) Addresses() []base.Address {
	ad := make([]base.Address, 2)

	ad[0] = it.contract
	ad[1] = it.customer

	return ad
}
