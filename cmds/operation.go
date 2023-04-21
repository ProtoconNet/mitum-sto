package cmds

type OperationCommand struct {
	CreateAccount         CreateAccountCommand         `cmd:"" name:"create-account" help:"create new account"`
	KeyUpdater            KeyUpdaterCommand            `cmd:"" name:"key-updater" help:"update account keys"`
	Transfer              TransferCommand              `cmd:"" name:"transfer" help:"transfer amounts to receiver"`
	CreateContractAccount CreateContractAccountCommand `cmd:"" name:"create-contract-account" help:"create new contract account"`
	Withdraw              WithdrawCommand              `cmd:"" name:"withdraw" help:"withdraw amounts from target contract account"`
	CreateSecurityTokens  CreateSecurityTokensCommand  `cmd:"" name:"create-security-token" help:"create security token in contract account"`
	CurrencyRegister      CurrencyRegisterCommand      `cmd:"" name:"currency-register" help:"register new currency"`
	CurrencyPolicyUpdater CurrencyPolicyUpdaterCommand `cmd:"" name:"currency-policy-updater" help:"update currency policy"`
	SuffrageInflation     SuffrageInflationCommand     `cmd:"" name:"suffrage-inflation" help:"suffrage inflation operation"`
	SuffrageCandidate     SuffrageCandidateCommand     `cmd:"" name:"suffrage-candidate" help:"suffrage candidate operation"`
	SuffrageJoin          SuffrageJoinCommand          `cmd:"" name:"suffrage-join" help:"suffrage join operation"`
	SuffrageDisjoin       SuffrageDisjoinCommand       `cmd:"" name:"suffrage-disjoin" help:"suffrage disjoin operation"` // revive:disable-line:line-length-limit
}

func NewOperationCommand() OperationCommand {
	return OperationCommand{
		CreateAccount:         NewCreateAccountCommand(),
		KeyUpdater:            NewKeyUpdaterCommand(),
		Transfer:              NewTransferCommand(),
		CreateContractAccount: NewCreateContractAccountCommand(),
		Withdraw:              NewWithdrawCommand(),
		CreateSecurityTokens:  NewCreateSecurityTokensCommand(),
		CurrencyRegister:      NewCurrencyRegisterCommand(),
		CurrencyPolicyUpdater: NewCurrencyPolicyUpdaterCommand(),
		SuffrageInflation:     NewSuffrageInflationCommand(),
		SuffrageCandidate:     NewSuffrageCandidateCommand(),
		SuffrageJoin:          NewSuffrageJoinCommand(),
		SuffrageDisjoin:       NewSuffrageDisjoinCommand(),
	}
}
