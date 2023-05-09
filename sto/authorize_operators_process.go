package sto

import (
	"context"
	"sync"

	extensioncurrency "github.com/ProtoconNet/mitum-currency-extension/v2/currency"
	"github.com/ProtoconNet/mitum-currency/v2/currency"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/pkg/errors"
)

var authorizeOperatorsItemProcessorPool = sync.Pool{
	New: func() interface{} {
		return new(AuthorizeOperatorsItemProcessor)
	},
}

var authorizeOperatorsProcessorPool = sync.Pool{
	New: func() interface{} {
		return new(AuthorizeOperatorsProcessor)
	},
}

func (AuthorizeOperators) Process(
	ctx context.Context, getStateFunc base.GetStateFunc,
) ([]base.StateMergeValue, base.OperationProcessReasonError, error) {
	return nil, nil, nil
}

type AuthorizeOperatorsItemProcessor struct {
	h            util.Hash
	sender       base.Address
	item         AuthorizeOperatorsItem
	operators    *[]base.Address
	tokenHolders *[]base.Address
}

func (ipp *AuthorizeOperatorsItemProcessor) PreProcess(
	ctx context.Context, op base.Operation, getStateFunc base.GetStateFunc,
) error {
	it := ipp.item

	if err := checkExistsState(extensioncurrency.StateKeyContractAccount(it.Contract()), getStateFunc); err != nil {
		return err
	}

	if err := checkExistsState(StateKeySTODesign(it.Contract(), it.STO()), getStateFunc); err != nil {
		return err
	}

	if err := checkExistsState(StateKeyPartitionBalance(it.Contract(), it.STO(), it.Partition()), getStateFunc); err != nil {
		return err
	}

	partitions, err := existsTokenHolderPartitions(it.Contract(), it.STO(), ipp.sender, getStateFunc)
	if err != nil {
		return err
	}

	for i, p := range partitions {
		if p == it.Partition() {
			break
		}

		if i == len(partitions)-1 {
			return errors.Errorf("partition not in tokenholder partition list, %s-%s, %s", ipp.sender, it.Contract(), it.STO(), it.Partition())
		}
	}

	balance, err := existsTokenHolderPartitionBalance(it.Contract(), it.STO(), ipp.sender, it.Partition(), getStateFunc)
	if err != nil {
		return err
	}

	if !balance.OverZero() {
		return errors.Errorf("zero tokenholder partition balance, %s-%s-%s", ipp.sender, it.Contract(), it.STO(), it.Partition())
	}

	for _, ad := range *ipp.operators {
		if ad.Equal(it.Operator()) {
			return errors.Errorf("operator is already in tokenholder operators, %q", ad)
		}
	}

	for _, ad := range *ipp.tokenHolders {
		if ad.Equal(ipp.sender) {
			return errors.Errorf("sender is already in operator tokenholders, %q", ad)
		}
	}

	if err := checkExistsState(currency.StateKeyCurrencyDesign(it.Currency()), getStateFunc); err != nil {
		return err
	}

	return nil
}

func (ipp *AuthorizeOperatorsItemProcessor) Process(
	ctx context.Context, op base.Operation, getStateFunc base.GetStateFunc,
) ([]base.StateMergeValue, error) {
	sts := make([]base.StateMergeValue, 1)

	it := ipp.item

	*ipp.operators = append(*ipp.operators, it.Operator())
	holders := append(*ipp.tokenHolders, ipp.sender)

	sts[0] = NewStateMergeValue(
		StateKeyOperatorTokenHolders(it.Contract(), it.STO(), it.Operator(), it.Partition()),
		NewOperatorTokenHoldersStateValue(holders),
	)

	return sts, nil
}

func (ipp *AuthorizeOperatorsItemProcessor) Close() error {
	ipp.h = nil
	ipp.sender = nil
	ipp.item = AuthorizeOperatorsItem{}
	ipp.operators = nil
	ipp.tokenHolders = nil

	authorizeOperatorsItemProcessorPool.Put(ipp)

	return nil
}

type AuthorizeOperatorsProcessor struct {
	*base.BaseOperationProcessor
}

func NewAuthorizeOperatorsProcessor() extensioncurrency.GetNewProcessor {
	return func(
		height base.Height,
		getStateFunc base.GetStateFunc,
		newPreProcessConstraintFunc base.NewOperationProcessorProcessFunc,
		newProcessConstraintFunc base.NewOperationProcessorProcessFunc,
	) (base.OperationProcessor, error) {
		e := util.StringErrorFunc("failed to create new AuthorizeOperatorsProcessor")

		nopp := authorizeOperatorsProcessorPool.Get()
		opp, ok := nopp.(*AuthorizeOperatorsProcessor)
		if !ok {
			return nil, e(nil, "expected AuthorizeOperatorsProcessor, not %T", nopp)
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

func (opp *AuthorizeOperatorsProcessor) PreProcess(
	ctx context.Context, op base.Operation, getStateFunc base.GetStateFunc,
) (context.Context, base.OperationProcessReasonError, error) {
	e := util.StringErrorFunc("failed to preprocess AuthorizeOperators")

	fact, ok := op.Fact().(AuthorizeOperatorsFact)
	if !ok {
		return ctx, nil, e(nil, "expected AuthorizeOperatorsFact, not %T", op.Fact())
	}

	if err := fact.IsValid(nil); err != nil {
		return ctx, nil, e(err, "")
	}

	if err := checkExistsState(currency.StateKeyAccount(fact.Sender()), getStateFunc); err != nil {
		return ctx, base.NewBaseOperationProcessReasonError("sender not found, %q: %w", fact.Sender(), err), nil
	}

	if err := checkNotExistsState(extensioncurrency.StateKeyContractAccount(fact.Sender()), getStateFunc); err != nil {
		return ctx, base.NewBaseOperationProcessReasonError("contract account cannot set its operators, %q: %w", fact.Sender(), err), nil
	}

	if err := checkFactSignsByState(fact.sender, op.Signs(), getStateFunc); err != nil {
		return ctx, base.NewBaseOperationProcessReasonError("invalid signing: %w", err), nil
	}

	operators := map[string]*[]base.Address{}
	holders := map[string]*[]base.Address{}

	for _, it := range fact.Items() {
		var ops, hds []base.Address

		k := StateKeyTokenHolderPartitionOperators(it.Contract(), it.STO(), fact.sender, it.Partition())
		if _, found := operators[k]; !found {
			switch st, found, err := getStateFunc(k); {
			case err != nil:
				return nil, base.NewBaseOperationProcessReasonError("failed to find tokenholder partition operators, %s: %w", k, err), nil
			case found:
				ops, err = StateTokenHolderPartitionOperatorsValue(st)
				if err != nil {
					return nil, base.NewBaseOperationProcessReasonError("failed to get tokenholder partition operators, %s: %w", k, err), nil
				}
			default:
				ops = []base.Address{}
			}
			operators[k] = &ops
		}

		k = StateKeyOperatorTokenHolders(it.Contract(), it.STO(), it.Operator(), it.Partition())
		if _, found := holders[k]; !found {
			switch st, found, err := getStateFunc(k); {
			case err != nil:
				return nil, base.NewBaseOperationProcessReasonError("failed to find operator tokenholders, %s: %w", k, err), nil
			case found:
				hds, err = StateOperatorTokenHoldersValue(st)
				if err != nil {
					return nil, base.NewBaseOperationProcessReasonError("failed to get operator tokenholders, %s: %w", k, err), nil
				}
			default:
				hds = []base.Address{}
			}

			holders[k] = &hds
		}
	}

	for _, it := range fact.Items() {
		ip := authorizeOperatorsItemProcessorPool.Get()
		ipc, ok := ip.(*AuthorizeOperatorsItemProcessor)
		if !ok {
			return nil, nil, e(nil, "expected AuthorizeOperatorsItemProcessor, not %T", ip)
		}

		ipc.h = op.Hash()
		ipc.sender = fact.Sender()
		ipc.item = it
		ipc.operators = operators[StateKeyTokenHolderPartitionOperators(it.Contract(), it.STO(), fact.sender, it.Partition())]
		ipc.tokenHolders = holders[StateKeyOperatorTokenHolders(it.Contract(), it.STO(), it.Operator(), it.Partition())]

		if err := ipc.PreProcess(ctx, op, getStateFunc); err != nil {
			return nil, base.NewBaseOperationProcessReasonError("fail to preprocess AuthorizeOperatorsItem: %w", err), nil
		}

		ipc.Close()
	}

	return ctx, nil, nil
}

func (opp *AuthorizeOperatorsProcessor) Process( // nolint:dupl
	ctx context.Context, op base.Operation, getStateFunc base.GetStateFunc) (
	[]base.StateMergeValue, base.OperationProcessReasonError, error,
) {
	e := util.StringErrorFunc("failed to process AuthorizeOperators")

	fact, ok := op.Fact().(AuthorizeOperatorsFact)
	if !ok {
		return nil, nil, e(nil, "expected AuthorizeOperatorsFact, not %T", op.Fact())
	}

	var sts []base.StateMergeValue // nolint:prealloc

	operators := map[string]*[]base.Address{}
	holders := map[string]*[]base.Address{}

	for _, it := range fact.Items() {
		var ops, hds []base.Address

		k := StateKeyTokenHolderPartitionOperators(it.Contract(), it.STO(), fact.sender, it.Partition())
		if _, found := operators[k]; !found {
			switch st, found, err := getStateFunc(k); {
			case err != nil:
				return nil, base.NewBaseOperationProcessReasonError("failed to find tokenholder partition operators, %s: %w", k, err), nil
			case found:
				ops, err = StateTokenHolderPartitionOperatorsValue(st)
				if err != nil {
					return nil, base.NewBaseOperationProcessReasonError("failed to get tokenholder partition operators, %s: %w", k, err), nil
				}
			default:
				ops = []base.Address{}
			}
			operators[k] = &ops
		}

		k = StateKeyOperatorTokenHolders(it.Contract(), it.STO(), it.Operator(), it.Partition())
		if _, found := holders[k]; !found {
			switch st, found, err := getStateFunc(k); {
			case err != nil:
				return nil, base.NewBaseOperationProcessReasonError("failed to find operator tokenholders, %s: %w", k, err), nil
			case found:
				hds, err = StateOperatorTokenHoldersValue(st)
				if err != nil {
					return nil, base.NewBaseOperationProcessReasonError("failed to get operator tokenholders, %s: %w", k, err), nil
				}
			default:
				hds = []base.Address{}
			}

			holders[k] = &hds
		}
	}

	ipcs := make([]*AuthorizeOperatorsItemProcessor, len(fact.items))
	for i, it := range fact.Items() {
		ip := authorizeOperatorsItemProcessorPool.Get()
		ipc, ok := ip.(*AuthorizeOperatorsItemProcessor)
		if !ok {
			return nil, nil, e(nil, "expected AuthorizeOperatorsItemProcessor, not %T", ip)
		}

		ipc.h = op.Hash()
		ipc.sender = fact.Sender()
		ipc.item = it
		ipc.operators = operators[StateKeyTokenHolderPartitionOperators(it.Contract(), it.STO(), fact.sender, it.Partition())]
		ipc.tokenHolders = holders[StateKeyOperatorTokenHolders(it.Contract(), it.STO(), it.Operator(), it.Partition())]

		s, err := ipc.Process(ctx, op, getStateFunc)
		if err != nil {
			return nil, base.NewBaseOperationProcessReasonError("failed to process AuthorizeOperatorsItem: %w", err), nil
		}
		sts = append(sts, s...)

		ipcs[i] = ipc
	}

	for k, v := range operators {
		sts = append(sts, NewStateMergeValue(
			k,
			NewTokenHolderPartitionOperatorsStateValue(*v),
		))
	}

	for _, ipc := range ipcs {
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

func (opp *AuthorizeOperatorsProcessor) Close() error {
	authorizeOperatorsProcessorPool.Put(opp)

	return nil
}
