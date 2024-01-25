package cmds

import (
	"context"

	"github.com/pkg/errors"

	crcycmds "github.com/ProtoconNet/mitum-currency/v3/cmds"
	"github.com/ProtoconNet/mitum-sto/operation/sto"
	"github.com/ProtoconNet/mitum2/base"
)

type CreateSecurityTokensCommand struct {
	BaseCommand
	crcycmds.OperationFlags
	Sender      crcycmds.AddressFlag    `arg:"" name:"sender" help:"sender address" required:"true"`
	Contract    crcycmds.AddressFlag    `arg:"" name:"contract" help:"contract account address" required:"true"`
	Granularity uint64                  `arg:"" name:"granularity" help:"granularity" required:"true"`
	Partition   PartitionFlag           `arg:"" name:"default-partition" help:"default partition" required:"true"`
	Currency    crcycmds.CurrencyIDFlag `arg:"" name:"currency-id" help:"currency id" required:"true"`
	sender      base.Address
	contract    base.Address
}

func (cmd *CreateSecurityTokensCommand) Run(pctx context.Context) error {
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

func (cmd *CreateSecurityTokensCommand) parseFlags() error {
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

	return nil
}

func (cmd *CreateSecurityTokensCommand) createOperation() (base.Operation, error) { // nolint:dupl
	var items []sto.CreateSecurityTokenItem

	item := sto.NewCreateSecurityTokenItem(
		cmd.contract,
		cmd.Granularity,
		cmd.Partition.Partition,
		//cmd.controllers,
		cmd.Currency.CID,
	)
	if err := item.IsValid(nil); err != nil {
		return nil, err
	}
	items = append(items, item)

	fact := sto.NewCreateSecurityTokenFact([]byte(cmd.Token), cmd.sender, items)

	op, err := sto.NewCreateSecurityToken(fact)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create security tokens operation")
	}
	err = op.HashSign(cmd.Privatekey, cmd.NetworkID.NetworkID())
	if err != nil {
		return nil, errors.Wrap(err, "failed to create security tokens operation")
	}

	return op, nil
}
