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
	SetDocumentsFactHint = hint.MustNewHint("mitum-sto-set-document-operation-fact-v0.0.1")
	SetDocumentsHint     = hint.MustNewHint("mitum-sto-set-document-operation-v0.0.1")
)

type SetDocumentsFact struct {
	base.BaseFact
	sender       base.Address
	stoID        currencyextension.ContractID // token id
	contract     base.Address                 // contract account
	title        string                       // document title
	uri          URI                          // document uri
	documentHash string                       // document hash
	currency     currency.CurrencyID          // fee
}

func NewSetDocumentsFact(token []byte, stoID currencyextension.ContractID, sender, contract base.Address, title string, uri URI, hash string, currency currency.CurrencyID) SetDocumentsFact {
	bf := base.NewBaseFact(SetDocumentsFactHint, token)
	fact := SetDocumentsFact{
		BaseFact:     bf,
		sender:       sender,
		stoID:        stoID,
		contract:     contract,
		title:        title,
		uri:          uri,
		documentHash: hash,
		currency:     currency,
	}
	fact.SetHash(fact.GenerateHash())

	return fact
}

func (fact SetDocumentsFact) Hash() util.Hash {
	return fact.BaseFact.Hash()
}

func (fact SetDocumentsFact) GenerateHash() util.Hash {
	return valuehash.NewSHA256(fact.Bytes())
}

func (fact SetDocumentsFact) Bytes() []byte {
	return util.ConcatBytesSlice(
		fact.Token(),
		fact.sender.Bytes(),
		fact.stoID.Bytes(),
		fact.contract.Bytes(),
		[]byte(fact.title),
		fact.uri.Bytes(),
		[]byte(fact.documentHash),
		fact.currency.Bytes(),
	)
}

func (fact SetDocumentsFact) IsValid(b []byte) error {
	if err := fact.BaseHinter.IsValid(nil); err != nil {
		return err
	}

	if err := currency.IsValidOperationFact(fact, b); err != nil {
		return err
	}

	if err := util.CheckIsValiders(nil, false, fact.sender, fact.stoID, fact.contract, fact.uri, fact.currency); err != nil {
		return err
	}

	return nil
}

func (fact SetDocumentsFact) Token() base.Token {
	return fact.BaseFact.Token()
}

func (fact SetDocumentsFact) Sender() base.Address {
	return fact.sender
}

func (fact SetDocumentsFact) Addresses() ([]base.Address, error) {
	as := make([]base.Address, 2)

	as[0] = fact.sender
	as[1] = fact.contract

	return as, nil
}

type SetDocuments struct {
	currency.BaseOperation
}

func NewSetDocuments(fact SetDocumentsFact) (SetDocuments, error) {
	return SetDocuments{BaseOperation: currency.NewBaseOperation(SetDocumentsHint, fact)}, nil
}

func (op *SetDocuments) HashSign(priv base.Privatekey, networkID base.NetworkID) error {
	err := op.Sign(priv, networkID)
	if err != nil {
		return err
	}
	return nil
}
