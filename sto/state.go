package sto

import (
	"bytes"
	"fmt"
	"sort"
	"strings"

	"github.com/ProtoconNet/mitum-currency/v2/currency"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
	"github.com/pkg/errors"
)

var MaxOperatorInOperators = 10
var MaxTokenHolderInTokenHolders = 10

var (
	TokenHolderPartitionsStateValueHint = hint.MustNewHint("token-holder-partitions-state-value-v0.0.1")
	StateKeyTokenHolderPrefix           = "token-holder:"
	StateKeyPartitionsSuffix            = ":partitions"
)

type TokenHolderPartitionsStateValue struct {
	hint.BaseHinter
	Partitions []Partition
}

func NewTokenHolderPartitionsStateValue(partitions []Partition) TokenHolderPartitionsStateValue {
	return TokenHolderPartitionsStateValue{
		BaseHinter: hint.NewBaseHinter(TokenHolderPartitionsStateValueHint),
		Partitions: partitions,
	}
}

func (sv TokenHolderPartitionsStateValue) Hint() hint.Hint {
	return sv.BaseHinter.Hint()
}

func (sv TokenHolderPartitionsStateValue) IsValid([]byte) error {
	e := util.ErrInvalid.Errorf("invalid TokenHolderPartitionsStateValue")

	if err := sv.BaseHinter.IsValid(TokenHolderPartitionsStateValueHint.Type().Bytes()); err != nil {
		return e.Wrap(err)
	}

	if len(sv.Partitions) == 0 {
		return errors.Errorf("empty partitions")
	}

	for _, partition := range sv.Partitions {
		if err := partition.IsValid(nil); err != nil {
			return e.Wrap(err)
		}
	}

	return nil
}

func (sv TokenHolderPartitionsStateValue) HashBytes() []byte {
	bs := make([][]byte, len(sv.Partitions))
	sort.Slice(sv.Partitions, func(i, j int) bool {
		return bytes.Compare(sv.Partitions[i].Bytes(), sv.Partitions[j].Bytes()) < 0
	})
	for i, partition := range sv.Partitions {
		bs[i] = partition.Bytes()
	}
	return util.ConcatBytesSlice(bs...)
}

func StateKeyTokenHolderPartitions(addr base.Address) string {
	return StateKeyTokenHolderPrefix + addr.String() + StateKeyPartitionsSuffix
}

var (
	TokenHolderPartitionBalanceStateValueHint = hint.MustNewHint("token-holder-partition-balance-state-value-v0.0.1")
	StateKeyPartitionBalanceSuffix            = ":partition-balance"
)

type TokenHolderPartitionBalanceStateValue struct {
	hint.BaseHinter
	Amount    currency.Amount
	Partition Partition
}

func NewTokenHolderPartitionBalanceStateValue(amount currency.Amount, partition Partition) TokenHolderPartitionBalanceStateValue {
	return TokenHolderPartitionBalanceStateValue{
		BaseHinter: hint.NewBaseHinter(TokenHolderPartitionBalanceStateValueHint),
		Amount:     amount,
		Partition:  partition,
	}
}

func (sv TokenHolderPartitionBalanceStateValue) Hint() hint.Hint {
	return sv.BaseHinter.Hint()
}

func (sv TokenHolderPartitionBalanceStateValue) IsValid([]byte) error {
	e := util.ErrInvalid.Errorf("invalid TokenHolderPartitionBalanceStateValue")

	if err := sv.BaseHinter.IsValid(TokenHolderPartitionBalanceStateValueHint.Type().Bytes()); err != nil {
		return e.Wrap(err)
	}

	if err := util.CheckIsValiders(nil, false, sv.Amount, sv.Partition); err != nil {
		return e.Wrap(err)
	}

	return nil
}

func (sv TokenHolderPartitionBalanceStateValue) HashBytes() []byte {
	return sv.Amount.Bytes()
}

func StateKeyTokenHolderPartitionBalance(address base.Address, partition Partition) string {
	return fmt.Sprintf("%s%s%s%s", StateKeyTokenHolderPrefix, address.String(), partition, StateKeyPartitionBalanceSuffix)
}

func IsStateTokenHolderPartitionBalanceKey(key string) bool {
	return strings.HasSuffix(key, StateKeyPartitionBalanceSuffix)
}

var (
	TokenHolderPartitionOperatorsStateValueHint = hint.MustNewHint("token-holder-partition-operators-state-value-v0.0.1")
	StateKeyPartitionOperatorsSuffix            = ":partition-operators"
)

type TokenHolderPartitionOperatorsStateValue struct {
	hint.BaseHinter
	Operators []currency.Account
}

func NewTokenHolderPartitionOperatorsStateValue(operators []currency.Account) TokenHolderPartitionOperatorsStateValue {
	return TokenHolderPartitionOperatorsStateValue{
		BaseHinter: hint.NewBaseHinter(TokenHolderPartitionOperatorsStateValueHint),
		Operators:  operators,
	}
}

func (sv TokenHolderPartitionOperatorsStateValue) Hint() hint.Hint {
	return sv.BaseHinter.Hint()
}

func (sv TokenHolderPartitionOperatorsStateValue) IsValid([]byte) error {
	e := util.ErrInvalid.Errorf("invalid TokenHolderPartitionOperatorsStateValue")

	if err := sv.BaseHinter.IsValid(TokenHolderPartitionOperatorsStateValueHint.Type().Bytes()); err != nil {
		return e.Wrap(err)
	}

	if n := len(sv.Operators); n < 1 {
		return util.ErrInvalid.Errorf("empty keys")
	} else if n > MaxOperatorInOperators {
		return util.ErrInvalid.Errorf("keys over %d, %d", MaxOperatorInOperators, n)
	}

	m := map[string]struct{}{}
	for i := range sv.Operators {
		k := sv.Operators[i]
		if err := util.CheckIsValiders(nil, false, k); err != nil {
			return err
		}

		if _, found := m[k.Address().String()]; found {
			return util.ErrInvalid.Errorf("duplicated Account found")
		}

		m[k.Address().String()] = struct{}{}
	}

	return nil
}

func (sv TokenHolderPartitionOperatorsStateValue) HashBytes() []byte {
	bs := make([][]byte, len(sv.Operators))
	sort.Slice(sv.Operators, func(i, j int) bool {
		return bytes.Compare(sv.Operators[i].Bytes(), sv.Operators[j].Bytes()) < 0
	})
	for i, partition := range sv.Operators {
		bs[i] = partition.Bytes()
	}
	return util.ConcatBytesSlice(bs...)
}

func StateKeyTokenHolderPartitionOperators(addr base.Address, partition Partition) string {
	return fmt.Sprintf("%s%s_%s%s", StateKeyTokenHolderPrefix, addr.String(), partition.String(), StateKeyPartitionOperatorsSuffix)
}

func IsStateTokenHolderPartitionOperatorsKey(key string) bool {
	return strings.HasSuffix(key, StateKeyPartitionOperatorsSuffix)
}

var (
	PartitionBalanceStateValueHint = hint.MustNewHint("partition-balance-state-value-v0.0.1")
	StateKeyPartitionPrefix        = "partition:"
)

type PartitionBalanceStateValue struct {
	hint.BaseHinter
	Amount currency.Amount
}

func NewPartitionBalanceStateValue(amount currency.Amount) PartitionBalanceStateValue {
	return PartitionBalanceStateValue{
		BaseHinter: hint.NewBaseHinter(PartitionBalanceStateValueHint),
		Amount:     amount,
	}
}

func (sv PartitionBalanceStateValue) Hint() hint.Hint {
	return sv.BaseHinter.Hint()
}

func (sv PartitionBalanceStateValue) IsValid([]byte) error {
	e := util.ErrInvalid.Errorf("invalid PartitionBalanceStateValue")

	if err := sv.BaseHinter.IsValid(PartitionBalanceStateValueHint.Type().Bytes()); err != nil {
		return e.Wrap(err)
	}

	return nil
}

func (sv PartitionBalanceStateValue) HashBytes() []byte {
	return sv.Amount.Bytes()
}

func StateKeyPartitionBalance(partition Partition) string {
	return fmt.Sprintf("%s%s%s", StateKeyPartitionPrefix, partition.String(), StateKeyPartitionBalanceSuffix)
}

func IsStatePartitionBalanceKey(key string) bool {
	return strings.HasSuffix(key, StateKeyPartitionBalanceSuffix)
}

var (
	PartitionControllersStateValueHint = hint.MustNewHint("partition-controllers-state-value-v0.0.1")
	StateKeyPartitionControllersSuffix = ":partition-controllers"
)

type PartitionControllersStateValue struct {
	hint.BaseHinter
	Controllers []base.Address
}

func NewPartitionControllersStateValue(controllers []base.Address) PartitionControllersStateValue {
	return PartitionControllersStateValue{
		BaseHinter:  hint.NewBaseHinter(PartitionControllersStateValueHint),
		Controllers: controllers,
	}
}

func (p PartitionControllersStateValue) Hint() hint.Hint {
	return p.BaseHinter.Hint()
}

func (p PartitionControllersStateValue) IsValid([]byte) error {
	e := util.ErrInvalid.Errorf("invalid PartitionControllersStateValue")

	if err := p.BaseHinter.IsValid(PartitionControllersStateValueHint.Type().Bytes()); err != nil {
		return e.Wrap(err)
	}

	if len(p.Controllers) == 0 {
		return errors.Errorf("empty controllers")
	}

	m := map[string]struct{}{}
	for _, controller := range p.Controllers {
		if _, found := m[controller.String()]; found {
			return util.ErrInvalid.Errorf("duplicated Address found")
		}
		m[controller.String()] = struct{}{}
	}

	return nil
}

func (p PartitionControllersStateValue) HashBytes() []byte {
	bs := make([][]byte, len(p.Controllers))
	sort.Slice(p.Controllers, func(i, j int) bool {
		return bytes.Compare(p.Controllers[i].Bytes(), p.Controllers[j].Bytes()) < 0
	})
	for i, controller := range p.Controllers {
		bs[i] = controller.Bytes()
	}
	return util.ConcatBytesSlice(bs...)
}

func StateKeyPartitionControllers(partition Partition) string {
	return fmt.Sprintf("%s%s%s", StateKeyPartitionPrefix, partition.String(), StateKeyPartitionControllersSuffix)
}

func IsStatePartitionControllersKey(key string) bool {
	return strings.HasSuffix(key, StateKeyPartitionControllersSuffix)
}

var (
	StateKeyOperatorPrefix             = "state:operator:"
	StateKeyOperatorTokenHoldersSuffix = ":token-holders"
)

type OperatorTokenHoldersStateValue struct {
	hint.BaseHinter
	TokenHolders []currency.Account
}

func NewOperatorTokenHoldersStateValue(tokenHolders []currency.Account) OperatorTokenHoldersStateValue {
	return OperatorTokenHoldersStateValue{
		BaseHinter:   hint.NewBaseHinter(TokenHolderPartitionOperatorsStateValueHint),
		TokenHolders: tokenHolders,
	}
}

func (o OperatorTokenHoldersStateValue) Hint() hint.Hint {
	return o.BaseHinter.Hint()
}

func (o OperatorTokenHoldersStateValue) IsValid([]byte) error {
	e := util.ErrInvalid.Errorf("invalid OperatorTokenHoldersStateValue")

	if err := o.BaseHinter.IsValid(TokenHolderPartitionOperatorsStateValueHint.Type().Bytes()); err != nil {
		return e.Wrap(err)
	}

	if n := len(o.TokenHolders); n < 1 {
		return util.ErrInvalid.Errorf("empty keys")
	} else if n > MaxTokenHolderInTokenHolders {
		return util.ErrInvalid.Errorf("keys over %d, %d", MaxTokenHolderInTokenHolders, n)
	}

	m := map[string]struct{}{}
	for i := range o.TokenHolders {
		k := o.TokenHolders[i]
		if err := util.CheckIsValiders(nil, false, k); err != nil {
			return err
		}

		if _, found := m[k.Address().String()]; found {
			return util.ErrInvalid.Errorf("duplicated Account found")
		}

		m[k.Address().String()] = struct{}{}
	}

	return nil
}

func (o OperatorTokenHoldersStateValue) HashBytes() []byte {
	bs := make([][]byte, len(o.TokenHolders))
	sort.Slice(o.TokenHolders, func(i, j int) bool {
		return bytes.Compare(o.TokenHolders[i].Bytes(), o.TokenHolders[j].Bytes()) < 0
	})
	for i, partition := range o.TokenHolders {
		bs[i] = partition.Bytes()
	}
	return util.ConcatBytesSlice(bs...)
}

func StateKeyOperatorTokenHolders(addr base.Address) string {
	return fmt.Sprintf("%s%s%s", StateKeyOperatorPrefix, addr.String(), StateKeyOperatorTokenHoldersSuffix)
}

func IsStateOperatorTokenHoldersKey(key string) bool {
	return strings.HasSuffix(key, StateKeyOperatorTokenHoldersSuffix)
}
