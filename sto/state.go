package sto

import (
	"bytes"
	"fmt"
	"sort"
	"strings"

	currencybase "github.com/ProtoconNet/mitum-currency/v3/base"
	currency "github.com/ProtoconNet/mitum-currency/v3/state/currency"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
	"github.com/pkg/errors"
)

var (
	STOPrefix            = "sto:"
	DesignStateValueHint = hint.MustNewHint("mitum-sto-design-state-value-v0.0.1")
	DesignSuffix         = ":design"
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
func StateKeySTOPrefix(addr base.Address, stoID currencybase.ContractID) string {
	return fmt.Sprintf("%s%s-%s", STOPrefix, addr.String(), stoID)
}

type DesignStateValue struct {
	hint.BaseHinter
	Design Design
}

func NewDesignStateValue(design Design) DesignStateValue {
	return DesignStateValue{
		BaseHinter: hint.NewBaseHinter(DesignStateValueHint),
		Design:     design,
	}
}

func (sd DesignStateValue) Hint() hint.Hint {
	return sd.BaseHinter.Hint()
}

func (sd DesignStateValue) IsValid([]byte) error {
	e := util.ErrInvalid.Errorf("invalid DesignStateValue")

	if err := sd.BaseHinter.IsValid(DesignStateValueHint.Type().Bytes()); err != nil {
		return e.Wrap(err)
	}

	if err := sd.Design.IsValid(nil); err != nil {
		return e.Wrap(err)
	}

	return nil
}

func (sd DesignStateValue) HashBytes() []byte {
	return sd.Design.Bytes()
}

func StateDesignValue(st base.State) (Design, error) {
	v := st.Value()
	if v == nil {
		return Design{}, util.ErrNotFound.Errorf("sto design not found in State")
	}

	d, ok := v.(DesignStateValue)
	if !ok {
		return Design{}, errors.Errorf("invalid sto design value found, %T", v)
	}

	return d.Design, nil
}

func IsStateDesignKey(key string) bool {
	return strings.HasPrefix(key, STOPrefix) && strings.HasSuffix(key, DesignSuffix)
}

// sto:address-stoID:design
func StateKeyDesign(addr base.Address, sid currencybase.ContractID) string {
	return fmt.Sprintf("%s%s", StateKeySTOPrefix(addr, sid), DesignSuffix)
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
	return strings.HasPrefix(key, STOPrefix) && strings.HasSuffix(key, TokenHolderPartitionsSuffix)
}

func StateKeyTokenHolderPartitions(caddr base.Address, sid currencybase.ContractID, uaddr base.Address) string {
	return fmt.Sprintf("%s-%s%s", StateKeySTOPrefix(caddr, sid), uaddr.String(), TokenHolderPartitionsSuffix)
}

func StateTokenHolderPartitionsValue(st base.State) ([]Partition, error) {
	v := st.Value()
	if v == nil {
		return []Partition{}, util.ErrNotFound.Errorf("tokenholder partitions not found in State")
	}

	p, ok := v.(TokenHolderPartitionsStateValue)
	if !ok {
		return []Partition{}, errors.Errorf("invalid tokenholder partitions value found, %T", v)
	}

	return p.Partitions, nil
}

var (
	TokenHolderPartitionBalanceStateValueHint = hint.MustNewHint("mitum-sto-tokenholder-partition-balance-state-value-v0.0.1")
	TokenHolderPartitionBalanceSuffix         = ":holder-partition-balance"
)

type TokenHolderPartitionBalanceStateValue struct {
	hint.BaseHinter
	Amount    currencybase.Big
	Partition Partition
}

func NewTokenHolderPartitionBalanceStateValue(amount currencybase.Big, partition Partition) TokenHolderPartitionBalanceStateValue {
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

func StateKeyTokenHolderPartitionBalance(caddr base.Address, stoID currencybase.ContractID, uaddr base.Address, partition Partition) string {
	return fmt.Sprintf("%s-%s-%s%s", StateKeySTOPrefix(caddr, stoID), uaddr.String(), partition, TokenHolderPartitionBalanceSuffix)
}

func IsStateTokenHolderPartitionBalanceKey(key string) bool {
	return strings.HasPrefix(key, STOPrefix) && strings.HasSuffix(key, TokenHolderPartitionBalanceSuffix)
}

func StateTokenHolderPartitionBalanceValue(st base.State) (currencybase.Big, error) {
	v := st.Value()
	if v == nil {
		return currencybase.Big{}, util.ErrNotFound.Errorf("tokenholder partition balance not found in State")
	}

	p, ok := v.(TokenHolderPartitionBalanceStateValue)
	if !ok {
		return currencybase.Big{}, errors.Errorf("invalid tokenholder partition balance value found, %T", v)
	}

	return p.Amount, nil
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

func StateKeyTokenHolderPartitionOperators(caddr base.Address, stoID currencybase.ContractID, uaddr base.Address, partition Partition) string {
	return fmt.Sprintf("%s-%s-%s%s", StateKeySTOPrefix(caddr, stoID), uaddr.String(), partition.String(), TokenHolderPartitionOperatorsSuffix)
}

func IsStateTokenHolderPartitionOperatorsKey(key string) bool {
	return strings.HasPrefix(key, STOPrefix) && strings.HasSuffix(key, TokenHolderPartitionOperatorsSuffix)
}

func StateTokenHolderPartitionOperatorsValue(st base.State) ([]base.Address, error) {
	v := st.Value()
	if v == nil {
		return []base.Address{}, util.ErrNotFound.Errorf("tokenholder partition operators not found in State")
	}

	addrs, ok := v.(TokenHolderPartitionOperatorsStateValue)
	if !ok {
		return []base.Address{}, errors.Errorf("invalid tokenholder partition operators value found, %T", v)
	}

	return addrs.Operators, nil
}

var (
	PartitionBalanceStateValueHint = hint.MustNewHint("mitum-sto-partition-balance-state-value-v0.0.1")
	PartitionBalanceSuffix         = ":partition-balance"
)

type PartitionBalanceStateValue struct {
	hint.BaseHinter
	Amount currencybase.Big
}

func NewPartitionBalanceStateValue(amount currencybase.Big) PartitionBalanceStateValue {
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

func StateKeyPartitionBalance(caddr base.Address, stoID currencybase.ContractID, partition Partition) string {
	return fmt.Sprintf("%s-%s%s", StateKeySTOPrefix(caddr, stoID), partition.String(), PartitionBalanceSuffix)
}

func IsStatePartitionBalanceKey(key string) bool {
	return strings.HasPrefix(key, STOPrefix) && strings.HasSuffix(key, PartitionBalanceSuffix)
}

func StatePartitionBalanceValue(st base.State) (currencybase.Big, error) {
	v := st.Value()
	if v == nil {
		return currencybase.Big{}, util.ErrNotFound.Errorf("partition balance not found in State")
	}

	pb, ok := v.(PartitionBalanceStateValue)
	if !ok {
		return currencybase.Big{}, errors.Errorf("invalid partition balance value found, %T", v)
	}

	return pb.Amount, nil
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

func StateKeyPartitionControllers(caddr base.Address, stoID currencybase.ContractID, partition Partition) string {
	return fmt.Sprintf("%s-%s%s", StateKeySTOPrefix(caddr, stoID), partition.String(), PartitionControllersSuffix)
}

func IsStatePartitionControllersKey(key string) bool {
	return strings.HasPrefix(key, STOPrefix) && strings.HasSuffix(key, PartitionControllersSuffix)
}

var (
	OperatorTokenHoldersStateValueHint = hint.MustNewHint("mitum-sto-operator-tokenholders-state-value-v0.0.1")
	OperatorTokenHoldersSuffix         = ":operator-holders"
)

type OperatorTokenHoldersStateValue struct {
	hint.BaseHinter
	TokenHolders []base.Address
}

func NewOperatorTokenHoldersStateValue(tokenHolders []base.Address) OperatorTokenHoldersStateValue {
	return OperatorTokenHoldersStateValue{
		BaseHinter:   hint.NewBaseHinter(OperatorTokenHoldersStateValueHint),
		TokenHolders: tokenHolders,
	}
}

func (o OperatorTokenHoldersStateValue) Hint() hint.Hint {
	return o.BaseHinter.Hint()
}

func (o OperatorTokenHoldersStateValue) IsValid([]byte) error {
	e := util.ErrInvalid.Errorf("invalid OperatorTokenHoldersStateValue")

	if err := o.BaseHinter.IsValid(OperatorTokenHoldersStateValueHint.Type().Bytes()); err != nil {
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

func StateKeyOperatorTokenHolders(caddr base.Address, stoID currencybase.ContractID, oaddr base.Address, partition Partition) string {
	return fmt.Sprintf("%s-%s-%s%s", StateKeySTOPrefix(caddr, stoID), oaddr.String(), partition.String(), OperatorTokenHoldersSuffix)
}

func IsStateOperatorTokenHoldersKey(key string) bool {
	return strings.HasPrefix(key, STOPrefix) && strings.HasSuffix(key, OperatorTokenHoldersSuffix)
}

func StateOperatorTokenHoldersValue(st base.State) ([]base.Address, error) {
	v := st.Value()
	if v == nil {
		return []base.Address{}, util.ErrNotFound.Errorf("operator tokenholders not found in State")
	}

	addrs, ok := v.(OperatorTokenHoldersStateValue)
	if !ok {
		return []base.Address{}, errors.Errorf("invalid operator tokenholders value found, %T", v)
	}

	return addrs.TokenHolders, nil
}

func checkExistsState(
	key string,
	getState base.GetStateFunc,
) error {
	switch _, found, err := getState(key); {
	case err != nil:
		return err
	case !found:
		return base.NewBaseOperationProcessReasonError("state, %q does not exist", key)
	default:
		return nil
	}
}

func checkNotExistsState(
	key string,
	getState base.GetStateFunc,
) error {
	switch _, found, err := getState(key); {
	case err != nil:
		return err
	case found:
		return base.NewBaseOperationProcessReasonError("state, %q exists", key)
	default:
		return nil
	}
}

func existsState(
	k,
	name string,
	getState base.GetStateFunc,
) (base.State, error) {
	switch st, found, err := getState(k); {
	case err != nil:
		return nil, err
	case !found:
		return nil, base.NewBaseOperationProcessReasonError("%s does not exist", name)
	default:
		return st, nil
	}
}

func notExistsState(
	k,
	name string,
	getState base.GetStateFunc,
) (base.State, error) {
	var st base.State
	switch _, found, err := getState(k); {
	case err != nil:
		return nil, err
	case found:
		return nil, base.NewBaseOperationProcessReasonError("%s already exists", name)
	case !found:
		st = currencybase.NewBaseState(base.NilHeight, k, nil, nil, nil)
	}
	return st, nil
}

func existsCurrencyPolicy(cid currencybase.CurrencyID, getStateFunc base.GetStateFunc) (currencybase.CurrencyPolicy, error) {
	var policy currencybase.CurrencyPolicy
	switch i, found, err := getStateFunc(currency.StateKeyCurrencyDesign(cid)); {
	case err != nil:
		return currencybase.CurrencyPolicy{}, err
	case !found:
		return currencybase.CurrencyPolicy{}, base.NewBaseOperationProcessReasonError("currency not found, %v", cid)
	default:
		currencydesign, ok := i.Value().(currency.CurrencyDesignStateValue) //nolint:forcetypeassert //...
		if !ok {
			return currencybase.CurrencyPolicy{}, errors.Errorf("expected CurrencyDesignStateValue, not %T", i.Value())
		}
		policy = currencydesign.CurrencyDesign.Policy()
	}
	return policy, nil
}

func existsPolicy(addr base.Address, sid currencybase.ContractID, getStateFunc base.GetStateFunc) (Policy, error) {
	var policy Policy
	switch i, found, err := getStateFunc(StateKeyDesign(addr, sid)); {
	case err != nil:
		return Policy{}, err
	case !found:
		return Policy{}, base.NewBaseOperationProcessReasonError("sto not found, %s-%s", addr, sid)
	default:
		design, ok := i.Value().(DesignStateValue) //nolint:forcetypeassert //...
		if !ok {
			return Policy{}, errors.Errorf("expected DesignStateValue, not %T", i.Value())
		}
		policy = design.Design.Policy()
	}
	return policy, nil
}

func existsTokenHolderPartitions(ca base.Address, sid currencybase.ContractID, holder base.Address, getStateFunc base.GetStateFunc) ([]Partition, error) {
	var partitions []Partition
	switch i, found, err := getStateFunc(StateKeyTokenHolderPartitions(ca, sid, holder)); {
	case err != nil:
		return nil, err
	case !found:
		return nil, base.NewBaseOperationProcessReasonError("tokenholder partitions not found, %s-%s-%s", ca, sid, holder)
	default:
		pts, ok := i.Value().(TokenHolderPartitionsStateValue) //nolint:forcetypeassert //...
		if !ok {
			return nil, errors.Errorf("expected TokenHolderPartitionsStateValue, not %T", i.Value())
		}
		partitions = pts.Partitions
	}
	return partitions, nil
}

func existsTokenHolderPartitionBalance(ca base.Address, sid currencybase.ContractID, holder base.Address, p Partition, getStateFunc base.GetStateFunc) (currencybase.Big, error) {
	var balance currencybase.Big
	switch i, found, err := getStateFunc(StateKeyTokenHolderPartitionBalance(ca, sid, holder, p)); {
	case err != nil:
		return currencybase.Big{}, err
	case !found:
		return currencybase.Big{}, base.NewBaseOperationProcessReasonError("tokenholder partition balance not found, %s-%s-%s-%s", ca, sid, p, holder)
	default:
		b, ok := i.Value().(TokenHolderPartitionBalanceStateValue) //nolint:forcetypeassert //...
		if !ok {
			return currencybase.Big{}, errors.Errorf("expected TokenHolderPartitionBalanceStateValue, not %T", i.Value())
		}
		balance = b.Amount
	}
	return balance, nil
}
