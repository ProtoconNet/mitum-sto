package sto

import (
	"fmt"

	"github.com/ProtoconNet/mitum-currency/v2/currency"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
	"github.com/ProtoconNet/mitum2/util/valuehash"
)

var MaxIssueSecurityTokensItems uint = 10
var (
	IssueSecurityTokensFactHint = hint.MustNewHint("mitum-sto-issue-security-tokens-operation-fact-v0.0.1")
	IssueSecurityTokensHint     = hint.MustNewHint("mitum-sto-issue-security-tokens-operation-v0.0.1")
)

type IssueSecurityTokensFact struct {
	base.BaseFact
	sender base.Address
	items  []IssueSecurityTokensItem
}

func NewIssueSecurityTokensFact(token []byte, sender base.Address, items []IssueSecurityTokensItem) IssueSecurityTokensFact {
	bf := base.NewBaseFact(IssueSecurityTokensFactHint, token)
	fact := IssueSecurityTokensFact{
		BaseFact: bf,
		sender:   sender,
		items:    items,
	}
	fact.SetHash(fact.GenerateHash())

	return fact
}

func (fact IssueSecurityTokensFact) Hash() util.Hash {
	return fact.BaseFact.Hash()
}

func (fact IssueSecurityTokensFact) GenerateHash() util.Hash {
	return valuehash.NewSHA256(fact.Bytes())
}

func (fact IssueSecurityTokensFact) Bytes() []byte {
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

func (fact IssueSecurityTokensFact) IsValid(b []byte) error {
	if err := fact.BaseHinter.IsValid(nil); err != nil {
		return err
	}

	if err := currency.IsValidOperationFact(fact, b); err != nil {
		return err
	}

	if n := len(fact.items); n < 1 {
		return util.ErrInvalid.Errorf("empty items")
	} else if n > int(MaxIssueSecurityTokensItems) {
		return util.ErrInvalid.Errorf("items, %d over max, %d", n, MaxIssueSecurityTokensItems)
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

		k := fmt.Sprintf("%s-%s-%s", it.contract.String(), it.stoID.String(), it.partition.String())

		if _, found := founds[k]; found {
			return util.ErrInvalid.Errorf("duplicated contract-sto-partition found, %s", k)
		}

		founds[k] = struct{}{}
	}

	return nil
}

func (fact IssueSecurityTokensFact) Token() base.Token {
	return fact.BaseFact.Token()
}

func (fact IssueSecurityTokensFact) Sender() base.Address {
	return fact.sender
}

func (fact IssueSecurityTokensFact) Items() []IssueSecurityTokensItem {
	return fact.items
}

func (fact IssueSecurityTokensFact) Addresses() ([]base.Address, error) {
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

type IssueSecurityTokens struct {
	currency.BaseOperation
}

func NewIssueSecurityTokens(fact IssueSecurityTokensFact) (IssueSecurityTokens, error) {
	return IssueSecurityTokens{BaseOperation: currency.NewBaseOperation(IssueSecurityTokensHint, fact)}, nil
}

func (op *IssueSecurityTokens) HashSign(priv base.Privatekey, networkID base.NetworkID) error {
	err := op.Sign(priv, networkID)
	if err != nil {
		return err
	}
	return nil
}
