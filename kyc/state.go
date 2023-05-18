package kyc

import (
	"fmt"
	"strings"

	extensioncurrency "github.com/ProtoconNet/mitum-currency-extension/v2/currency"
	"github.com/ProtoconNet/mitum-currency/v2/currency"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
	"github.com/pkg/errors"
)

var (
	KYCPrefix               = "kyc:"
	KYCDesignStateValueHint = hint.MustNewHint("mitum-kyc-design-state-value-v0.0.1")
	KYCDesignSuffix         = ":design"
	CustomerSuffix          = ":customer"
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

// kyc:{address}:{kycID}
func StateKeyKYCPrefix(addr base.Address, kycID extensioncurrency.ContractID) string {
	return fmt.Sprintf("%s%s:%s", KYCPrefix, addr.String(), kycID)
}

type KYCDesignStateValue struct {
	hint.BaseHinter
	Design KYCDesign
}

func NewKYCDesignStateValue(design KYCDesign) KYCDesignStateValue {
	return KYCDesignStateValue{
		BaseHinter: hint.NewBaseHinter(KYCDesignStateValueHint),
		Design:     design,
	}
}

func (sd KYCDesignStateValue) Hint() hint.Hint {
	return sd.BaseHinter.Hint()
}

func (sd KYCDesignStateValue) IsValid([]byte) error {
	e := util.ErrInvalid.Errorf("invalid STOKYCDesignStateValue")

	if err := sd.BaseHinter.IsValid(KYCDesignStateValueHint.Type().Bytes()); err != nil {
		return e.Wrap(err)
	}

	if err := sd.Design.IsValid(nil); err != nil {
		return e.Wrap(err)
	}

	return nil
}

func (sd KYCDesignStateValue) HashBytes() []byte {
	return sd.Design.Bytes()
}

func StateKYCDesignValue(st base.State) (KYCDesign, error) {
	v := st.Value()
	if v == nil {
		return KYCDesign{}, util.ErrNotFound.Errorf("kyc design not found in State")
	}

	d, ok := v.(KYCDesignStateValue)
	if !ok {
		return KYCDesign{}, errors.Errorf("invalid kyc design value found, %T", v)
	}

	return d.Design, nil
}

func IsStateKYCDesignKey(key string) bool {
	return strings.HasPrefix(key, KYCPrefix) && strings.HasSuffix(key, KYCDesignSuffix)
}

// kyc:{address}:{kycID}:design
func StateKeyKYCDesign(addr base.Address, sid extensioncurrency.ContractID) string {
	return fmt.Sprintf("%s%s", StateKeyKYCPrefix(addr, sid), KYCDesignSuffix)
}

type Status bool

var CustomerStateValueHint = hint.MustNewHint("mitum-kyc-customer-state-value-v0.0.1")

type CustomerStateValue struct {
	hint.BaseHinter
	status Status
}

func NewCustomerStateValue(status Status) CustomerStateValue {
	return CustomerStateValue{
		BaseHinter: hint.NewBaseHinter(CustomerStateValueHint),
		status:     status,
	}
}

func (sd CustomerStateValue) Hint() hint.Hint {
	return sd.BaseHinter.Hint()
}

func (sd CustomerStateValue) IsValid([]byte) error {
	e := util.ErrInvalid.Errorf("invalid kyc CustomerStateValue")

	if err := sd.BaseHinter.IsValid(KYCDesignStateValueHint.Type().Bytes()); err != nil {
		return e.Wrap(err)
	}

	return nil
}

func (sd CustomerStateValue) HashBytes() []byte {
	var v int8
	if sd.status {
		v = 1
	}
	return []byte{byte(v)}
}

func StateCustomerValue(st base.State) (*Status, error) {
	v := st.Value()
	if v == nil {
		return nil, util.ErrNotFound.Errorf("kyc customer not found in State")
	}

	d, ok := v.(CustomerStateValue)
	if !ok {
		return nil, errors.Errorf("invalid kyc customer value found, %T", v)
	}

	return &d.status, nil
}

func IsStateCustomerKey(key string) bool {
	return strings.HasPrefix(key, KYCPrefix) && strings.HasSuffix(key, CustomerSuffix)
}

// kyc:{address}:{kycID}:{address}:customer
func StateKeyCustomer(addr base.Address, sid extensioncurrency.ContractID, customer base.Address) string {
	return fmt.Sprintf("%s:%s%s", StateKeyKYCPrefix(addr, sid), customer.String(), CustomerSuffix)
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
		st = currency.NewBaseState(base.NilHeight, k, nil, nil, nil)
	}
	return st, nil
}

func existsCurrencyPolicy(cid currency.CurrencyID, getStateFunc base.GetStateFunc) (extensioncurrency.CurrencyPolicy, error) {
	var policy extensioncurrency.CurrencyPolicy
	switch i, found, err := getStateFunc(extensioncurrency.StateKeyCurrencyDesign(cid)); {
	case err != nil:
		return extensioncurrency.CurrencyPolicy{}, err
	case !found:
		return extensioncurrency.CurrencyPolicy{}, base.NewBaseOperationProcessReasonError("currency not found, %v", cid)
	default:
		currencydesign, ok := i.Value().(extensioncurrency.CurrencyDesignStateValue) //nolint:forcetypeassert //...
		if !ok {
			return extensioncurrency.CurrencyPolicy{}, errors.Errorf("expected CurrencyDesignStateValue, not %T", i.Value())
		}
		policy = currencydesign.CurrencyDesign.Policy()
	}
	return policy, nil
}

func existsKYCPolicy(addr base.Address, kycid extensioncurrency.ContractID, getStateFunc base.GetStateFunc) (KYCPolicy, error) {
	var policy KYCPolicy
	switch i, found, err := getStateFunc(StateKeyKYCDesign(addr, kycid)); {
	case err != nil:
		return KYCPolicy{}, err
	case !found:
		return KYCPolicy{}, base.NewBaseOperationProcessReasonError("kyc not found, %s-%s", addr, kycid)
	default:
		design, ok := i.Value().(KYCDesignStateValue) //nolint:forcetypeassert //...
		if !ok {
			return KYCPolicy{}, errors.Errorf("expected KYCDesignStateValue, not %T", i.Value())
		}
		policy = design.Design.Policy()
	}
	return policy, nil
}
