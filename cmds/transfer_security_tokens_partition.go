package cmds

import (
	"context"

	"github.com/pkg/errors"

	currencycmds "github.com/ProtoconNet/mitum-currency/v3/cmds"
	"github.com/ProtoconNet/mitum-sto/operation/sto"
	"github.com/ProtoconNet/mitum2/base"
)

type TransferSecurityTokensPartitionCommand struct {
	BaseCommand
	currencycmds.OperationFlags
	Sender      currencycmds.AddressFlag    `arg:"" name:"sender" help:"sender address" required:"true"`
	Contract    currencycmds.AddressFlag    `arg:"" name:"contract" help:"contract account address" required:"true"`
	TokenHolder currencycmds.AddressFlag    `arg:"" name:"tokenholder" help:"tokenholder" required:"true"`
	Receiver    currencycmds.AddressFlag    `arg:"" name:"receiver" help:"token receiver" required:"true"`
	Partition   PartitionFlag               `arg:"" name:"partition" help:"partition" required:"true"`
	Amount      currencycmds.BigFlag        `arg:"" name:"amount" help:"token amount" required:"true"`
	Currency    currencycmds.CurrencyIDFlag `arg:"" name:"currency-id" help:"currency id" required:"true"`
	sender      base.Address
	contract    base.Address
	holder      base.Address
	receiver    base.Address
}

func (cmd *TransferSecurityTokensPartitionCommand) Run(pctx context.Context) error {
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

func (cmd *TransferSecurityTokensPartitionCommand) parseFlags() error {
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

	holder, err := cmd.TokenHolder.Encode(enc)
	if err != nil {
		return errors.Wrapf(err, "invalid tokenholder format, %q", cmd.TokenHolder.String())
	}
	cmd.holder = holder

	receiver, err := cmd.Receiver.Encode(enc)
	if err != nil {
		return errors.Wrapf(err, "invalid receiver format, %q", cmd.Receiver.String())
	}
	cmd.receiver = receiver

	if !cmd.Amount.OverZero() {
		return errors.Wrap(nil, "amount must be over zero")
	}

	return nil
}

func (cmd *TransferSecurityTokensPartitionCommand) createOperation() (base.Operation, error) { // nolint:dupl
	var items []sto.TransferSecurityTokensPartitionItem

	item := sto.NewTransferSecurityTokensPartitionItem(
		cmd.contract,
		cmd.holder,
		cmd.receiver,
		cmd.Partition.Partition,
		cmd.Amount.Big,
		cmd.Currency.CID,
	)

	if err := item.IsValid(nil); err != nil {
		return nil, err
	}
	items = append(items, item)

	fact := sto.NewTransferSecurityTokensPartitionFact([]byte(cmd.Token), cmd.sender, items)

	op, err := sto.NewTransferSecurityTokensPartition(fact)
	if err != nil {
		return nil, errors.Wrap(err, "failed to transfer security tokens partition operation")
	}
	err = op.HashSign(cmd.Privatekey, cmd.NetworkID.NetworkID())
	if err != nil {
		return nil, errors.Wrap(err, "failed to transfer security tokens partition operation")
	}

	return op, nil
}
