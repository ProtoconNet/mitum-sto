package sto

import (
	"context"
	"sync"

	extensioncurrency "github.com/ProtoconNet/mitum-currency-extension/v2/currency"
	"github.com/ProtoconNet/mitum-currency/v2/currency"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
)

var createSecurityTokensItemProcessorPool = sync.Pool{
	New: func() interface{} {
		return new(CreateSecurityTokensItemProcessor)
	},
}

var createSecurityTokensProcessorPool = sync.Pool{
	New: func() interface{} {
		return new(CreateSecurityTokensProcessor)
	},
}

func (CreateSecurityTokens) Process(
	ctx context.Context, getStateFunc base.GetStateFunc,
) ([]base.StateMergeValue, base.OperationProcessReasonError, error) {
	return nil, nil, nil
}

type CreateSecurityTokensItemProcessor struct {
	h      util.Hash
	sender base.Address
	item   CreateSecurityTokensItem
}

func (ipp *CreateSecurityTokensItemProcessor) PreProcess(
	ctx context.Context, op base.Operation, getStateFunc base.GetStateFunc,
) error {
	it := ipp.item

	if err := checkExistsState(extensioncurrency.StateKeyContractAccount(it.Contract()), getStateFunc); err != nil {
		return err
	}

	if err := checkNotExistsState(StateKeySTODesign(it.Contract(), it.STO()), getStateFunc); err != nil {
		return err
	}

	if err := checkNotExistsState(StateKeyPartitionBalance(it.Contract(), it.STO(), it.DefaultPartition()), getStateFunc); err != nil {
		return err
	}

	for _, con := range it.Controllers() {
		if err := checkExistsState(currency.StateKeyAccount(con), getStateFunc); err != nil {
			return err
		}
	}

	if err := checkExistsState(currency.StateKeyCurrencyDesign(it.Currency()), getStateFunc); err != nil {
		return err
	}

	return nil
}

func (ipp *CreateSecurityTokensItemProcessor) Process(
	ctx context.Context, op base.Operation, getStateFunc base.GetStateFunc,
) ([]base.StateMergeValue, error) {
	sts := make([]base.StateMergeValue, 2)

	it := ipp.item

	partition := it.DefaultPartition()
	partitions := []Partition{partition}
	documents := []Document{}

	policy := NewSTOPolicy(partitions, currency.NewBig(0), it.Controllers(), documents)
	design := NewSTODesign(it.STO(), it.Granularity(), policy)

	if err := design.IsValid(nil); err != nil {
		return nil, err
	}

	sts[0] = NewStateMergeValue(
		StateKeySTODesign(it.Contract(), it.STO()),
		NewSTODesignStateValue(design),
	)
	sts[1] = NewStateMergeValue(
		StateKeyPartitionBalance(it.Contract(), it.STO(), it.DefaultPartition()),
		NewPartitionBalanceStateValue(currency.ZeroBig),
	)

	return sts, nil
}

func (ipp *CreateSecurityTokensItemProcessor) Close() error {
	ipp.h = nil
	ipp.sender = nil
	ipp.item = CreateSecurityTokensItem{}

	createSecurityTokensItemProcessorPool.Put(ipp)

	return nil
}

type CreateSecurityTokensProcessor struct {
	*base.BaseOperationProcessor
}

func NewCreateSecurityTokensProcessor() extensioncurrency.GetNewProcessor {
	return func(
		height base.Height,
		getStateFunc base.GetStateFunc,
		newPreProcessConstraintFunc base.NewOperationProcessorProcessFunc,
		newProcessConstraintFunc base.NewOperationProcessorProcessFunc,
	) (base.OperationProcessor, error) {
		e := util.StringErrorFunc("failed to create new CreateSecurityTokensProcessor")

		nopp := createSecurityTokensProcessorPool.Get()
		opp, ok := nopp.(*CreateSecurityTokensProcessor)
		if !ok {
			return nil, e(nil, "expected CreateSecurityTokensProcessor, not %T", nopp)
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

func (opp *CreateSecurityTokensProcessor) PreProcess(
	ctx context.Context, op base.Operation, getStateFunc base.GetStateFunc,
) (context.Context, base.OperationProcessReasonError, error) {
	e := util.StringErrorFunc("failed to preprocess CreateSecurityTokens")

	fact, ok := op.Fact().(CreateSecurityTokensFact)
	if !ok {
		return ctx, nil, e(nil, "expected CreateSecurityTokensFact, not %T", op.Fact())
	}

	if err := fact.IsValid(nil); err != nil {
		return ctx, nil, e(err, "")
	}

	if err := checkExistsState(currency.StateKeyAccount(fact.Sender()), getStateFunc); err != nil {
		return ctx, base.NewBaseOperationProcessReasonError("sender not found, %q: %w", fact.Sender(), err), nil
	}

	if err := checkNotExistsState(extensioncurrency.StateKeyContractAccount(fact.Sender()), getStateFunc); err != nil {
		return ctx, base.NewBaseOperationProcessReasonError("contract account cannot create security tokens, %q", fact.Sender()), nil
	}

	if err := checkFactSignsByState(fact.sender, op.Signs(), getStateFunc); err != nil {
		return ctx, base.NewBaseOperationProcessReasonError("invalid signing: %w", err), nil
	}

	for _, it := range fact.Items() {
		ip := createSecurityTokensItemProcessorPool.Get()
		ipc, ok := ip.(*CreateSecurityTokensItemProcessor)
		if !ok {
			return nil, nil, e(nil, "expected CreateSecurityTokensItemProcessor, not %T", ip)
		}

		ipc.h = op.Hash()
		ipc.sender = fact.Sender()
		ipc.item = it

		if err := ipc.PreProcess(ctx, op, getStateFunc); err != nil {
			return nil, base.NewBaseOperationProcessReasonError("fail to preprocess CreateSecurityTokensItem: %w", err), nil
		}

		ipc.Close()
	}

	return ctx, nil, nil
}

func (opp *CreateSecurityTokensProcessor) Process( // nolint:dupl
	ctx context.Context, op base.Operation, getStateFunc base.GetStateFunc) (
	[]base.StateMergeValue, base.OperationProcessReasonError, error,
) {
	e := util.StringErrorFunc("failed to process CreateSecurityTokens")

	fact, ok := op.Fact().(CreateSecurityTokensFact)
	if !ok {
		return nil, nil, e(nil, "expected CreateSecurityTokensFact, not %T", op.Fact())
	}

	var sts []base.StateMergeValue // nolint:prealloc

	for _, it := range fact.Items() {
		ip := createSecurityTokensItemProcessorPool.Get()
		ipc, ok := ip.(*CreateSecurityTokensItemProcessor)
		if !ok {
			return nil, nil, e(nil, "expected CreateSecurityTokensItemProcessor, not %T", ip)
		}

		ipc.h = op.Hash()
		ipc.sender = fact.Sender()
		ipc.item = it

		s, err := ipc.Process(ctx, op, getStateFunc)
		if err != nil {
			return nil, base.NewBaseOperationProcessReasonError("failed to process CreateSecurityTokensItem: %w", err), nil
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
	sb, err := currency.CheckEnoughBalance(fact.sender, required, getStateFunc)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError("failed to check enough balance: %w", err), nil
	}

	for i := range sb {
		v, ok := sb[i].Value().(currency.BalanceStateValue)
		if !ok {
			return nil, nil, e(nil, "expected BalanceStateValue, not %T", sb[i].Value())
		}
		stv := currency.NewBalanceStateValue(v.Amount.WithBig(v.Amount.Big().Sub(required[i][0])))
		sts = append(sts, currency.NewBalanceStateMergeValue(sb[i].Key(), stv))
	}

	return sts, nil, nil
}

func (opp *CreateSecurityTokensProcessor) Close() error {
	createSecurityTokensProcessorPool.Put(opp)

	return nil
}

func calculateSTOItemsFee(getStateFunc base.GetStateFunc, items []STOItem) (map[currency.CurrencyID][2]currency.Big, error) {
	required := map[currency.CurrencyID][2]currency.Big{}

	for _, item := range items {
		rq := [2]currency.Big{currency.ZeroBig, currency.ZeroBig}

		if k, found := required[item.Currency()]; found {
			rq = k
		}

		policy, err := existsCurrencyPolicy(item.Currency(), getStateFunc)
		if err != nil {
			return nil, err
		}

		switch k, err := policy.Feeer().Fee(currency.ZeroBig); {
		case err != nil:
			return nil, err
		case !k.OverZero():
			required[item.Currency()] = [2]currency.Big{rq[0], rq[1]}
		default:
			required[item.Currency()] = [2]currency.Big{rq[0].Add(k), rq[1].Add(k)}
		}

	}

	return required, nil

}
