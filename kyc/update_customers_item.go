package kyc

import (
	currencybase "github.com/ProtoconNet/mitum-currency/v3/base"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
)

var UpdateCustomersItemHint = hint.MustNewHint("mitum-kyc-update-customers-item-v0.0.1")

type UpdateCustomersItem struct {
	hint.BaseHinter
	contract base.Address
	kycID    currencybase.ContractID
	customer base.Address
	status   bool
	currency currencybase.CurrencyID
}

func NewUpdateCustomersItem(
	contract base.Address,
	kycID currencybase.ContractID,
	customer base.Address,
	status bool,
	currency currencybase.CurrencyID,
) UpdateCustomersItem {
	return UpdateCustomersItem{
		BaseHinter: hint.NewBaseHinter(UpdateCustomersItemHint),
		contract:   contract,
		kycID:      kycID,
		customer:   customer,
		status:     status,
		currency:   currency,
	}
}

func (it UpdateCustomersItem) Bytes() []byte {
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

func (it UpdateCustomersItem) IsValid([]byte) error {
	if err := util.CheckIsValiders(nil, false, it.BaseHinter, it.kycID, it.contract, it.customer, it.currency); err != nil {
		return err
	}

	if it.contract.Equal(it.customer) {
		return util.ErrInvalid.Errorf("contract address is same with customer, %q", it.contract)
	}

	return nil
}

func (it UpdateCustomersItem) KYC() currencybase.ContractID {
	return it.kycID
}

func (it UpdateCustomersItem) Contract() base.Address {
	return it.contract
}

func (it UpdateCustomersItem) Customer() base.Address {
	return it.customer
}

func (it UpdateCustomersItem) Status() bool {
	return it.status
}

func (it UpdateCustomersItem) Currency() currencybase.CurrencyID {
	return it.currency
}

func (it UpdateCustomersItem) Addresses() []base.Address {
	ad := make([]base.Address, 2)

	ad[0] = it.contract
	ad[1] = it.customer

	return ad
}
