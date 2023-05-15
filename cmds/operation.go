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
	TransferSecurityTokensPartition TransferSecurityTokensPartitionCommand             `cmd:"" name:"transfer-security-token" help:"transfer security tokens by partition"`
	RedeemTokens                    RedeemTokensCommand                                `cmd:"" name:"redeem-token" help:"redeem tokens from tokenholder"`
	AuthorizeOperators              AuthorizeOperatorsCommand                          `cmd:"" name:"authorize-operator" help:"authorize operator"`
	RevokeOperators                 RevokeOperatorsCommand                             `cmd:"" name:"revoke-operator" help:"revoke operator"`
	SetDocument                     SetDocumentCommand                                 `cmd:"" name:"set-document" help:"set sto documents"`
	CreateKYCService                CreateKYCServiceCommand                            `cmd:"" name:"create-kyc-service" help:"create kyc service to contract account"`
	AddControllers                  AddControllersCommand                              `cmd:"" name:"add-controllers" help:"add controllers to kyc service"`
	CurrencyRegister                currencycmds.CurrencyRegisterCommand               `cmd:"" name:"currency-register" help:"register new currency"`
	CurrencyPolicyUpdater           currencycmds.CurrencyPolicyUpdaterCommand          `cmd:"" name:"currency-policy-updater" help:"update currency policy"`
	SuffrageInflation               currencycmds.SuffrageInflationCommand              `cmd:"" name:"suffrage-inflation" help:"suffrage inflation operation"`
	SuffrageCandidate               currencycmds.SuffrageCandidateCommand              `cmd:"" name:"suffrage-candidate" help:"suffrage candidate operation"`
	SuffrageJoin                    currencycmds.SuffrageJoinCommand                   `cmd:"" name:"suffrage-join" help:"suffrage join operation"`
	SuffrageDisjoin                 currencycmds.SuffrageDisjoinCommand                `cmd:"" name:"suffrage-disjoin" help:"suffrage disjoin operation"` // revive:disable-line:line-length-limit
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
		SetDocument:                     NewSetDocumentCommand(),
		CreateKYCService:                NewCreateKYCServiceCommand(),
		AddControllers:                  NewAddControllersCommand(),
		CurrencyRegister:                currencycmds.NewCurrencyRegisterCommand(),
		CurrencyPolicyUpdater:           currencycmds.NewCurrencyPolicyUpdaterCommand(),
		SuffrageInflation:               currencycmds.NewSuffrageInflationCommand(),
		SuffrageCandidate:               currencycmds.NewSuffrageCandidateCommand(),
		SuffrageJoin:                    currencycmds.NewSuffrageJoinCommand(),
		SuffrageDisjoin:                 currencycmds.NewSuffrageDisjoinCommand(),
	}
}
