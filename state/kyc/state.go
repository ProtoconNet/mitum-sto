package kyc

import (
	"fmt"
	"strings"

	currencytypes "github.com/ProtoconNet/mitum-currency/v3/types"
	kyctypes "github.com/ProtoconNet/mitum-sto/types/kyc"
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

// kyc:{address}:{kycID}
func StateKeyKYCPrefix(addr base.Address, kycID currencytypes.ContractID) string {
	return fmt.Sprintf("%s%s:%s", KYCPrefix, addr.String(), kycID)
}

type DesignStateValue struct {
	hint.BaseHinter
	Design kyctypes.Design
}

func NewDesignStateValue(design kyctypes.Design) DesignStateValue {
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

func StateDesignValue(st base.State) (kyctypes.Design, error) {
	v := st.Value()
	if v == nil {
		return kyctypes.Design{}, util.ErrNotFound.Errorf("kyc design not found in State")
	}

	d, ok := v.(DesignStateValue)
	if !ok {
		return kyctypes.Design{}, errors.Errorf("invalid kyc design value found, %T", v)
	}

	return d.Design, nil
}

func IsStateDesignKey(key string) bool {
	return strings.HasPrefix(key, KYCPrefix) && strings.HasSuffix(key, DesignSuffix)
}

// kyc:{address}:{kycID}:design
func StateKeyDesign(addr base.Address, sid currencytypes.ContractID) string {
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

	if err := sd.BaseHinter.IsValid(CustomerStateValueHint.Type().Bytes()); err != nil {
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
func StateKeyCustomer(addr base.Address, sid currencytypes.ContractID, customer base.Address) string {
	return fmt.Sprintf("%s:%s%s", StateKeyKYCPrefix(addr, sid), customer.String(), CustomerSuffix)
}

func ExistsPolicy(addr base.Address, kycid currencytypes.ContractID, getStateFunc base.GetStateFunc) (kyctypes.Policy, error) {
	var policy kyctypes.Policy
	switch i, found, err := getStateFunc(StateKeyDesign(addr, kycid)); {
	case err != nil:
		return kyctypes.Policy{}, err
	case !found:
		return kyctypes.Policy{}, base.NewBaseOperationProcessReasonError("kyc not found, %s-%s", addr, kycid)
	default:
		design, ok := i.Value().(DesignStateValue) //nolint:forcetypeassert //...
		if !ok {
			return kyctypes.Policy{}, errors.Errorf("expected DesignStateValue, not %T", i.Value())
		}
		policy = design.Design.Policy()
	}
	return policy, nil
}
