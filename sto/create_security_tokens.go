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
	CreateSecurityTokensFactHint = hint.MustNewHint("mitum-sto-create-security-tokens-operation-fact-v0.0.1")
	CreateSecurityTokensHint     = hint.MustNewHint("mitum-sto-create-security-tokenss-operation-v0.0.1")
)

var MaxCreateSecurityTokensItems uint = 10

type CreateSecurityTokensItem interface {
	hint.Hinter
	util.IsValider
	Bytes() []byte
	STO() currencyextension.ContractID
	Granularity() uint64
	DefaultPartitions() []Partition
	Controllers() []base.Address
}

type CreateSecurityTokensFact struct {
	base.BaseFact
	sender base.Address
	items  []CreateSecurityTokensItem
}

func NewCreateTokenAccountsFact(token []byte, sender base.Address, items []CreateSecurityTokensItem) CreateSecurityTokensFact {
	bf := base.NewBaseFact(CreateSecurityTokensFactHint, token)
	fact := CreateSecurityTokensFact{
		BaseFact: bf,
		sender:   sender,
		items:    items,
	}
	fact.SetHash(fact.GenerateHash())

	return fact
}

func (fact CreateSecurityTokensFact) Hash() util.Hash {
	return fact.BaseFact.Hash()
}

func (fact CreateSecurityTokensFact) GenerateHash() util.Hash {
	return valuehash.NewSHA256(fact.Bytes())
}

func (fact CreateSecurityTokensFact) Bytes() []byte {
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

func (fact CreateSecurityTokensFact) IsValid(b []byte) error {
	if err := fact.BaseHinter.IsValid(nil); err != nil {
		return err
	}

	if err := currency.IsValidOperationFact(fact, b); err != nil {
		return err
	}

	if n := len(fact.items); n < 1 {
		return util.ErrInvalid.Errorf("empty items")
	} else if n > int(MaxCreateSecurityTokensItems) {
		return util.ErrInvalid.Errorf("items, %d over max, %d", n, MaxCreateSecurityTokensItems)
	}

	if err := util.CheckIsValiders(nil, false, fact.sender); err != nil {
		return err
	}

	return nil
}

func (fact CreateSecurityTokensFact) Token() base.Token {
	return fact.BaseFact.Token()
}

func (fact CreateSecurityTokensFact) Sender() base.Address {
	return fact.sender
}

func (fact CreateSecurityTokensFact) Items() []CreateSecurityTokensItem {
	return fact.items
}

func (fact CreateSecurityTokensFact) Addresses() ([]base.Address, error) {
	as := make([]base.Address, 1)
	as[0] = fact.sender

	return as, nil
}

type CreateSecurityTokens struct {
	currency.BaseOperation
}

func NewCreateSecurityTokens(fact CreateSecurityTokensFact) (CreateSecurityTokens, error) {
	return CreateSecurityTokens{BaseOperation: currency.NewBaseOperation(CreateSecurityTokensHint, fact)}, nil
}

func (op *CreateSecurityTokens) HashSign(priv base.Privatekey, networkID base.NetworkID) error {
	err := op.Sign(priv, networkID)
	if err != nil {
		return err
	}
	return nil
}
