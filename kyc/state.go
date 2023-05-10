package kyc

import (
	"fmt"
	"strings"

	extensioncurrency "github.com/ProtoconNet/mitum-currency-extension/v2/currency"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
	"github.com/pkg/errors"
)

var (
	KYCPrefix            = "kyc:"
	DesignStateValueHint = hint.MustNewHint("mitum-kyc-design-state-value-v0.0.1")
	DesignSuffix         = ":design"
	CustomerSuffix       = ":customer"
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
	e := util.ErrInvalid.Errorf("invalid STODesignStateValue")

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
		return Design{}, util.ErrNotFound.Errorf("kyc design not found in State")
	}

	d, ok := v.(DesignStateValue)
	if !ok {
		return Design{}, errors.Errorf("invalid kyc design value found, %T", v)
	}

	return d.Design, nil
}

func IsStateDesignKey(key string) bool {
	return strings.HasPrefix(key, KYCPrefix) && strings.HasSuffix(key, DesignSuffix)
}

// kyc:{address}:{kycID}:design
func StateKeyDesign(addr base.Address, sid extensioncurrency.ContractID) string {
	return fmt.Sprintf("%s%s", StateKeyKYCPrefix(addr, sid), DesignSuffix)
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

	if err := sd.BaseHinter.IsValid(DesignStateValueHint.Type().Bytes()); err != nil {
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
