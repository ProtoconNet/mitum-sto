package cmds

import (
	"context"

	crcycmds "github.com/ProtoconNet/mitum-currency/v3/cmds"
	"github.com/ProtoconNet/mitum-sto/operation/sto"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/pkg/errors"
)

type RedeemTokensCommand struct {
	BaseCommand
	crcycmds.OperationFlags
	Sender      crcycmds.AddressFlag    `arg:"" name:"sender" help:"sender address" required:"true"`
	Contract    crcycmds.AddressFlag    `arg:"" name:"contract" help:"contract account address" required:"true"`
	TokenHolder crcycmds.AddressFlag    `arg:"" name:"token holder" help:"token holder" required:"true"`
	Amount      crcycmds.BigFlag        `arg:"" name:"amount" help:"token amount" required:"true"`
	Partition   PartitionFlag           `arg:"" name:"partition" help:"partition" required:"true"`
	Currency    crcycmds.CurrencyIDFlag `arg:"" name:"currency-id" help:"currency id" required:"true"`
	sender      base.Address
	contract    base.Address
	holder      base.Address
}

func (cmd *RedeemTokensCommand) Run(pctx context.Context) error {
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

	crcycmds.PrettyPrint(cmd.Out, op)

	return nil
}

func (cmd *RedeemTokensCommand) parseFlags() error {
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
		return errors.Wrapf(err, "invalid token holder format, %q", cmd.TokenHolder.String())
	}
	cmd.holder = holder

	if !cmd.Amount.OverZero() {
		return errors.Wrap(nil, "amount must be over zero")
	}

	return nil
}

func (cmd *RedeemTokensCommand) createOperation() (base.Operation, error) { // nolint:dupl
	var items []sto.RedeemItem

	item := sto.NewRedeemItem(
		cmd.contract,
		cmd.holder,
		cmd.Amount.Big,
		cmd.Partition.Partition,
		cmd.Currency.CID,
	)

	if err := item.IsValid(nil); err != nil {
		return nil, err
	}
	items = append(items, item)

	fact := sto.NewRedeemFact([]byte(cmd.Token), cmd.sender, items)

	op, err := sto.NewRedeem(fact)
	if err != nil {
		return nil, errors.Wrap(err, "failed to redeem tokens operation")
	}
	err = op.HashSign(cmd.Privatekey, cmd.NetworkID.NetworkID())
	if err != nil {
		return nil, errors.Wrap(err, "failed to redeem tokens operation")
	}

	return op, nil
}
