package cmds

import (
	"context"

	"github.com/pkg/errors"

	crcycmds "github.com/ProtoconNet/mitum-currency/v3/cmds"
	"github.com/ProtoconNet/mitum-sto/operation/sto"
	"github.com/ProtoconNet/mitum2/base"
)

type AuthorizeOperatorsCommand struct {
	BaseCommand
	crcycmds.OperationFlags
	Sender    crcycmds.AddressFlag    `arg:"" name:"sender" help:"sender address" required:"true"`
	Contract  crcycmds.AddressFlag    `arg:"" name:"contract" help:"contract account address" required:"true"`
	Operator  crcycmds.AddressFlag    `arg:"" name:"operator" help:"operator" required:"true"`
	Partition PartitionFlag           `arg:"" name:"partition" help:"default partition" required:"true"`
	Currency  crcycmds.CurrencyIDFlag `arg:"" name:"currency-id" help:"currency id" required:"true"`
	sender    base.Address
	contract  base.Address
	operator  base.Address
}

func (cmd *AuthorizeOperatorsCommand) Run(pctx context.Context) error {
	if _, err := cmd.prepare(pctx); err != nil {
		return err
	}

	encs = cmd.Encoders
	enc = cmd.Encoder

	if err := cmd.parseFlags(); err != nil {
		return err
	}

	op, err := cmd.createOperation()
	if err != nil {
		return err
	}

	crcycmds.PrettyPrint(cmd.Out, op)

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
	var items []sto.AuthorizeOperatorItem

	item := sto.NewAuthorizeOperatorItem(
		cmd.contract,
		cmd.operator,
		cmd.Partition.Partition,
		cmd.Currency.CID,
	)
	if err := item.IsValid(nil); err != nil {
		return nil, err
	}
	items = append(items, item)

	fact := sto.NewAuthorizeOperatorFact([]byte(cmd.Token), cmd.sender, items)

	op, err := sto.NewAuthorizeOperator(fact)
	if err != nil {
		return nil, errors.Wrap(err, "failed to authorize operators operation")
	}
	err = op.HashSign(cmd.Privatekey, cmd.NetworkID.NetworkID())
	if err != nil {
		return nil, errors.Wrap(err, "failed to authorize operators operation")
	}

	return op, nil
}
