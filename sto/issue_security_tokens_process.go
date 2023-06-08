package sto

import (
	"context"
	"math/big"
	"sync"

	currencybase "github.com/ProtoconNet/mitum-currency/v3/base"
	currencyoperation "github.com/ProtoconNet/mitum-currency/v3/operation/currency"
	types "github.com/ProtoconNet/mitum-currency/v3/operation/type"
	currency "github.com/ProtoconNet/mitum-currency/v3/state/currency"
	extensioncurrency "github.com/ProtoconNet/mitum-currency/v3/state/extension"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/pkg/errors"
)

var issueSecurityTokensItemProcessorPool = sync.Pool{
	New: func() interface{} {
		return new(IssueSecurityTokensItemProcessor)
	},
}

var issueSecurityTokensProcessorPool = sync.Pool{
	New: func() interface{} {
		return new(IssueSecurityTokensProcessor)
	},
}

func (IssueSecurityTokens) Process(
	ctx context.Context, getStateFunc base.GetStateFunc,
) ([]base.StateMergeValue, base.OperationProcessReasonError, error) {
	return nil, nil, nil
}

type IssueSecurityTokensItemProcessor struct {
	h      util.Hash
	sender base.Address
	item   IssueSecurityTokensItem
}

func (ipp *IssueSecurityTokensItemProcessor) PreProcess(
	ctx context.Context, op base.Operation, getStateFunc base.GetStateFunc,
) error {
	it := ipp.item

	if err := checkExistsState(extensioncurrency.StateKeyContractAccount(it.Contract()), getStateFunc); err != nil {
		return err
	}

	if err := checkExistsState(currency.StateKeyAccount(it.Receiver()), getStateFunc); err != nil {
		return err
	}

	if err := checkNotExistsState(extensioncurrency.StateKeyContractAccount(it.Receiver()), getStateFunc); err != nil {
		return err
	}

	st, err := existsState(StateKeyDesign(it.Contract(), it.STO()), "key of sto design", getStateFunc)
	if err != nil {
		return err
	}

	design, err := StateDesignValue(st)
	if err != nil {
		return err
	}

	policy := design.Policy()

	controllers := policy.Controllers()
	if len(controllers) == 0 {
		return errors.Errorf("empty controllers, %s-%s", it.Contract(), it.STO())
	}

	for i, con := range controllers {
		if con.Equal(ipp.sender) {
			break
		}

		if i == len(controllers)-1 {
			return errors.Errorf("sender is not controller of sto, %q, %s-%s", ipp.sender, it.Contract(), it.STO())
		}
	}

	gn := new(big.Int)
	gn.SetUint64(design.Granularity())

	if mod := currencybase.NewBigFromBigInt(new(big.Int)).Mod(it.Amount().Int, gn); currencybase.NewBigFromBigInt(mod).OverZero() {
		return errors.Errorf("amount unit does not comply with sto granularity rule, %q, %q", it.Amount(), design.Granularity())
	}

	if err := checkExistsState(currency.StateKeyCurrencyDesign(it.Currency()), getStateFunc); err != nil {
		return err
	}

	return nil
}

func (ipp *IssueSecurityTokensItemProcessor) Process(
	ctx context.Context, op base.Operation, getStateFunc base.GetStateFunc,
) ([]base.StateMergeValue, error) {
	sts := make([]base.StateMergeValue, 4)

	it := ipp.item

	st, err := existsState(StateKeyDesign(it.Contract(), it.STO()), "key of sto design", getStateFunc)
	if err != nil {
		return nil, err
	}

	design, err := StateDesignValue(st)
	if err != nil {
		return nil, err
	}
	p := design.Policy()
	dps := p.Partitions()

	var pb currencybase.Big
	switch st, found, err := getStateFunc(StateKeyPartitionBalance(it.Contract(), it.STO(), it.Partition())); {
	case err != nil:
		return nil, err
	case found:
		pb, err = StatePartitionBalanceValue(st)
		if err != nil {
			return nil, err
		}
		pb = pb.Add(it.Amount())
	default:
		pb = it.Amount()
		dps = append(dps, it.Partition())
	}

	policy := NewPolicy(dps, it.Amount().Add(p.Aggregate()), p.Controllers(), p.Documents())
	if err := policy.IsValid(nil); err != nil {
		return nil, err
	}

	design = NewDesign(design.STO(), design.Granularity(), policy)
	if err := design.IsValid(nil); err != nil {
		return nil, err
	}

	sts[0] = NewStateMergeValue(
		StateKeyDesign(it.Contract(), it.STO()),
		NewDesignStateValue(design),
	)

	sts[1] = NewStateMergeValue(
		StateKeyPartitionBalance(it.Contract(), it.STO(), it.Partition()),
		NewPartitionBalanceStateValue(pb),
	)

	var ps []Partition
	switch st, found, err := getStateFunc(StateKeyTokenHolderPartitions(it.Contract(), it.STO(), it.Receiver())); {
	case err != nil:
		return nil, err
	case found:
		ps, err = StateTokenHolderPartitionsValue(st)
		if err != nil {
			return nil, err
		}
	default:
		ps = []Partition{}
	}

	if len(ps) == 0 {
		ps = append(ps, it.Partition())
	} else {
		for i, pt := range ps {
			if pt == it.Partition() {
				break
			}

			if i == len(ps)-1 {
				ps = append(ps, it.Partition())
			}
		}
	}

	sts[2] = NewStateMergeValue(
		StateKeyTokenHolderPartitions(it.Contract(), it.STO(), it.Receiver()),
		NewTokenHolderPartitionsStateValue(ps),
	)

	var am currencybase.Big
	switch st, found, err := getStateFunc(StateKeyTokenHolderPartitionBalance(it.Contract(), it.STO(), it.Receiver(), it.Partition())); {
	case err != nil:
		return nil, err
	case found:
		am, err = StateTokenHolderPartitionBalanceValue(st)
		if err != nil {
			return nil, err
		}
	default:
		am = currencybase.ZeroBig
	}

	am = am.Add(it.Amount())

	sts[3] = NewStateMergeValue(
		StateKeyTokenHolderPartitionBalance(it.Contract(), it.STO(), it.Receiver(), it.Partition()),
		NewTokenHolderPartitionBalanceStateValue(am, it.Partition()),
	)

	return sts, nil
}

func (ipp *IssueSecurityTokensItemProcessor) Close() error {
	ipp.h = nil
	ipp.sender = nil
	ipp.item = IssueSecurityTokensItem{}

	issueSecurityTokensItemProcessorPool.Put(ipp)

	return nil
}

type IssueSecurityTokensProcessor struct {
	*base.BaseOperationProcessor
}

func NewIssueSecurityTokensProcessor() types.GetNewProcessor {
	return func(
		height base.Height,
		getStateFunc base.GetStateFunc,
		newPreProcessConstraintFunc base.NewOperationProcessorProcessFunc,
		newProcessConstraintFunc base.NewOperationProcessorProcessFunc,
	) (base.OperationProcessor, error) {
		e := util.StringErrorFunc("failed to create new IssueSecurityTokensProcessor")

		nopp := issueSecurityTokensProcessorPool.Get()
		opp, ok := nopp.(*IssueSecurityTokensProcessor)
		if !ok {
			return nil, e(nil, "expected IssueSecurityTokensProcessor, not %T", nopp)
		}

		b, err := base.NewBaseOperationProcessor(
			height, getStateFunc, newPreProcessConstraintFunc, newProcessConstraintFunc)
		if err != nil {
			return nil, e(err, "")
		}

		opp.BaseOperationProcessor = b

		return opp, nil
	}
}

func (opp *IssueSecurityTokensProcessor) PreProcess(
	ctx context.Context, op base.Operation, getStateFunc base.GetStateFunc,
) (context.Context, base.OperationProcessReasonError, error) {
	e := util.StringErrorFunc("failed to preprocess IssueSecurityTokens")

	fact, ok := op.Fact().(IssueSecurityTokensFact)
	if !ok {
		return ctx, nil, e(nil, "expected IssueSecurityTokensFact, not %T", op.Fact())
	}

	if err := fact.IsValid(nil); err != nil {
		return ctx, nil, e(err, "")
	}

	if err := checkExistsState(currency.StateKeyAccount(fact.Sender()), getStateFunc); err != nil {
		return ctx, base.NewBaseOperationProcessReasonError("sender not found, %q: %w", fact.Sender(), err), nil
	}

	if err := checkNotExistsState(extensioncurrency.StateKeyContractAccount(fact.Sender()), getStateFunc); err != nil {
		return ctx, base.NewBaseOperationProcessReasonError("contract account cannot issue security tokens, %q: %w", fact.Sender(), err), nil
	}

	if err := checkFactSignsByState(fact.sender, op.Signs(), getStateFunc); err != nil {
		return ctx, base.NewBaseOperationProcessReasonError("invalid signing: %w", err), nil
	}

	for _, item := range fact.Items() {
		ip := issueSecurityTokensItemProcessorPool.Get()
		ipc, ok := ip.(*IssueSecurityTokensItemProcessor)
		if !ok {
			return nil, nil, e(nil, "expected IssueSecurityTokensItemProcessor, not %T", ip)
		}

		ipc.h = op.Hash()
		ipc.sender = fact.Sender()
		ipc.item = item

		if err := ipc.PreProcess(ctx, op, getStateFunc); err != nil {
			return nil, base.NewBaseOperationProcessReasonError("fail to preprocess IssueSecurityTokensItem: %w", err), nil
		}

		ipc.Close()
	}

	return ctx, nil, nil
}

func (opp *IssueSecurityTokensProcessor) Process( // nolint:dupl
	ctx context.Context, op base.Operation, getStateFunc base.GetStateFunc) (
	[]base.StateMergeValue, base.OperationProcessReasonError, error,
) {
	e := util.StringErrorFunc("failed to process IssueSecurityTokens")

	fact, ok := op.Fact().(IssueSecurityTokensFact)
	if !ok {
		return nil, nil, e(nil, "expected IssueSecurityTokensFact, not %T", op.Fact())
	}

	var sts []base.StateMergeValue // nolint:prealloc

	for _, item := range fact.Items() {
		ip := issueSecurityTokensItemProcessorPool.Get()
		ipc, ok := ip.(*IssueSecurityTokensItemProcessor)
		if !ok {
			return nil, nil, e(nil, "expected IssueSecurityTokensItemProcessor, not %T", ip)
		}

		ipc.h = op.Hash()
		ipc.sender = fact.Sender()
		ipc.item = item

		s, err := ipc.Process(ctx, op, getStateFunc)
		if err != nil {
			return nil, base.NewBaseOperationProcessReasonError("failed to process IssueSecurityTokensItem: %w", err), nil
		}
		sts = append(sts, s...)

		ipc.Close()
	}

	fitems := fact.Items()
	items := make([]STOItem, len(fitems))
	for i := range fact.Items() {
		items[i] = fitems[i]
	}

	required, err := calculateSTOItemsFee(getStateFunc, items)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError("failed to calculate fee: %w", err), nil
	}
	sb, err := currencyoperation.CheckEnoughBalance(fact.sender, required, getStateFunc)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError("failed to check enough balance: %w", err), nil
	}

	for i := range sb {
		v, ok := sb[i].Value().(currency.BalanceStateValue)
		if !ok {
			return nil, nil, e(nil, "expected BalanceStateValue, not %T", sb[i].Value())
		}
		stv := currency.NewBalanceStateValue(v.Amount.WithBig(v.Amount.Big().Sub(required[i][0])))
		sts = append(sts, NewStateMergeValue(sb[i].Key(), stv))
	}

	return sts, nil, nil
}

func (opp *IssueSecurityTokensProcessor) Close() error {
	issueSecurityTokensProcessorPool.Put(opp)

	return nil
}
