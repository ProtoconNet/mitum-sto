package sto

import (
	"github.com/ProtoconNet/mitum-currency/v2/currency"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
	"github.com/ProtoconNet/mitum2/util/valuehash"
)

var (
	SetDocumentsFactHint = hint.MustNewHint("mitum-sto-set-document-operation-fact-v0.0.1")
	SetDocumentsHint     = hint.MustNewHint("mitum-sto-set-document-operation-v0.0.1")
)

var MaxSetDocumentsItems uint = 10

type SetDocumentsItem interface {
	hint.Hinter
	util.IsValider
	Bytes() []byte
	Address() (base.Address, error)
	Rebuild() SetDocumentsItem
}

type SetDocumentsFact struct {
	base.BaseFact
	sender base.Address
	items  []SetDocumentsItem
}

func NewSetDocumentsFact(token []byte, sender base.Address, items []SetDocumentsItem) SetDocumentsFact {
	bf := base.NewBaseFact(SetDocumentsFactHint, token)
	fact := SetDocumentsFact{
		BaseFact: bf,
		sender:   sender,
		items:    items,
	}
	fact.SetHash(fact.GenerateHash())

	return fact
}

func (fact SetDocumentsFact) Hash() util.Hash {
	return fact.BaseFact.Hash()
}

func (fact SetDocumentsFact) GenerateHash() util.Hash {
	return valuehash.NewSHA256(fact.Bytes())
}

func (fact SetDocumentsFact) Bytes() []byte {
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

func (fact SetDocumentsFact) IsValid(b []byte) error {
	if err := fact.BaseHinter.IsValid(nil); err != nil {
		return err
	}

	if err := currency.IsValidOperationFact(fact, b); err != nil {
		return err
	}

	if n := len(fact.items); n < 1 {
		return util.ErrInvalid.Errorf("empty items")
	} else if n > int(MaxSetDocumentsItems) {
		return util.ErrInvalid.Errorf("items, %d over max, %d", n, MaxSetDocumentsItems)
	}

	if err := util.CheckIsValiders(nil, false, fact.sender); err != nil {
		return err
	}

	return nil
}

func (fact SetDocumentsFact) Token() base.Token {
	return fact.BaseFact.Token()
}

func (fact SetDocumentsFact) Sender() base.Address {
	return fact.sender
}

func (fact SetDocumentsFact) Items() []SetDocumentsItem {
	return fact.items
}

func (fact SetDocumentsFact) Targets() ([]base.Address, error) {
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

func (fact SetDocumentsFact) Addresses() ([]base.Address, error) {
	as := make([]base.Address, len(fact.items)+1)

	tas, err := fact.Targets()
	if err != nil {
		return nil, err
	}
	copy(as, tas)

	as[len(fact.items)] = fact.sender

	return as, nil
}

func (fact SetDocumentsFact) Rebuild() SetDocumentsFact {
	items := make([]SetDocumentsItem, len(fact.items))
	for i := range fact.items {
		it := fact.items[i]
		items[i] = it.Rebuild()
	}

	fact.items = items
	fact.SetHash(fact.GenerateHash())

	return fact
}

type SetDocuments struct {
	currency.BaseOperation
}

func NewSetDocuments(fact SetDocumentsFact) (SetDocuments, error) {
	return SetDocuments{BaseOperation: currency.NewBaseOperation(SetDocumentsHint, fact)}, nil
}

func (op *SetDocuments) HashSign(priv base.Privatekey, networkID base.NetworkID) error {
	err := op.Sign(priv, networkID)
	if err != nil {
		return err
	}
	return nil
}
