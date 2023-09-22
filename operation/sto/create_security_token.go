package sto

import (
	"fmt"

	"github.com/ProtoconNet/mitum-currency/v3/common"
	currencytypes "github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
	"github.com/ProtoconNet/mitum2/util/valuehash"
)

type STOItem interface {
	util.Byter
	util.IsValider
	Currency() currencytypes.CurrencyID
}

var (
	CreateSecurityTokenFactHint = hint.MustNewHint("mitum-sto-create-security-token-operation-fact-v0.0.1")
	CreateSecurityTokenHint     = hint.MustNewHint("mitum-sto-create-security-token-operation-v0.0.1")
)

var MaxCreateSecurityTokenItems uint = 10

type CreateSecurityTokenFact struct {
	base.BaseFact
	sender base.Address
	items  []CreateSecurityTokenItem
}

func NewCreateSecurityTokenFact(token []byte, sender base.Address, items []CreateSecurityTokenItem) CreateSecurityTokenFact {
	bf := base.NewBaseFact(CreateSecurityTokenFactHint, token)
	fact := CreateSecurityTokenFact{
		BaseFact: bf,
		sender:   sender,
		items:    items,
	}
	fact.SetHash(fact.GenerateHash())

	return fact
}

func (fact CreateSecurityTokenFact) Hash() util.Hash {
	return fact.BaseFact.Hash()
}

func (fact CreateSecurityTokenFact) GenerateHash() util.Hash {
	return valuehash.NewSHA256(fact.Bytes())
}

func (fact CreateSecurityTokenFact) Bytes() []byte {
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

func (fact CreateSecurityTokenFact) IsValid(b []byte) error {
	if err := fact.BaseHinter.IsValid(nil); err != nil {
		return err
	}

	if err := common.IsValidOperationFact(fact, b); err != nil {
		return err
	}

	if n := len(fact.items); n < 1 {
		return util.ErrInvalid.Errorf("empty items")
	} else if n > int(MaxCreateSecurityTokenItems) {
		return util.ErrInvalid.Errorf("items, %d over max, %d", n, MaxCreateSecurityTokenItems)
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

	return nil
}

func (fact CreateSecurityTokenFact) Token() base.Token {
	return fact.BaseFact.Token()
}

func (fact CreateSecurityTokenFact) Sender() base.Address {
	return fact.sender
}

func (fact CreateSecurityTokenFact) Items() []CreateSecurityTokenItem {
	return fact.items
}

func (fact CreateSecurityTokenFact) Addresses() ([]base.Address, error) {
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

type CreateSecurityToken struct {
	common.BaseOperation
}

func NewCreateSecurityToken(fact CreateSecurityTokenFact) (CreateSecurityToken, error) {
	return CreateSecurityToken{BaseOperation: common.NewBaseOperation(CreateSecurityTokenHint, fact)}, nil
}

func (op *CreateSecurityToken) HashSign(priv base.Privatekey, networkID base.NetworkID) error {
	err := op.Sign(priv, networkID)
	if err != nil {
		return err
	}
	return nil
}
