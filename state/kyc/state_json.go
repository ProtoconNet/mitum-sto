package kyc

import (
	"encoding/json"

	kyctypes "github.com/ProtoconNet/mitum-sto/types/kyc"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
	"github.com/ProtoconNet/mitum2/util/hint"
)

type DesignStateValueJSONMarshaler struct {
	hint.BaseHinter
	KYC kyctypes.Design `json:"kyc"`
}

func (de DesignStateValue) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(DesignStateValueJSONMarshaler{
		BaseHinter: de.BaseHinter,
		KYC:        de.Design,
	})
}

type DesignStateValueJSONUnmarshaler struct {
	KYC json.RawMessage `json:"kyc"`
}

func (de *DesignStateValue) DecodeJSON(b []byte, enc encoder.Encoder) error {
	e := util.StringError("failed to decode json of DesignStateValue")

	var u DesignStateValueJSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return e.Wrap(err)
	}

	var design kyctypes.Design

	if err := design.DecodeJSON(u.KYC, enc); err != nil {
		return e.Wrap(err)
	}

	de.Design = design

	return nil
}

type CustomerStateValueJSONMarshaler struct {
	hint.BaseHinter
	Status Status `json:"status"`
}

func (cm CustomerStateValue) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(CustomerStateValueJSONMarshaler{
		BaseHinter: cm.BaseHinter,
		Status:     cm.status,
	})
}

type CustomerStateValueJSONUnmarshaler struct {
	Status bool `json:"status"`
}

func (cm *CustomerStateValue) DecodeJSON(b []byte, enc encoder.Encoder) error {
	e := util.StringError("failed to decode json of CustomerStateValue")

	var u CustomerStateValueJSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return e.Wrap(err)
	}
	cm.status = Status(u.Status)

	return nil
}
