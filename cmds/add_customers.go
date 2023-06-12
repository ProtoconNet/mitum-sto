package cmds

import (
	"context"

	"github.com/pkg/errors"

	"github.com/ProtoconNet/mitum-sto/operation/kyc"
	"github.com/ProtoconNet/mitum2/base"
)

type AddCustomersCommand struct {
	baseCommand
	OperationFlags
	Sender   AddressFlag    `arg:"" name:"sender" help:"sender address" required:"true"`
	Contract AddressFlag    `arg:"" name:"contract" help:"contract account address" required:"true"`
	KYC      ContractIDFlag `arg:"" name:"kyc-id" help:"kyc id" required:"true"`
	Customer AddressFlag    `arg:"" name:"customer" help:"customer" required:"true"`
	Status   bool           `arg:"" name:"status" help:"customer status" required:"true"`
	Currency CurrencyIDFlag `arg:"" name:"currency-id" help:"currency id" required:"true"`
	sender   base.Address
	contract base.Address
	customer base.Address
}

func NewAddCustomersCommand() AddCustomersCommand {
	cmd := NewbaseCommand()
	return AddCustomersCommand{
		baseCommand: *cmd,
	}
}

func (cmd *AddCustomersCommand) Run(pctx context.Context) error {
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

func (cmd *AddCustomersCommand) parseFlags() error {
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

	customer, err := cmd.Customer.Encode(enc)
	if err != nil {
		return errors.Wrapf(err, "invalid customer account format, %q", cmd.Customer.String())
	}
	cmd.customer = customer

	return nil
}

func (cmd *AddCustomersCommand) createOperation() (base.Operation, error) { // nolint:dupl
	var items []kyc.AddCustomersItem

	item := kyc.NewAddCustomersItem(
		cmd.contract,
		cmd.KYC.ID,
		cmd.customer,
		cmd.Status,
		cmd.Currency.CID,
	)
	if err := item.IsValid(nil); err != nil {
		return nil, err
	}
	items = append(items, item)

	fact := kyc.NewAddCustomersFact([]byte(cmd.Token), cmd.sender, items)

	op, err := kyc.NewAddCustomers(fact)
	if err != nil {
		return nil, errors.Wrap(err, "failed to add customers operation")
	}
	err = op.HashSign(cmd.Privatekey, cmd.NetworkID.NetworkID())
	if err != nil {
		return nil, errors.Wrap(err, "failed to add customers operation")
	}

	return op, nil
}
