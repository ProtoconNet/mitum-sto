package processor

import (
	"github.com/ProtoconNet/mitum-currency/v3/operation/currency"
	extensioncurrency "github.com/ProtoconNet/mitum-currency/v3/operation/extension"
	currencyprocessor "github.com/ProtoconNet/mitum-currency/v3/operation/processor"
	currencytypes "github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum-sto/operation/kyc"
	"github.com/ProtoconNet/mitum-sto/operation/sto"
	mitumbase "github.com/ProtoconNet/mitum2/base"
	"github.com/pkg/errors"
)

const (
	DuplicationTypeSender   currencytypes.DuplicationType = "sender"
	DuplicationTypeCurrency currencytypes.DuplicationType = "currency"
	DuplicationTypeContract currencytypes.DuplicationType = "contract"
)

func CheckDuplication(opr *currencyprocessor.OperationProcessor, op mitumbase.Operation) error {
	opr.Lock()
	defer opr.Unlock()

	var duplicationTypeSenderID string
	var duplicationTypeCurrencyID string
	var duplicationTypeContract string
	var newAddresses []mitumbase.Address

	switch t := op.(type) {
	case currency.CreateAccount:
		fact, ok := t.Fact().(currency.CreateAccountFact)
		if !ok {
			return errors.Errorf("expected CreateAccountFact, not %T", t.Fact())
		}
		as, err := fact.Targets()
		if err != nil {
			return errors.Errorf("failed to get Addresses")
		}
		newAddresses = as
		duplicationTypeSenderID = fact.Sender().String()
	case currency.UpdateKey:
		fact, ok := t.Fact().(currency.UpdateKeyFact)
		if !ok {
			return errors.Errorf("expected UpdateKeyFact, not %T", t.Fact())
		}
		duplicationTypeSenderID = fact.Target().String()
	case currency.Transfer:
		fact, ok := t.Fact().(currency.TransferFact)
		if !ok {
			return errors.Errorf("expected TransferFact, not %T", t.Fact())
		}
		duplicationTypeSenderID = fact.Sender().String()
	case currency.RegisterCurrency:
		fact, ok := t.Fact().(currency.RegisterCurrencyFact)
		if !ok {
			return errors.Errorf("expected RegisterCurrencyFact, not %T", t.Fact())
		}
		duplicationTypeCurrencyID = fact.Currency().Currency().String()
	case currency.UpdateCurrency:
		fact, ok := t.Fact().(currency.UpdateCurrencyFact)
		if !ok {
			return errors.Errorf("expected UpdateCurrencyFact, not %T", t.Fact())
		}
		duplicationTypeSenderID = fact.Currency().String()
	case currency.Mint:
	case extensioncurrency.CreateContractAccount:
		fact, ok := t.Fact().(extensioncurrency.CreateContractAccountFact)
		if !ok {
			return errors.Errorf("expected CreateContractAccountFact, not %T", t.Fact())
		}
		as, err := fact.Targets()
		if err != nil {
			return errors.Errorf("failed to get Addresses")
		}
		newAddresses = as
		duplicationTypeSenderID = fact.Sender().String()
	case extensioncurrency.Withdraw:
		fact, ok := t.Fact().(extensioncurrency.WithdrawFact)
		if !ok {
			return errors.Errorf("expected WithdrawFact, not %T", t.Fact())
		}
		duplicationTypeSenderID = fact.Sender().String()
	case sto.AuthorizeOperator:
		fact, ok := t.Fact().(sto.AuthorizeOperatorFact)
		if !ok {
			return errors.Errorf("expected AuthorizeOperator, not %T", t.Fact())
		}
		duplicationTypeSenderID = fact.Sender().String()
	case sto.CreateSecurityToken:
		fact, ok := t.Fact().(sto.CreateSecurityTokenFact)
		if !ok {
			return errors.Errorf("expected CreateSecurityTokenFact, not %T", t.Fact())
		}
		duplicationTypeSenderID = fact.Sender().String()
	case sto.Issue:
		fact, ok := t.Fact().(sto.IssueFact)
		if !ok {
			return errors.Errorf("expected IssueFact, not %T", t.Fact())
		}
		duplicationTypeSenderID = fact.Sender().String()
	case sto.Redeem:
		fact, ok := t.Fact().(sto.RedeemFact)
		if !ok {
			return errors.Errorf("expected RedeemFact, not %T", t.Fact())
		}
		duplicationTypeSenderID = fact.Sender().String()
	case sto.RevokeOperator:
		fact, ok := t.Fact().(sto.RevokeOperatorFact)
		if !ok {
			return errors.Errorf("expected RevokeOperatorFact, not %T", t.Fact())
		}
		duplicationTypeSenderID = fact.Sender().String()
	case sto.SetDocument:
		fact, ok := t.Fact().(sto.SetDocumentFact)
		if !ok {
			return errors.Errorf("expected SetDocumentFact, not %T", t.Fact())
		}
		duplicationTypeSenderID = fact.Sender().String()
	case sto.TransferByPartition:
		fact, ok := t.Fact().(sto.TransferByPartitionFact)
		if !ok {
			return errors.Errorf("expected TransferByPartitionFact, not %T", t.Fact())
		}
		duplicationTypeSenderID = fact.Sender().String()
	case kyc.AddController:
		fact, ok := t.Fact().(kyc.AddControllerFact)
		if !ok {
			return errors.Errorf("expected AddControllerFact, not %T", t.Fact())
		}
		duplicationTypeSenderID = fact.Sender().String()
	case kyc.AddCustomer:
		fact, ok := t.Fact().(kyc.AddCustomerFact)
		if !ok {
			return errors.Errorf("expected AddCustomerFact, not %T", t.Fact())
		}
		duplicationTypeSenderID = fact.Sender().String()
	case kyc.CreateService:
		fact, ok := t.Fact().(kyc.CreateServiceFact)
		if !ok {
			return errors.Errorf("expected CreateServiceFact, not %T", t.Fact())
		}
		duplicationTypeSenderID = fact.Sender().String()
	case kyc.RemoveController:
		fact, ok := t.Fact().(kyc.RemoveControllerFact)
		if !ok {
			return errors.Errorf("expected RemoveControllerFact, not %T", t.Fact())
		}
		duplicationTypeSenderID = fact.Sender().String()
	case kyc.UpdateCustomers:
		fact, ok := t.Fact().(kyc.UpdateCustomersFact)
		if !ok {
			return errors.Errorf("expected UpdateCustomersFact, not %T", t.Fact())
		}
		duplicationTypeSenderID = fact.Sender().String()
	default:
		return nil
	}

	if len(duplicationTypeSenderID) > 0 {
		if _, found := opr.Duplicated[duplicationTypeSenderID]; found {
			return errors.Errorf("proposal cannot have duplicate sender, %v", duplicationTypeSenderID)
		}

		opr.Duplicated[duplicationTypeSenderID] = DuplicationTypeSender
	}
	if len(duplicationTypeCurrencyID) > 0 {
		if _, found := opr.Duplicated[duplicationTypeCurrencyID]; found {
			return errors.Errorf(
				"cannot register duplicate currency id, %v within a proposal",
				duplicationTypeCurrencyID,
			)
		}

		opr.Duplicated[duplicationTypeCurrencyID] = DuplicationTypeCurrency
	}
	if len(duplicationTypeContract) > 0 {
		if _, found := opr.Duplicated[duplicationTypeContract]; found {
			return errors.Errorf(
				"cannot use a duplicated contract for registering in contract model , %v within a proposal",
				duplicationTypeSenderID,
			)
		}

		opr.Duplicated[duplicationTypeContract] = DuplicationTypeContract
	}

	if len(newAddresses) > 0 {
		if err := opr.CheckNewAddressDuplication(newAddresses); err != nil {
			return err
		}
	}

	return nil
}

func GetNewProcessor(opr *currencyprocessor.OperationProcessor, op mitumbase.Operation) (mitumbase.OperationProcessor, bool, error) {
	switch i, err := opr.GetNewProcessorFromHintset(op); {
	case err != nil:
		return nil, false, err
	case i != nil:
		return i, true, nil
	}

	switch t := op.(type) {
	case currency.CreateAccount,
		currency.UpdateKey,
		currency.Transfer,
		extensioncurrency.CreateContractAccount,
		extensioncurrency.Withdraw,
		currency.RegisterCurrency,
		currency.UpdateCurrency,
		currency.Mint,
		sto.AuthorizeOperator,
		sto.CreateSecurityToken,
		sto.Issue,
		sto.Redeem,
		sto.RevokeOperator,
		sto.SetDocument,
		sto.TransferByPartition:
		return nil, false, errors.Errorf("%T needs SetProcessor", t)
	default:
		return nil, false, nil
	}
}
