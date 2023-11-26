package kyc

import (
	"github.com/ProtoconNet/mitum-currency/v3/common"
	"github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
	"github.com/ProtoconNet/mitum2/util/valuehash"
)

var (
	CreateServiceFactHint = hint.MustNewHint("mitum-kyc-create-service-operation-fact-v0.0.1")
	CreateServiceHint     = hint.MustNewHint("mitum-kyc-create-service-operation-v0.0.1")
)

type CreateServiceFact struct {
	base.BaseFact
	sender      base.Address
	contract    base.Address // contract account
	controllers []base.Address
	currency    types.CurrencyID
}

func NewCreateServiceFact(
	token []byte,
	sender base.Address,
	contract base.Address,
	controllers []base.Address,
	currency types.CurrencyID,
) CreateServiceFact {
	bf := base.NewBaseFact(CreateServiceFactHint, token)
	fact := CreateServiceFact{
		BaseFact:    bf,
		sender:      sender,
		contract:    contract,
		controllers: controllers,
		currency:    currency,
	}
	fact.SetHash(fact.GenerateHash())

	return fact
}

func (fact CreateServiceFact) Hash() util.Hash {
	return fact.BaseFact.Hash()
}

func (fact CreateServiceFact) GenerateHash() util.Hash {
	return valuehash.NewSHA256(fact.Bytes())
}

func (fact CreateServiceFact) Bytes() []byte {
	bc := make([][]byte, len(fact.controllers))

	for i, con := range fact.controllers {
		bc[i] = con.Bytes()
	}

	return util.ConcatBytesSlice(
		fact.Token(),
		fact.sender.Bytes(),
		fact.contract.Bytes(),
		util.ConcatBytesSlice(bc...),
		fact.currency.Bytes(),
	)
}

func (fact CreateServiceFact) IsValid(b []byte) error {
	if err := fact.BaseHinter.IsValid(nil); err != nil {
		return err
	}

	if err := util.CheckIsValiders(nil, false, fact.sender, fact.contract, fact.currency); err != nil {
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

	if err := common.IsValidOperationFact(fact, b); err != nil {
		return err
	}

	return nil
}

func (fact CreateServiceFact) Token() base.Token {
	return fact.BaseFact.Token()
}

func (fact CreateServiceFact) Sender() base.Address {
	return fact.sender
}

func (fact CreateServiceFact) Contract() base.Address {
	return fact.contract
}

func (fact CreateServiceFact) Controllers() []base.Address {
	return append([]base.Address{}, fact.controllers...)
}

func (fact CreateServiceFact) Currency() types.CurrencyID {
	return fact.currency
}

func (fact CreateServiceFact) Addresses() ([]base.Address, error) {
	as := make([]base.Address, 2+len(fact.controllers))

	as[0] = fact.sender
	as[1] = fact.contract

	for i, con := range fact.controllers {
		as[i+2] = con
	}

	return as, nil
}

type CreateService struct {
	common.BaseOperation
}

func NewCreateService(fact CreateServiceFact) (CreateService, error) {
	return CreateService{BaseOperation: common.NewBaseOperation(CreateServiceHint, fact)}, nil
}

func (op *CreateService) HashSign(priv base.Privatekey, networkID base.NetworkID) error {
	err := op.Sign(priv, networkID)
	if err != nil {
		return err
	}
	return nil
}
