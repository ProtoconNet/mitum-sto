package kyc

import (
	currencybase "github.com/ProtoconNet/mitum-currency/v3/base"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
	"github.com/ProtoconNet/mitum2/util/valuehash"
)

var (
	UpdateCustomersFactHint = hint.MustNewHint("mitum-kyc-update-customers-operation-fact-v0.0.1")
	UpdateCustomersHint     = hint.MustNewHint("mitum-kyc-update-customers-operation-v0.0.1")
)

var MaxUpdateCustomersItems uint = 10

type UpdateCustomersFact struct {
	base.BaseFact
	sender base.Address
	items  []UpdateCustomersItem
}

func NewUpdateCustomersFact(token []byte, sender base.Address, items []UpdateCustomersItem) UpdateCustomersFact {
	bf := base.NewBaseFact(UpdateCustomersFactHint, token)
	fact := UpdateCustomersFact{
		BaseFact: bf,
		sender:   sender,
		items:    items,
	}
	fact.SetHash(fact.GenerateHash())

	return fact
}

func (fact UpdateCustomersFact) Hash() util.Hash {
	return fact.BaseFact.Hash()
}

func (fact UpdateCustomersFact) GenerateHash() util.Hash {
	return valuehash.NewSHA256(fact.Bytes())
}

func (fact UpdateCustomersFact) Bytes() []byte {
	is := make([][]byte, len(fact.items))
	for i := range fact.items {
		is[i] = fact.items[i].Bytes()
	}

	return util.ConcatBytesSlice(
		fact.Token(),
		fact.sender.Bytes(),
		util.ConcatBytesSlice(is...),
	)
}

func (fact UpdateCustomersFact) IsValid(b []byte) error {
	if err := fact.BaseHinter.IsValid(nil); err != nil {
		return err
	}

	if err := currencybase.IsValidOperationFact(fact, b); err != nil {
		return err
	}

	if n := len(fact.items); n < 1 {
		return util.ErrInvalid.Errorf("empty items")
	} else if n > int(MaxUpdateCustomersItems) {
		return util.ErrInvalid.Errorf("items, %d over max, %d", n, MaxUpdateCustomersItems)
	}

	if err := fact.sender.IsValid(nil); err != nil {
		return err
	}

	founds := map[string]struct{}{}
	for _, it := range fact.items {
		if err := it.IsValid(nil); err != nil {
			return err
		}

		if it.contract.Equal(fact.sender) {
			return util.ErrInvalid.Errorf("contract address is same with sender, %q", fact.sender)
		}

		if _, found := founds[it.Customer().String()]; found {
			return util.ErrInvalid.Errorf("duplicate customer found, %s", it.Customer())
		}

		founds[it.customer.String()] = struct{}{}
	}

	return nil
}

func (fact UpdateCustomersFact) Token() base.Token {
	return fact.BaseFact.Token()
}

func (fact UpdateCustomersFact) Sender() base.Address {
	return fact.sender
}

func (fact UpdateCustomersFact) Items() []UpdateCustomersItem {
	return fact.items
}

func (fact UpdateCustomersFact) Addresses() ([]base.Address, error) {
	as := []base.Address{}

	adrMap := make(map[string]struct{})
	for i := range fact.items {
		for j := range fact.items[i].Addresses() {
			if _, found := adrMap[fact.items[i].Addresses()[j].String()]; !found {
				adrMap[fact.items[i].Addresses()[j].String()] = struct{}{}
				as = append(as, fact.items[i].Addresses()[j])
			}
		}
	}
	as = append(as, fact.sender)

	return as, nil
}

type UpdateCustomers struct {
	currencybase.BaseOperation
}

func NewUpdateCustomers(fact UpdateCustomersFact) (UpdateCustomers, error) {
	return UpdateCustomers{BaseOperation: currencybase.NewBaseOperation(UpdateCustomersHint, fact)}, nil
}

func (op *UpdateCustomers) HashSign(priv base.Privatekey, networkID base.NetworkID) error {
	err := op.Sign(priv, networkID)
	if err != nil {
		return err
	}
	return nil
}
