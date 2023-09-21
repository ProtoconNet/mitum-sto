package cmds

import (
	"context"

	"github.com/pkg/errors"

	currencycmds "github.com/ProtoconNet/mitum-currency/v3/cmds"
	"github.com/ProtoconNet/mitum-sto/operation/sto"
	stotypes "github.com/ProtoconNet/mitum-sto/types/sto"
	"github.com/ProtoconNet/mitum2/base"
)

type SetDocumentCommand struct {
	BaseCommand
	currencycmds.OperationFlags
	Sender       currencycmds.AddressFlag    `arg:"" name:"sender" help:"sender address" required:"true"`
	Contract     currencycmds.AddressFlag    `arg:"" name:"contract" help:"contract address of sto" required:"true"`
	Title        string                      `arg:"" name:"title" help:"sto document title" required:"true"`
	URI          string                      `arg:"" name:"uri" help:"sto document uri" required:"true"`
	DocumentHash string                      `arg:"" name:"document-hash" help:"sto document hash" required:"true"`
	Currency     currencycmds.CurrencyIDFlag `arg:"" name:"currency-id" help:"currency id" required:"true"`
	sender       base.Address
	contract     base.Address
}

func (cmd *SetDocumentCommand) Run(pctx context.Context) error { // nolint:dupl
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

func (cmd *SetDocumentCommand) parseFlags() error {
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

	return nil
}

func (cmd *SetDocumentCommand) createOperation() (base.Operation, error) { // nolint:dupl}
	fact := sto.NewSetDocumentFact(
		[]byte(cmd.Token),
		cmd.sender,
		cmd.contract,
		cmd.Title,
		stotypes.URI(cmd.URI),
		cmd.DocumentHash,
		cmd.Currency.CID,
	)

	op, err := sto.NewSetDocument(fact)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create set-document operation")
	}
	err = op.HashSign(cmd.Privatekey, cmd.NetworkID.NetworkID())
	if err != nil {
		return nil, errors.Wrap(err, "failed to create set-document operation")
	}

	return op, nil
}
