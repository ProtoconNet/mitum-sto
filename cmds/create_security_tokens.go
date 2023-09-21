package cmds

import (
	"context"

	"github.com/pkg/errors"

	currencycmds "github.com/ProtoconNet/mitum-currency/v3/cmds"
	"github.com/ProtoconNet/mitum-sto/operation/sto"
	"github.com/ProtoconNet/mitum2/base"
)

type CreateSecurityTokensCommand struct {
	BaseCommand
	currencycmds.OperationFlags
	Sender      currencycmds.AddressFlag    `arg:"" name:"sender" help:"sender address" required:"true"`
	Contract    currencycmds.AddressFlag    `arg:"" name:"contract" help:"contract account address" required:"true"`
	Granularity uint64                      `arg:"" name:"granularity" help:"granularity" required:"true"`
	Partition   PartitionFlag               `arg:"" name:"default-partition" help:"default partition" required:"true"`
	Currency    currencycmds.CurrencyIDFlag `arg:"" name:"currency-id" help:"currency id" required:"true"`
	Controller  currencycmds.AddressFlag    `name:"controller" help:"controller"`
	sender      base.Address
	contract    base.Address
	controllers []base.Address
}

func (cmd *CreateSecurityTokensCommand) Run(pctx context.Context) error {
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

	currencycmds.PrettyPrint(cmd.Out, op)

	return nil
}

func (cmd *CreateSecurityTokensCommand) parseFlags() error {
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

	if cmd.Controller.String() != "" {
		controller, err := cmd.Controller.Encode(enc)
		if err != nil {
			return errors.Wrapf(err, "invalid controller format, %q", controller)
		}
		cmd.controllers = []base.Address{controller}
	}

	return nil
}

func (cmd *CreateSecurityTokensCommand) createOperation() (base.Operation, error) { // nolint:dupl
	var items []sto.CreateSecurityTokensItem

	item := sto.NewCreateSecurityTokensItem(
		cmd.contract,
		cmd.Granularity,
		cmd.Partition.Partition,
		cmd.controllers,
		cmd.Currency.CID,
	)
	if err := item.IsValid(nil); err != nil {
		return nil, err
	}
	items = append(items, item)

	fact := sto.NewCreateSecurityTokensFact([]byte(cmd.Token), cmd.sender, items)

	op, err := sto.NewCreateSecurityTokens(fact)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create security tokens operation")
	}
	err = op.HashSign(cmd.Privatekey, cmd.NetworkID.NetworkID())
	if err != nil {
		return nil, errors.Wrap(err, "failed to create security tokens operation")
	}

	return op, nil
}
