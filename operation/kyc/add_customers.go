package kyc

import (
	"github.com/ProtoconNet/mitum-currency/v3/common"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
	"github.com/ProtoconNet/mitum2/util/valuehash"
)

var (
	AddCustomersFactHint = hint.MustNewHint("mitum-kyc-add-customers-operation-fact-v0.0.1")
	AddCustomersHint     = hint.MustNewHint("mitum-kyc-add-customers-operation-v0.0.1")
)

var MaxAddCustomersItems uint = 10

type AddCustomersFact struct {
	base.BaseFact
	sender base.Address
	items  []AddCustomersItem
}

func NewAddCustomersFact(token []byte, sender base.Address, items []AddCustomersItem) AddCustomersFact {
	bf := base.NewBaseFact(AddCustomersFactHint, token)
	fact := AddCustomersFact{
		BaseFact: bf,
		sender:   sender,
		items:    items,
	}
	fact.SetHash(fact.GenerateHash())

	return fact
}

func (fact AddCustomersFact) Hash() util.Hash {
	return fact.BaseFact.Hash()
}

func (fact AddCustomersFact) GenerateHash() util.Hash {
	return valuehash.NewSHA256(fact.Bytes())
}

func (fact AddCustomersFact) Bytes() []byte {
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

func (fact AddCustomersFact) IsValid(b []byte) error {
	if err := fact.BaseHinter.IsValid(nil); err != nil {
		return err
	}

	if err := common.IsValidOperationFact(fact, b); err != nil {
		return err
	}

	if n := len(fact.items); n < 1 {
		return util.ErrInvalid.Errorf("empty items")
	} else if n > int(MaxAddCustomersItems) {
		return util.ErrInvalid.Errorf("items, %d over max, %d", n, MaxAddCustomersItems)
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

func (fact AddCustomersFact) Token() base.Token {
	return fact.BaseFact.Token()
}

func (fact AddCustomersFact) Sender() base.Address {
	return fact.sender
}

func (fact AddCustomersFact) Items() []AddCustomersItem {
	return fact.items
}

func (fact AddCustomersFact) Addresses() ([]base.Address, error) {
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

type AddCustomers struct {
	common.BaseOperation
}

func NewAddCustomers(fact AddCustomersFact) (AddCustomers, error) {
	return AddCustomers{BaseOperation: common.NewBaseOperation(AddCustomersHint, fact)}, nil
}
