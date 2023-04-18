package sto

import (
	currencyextension "github.com/ProtoconNet/mitum-currency-extension/v2/currency"
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

type AuthorizeOperatorsItem interface {
	hint.Hinter
	util.IsValider
	Bytes() []byte
	STO() currencyextension.ContractID
	Operator() base.Address
	Addresses() []base.Address
}

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

	if err := util.CheckIsValiders(nil, false, fact.sender); err != nil {
		return err
	}

	foundAddrs := map[string]struct{}{}
	for i := range fact.items {
		if err := util.CheckIsValiders(nil, false, fact.items[i]); err != nil {
			return err
		}

		it := fact.items[i]
		addr := it.Operator()

		if _, found := foundAddrs[addr.String()]; found {
			return util.ErrInvalid.Errorf("duplicate address found, %s", addr)
		}

		foundAddrs[addr.String()] = struct{}{}
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
