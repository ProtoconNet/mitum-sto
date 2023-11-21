package sto

import (
	"github.com/ProtoconNet/mitum-currency/v3/common"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
	"github.com/ProtoconNet/mitum2/util/valuehash"
)

var (
	RedeemFactHint = hint.MustNewHint("mitum-sto-redeem-operation-fact-v0.0.1")
	RedeemHint     = hint.MustNewHint("mitum-sto-redeem-operation-v0.0.1")
)

type RedeemFact struct {
	base.BaseFact
	sender base.Address
	items  []RedeemItem
}

func NewRedeemFact(token []byte, receiver base.Address, items []RedeemItem) RedeemFact {
	bf := base.NewBaseFact(RedeemFactHint, token)
	fact := RedeemFact{
		BaseFact: bf,
		sender:   receiver,
		items:    items,
	}
	fact.SetHash(fact.GenerateHash())

	return fact
}

func (fact RedeemFact) Hash() util.Hash {
	return fact.BaseFact.Hash()
}

func (fact RedeemFact) GenerateHash() util.Hash {
	return valuehash.NewSHA256(fact.Bytes())
}

func (fact RedeemFact) Bytes() []byte {
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

func (fact RedeemFact) IsValid(b []byte) error {
	if err := fact.BaseHinter.IsValid(nil); err != nil {
		return err
	}

	if n := len(fact.items); n < 1 {
		return util.ErrInvalid.Errorf("empty items")
	}

	if err := fact.sender.IsValid(nil); err != nil {
		return err
	}

	founds := map[string]struct{}{}
	for _, it := range fact.items {
		if err := it.IsValid(nil); err != nil {
			return err
		}

		addr := it.tokenHolder

		if _, found := founds[addr.String()]; found {
			return util.ErrInvalid.Errorf("duplicate address found, %s", addr)
		}

		founds[addr.String()] = struct{}{}
	}

	if err := common.IsValidOperationFact(fact, b); err != nil {
		return err
	}

	return nil
}

func (fact RedeemFact) Token() base.Token {
	return fact.BaseFact.Token()
}

func (fact RedeemFact) Sender() base.Address {
	return fact.sender
}

func (fact RedeemFact) Items() []RedeemItem {
	return fact.items
}

func (fact RedeemFact) Addresses() ([]base.Address, error) {
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

type Redeem struct {
	common.BaseOperation
}

func NewRedeem(fact RedeemFact) (Redeem, error) {
	return Redeem{BaseOperation: common.NewBaseOperation(RedeemHint, fact)}, nil
}

func (op *Redeem) HashSign(priv base.Privatekey, networkID base.NetworkID) error {
	err := op.Sign(priv, networkID)
	if err != nil {
		return err
	}

	return nil
}
