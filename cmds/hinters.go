package cmds

import (
	"github.com/ProtoconNet/mitum-currency/v3/common"
	"github.com/ProtoconNet/mitum-currency/v3/digest"
	digestisaac "github.com/ProtoconNet/mitum-currency/v3/digest/isaac"
	mitumcurrency "github.com/ProtoconNet/mitum-currency/v3/operation/currency"
	extensioncurrency "github.com/ProtoconNet/mitum-currency/v3/operation/extension"
	isaacoperation "github.com/ProtoconNet/mitum-currency/v3/operation/isaac"
	currencystate "github.com/ProtoconNet/mitum-currency/v3/state/currency"
	extensionstate "github.com/ProtoconNet/mitum-currency/v3/state/extension"
	currencytypes "github.com/ProtoconNet/mitum-currency/v3/types"
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

var hinters = []encoder.DecodeDetail{
	// revive:disable-next-line:line-length-limit
	{Hint: common.BaseStateHint, Instance: common.BaseState{}},
	{Hint: common.NodeHint, Instance: common.BaseNode{}},
	{Hint: currencytypes.AccountHint, Instance: currencytypes.Account{}},
	{Hint: currencytypes.AddressHint, Instance: currencytypes.Address{}},
	{Hint: currencytypes.AmountHint, Instance: currencytypes.Amount{}},
	{Hint: currencytypes.AccountKeysHint, Instance: currencytypes.BaseAccountKeys{}},
	{Hint: currencytypes.AccountKeyHint, Instance: currencytypes.BaseAccountKey{}},
	{Hint: mitumcurrency.CreateAccountsItemMultiAmountsHint, Instance: mitumcurrency.CreateAccountsItemMultiAmounts{}},
	{Hint: mitumcurrency.CreateAccountsItemSingleAmountHint, Instance: mitumcurrency.CreateAccountsItemSingleAmount{}},
	{Hint: mitumcurrency.CreateAccountsHint, Instance: mitumcurrency.CreateAccounts{}},
	{Hint: mitumcurrency.KeyUpdaterHint, Instance: mitumcurrency.KeyUpdater{}},
	{Hint: mitumcurrency.TransfersItemMultiAmountsHint, Instance: mitumcurrency.TransfersItemMultiAmounts{}},
	{Hint: mitumcurrency.TransfersItemSingleAmountHint, Instance: mitumcurrency.TransfersItemSingleAmount{}},
	{Hint: mitumcurrency.TransfersHint, Instance: mitumcurrency.Transfers{}},
	{Hint: mitumcurrency.SuffrageInflationHint, Instance: mitumcurrency.SuffrageInflation{}},
	{Hint: currencystate.AccountStateValueHint, Instance: currencystate.AccountStateValue{}},
	{Hint: currencystate.BalanceStateValueHint, Instance: currencystate.BalanceStateValue{}},

	{Hint: currencytypes.CurrencyDesignHint, Instance: currencytypes.CurrencyDesign{}},
	{Hint: currencytypes.CurrencyPolicyHint, Instance: currencytypes.CurrencyPolicy{}},
	{Hint: mitumcurrency.CurrencyRegisterHint, Instance: mitumcurrency.CurrencyRegister{}},
	{Hint: mitumcurrency.CurrencyPolicyUpdaterHint, Instance: mitumcurrency.CurrencyPolicyUpdater{}},
	{Hint: currencytypes.ContractAccountKeysHint, Instance: currencytypes.ContractAccountKeys{}},
	{Hint: extensioncurrency.CreateContractAccountsItemMultiAmountsHint, Instance: extensioncurrency.CreateContractAccountsItemMultiAmounts{}},
	{Hint: extensioncurrency.CreateContractAccountsItemSingleAmountHint, Instance: extensioncurrency.CreateContractAccountsItemSingleAmount{}},
	{Hint: extensioncurrency.CreateContractAccountsHint, Instance: extensioncurrency.CreateContractAccounts{}},
	{Hint: extensioncurrency.WithdrawsItemMultiAmountsHint, Instance: extensioncurrency.WithdrawsItemMultiAmounts{}},
	{Hint: extensioncurrency.WithdrawsItemSingleAmountHint, Instance: extensioncurrency.WithdrawsItemSingleAmount{}},
	{Hint: extensioncurrency.WithdrawsHint, Instance: extensioncurrency.Withdraws{}},
	{Hint: mitumcurrency.GenesisCurrenciesFactHint, Instance: mitumcurrency.GenesisCurrenciesFact{}},
	{Hint: mitumcurrency.GenesisCurrenciesHint, Instance: mitumcurrency.GenesisCurrencies{}},
	{Hint: currencytypes.NilFeeerHint, Instance: currencytypes.NilFeeer{}},
	{Hint: currencytypes.FixedFeeerHint, Instance: currencytypes.FixedFeeer{}},
	{Hint: currencytypes.RatioFeeerHint, Instance: currencytypes.RatioFeeer{}},
	{Hint: extensionstate.ContractAccountStateValueHint, Instance: extensionstate.ContractAccountStateValue{}},
	{Hint: currencystate.CurrencyDesignStateValueHint, Instance: currencystate.CurrencyDesignStateValue{}},

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

	{Hint: digestisaac.ManifestHint, Instance: digestisaac.Manifest{}},
	{Hint: digest.AccountValueHint, Instance: digest.AccountValue{}},
	{Hint: digest.OperationValueHint, Instance: digest.OperationValue{}},

	{Hint: isaacoperation.GenesisNetworkPolicyHint, Instance: isaacoperation.GenesisNetworkPolicy{}},
	{Hint: isaacoperation.SuffrageCandidateHint, Instance: isaacoperation.SuffrageCandidate{}},
	{Hint: isaacoperation.SuffrageGenesisJoinHint, Instance: isaacoperation.SuffrageGenesisJoin{}},
	{Hint: isaacoperation.SuffrageDisjoinHint, Instance: isaacoperation.SuffrageDisjoin{}},
	{Hint: isaacoperation.SuffrageJoinHint, Instance: isaacoperation.SuffrageJoin{}},
	{Hint: isaacoperation.NetworkPolicyHint, Instance: isaacoperation.NetworkPolicy{}},
	{Hint: isaacoperation.NetworkPolicyStateValueHint, Instance: isaacoperation.NetworkPolicyStateValue{}},
	{Hint: isaacoperation.FixedSuffrageCandidateLimiterRuleHint, Instance: isaacoperation.FixedSuffrageCandidateLimiterRule{}},
	{Hint: isaacoperation.MajoritySuffrageCandidateLimiterRuleHint, Instance: isaacoperation.MajoritySuffrageCandidateLimiterRule{}},
}

var supportedProposalOperationFactHinters = []encoder.DecodeDetail{
	{Hint: mitumcurrency.CreateAccountsFactHint, Instance: mitumcurrency.CreateAccountsFact{}},
	{Hint: mitumcurrency.KeyUpdaterFactHint, Instance: mitumcurrency.KeyUpdaterFact{}},
	{Hint: mitumcurrency.TransfersFactHint, Instance: mitumcurrency.TransfersFact{}},
	{Hint: mitumcurrency.SuffrageInflationFactHint, Instance: mitumcurrency.SuffrageInflationFact{}},

	{Hint: mitumcurrency.CurrencyRegisterFactHint, Instance: mitumcurrency.CurrencyRegisterFact{}},
	{Hint: mitumcurrency.CurrencyPolicyUpdaterFactHint, Instance: mitumcurrency.CurrencyPolicyUpdaterFact{}},
	{Hint: extensioncurrency.CreateContractAccountsFactHint, Instance: extensioncurrency.CreateContractAccountsFact{}},
	{Hint: extensioncurrency.WithdrawsFactHint, Instance: extensioncurrency.WithdrawsFact{}},

	{Hint: isaacoperation.GenesisNetworkPolicyFactHint, Instance: isaacoperation.GenesisNetworkPolicyFact{}},
	{Hint: isaacoperation.SuffrageCandidateFactHint, Instance: isaacoperation.SuffrageCandidateFact{}},
	{Hint: isaacoperation.SuffrageDisjoinFactHint, Instance: isaacoperation.SuffrageDisjoinFact{}},
	{Hint: isaacoperation.SuffrageJoinFactHint, Instance: isaacoperation.SuffrageJoinFact{}},
	{Hint: isaacoperation.SuffrageGenesisJoinFactHint, Instance: isaacoperation.SuffrageGenesisJoinFact{}},

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
	Hinters = make([]encoder.DecodeDetail, len(launch.Hinters)+len(hinters))
	copy(Hinters, launch.Hinters)
	copy(Hinters[len(launch.Hinters):], hinters)

	SupportedProposalOperationFactHinters = make([]encoder.DecodeDetail, len(launch.SupportedProposalOperationFactHinters)+len(supportedProposalOperationFactHinters))
	copy(SupportedProposalOperationFactHinters, launch.SupportedProposalOperationFactHinters)
	copy(SupportedProposalOperationFactHinters[len(launch.SupportedProposalOperationFactHinters):], supportedProposalOperationFactHinters)
}

func LoadHinters(enc encoder.Encoder) error {
	for i := range Hinters {
		if err := enc.Add(Hinters[i]); err != nil {
			return errors.Wrap(err, "failed to add to encoder")
		}
	}

	for i := range SupportedProposalOperationFactHinters {
		if err := enc.Add(SupportedProposalOperationFactHinters[i]); err != nil {
			return errors.Wrap(err, "failed to add to encoder")
		}
	}

	return nil
}
