package cmds

import (
	"context"

	"github.com/pkg/errors"

	currencycmds "github.com/ProtoconNet/mitum-currency/v3/cmds"
	"github.com/ProtoconNet/mitum-sto/operation/kyc"
	"github.com/ProtoconNet/mitum2/base"
)

type CreateKYCServiceCommand struct {
	BaseCommand
	currencycmds.OperationFlags
	Sender      currencycmds.AddressFlag    `arg:"" name:"sender" help:"sender address" required:"true"`
	Contract    currencycmds.AddressFlag    `arg:"" name:"contract" help:"contract address of kyc" required:"true"`
	Controller  currencycmds.AddressFlag    `arg:"" name:"controller" help:"controller" required:"true"`
	Currency    currencycmds.CurrencyIDFlag `arg:"" name:"currency-id" help:"currency id" required:"true"`
	sender      base.Address
	contract    base.Address
	controllers []base.Address
}

func (cmd *CreateKYCServiceCommand) Run(pctx context.Context) error { // nolint:dupl
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

func (cmd *CreateKYCServiceCommand) parseFlags() error {
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
	cmd.controllers = []base.Address{controller}

	return nil
}

func (cmd *CreateKYCServiceCommand) createOperation() (base.Operation, error) { // nolint:dupl}
	fact := kyc.NewCreateServiceFact(
		[]byte(cmd.Token),
		cmd.sender,
		cmd.contract,
		cmd.controllers,
		cmd.Currency.CID,
	)

	op, err := kyc.NewCreateService(fact)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create create-kyc-service operation")
	}
	err = op.HashSign(cmd.Privatekey, cmd.NetworkID.NetworkID())
	if err != nil {
		return nil, errors.Wrap(err, "failed to create create-kyc-service operation")
	}

	return op, nil
}
