package sto

import (
	"github.com/ProtoconNet/mitum-currency/v3/common"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
	"github.com/ProtoconNet/mitum2/util/valuehash"
)

var (
	TransferByPartitionFactHint = hint.MustNewHint("mitum-sto-transfer-by-partition-operation-fact-v0.0.1")
	TransferByPartitionHint     = hint.MustNewHint("mitum-sto-transfer-by-partition-operation-v0.0.1")
)

var MaxTransferByPartitionItems uint = 10

type TransferByPartitionFact struct {
	base.BaseFact
	sender base.Address
	items  []TransferByPartitionItem
}

func NewTransferByPartitionFact(token []byte, sender base.Address, items []TransferByPartitionItem) TransferByPartitionFact {
	bf := base.NewBaseFact(TransferByPartitionFactHint, token)
	fact := TransferByPartitionFact{
		BaseFact: bf,
		sender:   sender,
		items:    items,
	}
	fact.SetHash(fact.GenerateHash())

	return fact
}

func (fact TransferByPartitionFact) Hash() util.Hash {
	return fact.BaseFact.Hash()
}

func (fact TransferByPartitionFact) GenerateHash() util.Hash {
	return valuehash.NewSHA256(fact.Bytes())
}

func (fact TransferByPartitionFact) Bytes() []byte {
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

func (fact TransferByPartitionFact) IsValid(b []byte) error {
	if err := fact.BaseHinter.IsValid(nil); err != nil {
		return err
	}

	if n := len(fact.items); n < 1 {
		return util.ErrInvalid.Errorf("empty items")
	} else if n > int(MaxTransferByPartitionItems) {
		return util.ErrInvalid.Errorf("items, %d over max, %d", n, MaxTransferByPartitionItems)
	}

	if err := fact.sender.IsValid(nil); err != nil {
		return err
	}

	founds := map[string]struct{}{}
	for _, it := range fact.items {
		if err := it.IsValid(nil); err != nil {
			return err
		}

		k := it.tokenholder.String() + it.receiver.String() + it.partition.String()
		if _, found := founds[k]; found {
			return util.ErrInvalid.Errorf("duplicate token holder-receiver-partition found, %s", k)
		}

		founds[k] = struct{}{}
	}

	if err := common.IsValidOperationFact(fact, b); err != nil {
		return err
	}

	return nil
}

func (fact TransferByPartitionFact) Token() base.Token {
	return fact.BaseFact.Token()
}

func (fact TransferByPartitionFact) Sender() base.Address {
	return fact.sender
}

func (fact TransferByPartitionFact) Items() []TransferByPartitionItem {
	return fact.items
}

func (fact TransferByPartitionFact) Addresses() ([]base.Address, error) {
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

type TransferByPartition struct {
	common.BaseOperation
}

func NewTransferByPartition(fact TransferByPartitionFact) (TransferByPartition, error) {
	return TransferByPartition{BaseOperation: common.NewBaseOperation(TransferByPartitionHint, fact)}, nil
}

func (op *TransferByPartition) HashSign(priv base.Privatekey, networkID base.NetworkID) error {
	err := op.Sign(priv, networkID)
	if err != nil {
		return err
	}
	return nil
}
