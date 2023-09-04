package kyc

import (
	"github.com/ProtoconNet/mitum-currency/v3/common"
	currencytypes "github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
	"github.com/ProtoconNet/mitum2/util/valuehash"
)

var (
	CreateKYCServiceFactHint = hint.MustNewHint("mitum-kyc-create-key-service-operation-fact-v0.0.1")
	CreateKYCServiceHint     = hint.MustNewHint("mitum-kyc-create-key-service-operation-v0.0.1")
)

type CreateKYCServiceFact struct {
	base.BaseFact
	sender      base.Address
	contract    base.Address             // contract account
	kycID       currencytypes.ContractID // kyc id
	controllers []base.Address
	currency    currencytypes.CurrencyID
}

func NewCreateKYCServiceFact(
	token []byte,
	sender base.Address,
	contract base.Address,
	kycID currencytypes.ContractID,
	controllers []base.Address,
	currency currencytypes.CurrencyID,
) CreateKYCServiceFact {
	bf := base.NewBaseFact(CreateKYCServiceFactHint, token)
	fact := CreateKYCServiceFact{
		BaseFact:    bf,
		sender:      sender,
		contract:    contract,
		kycID:       kycID,
		controllers: controllers,
		currency:    currency,
	}
	fact.SetHash(fact.GenerateHash())

	return fact
}

func (fact CreateKYCServiceFact) Hash() util.Hash {
	return fact.BaseFact.Hash()
}

func (fact CreateKYCServiceFact) GenerateHash() util.Hash {
	return valuehash.NewSHA256(fact.Bytes())
}

func (fact CreateKYCServiceFact) Bytes() []byte {
	bc := make([][]byte, len(fact.controllers))

	for i, con := range fact.controllers {
		bc[i] = con.Bytes()
	}

	return util.ConcatBytesSlice(
		fact.Token(),
		fact.sender.Bytes(),
		fact.contract.Bytes(),
		fact.kycID.Bytes(),
		util.ConcatBytesSlice(bc...),
		fact.currency.Bytes(),
	)
}

func (fact CreateKYCServiceFact) IsValid(b []byte) error {
	if err := fact.BaseHinter.IsValid(nil); err != nil {
		return err
	}

	if err := common.IsValidOperationFact(fact, b); err != nil {
		return err
	}

	if err := util.CheckIsValiders(nil, false, fact.sender, fact.kycID, fact.contract, fact.currency); err != nil {
		return err
	}

	if fact.sender.Equal(fact.contract) {
		return util.ErrInvalid.Errorf("contract address is same with sender, %q", fact.sender)
	}

	founds := map[string]struct{}{}
	for _, con := range fact.controllers {
		if err := con.IsValid(nil); err != nil {
			return err
		}

		if con.Equal(fact.sender) {
			return util.ErrInvalid.Errorf("controller address is same with sender, %q", fact.sender)
		}

		if _, found := founds[con.String()]; found {
			return util.ErrInvalid.Errorf("duplicate controller found, %q", con)
		}

		founds[con.String()] = struct{}{}
	}

	return nil
}

func (fact CreateKYCServiceFact) Token() base.Token {
	return fact.BaseFact.Token()
}

func (fact CreateKYCServiceFact) Sender() base.Address {
	return fact.sender
}

func (fact CreateKYCServiceFact) Contract() base.Address {
	return fact.contract
}

func (fact CreateKYCServiceFact) KYC() currencytypes.ContractID {
	return fact.kycID
}

func (fact CreateKYCServiceFact) Controllers() []base.Address {
	return append([]base.Address{}, fact.controllers...)
}

func (fact CreateKYCServiceFact) Currency() currencytypes.CurrencyID {
	return fact.currency
}

func (fact CreateKYCServiceFact) Addresses() ([]base.Address, error) {
	as := make([]base.Address, 2+len(fact.controllers))

	as[0] = fact.sender
	as[1] = fact.contract

	for i, con := range fact.controllers {
		as[i+2] = con
	}

	return as, nil
}

type CreateKYCService struct {
	common.BaseOperation
}

func NewCreateKYCService(fact CreateKYCServiceFact) (CreateKYCService, error) {
	return CreateKYCService{BaseOperation: common.NewBaseOperation(CreateKYCServiceHint, fact)}, nil
}

func (op *CreateKYCService) HashSign(priv base.Privatekey, networkID base.NetworkID) error {
	err := op.Sign(priv, networkID)
	if err != nil {
		return err
	}
	return nil
}
