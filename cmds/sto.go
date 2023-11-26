package cmds

type STOCommand struct {
	CreateSecurityTokens            CreateSecurityTokensCommand            `cmd:"" name:"create-security-token" help:"create security token in contract account"`
	IssueSecurityTokens             IssueSecurityTokensCommand             `cmd:"" name:"issue-security-token" help:"issue security token in partition"`
	TransferSecurityTokensPartition TransferSecurityTokensPartitionCommand `cmd:"" name:"transfer-security-token" help:"transfer security tokens by partition"`
	RedeemTokens                    RedeemTokensCommand                    `cmd:"" name:"redeem-token" help:"redeem tokens from token holder"`
	AuthorizeOperators              AuthorizeOperatorsCommand              `cmd:"" name:"authorize-operator" help:"authorize operator"`
	RevokeOperators                 RevokeOperatorsCommand                 `cmd:"" name:"revoke-operator" help:"revoke operator"`
	SetDocument                     SetDocumentCommand                     `cmd:"" name:"set-document" help:"set sto documents"`
}
