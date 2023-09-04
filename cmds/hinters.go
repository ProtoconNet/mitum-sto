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
	{Hint: sto.CreateSecurityTokensItemHint, Instance: sto.CreateSecurityTokensItem{}},
	{Hint: sto.CreateSecurityTokensHint, Instance: sto.CreateSecurityTokens{}},
	{Hint: sto.IssueSecurityTokensItemHint, Instance: sto.IssueSecurityTokensItem{}},
	{Hint: sto.IssueSecurityTokensHint, Instance: sto.IssueSecurityTokens{}},
	{Hint: sto.TransferSecurityTokensPartitionItemHint, Instance: sto.TransferSecurityTokensPartitionItem{}},
	{Hint: sto.TransferSecurityTokensPartitionHint, Instance: sto.TransferSecurityTokensPartition{}},
	{Hint: sto.RedeemTokensItemHint, Instance: sto.RedeemTokensItem{}},
	{Hint: sto.RedeemTokensHint, Instance: sto.RedeemTokens{}},
	{Hint: sto.AuthorizeOperatorsItemHint, Instance: sto.AuthorizeOperatorsItem{}},
	{Hint: sto.AuthorizeOperatorsHint, Instance: sto.AuthorizeOperators{}},
	{Hint: sto.RevokeOperatorsItemHint, Instance: sto.RevokeOperatorsItem{}},
	{Hint: sto.RevokeOperatorsHint, Instance: sto.RevokeOperators{}},
	{Hint: sto.SetDocumentHint, Instance: sto.SetDocument{}},

	{Hint: kyctypes.DesignHint, Instance: kyctypes.Design{}},
	{Hint: kycstate.DesignStateValueHint, Instance: kycstate.DesignStateValue{}},
	{Hint: kyctypes.PolicyHint, Instance: kyctypes.Policy{}},
	{Hint: kycstate.CustomerStateValueHint, Instance: kycstate.CustomerStateValue{}},
	{Hint: kyc.CreateKYCServiceHint, Instance: kyc.CreateKYCService{}},
	{Hint: kyc.AddControllersItemHint, Instance: kyc.AddControllersItem{}},
	{Hint: kyc.AddControllersHint, Instance: kyc.AddControllers{}},
	{Hint: kyc.RemoveControllersItemHint, Instance: kyc.RemoveControllersItem{}},
	{Hint: kyc.RemoveControllersHint, Instance: kyc.RemoveControllers{}},
	{Hint: kyc.AddCustomersItemHint, Instance: kyc.AddCustomersItem{}},
	{Hint: kyc.AddCustomersHint, Instance: kyc.AddCustomers{}},
	{Hint: kyc.UpdateCustomersItemHint, Instance: kyc.UpdateCustomersItem{}},
	{Hint: kyc.UpdateCustomersHint, Instance: kyc.UpdateCustomers{}},
}

var AddedSupportedHinters = []encoder.DecodeDetail{
	{Hint: sto.CreateSecurityTokensFactHint, Instance: sto.CreateSecurityTokensFact{}},
	{Hint: sto.IssueSecurityTokensFactHint, Instance: sto.IssueSecurityTokensFact{}},
	{Hint: sto.TransferSecurityTokensPartitionFactHint, Instance: sto.TransferSecurityTokensPartitionFact{}},
	{Hint: sto.RedeemTokensFactHint, Instance: sto.RedeemTokensFact{}},
	{Hint: sto.AuthorizeOperatorsFactHint, Instance: sto.AuthorizeOperatorsFact{}},
	{Hint: sto.RevokeOperatorsFactHint, Instance: sto.RevokeOperatorsFact{}},
	{Hint: sto.SetDocumentFactHint, Instance: sto.SetDocumentFact{}},

	{Hint: kyc.CreateKYCServiceFactHint, Instance: kyc.CreateKYCServiceFact{}},
	{Hint: kyc.AddControllersFactHint, Instance: kyc.AddControllersFact{}},
	{Hint: kyc.RemoveControllersFactHint, Instance: kyc.RemoveControllers{}},
	{Hint: kyc.AddCustomersFactHint, Instance: kyc.AddCustomersFact{}},
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
