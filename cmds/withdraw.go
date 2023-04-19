package cmds

import (
	"context"

	"github.com/pkg/errors"

	"github.com/ProtoconNet/mitum-currency-extension/v2/currency"

	mitumcurrency "github.com/ProtoconNet/mitum-currency/v2/currency"
	"github.com/ProtoconNet/mitum2/base"
)

type WithdrawCommand struct {
	baseCommand
	OperationFlags
	Sender  AddressFlag          `arg:"" name:"sender" help:"sender address" required:"true"`
	Target  AddressFlag          `arg:"" name:"target" help:"target contract account address" required:"true"`
	Amounts []CurrencyAmountFlag `arg:"" name:"currency-amount" help:"amount (ex: \"<currency>,<amount>\")"`
	sender  base.Address
	target  base.Address
}

func NewWithdrawCommand() WithdrawCommand {
	cmd := NewbaseCommand()
	return WithdrawCommand{
		baseCommand: *cmd,
	}
}

func (cmd *WithdrawCommand) Run(pctx context.Context) error {
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

func (cmd *WithdrawCommand) parseFlags() error {
	if err := cmd.OperationFlags.IsValid(nil); err != nil {
		return err
	}

	if len(cmd.Amounts) < 1 {
		return errors.Errorf("empty currency-amount, must be given at least one")
	}

	if sender, err := cmd.Sender.Encode(enc); err != nil {
		return errors.Wrapf(err, "invalid sender format, %q", cmd.Sender.String())
	} else if target, err := cmd.Target.Encode(enc); err != nil {
		return errors.Wrapf(err, "invalid target format, %q", cmd.Target.String())
	} else {
		cmd.sender = sender
		cmd.target = target
	}

	return nil
}

func (cmd *WithdrawCommand) createOperation() (base.Operation, error) { // nolint:dupl
	var items []currency.WithdrawsItem

	ams := make([]mitumcurrency.Amount, len(cmd.Amounts))
	for i := range cmd.Amounts {
		a := cmd.Amounts[i]
		am := mitumcurrency.NewAmount(a.Big, a.CID)
		if err := am.IsValid(nil); err != nil {
			return nil, err
		}

		ams[i] = am
	}

	item := currency.NewWithdrawsItemMultiAmounts(cmd.target, ams)
	if err := item.IsValid(nil); err != nil {
		return nil, err
	}
	items = append(items, item)

	fact := currency.NewWithdrawsFact([]byte(cmd.Token), cmd.sender, items)

	op, err := currency.NewWithdraws(fact)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create withdraws operation")
	}
	err = op.HashSign(cmd.Privatekey, cmd.NetworkID.NetworkID())
	if err != nil {
		return nil, errors.Wrap(err, "failed to create withdraws operation")
	}

	return op, nil
}
