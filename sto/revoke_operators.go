package sto

import (
	"github.com/ProtoconNet/mitum-currency/v2/currency"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
	"github.com/ProtoconNet/mitum2/util/valuehash"
)

var (
	RevokeOperatorsFactHint = hint.MustNewHint("mitum-sto-revoke-operator-operation-fact-v0.0.1")
	RevokeOperatorsHint     = hint.MustNewHint("mitum-sto-revoke-operator-operation-v0.0.1")
)

var MaxRevokeOperatorsItems uint = 10

type RevokeOperatorsItem interface {
	hint.Hinter
	util.IsValider
	Bytes() []byte
	Address() (base.Address, error)
	Rebuild() RevokeOperatorsItem
}

type RevokeOperatorsFact struct {
	base.BaseFact
	sender base.Address
	items  []RevokeOperatorsItem
}

func NewRevokeOperatorsFact(token []byte, sender base.Address, items []RevokeOperatorsItem) RevokeOperatorsFact {
	bf := base.NewBaseFact(RevokeOperatorsFactHint, token)
	fact := RevokeOperatorsFact{
		BaseFact: bf,
		sender:   sender,
		items:    items,
	}
	fact.SetHash(fact.GenerateHash())

	return fact
}

func (fact RevokeOperatorsFact) Hash() util.Hash {
	return fact.BaseFact.Hash()
}

func (fact RevokeOperatorsFact) GenerateHash() util.Hash {
	return valuehash.NewSHA256(fact.Bytes())
}

func (fact RevokeOperatorsFact) Bytes() []byte {
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

func (fact RevokeOperatorsFact) IsValid(b []byte) error {
	if err := fact.BaseHinter.IsValid(nil); err != nil {
		return err
	}

	if err := currency.IsValidOperationFact(fact, b); err != nil {
		return err
	}

	if n := len(fact.items); n < 1 {
		return util.ErrInvalid.Errorf("empty items")
	} else if n > int(MaxRevokeOperatorsItems) {
		return util.ErrInvalid.Errorf("items, %d over max, %d", n, MaxRevokeOperatorsItems)
	}

	if err := util.CheckIsValiders(nil, false, fact.sender); err != nil {
		return err
	}

	for i := range fact.items {
		if err := util.CheckIsValiders(nil, false, fact.items[i]); err != nil {
			return err
		}

		if fact.sender.Equal(fact.sender) {
			return util.ErrInvalid.Errorf("target address is same with sender, %q", fact.sender)
		}
	}

	return nil
}

func (fact RevokeOperatorsFact) Token() base.Token {
	return fact.BaseFact.Token()
}

func (fact RevokeOperatorsFact) Sender() base.Address {
	return fact.sender
}

func (fact RevokeOperatorsFact) Items() []RevokeOperatorsItem {
	return fact.items
}

func (fact RevokeOperatorsFact) Targets() ([]base.Address, error) {
	as := make([]base.Address, len(fact.items))
	for i := range fact.items {
		a, err := fact.items[i].Address()
		if err != nil {
			return nil, err
		}
		as[i] = a
	}

	return as, nil
}

func (fact RevokeOperatorsFact) Addresses() ([]base.Address, error) {
	as := make([]base.Address, len(fact.items)+1)

	tas, err := fact.Targets()
	if err != nil {
		return nil, err
	}
	copy(as, tas)

	as[len(fact.items)] = fact.sender

	return as, nil
}

func (fact RevokeOperatorsFact) Rebuild() RevokeOperatorsFact {
	items := make([]RevokeOperatorsItem, len(fact.items))
	for i := range fact.items {
		it := fact.items[i]
		items[i] = it.Rebuild()
	}

	fact.items = items
	fact.SetHash(fact.GenerateHash())

	return fact
}

type RevokeOperators struct {
	currency.BaseOperation
}

func NewRevokeOperators(fact RevokeOperatorsFact) (RevokeOperators, error) {
	return RevokeOperators{BaseOperation: currency.NewBaseOperation(RevokeOperatorsHint, fact)}, nil
}

func (op *RevokeOperators) HashSign(priv base.Privatekey, networkID base.NetworkID) error {
	err := op.Sign(priv, networkID)
	if err != nil {
		return err
	}
	return nil
}
