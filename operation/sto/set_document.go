package sto

import (
	"github.com/ProtoconNet/mitum-currency/v3/common"
	currencytypes "github.com/ProtoconNet/mitum-currency/v3/types"
	stotypes "github.com/ProtoconNet/mitum-sto/types/sto"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
	"github.com/ProtoconNet/mitum2/util/valuehash"
)

var (
	SetDocumentFactHint = hint.MustNewHint("mitum-sto-set-document-operation-fact-v0.0.1")
	SetDocumentHint     = hint.MustNewHint("mitum-sto-set-document-operation-v0.0.1")
)

type SetDocumentFact struct {
	base.BaseFact
	sender       base.Address
	contract     base.Address             // contract account
	stoID        currencytypes.ContractID // token id
	title        string                   // document title
	uri          stotypes.URI             // document uri
	documentHash string                   // document hash
	currency     currencytypes.CurrencyID // fee
}

func NewSetDocumentFact(
	token []byte,
	sender base.Address,
	contract base.Address,
	stoID currencytypes.ContractID,
	title string,
	uri stotypes.URI,
	hash string,
	currency currencytypes.CurrencyID,
) SetDocumentFact {
	bf := base.NewBaseFact(SetDocumentFactHint, token)
	fact := SetDocumentFact{
		BaseFact:     bf,
		sender:       sender,
		contract:     contract,
		stoID:        stoID,
		title:        title,
		uri:          uri,
		documentHash: hash,
		currency:     currency,
	}
	fact.SetHash(fact.GenerateHash())

	return fact
}

func (fact SetDocumentFact) Hash() util.Hash {
	return fact.BaseFact.Hash()
}

func (fact SetDocumentFact) GenerateHash() util.Hash {
	return valuehash.NewSHA256(fact.Bytes())
}

func (fact SetDocumentFact) Bytes() []byte {
	return util.ConcatBytesSlice(
		fact.Token(),
		fact.sender.Bytes(),
		fact.contract.Bytes(),
		fact.stoID.Bytes(),
		[]byte(fact.title),
		fact.uri.Bytes(),
		[]byte(fact.documentHash),
		fact.currency.Bytes(),
	)
}

func (fact SetDocumentFact) IsValid(b []byte) error {
	if err := fact.BaseHinter.IsValid(nil); err != nil {
		return err
	}

	if err := common.IsValidOperationFact(fact, b); err != nil {
		return err
	}

	if err := util.CheckIsValiders(nil, false, fact.sender, fact.stoID, fact.contract, fact.uri, fact.currency); err != nil {
		return err
	}

	return nil
}

func (fact SetDocumentFact) Token() base.Token {
	return fact.BaseFact.Token()
}

func (fact SetDocumentFact) Sender() base.Address {
	return fact.sender
}

func (fact SetDocumentFact) Contract() base.Address {
	return fact.contract
}

func (fact SetDocumentFact) STO() currencytypes.ContractID {
	return fact.stoID
}

func (fact SetDocumentFact) Title() string {
	return fact.title
}

func (fact SetDocumentFact) URI() stotypes.URI {
	return fact.uri
}

func (fact SetDocumentFact) DocumentHash() string {
	return fact.documentHash
}

func (fact SetDocumentFact) Currency() currencytypes.CurrencyID {
	return fact.currency
}

func (fact SetDocumentFact) Addresses() ([]base.Address, error) {
	as := make([]base.Address, 2)

	as[0] = fact.sender
	as[1] = fact.contract

	return as, nil
}

type SetDocument struct {
	common.BaseOperation
}

func NewSetDocument(fact SetDocumentFact) (SetDocument, error) {
	return SetDocument{BaseOperation: common.NewBaseOperation(SetDocumentHint, fact)}, nil
}

func (op *SetDocument) HashSign(priv base.Privatekey, networkID base.NetworkID) error {
	err := op.Sign(priv, networkID)
	if err != nil {
		return err
	}
	return nil
}
