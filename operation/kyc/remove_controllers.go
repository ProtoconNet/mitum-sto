package kyc

import (
	"github.com/ProtoconNet/mitum-currency/v3/common"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
	"github.com/ProtoconNet/mitum2/util/valuehash"
)

var (
	RemoveControllersFactHint = hint.MustNewHint("mitum-kyc-remove-controllers-operation-fact-v0.0.1")
	RemoveControllersHint     = hint.MustNewHint("mitum-kyc-remove-controllers-operation-v0.0.1")
)

var MaxRemoveControllersItems uint = 10

type RemoveControllersFact struct {
	base.BaseFact
	sender base.Address
	items  []RemoveControllersItem
}

func NewRemoveControllersFact(token []byte, sender base.Address, items []RemoveControllersItem) RemoveControllersFact {
	bf := base.NewBaseFact(RemoveControllersFactHint, token)
	fact := RemoveControllersFact{
		BaseFact: bf,
		sender:   sender,
		items:    items,
	}
	fact.SetHash(fact.GenerateHash())

	return fact
}

func (fact RemoveControllersFact) Hash() util.Hash {
	return fact.BaseFact.Hash()
}

func (fact RemoveControllersFact) GenerateHash() util.Hash {
	return valuehash.NewSHA256(fact.Bytes())
}

func (fact RemoveControllersFact) Bytes() []byte {
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

func (fact RemoveControllersFact) IsValid(b []byte) error {
	if err := fact.BaseHinter.IsValid(nil); err != nil {
		return err
	}

	if err := common.IsValidOperationFact(fact, b); err != nil {
		return err
	}

	if n := len(fact.items); n < 1 {
		return util.ErrInvalid.Errorf("empty items")
	} else if n > int(MaxRemoveControllersItems) {
		return util.ErrInvalid.Errorf("items, %d over max, %d", n, MaxRemoveControllersItems)
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

	return nil
}

func (fact RemoveControllersFact) Token() base.Token {
	return fact.BaseFact.Token()
}

func (fact RemoveControllersFact) Sender() base.Address {
	return fact.sender
}

func (fact RemoveControllersFact) Items() []RemoveControllersItem {
	return fact.items
}

func (fact RemoveControllersFact) Addresses() ([]base.Address, error) {
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

type RemoveControllers struct {
	common.BaseOperation
}

func NewRemoveControllers(fact RemoveControllersFact) (RemoveControllers, error) {
	return RemoveControllers{BaseOperation: common.NewBaseOperation(RemoveControllersHint, fact)}, nil
}

func (op *RemoveControllers) HashSign(priv base.Privatekey, networkID base.NetworkID) error {
	err := op.Sign(priv, networkID)
	if err != nil {
		return err
	}
	return nil
}
