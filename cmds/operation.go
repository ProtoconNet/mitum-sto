package cmds

import (
	extensioncurrencycmds "github.com/ProtoconNet/mitum-currency-extension/v2/cmds"
	currencycmds "github.com/ProtoconNet/mitum-currency/v2/cmds"
)

type OperationCommand struct {
	CreateAccount                   currencycmds.CreateAccountCommand                  `cmd:"" name:"create-account" help:"create new account"`
	KeyUpdater                      currencycmds.KeyUpdaterCommand                     `cmd:"" name:"key-updater" help:"update account keys"`
	Transfer                        currencycmds.TransferCommand                       `cmd:"" name:"transfer" help:"transfer amounts to receiver"`
	CreateContractAccount           extensioncurrencycmds.CreateContractAccountCommand `cmd:"" name:"create-contract-account" help:"create new contract account"`
	Withdraw                        extensioncurrencycmds.WithdrawCommand              `cmd:"" name:"withdraw" help:"withdraw amounts from target contract account"`
	CreateSecurityTokens            CreateSecurityTokensCommand                        `cmd:"" name:"create-security-token" help:"create security token in contract account"`
	IssueSecurityTokens             IssueSecurityTokensCommand                         `cmd:"" name:"issue-security-token" help:"issue security token in partition"`
	TransferSecurityTokensPartition TransferSecurityTokensPartitionCommand             `cmd:"transfer-security-token" help:"transfer security tokens by partition"`
	RedeemTokens                    RedeemTokensCommand                                `cmd:"redeem-tokens" help:"redeem tokens from tokenholder"`
	AuthorizeOperators              AuthorizeOperatorsCommand                          `cmd:"" name:"authorize-operator" help:"authorize operator"`
	RevokeOperators                 RevokeOperatorsCommand                             `cmd:"" name:"revoke-operator" help:"revoke operator"`
	SetDocuments                    SetDocumentsCommand                                `cmd:"" name:"set-documents" help:"set sto documents"`
	CurrencyRegister                currencycmds.CurrencyRegisterCommand               `cmd:"" name:"currency-register" help:"register new currency"`
	CurrencyPolicyUpdater           currencycmds.CurrencyPolicyUpdaterCommand          `cmd:"" name:"currency-policy-updater" help:"update currency policy"`
	SuffrageInflation               SuffrageInflationCommand                           `cmd:"" name:"suffrage-inflation" help:"suffrage inflation operation"`
	SuffrageCandidate               SuffrageCandidateCommand                           `cmd:"" name:"suffrage-candidate" help:"suffrage candidate operation"`
	SuffrageJoin                    SuffrageJoinCommand                                `cmd:"" name:"suffrage-join" help:"suffrage join operation"`
	SuffrageDisjoin                 SuffrageDisjoinCommand                             `cmd:"" name:"suffrage-disjoin" help:"suffrage disjoin operation"` // revive:disable-line:line-length-limit
}

func NewOperationCommand() OperationCommand {
	return OperationCommand{
		CreateAccount:                   currencycmds.NewCreateAccountCommand(),
		KeyUpdater:                      currencycmds.NewKeyUpdaterCommand(),
		Transfer:                        currencycmds.NewTransferCommand(),
		CreateContractAccount:           extensioncurrencycmds.NewCreateContractAccountCommand(),
		Withdraw:                        extensioncurrencycmds.NewWithdrawCommand(),
		CreateSecurityTokens:            NewCreateSecurityTokensCommand(),
		IssueSecurityTokens:             NewIssueSecurityTokensCommand(),
		TransferSecurityTokensPartition: NewTransferSecurityTokensPartitionCommand(),
		RedeemTokens:                    NewRedeemTokensCommand(),
		AuthorizeOperators:              NewAuthorizeOperatorsCommand(),
		RevokeOperators:                 NewRevokeOperatorsCommand(),
		SetDocuments:                    NewSetDocumentsCommand(),
		CurrencyRegister:                currencycmds.NewCurrencyRegisterCommand(),
		CurrencyPolicyUpdater:           currencycmds.NewCurrencyPolicyUpdaterCommand(),
		SuffrageInflation:               NewSuffrageInflationCommand(),
		SuffrageCandidate:               NewSuffrageCandidateCommand(),
		SuffrageJoin:                    NewSuffrageJoinCommand(),
		SuffrageDisjoin:                 NewSuffrageDisjoinCommand(),
	}
}
