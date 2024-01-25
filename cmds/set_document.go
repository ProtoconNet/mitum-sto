package cmds

import (
	"context"

	crcycmds "github.com/ProtoconNet/mitum-currency/v3/cmds"
	"github.com/ProtoconNet/mitum-sto/operation/sto"
	typesto "github.com/ProtoconNet/mitum-sto/types/sto"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/pkg/errors"
)

type SetDocumentCommand struct {
	BaseCommand
	crcycmds.OperationFlags
	Sender       crcycmds.AddressFlag    `arg:"" name:"sender" help:"sender address" required:"true"`
	Contract     crcycmds.AddressFlag    `arg:"" name:"contract" help:"contract address of sto" required:"true"`
	Title        string                  `arg:"" name:"title" help:"sto document title" required:"true"`
	URI          string                  `arg:"" name:"uri" help:"sto document uri" required:"true"`
	DocumentHash string                  `arg:"" name:"document-hash" help:"sto document hash" required:"true"`
	Currency     crcycmds.CurrencyIDFlag `arg:"" name:"currency-id" help:"currency id" required:"true"`
	sender       base.Address
	contract     base.Address
}

func (cmd *SetDocumentCommand) Run(pctx context.Context) error { // nolint:dupl
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

func (cmd *SetDocumentCommand) parseFlags() error {
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

func (cmd *SetDocumentCommand) createOperation() (base.Operation, error) { // nolint:dupl}
	fact := sto.NewSetDocumentFact(
		[]byte(cmd.Token),
		cmd.sender,
		cmd.contract,
		cmd.Title,
		typesto.URI(cmd.URI),
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
