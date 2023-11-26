package cmds

import (
	crcycmds "github.com/ProtoconNet/mitum-currency/v3/cmds"
	"github.com/ProtoconNet/mitum-sto/operation/kyc"
	"github.com/ProtoconNet/mitum-sto/operation/sto"
	stkyc "github.com/ProtoconNet/mitum-sto/state/kyc"
	ststo "github.com/ProtoconNet/mitum-sto/state/sto"
	typekyc "github.com/ProtoconNet/mitum-sto/types/kyc"
	typesto "github.com/ProtoconNet/mitum-sto/types/sto"
	"github.com/ProtoconNet/mitum2/launch"
	"github.com/ProtoconNet/mitum2/util/encoder"
	"github.com/pkg/errors"
)

var Hinters []encoder.DecodeDetail
var SupportedProposalOperationFactHinters []encoder.DecodeDetail

var AddedHinters = []encoder.DecodeDetail{
	// revive:disable-next-line:line-length-limit
	{Hint: ststo.DesignStateValueHint, Instance: ststo.DesignStateValue{}},
	{Hint: ststo.TokenHolderPartitionsStateValueHint, Instance: ststo.TokenHolderPartitionsStateValue{}},
	{Hint: ststo.TokenHolderPartitionBalanceStateValueHint, Instance: ststo.TokenHolderPartitionBalanceStateValue{}},
	{Hint: ststo.TokenHolderPartitionOperatorsStateValueHint, Instance: ststo.TokenHolderPartitionOperatorsStateValue{}},
	{Hint: ststo.PartitionBalanceStateValueHint, Instance: ststo.PartitionBalanceStateValue{}},
	{Hint: ststo.OperatorTokenHoldersStateValueHint, Instance: ststo.OperatorTokenHoldersStateValue{}},
	{Hint: typesto.DesignHint, Instance: typesto.Design{}},
	{Hint: typesto.DocumentHint, Instance: typesto.Document{}},
	{Hint: typesto.PolicyHint, Instance: typesto.Policy{}},
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

	{Hint: typekyc.DesignHint, Instance: typekyc.Design{}},
	{Hint: stkyc.DesignStateValueHint, Instance: stkyc.DesignStateValue{}},
	{Hint: typekyc.PolicyHint, Instance: typekyc.Policy{}},
	{Hint: stkyc.CustomerStateValueHint, Instance: stkyc.CustomerStateValue{}},
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
	currencyExtendedLen := defaultLen + len(crcycmds.AddedHinters)
	allExtendedLen := currencyExtendedLen + len(AddedHinters)

	Hinters = make([]encoder.DecodeDetail, allExtendedLen)
	copy(Hinters, launch.Hinters)
	copy(Hinters[defaultLen:currencyExtendedLen], crcycmds.AddedHinters)
	copy(Hinters[currencyExtendedLen:], AddedHinters)

	defaultSupportedLen := len(launch.SupportedProposalOperationFactHinters)
	currencySupportedExtendedLen := defaultSupportedLen + len(crcycmds.AddedSupportedHinters)
	allSupportedExtendedLen := currencySupportedExtendedLen + len(AddedSupportedHinters)

	SupportedProposalOperationFactHinters = make(
		[]encoder.DecodeDetail,
		allSupportedExtendedLen)
	copy(SupportedProposalOperationFactHinters, launch.SupportedProposalOperationFactHinters)
	copy(SupportedProposalOperationFactHinters[defaultSupportedLen:currencySupportedExtendedLen], crcycmds.AddedSupportedHinters)
	copy(SupportedProposalOperationFactHinters[currencySupportedExtendedLen:], AddedSupportedHinters)
}

func LoadHinters(encs *encoder.Encoders) error {
	for i := range Hinters {
		if err := encs.AddDetail(Hinters[i]); err != nil {
			return errors.Wrap(err, "add hinter to encoder")
		}
	}

	for i := range SupportedProposalOperationFactHinters {
		if err := encs.AddDetail(SupportedProposalOperationFactHinters[i]); err != nil {
			return errors.Wrap(err, "add supported proposal operation fact hinter to encoder")
		}
	}

	return nil
}
