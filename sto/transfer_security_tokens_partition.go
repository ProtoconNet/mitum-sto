package sto

import (
	"github.com/ProtoconNet/mitum-currency/v2/currency"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
	"github.com/ProtoconNet/mitum2/util/valuehash"
)

var (
	TransferSecurityTokensPartitionFactHint = hint.MustNewHint("mitum-sto-transfer-security-tokens-partition-operation-fact-v0.0.1")
	TransferSecurityTokensPartitionHint     = hint.MustNewHint("mitum-sto-transfer-security-tokens-partition-operation-v0.0.1")
)

var MaxTransferSecurityTokensPartitionItems uint = 10

type TransferSecurityTokensPartitionFact struct {
	base.BaseFact
	sender base.Address
	items  []TransferSecurityTokensPartitionItem
}

func NewTransferSecurityTokensPartitionFact(token []byte, sender base.Address, items []TransferSecurityTokensPartitionItem) TransferSecurityTokensPartitionFact {
	bf := base.NewBaseFact(TransferSecurityTokensPartitionFactHint, token)
	fact := TransferSecurityTokensPartitionFact{
		BaseFact: bf,
		sender:   sender,
		items:    items,
	}
	fact.SetHash(fact.GenerateHash())

	return fact
}

func (fact TransferSecurityTokensPartitionFact) Hash() util.Hash {
	return fact.BaseFact.Hash()
}

func (fact TransferSecurityTokensPartitionFact) GenerateHash() util.Hash {
	return valuehash.NewSHA256(fact.Bytes())
}

func (fact TransferSecurityTokensPartitionFact) Bytes() []byte {
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

func (fact TransferSecurityTokensPartitionFact) IsValid(b []byte) error {
	if err := fact.BaseHinter.IsValid(nil); err != nil {
		return err
	}

	if err := currency.IsValidOperationFact(fact, b); err != nil {
		return err
	}

	if n := len(fact.items); n < 1 {
		return util.ErrInvalid.Errorf("empty items")
	} else if n > int(MaxTransferSecurityTokensPartitionItems) {
		return util.ErrInvalid.Errorf("items, %d over max, %d", n, MaxTransferSecurityTokensPartitionItems)
	}

	if err := fact.sender.IsValid(nil); err != nil {
		return err
	}

	founds := map[string]struct{}{}
	for _, it := range fact.items {
		if err := it.IsValid(nil); err != nil {
			return err
		}

		k := it.receiver.String() + it.partition.String()
		if _, found := founds[k]; found {
			return util.ErrInvalid.Errorf("duplicate receiver-partition found, %s", k)
		}

		if fact.sender.Equal(it.receiver) {
			return util.ErrInvalid.Errorf("receiver address is same with sender, %q", fact.sender)
		}

		founds[k] = struct{}{}
	}

	return nil
}

func (fact TransferSecurityTokensPartitionFact) Token() base.Token {
	return fact.BaseFact.Token()
}

func (fact TransferSecurityTokensPartitionFact) Sender() base.Address {
	return fact.sender
}

func (fact TransferSecurityTokensPartitionFact) Items() []TransferSecurityTokensPartitionItem {
	return fact.items
}

func (fact TransferSecurityTokensPartitionFact) Addresses() ([]base.Address, error) {
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

type TransferSecurityTokensPartition struct {
	currency.BaseOperation
}

func NewTransferSecurityTokensPartition(fact TransferSecurityTokensPartitionFact) (TransferSecurityTokensPartition, error) {
	return TransferSecurityTokensPartition{BaseOperation: currency.NewBaseOperation(TransferSecurityTokensPartitionHint, fact)}, nil
}

func (op *TransferSecurityTokensPartition) HashSign(priv base.Privatekey, networkID base.NetworkID) error {
	err := op.Sign(priv, networkID)
	if err != nil {
		return err
	}
	return nil
}
