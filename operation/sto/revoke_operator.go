package sto

import (
	"github.com/ProtoconNet/mitum-currency/v3/common"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
	"github.com/ProtoconNet/mitum2/util/valuehash"
)

var (
	RevokeOperatorFactHint = hint.MustNewHint("mitum-sto-revoke-operator-operation-fact-v0.0.1")
	RevokeOperatorHint     = hint.MustNewHint("mitum-sto-revoke-operator-operation-v0.0.1")
)

var MaxRevokeOperatorItems uint = 10

type RevokeOperatorFact struct {
	base.BaseFact
	sender base.Address
	items  []RevokeOperatorItem
}

func NewRevokeOperatorFact(token []byte, sender base.Address, items []RevokeOperatorItem) RevokeOperatorFact {
	bf := base.NewBaseFact(RevokeOperatorFactHint, token)
	fact := RevokeOperatorFact{
		BaseFact: bf,
		sender:   sender,
		items:    items,
	}
	fact.SetHash(fact.GenerateHash())

	return fact
}

func (fact RevokeOperatorFact) Hash() util.Hash {
	return fact.BaseFact.Hash()
}

func (fact RevokeOperatorFact) GenerateHash() util.Hash {
	return valuehash.NewSHA256(fact.Bytes())
}

func (fact RevokeOperatorFact) Bytes() []byte {
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

func (fact RevokeOperatorFact) IsValid(b []byte) error {
	if err := fact.BaseHinter.IsValid(nil); err != nil {
		return err
	}

	if n := len(fact.items); n < 1 {
		return util.ErrInvalid.Errorf("empty items")
	} else if n > int(MaxRevokeOperatorItems) {
		return util.ErrInvalid.Errorf("items, %d over max, %d", n, MaxRevokeOperatorItems)
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

		if _, found := founds[it.Operator().String()]; found {
			return util.ErrInvalid.Errorf("duplicate operator found, %s", it.Operator())
		}

		founds[it.operator.String()] = struct{}{}
	}

	if err := common.IsValidOperationFact(fact, b); err != nil {
		return err
	}

	return nil
}

func (fact RevokeOperatorFact) Token() base.Token {
	return fact.BaseFact.Token()
}

func (fact RevokeOperatorFact) Sender() base.Address {
	return fact.sender
}

func (fact RevokeOperatorFact) Items() []RevokeOperatorItem {
	return fact.items
}

func (fact RevokeOperatorFact) Addresses() ([]base.Address, error) {
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

type RevokeOperator struct {
	common.BaseOperation
}

func NewRevokeOperator(fact RevokeOperatorFact) (RevokeOperator, error) {
	return RevokeOperator{BaseOperation: common.NewBaseOperation(RevokeOperatorHint, fact)}, nil
}

func (op *RevokeOperator) HashSign(priv base.Privatekey, networkID base.NetworkID) error {
	err := op.Sign(priv, networkID)
	if err != nil {
		return err
	}
	return nil
}
