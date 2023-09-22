package cmds

import (
	currencycmds "github.com/ProtoconNet/mitum-currency/v3/cmds"
	"github.com/ProtoconNet/mitum-sto/operation/kyc"
	"github.com/ProtoconNet/mitum-sto/operation/sto"
	kycstate "github.com/ProtoconNet/mitum-sto/state/kyc"
	stostate "github.com/ProtoconNet/mitum-sto/state/sto"
	kyctypes "github.com/ProtoconNet/mitum-sto/types/kyc"
	stotypes "github.com/ProtoconNet/mitum-sto/types/sto"

	"github.com/ProtoconNet/mitum2/launch"
	"github.com/ProtoconNet/mitum2/util/encoder"
	"github.com/pkg/errors"
)

var Hinters []encoder.DecodeDetail
var SupportedProposalOperationFactHinters []encoder.DecodeDetail

var AddedHinters = []encoder.DecodeDetail{
	// revive:disable-next-line:line-length-limit
	{Hint: stostate.DesignStateValueHint, Instance: stostate.DesignStateValue{}},
	{Hint: stostate.TokenHolderPartitionsStateValueHint, Instance: stostate.TokenHolderPartitionsStateValue{}},
	{Hint: stostate.TokenHolderPartitionBalanceStateValueHint, Instance: stostate.TokenHolderPartitionBalanceStateValue{}},
	{Hint: stostate.TokenHolderPartitionOperatorsStateValueHint, Instance: stostate.TokenHolderPartitionOperatorsStateValue{}},
	{Hint: stostate.PartitionBalanceStateValueHint, Instance: stostate.PartitionBalanceStateValue{}},
	{Hint: stostate.OperatorTokenHoldersStateValueHint, Instance: stostate.OperatorTokenHoldersStateValue{}},
	{Hint: stotypes.DesignHint, Instance: stotypes.Design{}},
	{Hint: stotypes.DocumentHint, Instance: stotypes.Document{}},
	{Hint: stotypes.PolicyHint, Instance: stotypes.Policy{}},
	{Hint: sto.CreateSecurityTokenItemHint, Instance: sto.CreateSecurityTokenItem{}},
	{Hint: sto.CreateSecurityTokenHint, Instance: sto.CreateSecurityToken{}},
	{Hint: sto.IssueItemHint, Instance: sto.IssueItem{}},
	{Hint: sto.IssueHint, Instance: sto.Issue{}},
	{Hint: sto.TransferByPartitionItemHint, Instance: sto.TransferByPartitionItem{}},
	{Hint: sto.TransferByPartitionHint, Instance: sto.TransferByPartition{}},
	{Hint: sto.RedeemItemHint, Instance: sto.RedeemItem{}},
	{Hint: sto.RedeemHint, Instance: sto.Redeem{}},
	{Hint: sto.AuthorizeOperatorItemHint, Instance: sto.AuthorizeOperatorItem{}},
	{Hint: sto.AuthorizeOperatorHint, Instance: sto.AuthorizeOperator{}},
	{Hint: sto.RevokeOperatorItemHint, Instance: sto.RevokeOperatorItem{}},
	{Hint: sto.RevokeOperatorHint, Instance: sto.RevokeOperator{}},
	{Hint: sto.SetDocumentHint, Instance: sto.SetDocument{}},

	{Hint: kyctypes.DesignHint, Instance: kyctypes.Design{}},
	{Hint: kycstate.DesignStateValueHint, Instance: kycstate.DesignStateValue{}},
	{Hint: kyctypes.PolicyHint, Instance: kyctypes.Policy{}},
	{Hint: kycstate.CustomerStateValueHint, Instance: kycstate.CustomerStateValue{}},
	{Hint: kyc.CreateServiceHint, Instance: kyc.CreateService{}},
	{Hint: kyc.AddControllerItemHint, Instance: kyc.AddControllerItem{}},
	{Hint: kyc.AddControllerHint, Instance: kyc.AddController{}},
	{Hint: kyc.RemoveControllerItemHint, Instance: kyc.RemoveControllerItem{}},
	{Hint: kyc.RemoveControllerHint, Instance: kyc.RemoveController{}},
	{Hint: kyc.AddCustomerItemHint, Instance: kyc.AddCustomerItem{}},
	{Hint: kyc.AddCustomerHint, Instance: kyc.AddCustomer{}},
	{Hint: kyc.UpdateCustomersItemHint, Instance: kyc.UpdateCustomersItem{}},
	{Hint: kyc.UpdateCustomersHint, Instance: kyc.UpdateCustomers{}},
}

var AddedSupportedHinters = []encoder.DecodeDetail{
	{Hint: sto.CreateSecurityTokenFactHint, Instance: sto.CreateSecurityTokenFact{}},
	{Hint: sto.IssueFactHint, Instance: sto.IssueFact{}},
	{Hint: sto.TransferByPartitionFactHint, Instance: sto.TransferByPartitionFact{}},
	{Hint: sto.RedeemFactHint, Instance: sto.RedeemFact{}},
	{Hint: sto.AuthorizeOperatorFactHint, Instance: sto.AuthorizeOperatorFact{}},
	{Hint: sto.RevokeOperatorFactHint, Instance: sto.RevokeOperatorFact{}},
	{Hint: sto.SetDocumentFactHint, Instance: sto.SetDocumentFact{}},

	{Hint: kyc.CreateServiceFactHint, Instance: kyc.CreateServiceFact{}},
	{Hint: kyc.AddControllerFactHint, Instance: kyc.AddControllerFact{}},
	{Hint: kyc.RemoveControllerFactHint, Instance: kyc.RemoveController{}},
	{Hint: kyc.AddCustomerFactHint, Instance: kyc.AddCustomerFact{}},
	{Hint: kyc.UpdateCustomersFactHint, Instance: kyc.UpdateCustomersFact{}},
}

func init() {
	defaultLen := len(launch.Hinters)
	currencyExtendedLen := defaultLen + len(currencycmds.AddedHinters)
	allExtendedLen := currencyExtendedLen + len(AddedHinters)

	Hinters = make([]encoder.DecodeDetail, allExtendedLen)
	copy(Hinters, launch.Hinters)
	copy(Hinters[defaultLen:currencyExtendedLen], currencycmds.AddedHinters)
	copy(Hinters[currencyExtendedLen:], AddedHinters)

	defaultSupportedLen := len(launch.SupportedProposalOperationFactHinters)
	currencySupportedExtendedLen := defaultSupportedLen + len(currencycmds.AddedSupportedHinters)
	allSupportedExtendedLen := currencySupportedExtendedLen + len(AddedSupportedHinters)

	SupportedProposalOperationFactHinters = make(
		[]encoder.DecodeDetail,
		allSupportedExtendedLen)
	copy(SupportedProposalOperationFactHinters, launch.SupportedProposalOperationFactHinters)
	copy(SupportedProposalOperationFactHinters[defaultSupportedLen:currencySupportedExtendedLen], currencycmds.AddedSupportedHinters)
	copy(SupportedProposalOperationFactHinters[currencySupportedExtendedLen:], AddedSupportedHinters)
}

func LoadHinters(enc encoder.Encoder) error {
	for _, hinter := range Hinters {
		if err := enc.Add(hinter); err != nil {
			return errors.Wrap(err, "failed to add to encoder")
		}
	}

	for _, hinter := range SupportedProposalOperationFactHinters {
		if err := enc.Add(hinter); err != nil {
			return errors.Wrap(err, "failed to add to encoder")
		}
	}

	return nil
}
