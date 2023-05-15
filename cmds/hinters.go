package cmds

import (
	"github.com/ProtoconNet/mitum-currency-extension/v2/currency"
	"github.com/ProtoconNet/mitum-currency-extension/v2/digest"
	mitumcurrency "github.com/ProtoconNet/mitum-currency/v2/currency"
	digestisaac "github.com/ProtoconNet/mitum-currency/v2/digest/isaac"
	isaacoperation "github.com/ProtoconNet/mitum-currency/v2/isaac"
	"github.com/ProtoconNet/mitum-sto/kyc"
	"github.com/ProtoconNet/mitum-sto/sto"
	"github.com/ProtoconNet/mitum2/launch"
	"github.com/ProtoconNet/mitum2/util/encoder"
	"github.com/pkg/errors"
)

var Hinters []encoder.DecodeDetail
var SupportedProposalOperationFactHinters []encoder.DecodeDetail

var hinters = []encoder.DecodeDetail{
	// revive:disable-next-line:line-length-limit
	{Hint: mitumcurrency.BaseStateHint, Instance: mitumcurrency.BaseState{}},
	{Hint: mitumcurrency.NodeHint, Instance: mitumcurrency.BaseNode{}},
	{Hint: mitumcurrency.AccountHint, Instance: mitumcurrency.Account{}},
	{Hint: mitumcurrency.AddressHint, Instance: mitumcurrency.Address{}},
	{Hint: mitumcurrency.AmountHint, Instance: mitumcurrency.Amount{}},
	{Hint: mitumcurrency.AccountKeysHint, Instance: mitumcurrency.BaseAccountKeys{}},
	{Hint: mitumcurrency.AccountKeyHint, Instance: mitumcurrency.BaseAccountKey{}},
	{Hint: mitumcurrency.CreateAccountsItemMultiAmountsHint, Instance: mitumcurrency.CreateAccountsItemMultiAmounts{}},
	{Hint: mitumcurrency.CreateAccountsItemSingleAmountHint, Instance: mitumcurrency.CreateAccountsItemSingleAmount{}},
	{Hint: mitumcurrency.CreateAccountsHint, Instance: mitumcurrency.CreateAccounts{}},
	{Hint: mitumcurrency.KeyUpdaterHint, Instance: mitumcurrency.KeyUpdater{}},
	{Hint: mitumcurrency.TransfersItemMultiAmountsHint, Instance: mitumcurrency.TransfersItemMultiAmounts{}},
	{Hint: mitumcurrency.TransfersItemSingleAmountHint, Instance: mitumcurrency.TransfersItemSingleAmount{}},
	{Hint: mitumcurrency.TransfersHint, Instance: mitumcurrency.Transfers{}},
	{Hint: mitumcurrency.SuffrageInflationHint, Instance: mitumcurrency.SuffrageInflation{}},
	{Hint: mitumcurrency.AccountStateValueHint, Instance: mitumcurrency.AccountStateValue{}},
	{Hint: mitumcurrency.BalanceStateValueHint, Instance: mitumcurrency.BalanceStateValue{}},

	{Hint: currency.CurrencyDesignHint, Instance: currency.CurrencyDesign{}},
	{Hint: currency.CurrencyPolicyHint, Instance: currency.CurrencyPolicy{}},
	{Hint: currency.CurrencyRegisterHint, Instance: currency.CurrencyRegister{}},
	{Hint: currency.CurrencyPolicyUpdaterHint, Instance: currency.CurrencyPolicyUpdater{}},
	{Hint: currency.ContractAccountKeysHint, Instance: currency.ContractAccountKeys{}},
	{Hint: currency.CreateContractAccountsItemMultiAmountsHint, Instance: currency.CreateContractAccountsItemMultiAmounts{}},
	{Hint: currency.CreateContractAccountsItemSingleAmountHint, Instance: currency.CreateContractAccountsItemSingleAmount{}},
	{Hint: currency.CreateContractAccountsHint, Instance: currency.CreateContractAccounts{}},
	{Hint: currency.WithdrawsItemMultiAmountsHint, Instance: currency.WithdrawsItemMultiAmounts{}},
	{Hint: currency.WithdrawsItemSingleAmountHint, Instance: currency.WithdrawsItemSingleAmount{}},
	{Hint: currency.WithdrawsHint, Instance: currency.Withdraws{}},
	{Hint: currency.GenesisCurrenciesFactHint, Instance: currency.GenesisCurrenciesFact{}},
	{Hint: currency.GenesisCurrenciesHint, Instance: currency.GenesisCurrencies{}},
	{Hint: currency.NilFeeerHint, Instance: currency.NilFeeer{}},
	{Hint: currency.FixedFeeerHint, Instance: currency.FixedFeeer{}},
	{Hint: currency.RatioFeeerHint, Instance: currency.RatioFeeer{}},
	{Hint: currency.ContractAccountStateValueHint, Instance: currency.ContractAccountStateValue{}},
	{Hint: currency.CurrencyDesignStateValueHint, Instance: currency.CurrencyDesignStateValue{}},

	{Hint: sto.DesignStateValueHint, Instance: sto.STODesignStateValue{}},
	{Hint: sto.TokenHolderPartitionsStateValueHint, Instance: sto.TokenHolderPartitionsStateValue{}},
	{Hint: sto.TokenHolderPartitionBalanceStateValueHint, Instance: sto.TokenHolderPartitionBalanceStateValue{}},
	{Hint: sto.TokenHolderPartitionOperatorsStateValueHint, Instance: sto.TokenHolderPartitionOperatorsStateValue{}},
	{Hint: sto.PartitionBalanceStateValueHint, Instance: sto.PartitionBalanceStateValue{}},
	{Hint: sto.OperatorTokenHoldersStateValueHint, Instance: sto.OperatorTokenHoldersStateValue{}},
	{Hint: sto.STODesignHint, Instance: sto.STODesign{}},
	{Hint: sto.DocumentHint, Instance: sto.Document{}},
	{Hint: sto.STOPolicyHint, Instance: sto.STOPolicy{}},
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
	{Hint: kyc.CreateKYCServiceHint, Instance: kyc.CreateKYCService{}},
	{Hint: kyc.AddControllersItemHint, Instance: kyc.AddControllersItem{}},
	{Hint: kyc.AddControllersHint, Instance: kyc.AddControllers{}},

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

	{Hint: currency.CurrencyRegisterFactHint, Instance: currency.CurrencyRegisterFact{}},
	{Hint: currency.CurrencyPolicyUpdaterFactHint, Instance: currency.CurrencyPolicyUpdaterFact{}},
	{Hint: currency.CreateContractAccountsFactHint, Instance: currency.CreateContractAccountsFact{}},
	{Hint: currency.WithdrawsFactHint, Instance: currency.WithdrawsFact{}},

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
