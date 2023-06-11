package cmds

import (
	"context"

	"github.com/pkg/errors"

	"github.com/ProtoconNet/mitum-sto/operation/kyc"
	"github.com/ProtoconNet/mitum2/base"
)

type AddControllersCommand struct {
	baseCommand
	OperationFlags
	Sender     AddressFlag    `arg:"" name:"sender" help:"sender address" required:"true"`
	Contract   AddressFlag    `arg:"" name:"contract" help:"contract account address" required:"true"`
	KYC        ContractIDFlag `arg:"" name:"kyc-id" help:"kyc id" required:"true"`
	Controller AddressFlag    `arg:"" name:"controller" help:"controller" required:"true"`
	Currency   CurrencyIDFlag `arg:"" name:"currency-id" help:"currency id" required:"true"`
	sender     base.Address
	contract   base.Address
	controller base.Address
}

func NewAddControllersCommand() AddControllersCommand {
	cmd := NewbaseCommand()
	return AddControllersCommand{
		baseCommand: *cmd,
	}
}

func (cmd *AddControllersCommand) Run(pctx context.Context) error {
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

func (cmd *AddControllersCommand) parseFlags() error {
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

func (cmd *AddControllersCommand) createOperation() (base.Operation, error) { // nolint:dupl
	var items []kyc.AddControllersItem

	item := kyc.NewAddControllersItem(
		cmd.contract,
		cmd.KYC.ID,
		cmd.controller,
		cmd.Currency.CID,
	)
	if err := item.IsValid(nil); err != nil {
		return nil, err
	}
	items = append(items, item)

	fact := kyc.NewAddControllersFact([]byte(cmd.Token), cmd.sender, items)

	op, err := kyc.NewAddControllers(fact)
	if err != nil {
		return nil, errors.Wrap(err, "failed to add controllers operation")
	}
	err = op.HashSign(cmd.Privatekey, cmd.NetworkID.NetworkID())
	if err != nil {
		return nil, errors.Wrap(err, "failed to add controllers operation")
	}

	return op, nil
}
