package sto

import (
	"fmt"

	"github.com/ProtoconNet/mitum-currency/v3/common"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
	"github.com/ProtoconNet/mitum2/util/valuehash"
)

var MaxIssueItems uint = 10
var (
	IssueFactHint = hint.MustNewHint("mitum-sto-issue-operation-fact-v0.0.1")
	IssueHint     = hint.MustNewHint("mitum-sto-issue-operation-v0.0.1")
)

type IssueFact struct {
	base.BaseFact
	sender base.Address
	items  []IssueItem
}

func NewIssueFact(token []byte, sender base.Address, items []IssueItem) IssueFact {
	bf := base.NewBaseFact(IssueFactHint, token)
	fact := IssueFact{
		BaseFact: bf,
		sender:   sender,
		items:    items,
	}
	fact.SetHash(fact.GenerateHash())

	return fact
}

func (fact IssueFact) Hash() util.Hash {
	return fact.BaseFact.Hash()
}

func (fact IssueFact) GenerateHash() util.Hash {
	return valuehash.NewSHA256(fact.Bytes())
}

func (fact IssueFact) Bytes() []byte {
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

func (fact IssueFact) IsValid(b []byte) error {
	if err := fact.BaseHinter.IsValid(nil); err != nil {
		return err
	}

	if n := len(fact.items); n < 1 {
		return util.ErrInvalid.Errorf("empty items")
	} else if n > int(MaxIssueItems) {
		return util.ErrInvalid.Errorf("items, %d over max, %d", n, MaxIssueItems)
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

		k := fmt.Sprintf("%s", it.contract.String())

		if _, found := founds[k]; found {
			return util.ErrInvalid.Errorf("duplicated contract-sto found, %s", k)
		}

		founds[k] = struct{}{}
	}

	if err := common.IsValidOperationFact(fact, b); err != nil {
		return err
	}

	return nil
}

func (fact IssueFact) Token() base.Token {
	return fact.BaseFact.Token()
}

func (fact IssueFact) Sender() base.Address {
	return fact.sender
}

func (fact IssueFact) Items() []IssueItem {
	return fact.items
}

func (fact IssueFact) Addresses() ([]base.Address, error) {
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

type Issue struct {
	common.BaseOperation
}

func NewIssue(fact IssueFact) (Issue, error) {
	return Issue{BaseOperation: common.NewBaseOperation(IssueHint, fact)}, nil
}

func (op *Issue) HashSign(priv base.Privatekey, networkID base.NetworkID) error {
	err := op.Sign(priv, networkID)
	if err != nil {
		return err
	}
	return nil
}
