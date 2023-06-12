package cmds

import (
	"context"

	"github.com/pkg/errors"

	"github.com/ProtoconNet/mitum-sto/operation/sto"
	"github.com/ProtoconNet/mitum2/base"
)

type AuthorizeOperatorsCommand struct {
	baseCommand
	OperationFlags
	Sender    AddressFlag    `arg:"" name:"sender" help:"sender address" required:"true"`
	Contract  AddressFlag    `arg:"" name:"contract" help:"contract account address" required:"true"`
	STO       ContractIDFlag `arg:"" name:"sto-id" help:"sto id" required:"true"`
	Operator  AddressFlag    `arg:"" name:"operator" help:"operator" required:"true"`
	Partition PartitionFlag  `arg:"" name:"partition" help:"default partition" required:"true"`
	Currency  CurrencyIDFlag `arg:"" name:"currency-id" help:"currency id" required:"true"`
	sender    base.Address
	contract  base.Address
	operator  base.Address
}

func NewAuthorizeOperatorsCommand() AuthorizeOperatorsCommand {
	cmd := NewbaseCommand()
	return AuthorizeOperatorsCommand{
		baseCommand: *cmd,
	}
}

func (cmd *AuthorizeOperatorsCommand) Run(pctx context.Context) error {
	if _, err := cmd.prepare(pctx); err != nil {
		return err
	}

	encs = cmd.encs
	enc = cmd.enc

	if err := cmd.parseFlags(); err != nil {
		return err
	}

	op, err := cmd.createOperation()
	if err != nil {
		return err
	}

	PrettyPrint(cmd.Out, op)

	return nil
}

func (cmd *AuthorizeOperatorsCommand) parseFlags() error {
	if err := cmd.OperationFlags.IsValid(nil); err != nil {
		return err
	}

	sender, err := cmd.Sender.Encode(enc)
	if err != nil {
		return errors.Wrapf(err, "invalid sender format, %q", cmd.Sender.String())
	}
	cmd.sender = sender

	contract, err := cmd.Contract.Encode(enc)
	if err != nil {
		return errors.Wrapf(err, "invalid contract account format, %q", cmd.Contract.String())
	}
	cmd.contract = contract

	operator, err := cmd.Operator.Encode(enc)
	if err != nil {
		return errors.Wrapf(err, "invalid operator account format, %q", cmd.Operator.String())
	}
	cmd.operator = operator

	return nil
}

func (cmd *AuthorizeOperatorsCommand) createOperation() (base.Operation, error) { // nolint:dupl
	var items []sto.AuthorizeOperatorsItem

	item := sto.NewAuthorizeOperatorsItem(
		cmd.contract,
		cmd.STO.ID,
		cmd.operator,
		cmd.Partition.Partition,
		cmd.Currency.CID,
	)
	if err := item.IsValid(nil); err != nil {
		return nil, err
	}
	items = append(items, item)

	fact := sto.NewAuthorizeOperatorsFact([]byte(cmd.Token), cmd.sender, items)

	op, err := sto.NewAuthorizeOperators(fact)
	if err != nil {
		return nil, errors.Wrap(err, "failed to authorize operators operation")
	}
	err = op.HashSign(cmd.Privatekey, cmd.NetworkID.NetworkID())
	if err != nil {
		return nil, errors.Wrap(err, "failed to authorize operators operation")
	}

	return op, nil
}
