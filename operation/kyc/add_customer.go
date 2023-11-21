package kyc

import (
	"github.com/ProtoconNet/mitum-currency/v3/common"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
	"github.com/ProtoconNet/mitum2/util/valuehash"
)

var (
	AddCustomerFactHint = hint.MustNewHint("mitum-kyc-add-customer-operation-fact-v0.0.1")
	AddCustomerHint     = hint.MustNewHint("mitum-kyc-add-customer-operation-v0.0.1")
)

var MaxAddCustomerItems uint = 10

type AddCustomerFact struct {
	base.BaseFact
	sender base.Address
	items  []AddCustomerItem
}

func NewAddCustomerFact(token []byte, sender base.Address, items []AddCustomerItem) AddCustomerFact {
	bf := base.NewBaseFact(AddCustomerFactHint, token)
	fact := AddCustomerFact{
		BaseFact: bf,
		sender:   sender,
		items:    items,
	}
	fact.SetHash(fact.GenerateHash())

	return fact
}

func (fact AddCustomerFact) Hash() util.Hash {
	return fact.BaseFact.Hash()
}

func (fact AddCustomerFact) GenerateHash() util.Hash {
	return valuehash.NewSHA256(fact.Bytes())
}

func (fact AddCustomerFact) Bytes() []byte {
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

func (fact AddCustomerFact) IsValid(b []byte) error {
	if err := fact.BaseHinter.IsValid(nil); err != nil {
		return err
	}

	if n := len(fact.items); n < 1 {
		return util.ErrInvalid.Errorf("empty items")
	} else if n > int(MaxAddCustomerItems) {
		return util.ErrInvalid.Errorf("items, %d over max, %d", n, MaxAddCustomerItems)
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

	if err := common.IsValidOperationFact(fact, b); err != nil {
		return err
	}

	return nil
}

func (fact AddCustomerFact) Token() base.Token {
	return fact.BaseFact.Token()
}

func (fact AddCustomerFact) Sender() base.Address {
	return fact.sender
}

func (fact AddCustomerFact) Items() []AddCustomerItem {
	return fact.items
}

func (fact AddCustomerFact) Addresses() ([]base.Address, error) {
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

type AddCustomer struct {
	common.BaseOperation
}

func NewAddCustomer(fact AddCustomerFact) (AddCustomer, error) {
	return AddCustomer{BaseOperation: common.NewBaseOperation(AddCustomerHint, fact)}, nil
}
