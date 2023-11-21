package kyc

import (
	"github.com/ProtoconNet/mitum-currency/v3/common"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
	"github.com/ProtoconNet/mitum2/util/valuehash"
)

var (
	AddControllerFactHint = hint.MustNewHint("mitum-kyc-add-controller-operation-fact-v0.0.1")
	AddControllerHint     = hint.MustNewHint("mitum-kyc-add-controller-operation-v0.0.1")
)

var MaxAddControllersItems uint = 10

type AddControllerFact struct {
	base.BaseFact
	sender base.Address
	items  []AddControllerItem
}

func NewAddControllerFact(token []byte, sender base.Address, items []AddControllerItem) AddControllerFact {
	bf := base.NewBaseFact(AddControllerFactHint, token)
	fact := AddControllerFact{
		BaseFact: bf,
		sender:   sender,
		items:    items,
	}
	fact.SetHash(fact.GenerateHash())

	return fact
}

func (fact AddControllerFact) Hash() util.Hash {
	return fact.BaseFact.Hash()
}

func (fact AddControllerFact) GenerateHash() util.Hash {
	return valuehash.NewSHA256(fact.Bytes())
}

func (fact AddControllerFact) Bytes() []byte {
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

func (fact AddControllerFact) IsValid(b []byte) error {
	if err := fact.BaseHinter.IsValid(nil); err != nil {
		return err
	}

	if n := len(fact.items); n < 1 {
		return util.ErrInvalid.Errorf("empty items")
	} else if n > int(MaxAddControllersItems) {
		return util.ErrInvalid.Errorf("items, %d over max, %d", n, MaxAddControllersItems)
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

func (fact AddControllerFact) Token() base.Token {
	return fact.BaseFact.Token()
}

func (fact AddControllerFact) Sender() base.Address {
	return fact.sender
}

func (fact AddControllerFact) Items() []AddControllerItem {
	return fact.items
}

func (fact AddControllerFact) Addresses() ([]base.Address, error) {
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

type AddController struct {
	common.BaseOperation
}

func NewAddControllers(fact AddControllerFact) (AddController, error) {
	return AddController{BaseOperation: common.NewBaseOperation(AddControllerHint, fact)}, nil
}
