package sto

import (
	"bytes"
	"fmt"
	"sort"
	"strings"

	currencyextension "github.com/ProtoconNet/mitum-currency-extension/v2/currency"
	"github.com/ProtoconNet/mitum-currency/v2/currency"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
	"github.com/pkg/errors"
)

var (
	STOPrefix               = "sto:"
	STODesignStateValueHint = hint.MustNewHint("mitum-sto-design-state-value-v0.0.1")
	STODesignSuffix         = ":design"
)

type StateValueMerger struct {
	*base.BaseStateValueMerger
}

func NewStateValueMerger(height base.Height, key string, st base.State) *StateValueMerger {
	s := &StateValueMerger{
		BaseStateValueMerger: base.NewBaseStateValueMerger(height, key, st),
	}

	return s
}

func NewStateMergeValue(key string, stv base.StateValue) base.StateMergeValue {
	StateValueMergerFunc := func(height base.Height, st base.State) base.StateValueMerger {
		return NewStateValueMerger(height, key, st)
	}

	return base.NewBaseStateMergeValue(
		key,
		stv,
		StateValueMergerFunc,
	)
}

// sto:address-stoID
func StateKeySTOPrefix(addr base.Address, stoID currencyextension.ContractID) string {
	return fmt.Sprintf("%s%s-%s", STOPrefix, addr.String(), stoID)
}

type STODesignStateValue struct {
	hint.BaseHinter
	Design STODesign
}

func NewSTODesignStateValue(design STODesign) STODesignStateValue {
	return STODesignStateValue{
		BaseHinter: hint.NewBaseHinter(STODesignStateValueHint),
		Design:     design,
	}
}

func (sd STODesignStateValue) Hint() hint.Hint {
	return sd.BaseHinter.Hint()
}

func (sd STODesignStateValue) IsValid([]byte) error {
	e := util.ErrInvalid.Errorf("invalid STODesignStateValue")

	if err := sd.BaseHinter.IsValid(STODesignStateValueHint.Type().Bytes()); err != nil {
		return e.Wrap(err)
	}

	if err := sd.Design.IsValid(nil); err != nil {
		return e.Wrap(err)
	}

	return nil
}

func (sd STODesignStateValue) HashBytes() []byte {
	return sd.Design.Bytes()
}

func StateSTODesignValue(st base.State) (STODesign, error) {
	v := st.Value()
	if v == nil {
		return STODesign{}, util.ErrNotFound.Errorf("sto design not found in State")
	}

	d, ok := v.(STODesignStateValue)
	if !ok {
		return STODesign{}, errors.Errorf("invalid sto design value found, %T", v)
	}

	return d.Design, nil
}

func IsStateSTODesignKey(key string) bool {
	return strings.HasSuffix(key, STODesignSuffix)
}

// sto:address-stoID:design
func StateKeySTODesign(addr base.Address, sid currencyextension.ContractID) string {
	return fmt.Sprintf("%s%s", StateKeySTOPrefix(addr, sid), STODesignSuffix)
}

var MaxOperatorInOperators = 10
var MaxTokenHolderInTokenHolders = 10

var (
	TokenHolderPartitionsStateValueHint = hint.MustNewHint("mitum-sto-tokenholder-partitions-state-value-v0.0.1")
	TokenHolderPartitionsSuffix         = ":holder-partitions"
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

func IsStateTokenHolderPartitionsKey(key string) bool {
	return strings.HasSuffix(key, TokenHolderPartitionsSuffix)
}

func StateKeyTokenHolderPartitions(caddr base.Address, sid currencyextension.ContractID, uaddr base.Address) string {
	return fmt.Sprintf("%s-%s%s", StateKeySTOPrefix(caddr, sid), uaddr.String(), TokenHolderPartitionsSuffix)
}

var (
	TokenHolderPartitionBalanceStateValueHint = hint.MustNewHint("mitum-sto-tokenholder-partition-balance-state-value-v0.0.1")
	TokenHolderPartitionBalanceSuffix         = ":holder-partition-balance"
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

func StateKeyTokenHolderPartitionBalance(caddr base.Address, stoID currencyextension.ContractID, uaddr base.Address, partition Partition) string {
	return fmt.Sprintf("%s-%s-%s%s", StateKeySTOPrefix(caddr, stoID), uaddr.String(), partition, TokenHolderPartitionBalanceSuffix)
}

func IsStateTokenHolderPartitionBalanceKey(key string) bool {
	return strings.HasSuffix(key, TokenHolderPartitionBalanceSuffix)
}

var (
	TokenHolderPartitionOperatorsStateValueHint = hint.MustNewHint("mitum-sto-tokenholder-partition-operators-state-value-v0.0.1")
	TokenHolderPartitionOperatorsSuffix         = ":holder-partition-operators"
)

type TokenHolderPartitionOperatorsStateValue struct {
	hint.BaseHinter
	Operators []base.Address
}

func NewTokenHolderPartitionOperatorsStateValue(operators []base.Address) TokenHolderPartitionOperatorsStateValue {
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

		if _, found := m[k.String()]; found {
			return util.ErrInvalid.Errorf("duplicated Account found")
		}

		m[k.String()] = struct{}{}
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

func StateKeyTokenHolderPartitionOperators(caddr base.Address, stoID currencyextension.ContractID, uaddr base.Address, partition Partition) string {
	return fmt.Sprintf("%s-%s-%s%s", StateKeySTOPrefix(caddr, stoID), uaddr.String(), partition.String(), TokenHolderPartitionOperatorsSuffix)
}

func IsStateTokenHolderPartitionOperatorsKey(key string) bool {
	return strings.HasSuffix(key, TokenHolderPartitionOperatorsSuffix)
}

var (
	PartitionBalanceStateValueHint = hint.MustNewHint("mitum-sto-partition-balance-state-value-v0.0.1")
	PartitionBalanceSuffix         = ":partition-balance"
)

type PartitionBalanceStateValue struct {
	hint.BaseHinter
	Amount currency.Big
}

func NewPartitionBalanceStateValue(amount currency.Big) PartitionBalanceStateValue {
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

func StateKeyPartitionBalance(caddr base.Address, stoID currencyextension.ContractID, partition Partition) string {
	return fmt.Sprintf("%s-%s%s", StateKeySTOPrefix(caddr, stoID), partition.String(), PartitionBalanceSuffix)
}

func IsStatePartitionBalanceKey(key string) bool {
	return strings.HasSuffix(key, PartitionBalanceSuffix)
}

var (
	PartitionControllersStateValueHint = hint.MustNewHint("mitum-sto-partition-controllers-state-value-v0.0.1")
	PartitionControllersSuffix         = ":partition-controllers"
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

func StateKeyPartitionControllers(caddr base.Address, stoID currencyextension.ContractID, partition Partition) string {
	return fmt.Sprintf("%s-%s%s", StateKeySTOPrefix(caddr, stoID), partition.String(), PartitionControllersSuffix)
}

func IsStatePartitionControllersKey(key string) bool {
	return strings.HasSuffix(key, PartitionControllersSuffix)
}

var (
	OperatorTokenHoldersStateValueHint = hint.MustNewHint("mitum-sto-operator-tokenHolders-state-value-v0.0.1")
	OperatorTokenHoldersSuffix         = ":operator-holders"
)

type OperatorTokenHoldersStateValue struct {
	hint.BaseHinter
	TokenHolders []base.Address
}

func NewOperatorTokenHoldersStateValue(tokenHolders []base.Address) OperatorTokenHoldersStateValue {
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

		if _, found := m[k.String()]; found {
			return util.ErrInvalid.Errorf("duplicated Account found")
		}

		m[k.String()] = struct{}{}
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

func StateKeyOperatorTokenHolders(caddr base.Address, stoID currencyextension.ContractID, oaddr base.Address) string {
	return fmt.Sprintf("%s-%s%s", StateKeySTOPrefix(caddr, stoID), oaddr.String(), OperatorTokenHoldersSuffix)
}

func IsStateOperatorTokenHoldersKey(key string) bool {
	return strings.HasSuffix(key, OperatorTokenHoldersSuffix)
}
