package cmds

import (
	"context"

	"github.com/pkg/errors"

	"github.com/ProtoconNet/mitum-sto/sto"
	"github.com/ProtoconNet/mitum2/base"
)

type SetDocumentCommand struct {
	baseCommand
	OperationFlags
	Sender       AddressFlag    `arg:"" name:"sender" help:"sender address" required:"true"`
	Contract     AddressFlag    `arg:"" name:"contract" help:"contract address of sto" required:"true"`
	STO          ContractIDFlag `arg:"" name:"sto-id" help:"sto id" required:"true"`
	Title        string         `arg:"" name:"title" help:"sto document title" required:"true"`
	URI          string         `arg:"" name:"uri" help:"sto document uri" required:"true"`
	DocumentHash string         `arg:"" name:"document-hash" help:"sto document hash" required:"true"`
	Currency     CurrencyIDFlag `arg:"" name:"currency-id" help:"currency id" required:"true"`
	sender       base.Address
	contract     base.Address
}

func NewSetDocumentCommand() SetDocumentCommand {
	cmd := NewbaseCommand()
	return SetDocumentCommand{
		baseCommand: *cmd,
	}
}

func (cmd *SetDocumentCommand) Run(pctx context.Context) error { // nolint:dupl
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
	fact := sto.NewSetDocumentFact([]byte(cmd.Token), cmd.sender, cmd.contract, cmd.STO.ID, cmd.Title, sto.URI(cmd.URI), cmd.DocumentHash, cmd.Currency.CID)

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
