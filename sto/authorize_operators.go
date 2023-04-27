package sto

import (
	"fmt"

	"github.com/ProtoconNet/mitum-currency/v2/currency"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
	"github.com/ProtoconNet/mitum2/util/valuehash"
)

var (
	AuthorizeOperatorsFactHint = hint.MustNewHint("mitum-sto-authorize-operator-operation-fact-v0.0.1")
	AuthorizeOperatorsHint     = hint.MustNewHint("mitum-sto-authorize-operator-operation-v0.0.1")
)

var MaxAuthorizeOperatorsItems uint = 10

type AuthorizeOperatorsFact struct {
	base.BaseFact
	sender base.Address
	items  []AuthorizeOperatorsItem
}

func NewAuthorizeOperatorsFact(token []byte, sender base.Address, items []AuthorizeOperatorsItem) AuthorizeOperatorsFact {
	bf := base.NewBaseFact(AuthorizeOperatorsFactHint, token)
	fact := AuthorizeOperatorsFact{
		BaseFact: bf,
		sender:   sender,
		items:    items,
	}
	fact.SetHash(fact.GenerateHash())

	return fact
}

func (fact AuthorizeOperatorsFact) Hash() util.Hash {
	return fact.BaseFact.Hash()
}

func (fact AuthorizeOperatorsFact) GenerateHash() util.Hash {
	return valuehash.NewSHA256(fact.Bytes())
}

func (fact AuthorizeOperatorsFact) Bytes() []byte {
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

func (fact AuthorizeOperatorsFact) IsValid(b []byte) error {
	if err := fact.BaseHinter.IsValid(nil); err != nil {
		return err
	}

	if err := currency.IsValidOperationFact(fact, b); err != nil {
		return err
	}

	if n := len(fact.items); n < 1 {
		return util.ErrInvalid.Errorf("empty items")
	} else if n > int(MaxAuthorizeOperatorsItems) {
		return util.ErrInvalid.Errorf("items, %d over max, %d", n, MaxAuthorizeOperatorsItems)
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

		k := fmt.Sprintf("%s-%s-%s", it.contract, it.stoID, it.partition)
		if _, found := founds[k]; found {
			return util.ErrInvalid.Errorf("duplicate contract-sto-partition found, %s", k)
		}

		founds[k] = struct{}{}
	}

	return nil
}

func (fact AuthorizeOperatorsFact) Token() base.Token {
	return fact.BaseFact.Token()
}

func (fact AuthorizeOperatorsFact) Sender() base.Address {
	return fact.sender
}

func (fact AuthorizeOperatorsFact) Items() []AuthorizeOperatorsItem {
	return fact.items
}

func (fact AuthorizeOperatorsFact) Addresses() ([]base.Address, error) {
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

type AuthorizeOperators struct {
	currency.BaseOperation
}

func NewAuthorizeOperators(fact AuthorizeOperatorsFact) (AuthorizeOperators, error) {
	return AuthorizeOperators{BaseOperation: currency.NewBaseOperation(AuthorizeOperatorsHint, fact)}, nil
}

func (op *AuthorizeOperators) HashSign(priv base.Privatekey, networkID base.NetworkID) error {
	err := op.Sign(priv, networkID)
	if err != nil {
		return err
	}
	return nil
}
