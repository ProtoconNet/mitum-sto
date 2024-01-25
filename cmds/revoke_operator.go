package cmds

import (
	"context"

	crcycmds "github.com/ProtoconNet/mitum-currency/v3/cmds"
	"github.com/ProtoconNet/mitum-sto/operation/sto"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/pkg/errors"
)

type RevokeOperatorsCommand struct {
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

func (cmd *RevokeOperatorsCommand) Run(pctx context.Context) error {
	if _, err := cmd.prepare(pctx); err != nil {
		return err
	}

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

func (cmd *RevokeOperatorsCommand) parseFlags() error {
	if err := cmd.OperationFlags.IsValid(nil); err != nil {
		return err
	}

	sender, err := cmd.Sender.Encode(cmd.Encoders.JSON())
	if err != nil {
		return errors.Wrapf(err, "invalid sender format, %q", cmd.Sender.String())
	}
	cmd.sender = sender

	contract, err := cmd.Contract.Encode(cmd.Encoders.JSON())
	if err != nil {
		return errors.Wrapf(err, "invalid contract account format, %q", cmd.Contract.String())
	}
	cmd.contract = contract

	operator, err := cmd.Operator.Encode(cmd.Encoders.JSON())
	if err != nil {
		return errors.Wrapf(err, "invalid operator account format, %q", cmd.Operator.String())
	}
	cmd.operator = operator

	return nil
}

func (cmd *RevokeOperatorsCommand) createOperation() (base.Operation, error) { // nolint:dupl
	var items []sto.RevokeOperatorItem

	item := sto.NewRevokeOperatorItem(
		cmd.contract,
		cmd.operator,
		cmd.Partition.Partition,
		cmd.Currency.CID,
	)
	if err := item.IsValid(nil); err != nil {
		return nil, err
	}
	items = append(items, item)

	fact := sto.NewRevokeOperatorFact([]byte(cmd.Token), cmd.sender, items)

	op, err := sto.NewRevokeOperator(fact)
	if err != nil {
		return nil, errors.Wrap(err, "failed to revoke operators operation")
	}
	err = op.HashSign(cmd.Privatekey, cmd.NetworkID.NetworkID())
	if err != nil {
		return nil, errors.Wrap(err, "failed to revoke operators operation")
	}

	return op, nil
}
