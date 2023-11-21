package kyc

import (
	"github.com/ProtoconNet/mitum-currency/v3/common"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
	"github.com/ProtoconNet/mitum2/util/valuehash"
)

var (
	RemoveControllerFactHint = hint.MustNewHint("mitum-kyc-remove-controller-operation-fact-v0.0.1")
	RemoveControllerHint     = hint.MustNewHint("mitum-kyc-remove-controller-operation-v0.0.1")
)

var MaxRemoveControllerItems uint = 10

type RemoveControllerFact struct {
	base.BaseFact
	sender base.Address
	items  []RemoveControllerItem
}

func NewRemoveControllerFact(token []byte, sender base.Address, items []RemoveControllerItem) RemoveControllerFact {
	bf := base.NewBaseFact(RemoveControllerFactHint, token)
	fact := RemoveControllerFact{
		BaseFact: bf,
		sender:   sender,
		items:    items,
	}
	fact.SetHash(fact.GenerateHash())

	return fact
}

func (fact RemoveControllerFact) Hash() util.Hash {
	return fact.BaseFact.Hash()
}

func (fact RemoveControllerFact) GenerateHash() util.Hash {
	return valuehash.NewSHA256(fact.Bytes())
}

func (fact RemoveControllerFact) Bytes() []byte {
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

func (fact RemoveControllerFact) IsValid(b []byte) error {
	if err := fact.BaseHinter.IsValid(nil); err != nil {
		return err
	}

	if n := len(fact.items); n < 1 {
		return util.ErrInvalid.Errorf("empty items")
	} else if n > int(MaxRemoveControllerItems) {
		return util.ErrInvalid.Errorf("items, %d over max, %d", n, MaxRemoveControllerItems)
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

		if _, found := founds[it.Controller().String()]; found {
			return util.ErrInvalid.Errorf("duplicate controller found, %s", it.Controller())
		}

		founds[it.controller.String()] = struct{}{}
	}

	if err := common.IsValidOperationFact(fact, b); err != nil {
		return err
	}

	return nil
}

func (fact RemoveControllerFact) Token() base.Token {
	return fact.BaseFact.Token()
}

func (fact RemoveControllerFact) Sender() base.Address {
	return fact.sender
}

func (fact RemoveControllerFact) Items() []RemoveControllerItem {
	return fact.items
}

func (fact RemoveControllerFact) Addresses() ([]base.Address, error) {
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

type RemoveController struct {
	common.BaseOperation
}

func NewRemoveController(fact RemoveControllerFact) (RemoveController, error) {
	return RemoveController{BaseOperation: common.NewBaseOperation(RemoveControllerHint, fact)}, nil
}

func (op *RemoveController) HashSign(priv base.Privatekey, networkID base.NetworkID) error {
	err := op.Sign(priv, networkID)
	if err != nil {
		return err
	}
	return nil
}
