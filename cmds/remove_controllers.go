package cmds

import (
	"context"

	"github.com/pkg/errors"

	currencycmds "github.com/ProtoconNet/mitum-currency/v3/cmds"
	"github.com/ProtoconNet/mitum-sto/operation/kyc"
	"github.com/ProtoconNet/mitum2/base"
)

type RemoveControllersCommand struct {
	BaseCommand
	currencycmds.OperationFlags
	Sender     currencycmds.AddressFlag    `arg:"" name:"sender" help:"sender address" required:"true"`
	Contract   currencycmds.AddressFlag    `arg:"" name:"contract" help:"contract account address" required:"true"`
	Controller currencycmds.AddressFlag    `arg:"" name:"controller" help:"controller" required:"true"`
	Currency   currencycmds.CurrencyIDFlag `arg:"" name:"currency-id" help:"currency id" required:"true"`
	sender     base.Address
	contract   base.Address
	controller base.Address
}

func (cmd *RemoveControllersCommand) Run(pctx context.Context) error {
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

func (cmd *RemoveControllersCommand) parseFlags() error {
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

	controller, err := cmd.Controller.Encode(enc)
	if err != nil {
		return errors.Wrapf(err, "invalid controller account format, %q", cmd.Controller.String())
	}
	cmd.controller = controller

	return nil
}

func (cmd *RemoveControllersCommand) createOperation() (base.Operation, error) { // nolint:dupl
	var items []kyc.RemoveControllersItem

	item := kyc.NewRemoveControllersItem(
		cmd.contract,
		cmd.controller,
		cmd.Currency.CID,
	)
	if err := item.IsValid(nil); err != nil {
		return nil, err
	}
	items = append(items, item)

	fact := kyc.NewRemoveControllersFact([]byte(cmd.Token), cmd.sender, items)

	op, err := kyc.NewRemoveControllers(fact)
	if err != nil {
		return nil, errors.Wrap(err, "failed to remove controllers operation")
	}
	err = op.HashSign(cmd.Privatekey, cmd.NetworkID.NetworkID())
	if err != nil {
		return nil, errors.Wrap(err, "failed to remove controllers operation")
	}

	return op, nil
}
