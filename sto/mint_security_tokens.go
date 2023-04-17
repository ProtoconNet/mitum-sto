package sto

import (
	"github.com/ProtoconNet/mitum-currency/v2/currency"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
	"github.com/ProtoconNet/mitum2/util/valuehash"
)

var MaxMintSecurityTokensItems uint = 10
var (
	MintSecurityTokensFactHint = hint.MustNewHint("mitum-sto-mint-security-tokens-operation-fact-v0.0.1")
	MintSecurityTokensHint     = hint.MustNewHint("mitum-sto-mint-security-tokenss-operation-v0.0.1")
)

type MintSecurityTokensItem interface {
	hint.Hinter
	util.IsValider
	Bytes() []byte
	Address() (base.Address, error)
	Rebuild() CreateSecurityTokensItem
}

type MintSecurityTokensFact struct {
	base.BaseFact
	sender base.Address
	token  base.Token
	items  []MintSecurityTokensItem
}

func NewMintSecurityTokensFact(token []byte, sender, target base.Address, items []MintSecurityTokensItem) MintSecurityTokensFact {
	bf := base.NewBaseFact(MintSecurityTokensFactHint, token)
	fact := MintSecurityTokensFact{
		BaseFact: bf,
		sender:   sender,
		items:    items,
	}
	fact.SetHash(fact.GenerateHash())

	return fact
}

func (fact MintSecurityTokensFact) GenerateHash() util.Hash {
	return valuehash.NewSHA256(fact.Bytes())
}

func (fact MintSecurityTokensFact) Bytes() []byte {
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

func (fact MintSecurityTokensFact) IsValid(b []byte) error {
	if err := fact.BaseHinter.IsValid(nil); err != nil {
		return err
	}

	if err := currency.IsValidOperationFact(fact, b); err != nil {
		return err
	}

	if n := len(fact.items); n < 1 {
		return util.ErrInvalid.Errorf("empty items")
	} else if n > int(MaxMintSecurityTokensItems) {
		return util.ErrInvalid.Errorf("items, %d over max, %d", n, MaxMintSecurityTokensItems)
	}

	if err := util.CheckIsValiders(nil, false, fact.sender); err != nil {
		return err
	}

	for i := range fact.items {
		if err := util.CheckIsValiders(nil, false, fact.items[i]); err != nil {
			return err
		}
	}

	return nil
}

type MintSecurityTokens struct {
	currency.BaseOperation
}

func NewMintSecurityTokens(fact MintSecurityTokensFact) (MintSecurityTokens, error) {
	return MintSecurityTokens{BaseOperation: currency.NewBaseOperation(MintSecurityTokensHint, fact)}, nil
}

func (op *MintSecurityTokens) HashSign(priv base.Privatekey, networkID base.NetworkID) error {
	err := op.Sign(priv, networkID)
	if err != nil {
		return err
	}
	return nil
}
