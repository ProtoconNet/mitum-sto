package digest

import (
	"github.com/ProtoconNet/mitum-currency/v3/common"
	mongodbstorage "github.com/ProtoconNet/mitum-currency/v3/digest/mongodb"
	bsonenc "github.com/ProtoconNet/mitum-currency/v3/digest/util/bson"
	stostate "github.com/ProtoconNet/mitum-sto/state/sto"
	stotypes "github.com/ProtoconNet/mitum-sto/types/sto"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util/encoder"
)

type DesignDoc struct {
	mongodbstorage.BaseDoc
	st base.State
	de stotypes.Design
}

func NewDesignDoc(st base.State, enc encoder.Encoder) (DesignDoc, error) {
	de, err := stostate.StateDesignValue(st)
	if err != nil {
		return DesignDoc{}, err
	}
	b, err := mongodbstorage.NewBaseDoc(nil, st, enc)
	if err != nil {
		return DesignDoc{}, err
	}

	return DesignDoc{
		BaseDoc: b,
		st:      st,
		de:      de,
	}, nil
}

func (doc DesignDoc) MarshalBSON() ([]byte, error) {
	m, err := doc.BaseDoc.M()
	if err != nil {
		return nil, err
	}

	parsedKey, err := stostate.ParseStateKey(doc.st.Key(), stostate.STOPrefix)
	m["contract"] = parsedKey[1]
	m["height"] = doc.st.Height()
	m["design"] = doc.de

	return bsonenc.Marshal(m)
}

type HolderPartitionsDoc struct {
	mongodbstorage.BaseDoc
	st  base.State
	pts []stotypes.Partition
}

func NewHolderPartitionsDoc(st base.State, enc encoder.Encoder) (HolderPartitionsDoc, error) {
	pts, err := stostate.StateTokenHolderPartitionsValue(st)
	if err != nil {
		return HolderPartitionsDoc{}, err
	}
	b, err := mongodbstorage.NewBaseDoc(nil, st, enc)
	if err != nil {
		return HolderPartitionsDoc{}, err
	}

	return HolderPartitionsDoc{
		BaseDoc: b,
		st:      st,
		pts:     pts,
	}, nil
}

func (doc HolderPartitionsDoc) MarshalBSON() ([]byte, error) {
	m, err := doc.BaseDoc.M()
	if err != nil {
		return nil, err
	}

	parsedKey, err := stostate.ParseStateKey(doc.st.Key(), stostate.STOPrefix)
	m["contract"] = parsedKey[1]
	m["holder"] = parsedKey[2]
	m["height"] = doc.st.Height()
	m["partitions"] = doc.pts

	return bsonenc.Marshal(m)
}

type HolderPartitionBalanceDoc struct {
	mongodbstorage.BaseDoc
	st base.State
	am common.Big
}

func NewHolderPartitionBalanceDoc(st base.State, enc encoder.Encoder) (HolderPartitionBalanceDoc, error) {
	am, err := stostate.StateTokenHolderPartitionBalanceValue(st)
	if err != nil {
		return HolderPartitionBalanceDoc{}, err
	}
	b, err := mongodbstorage.NewBaseDoc(nil, st, enc)
	if err != nil {
		return HolderPartitionBalanceDoc{}, err
	}

	return HolderPartitionBalanceDoc{
		BaseDoc: b,
		st:      st,
		am:      am,
	}, nil
}

func (doc HolderPartitionBalanceDoc) MarshalBSON() ([]byte, error) {
	m, err := doc.BaseDoc.M()
	if err != nil {
		return nil, err
	}

	parsedKey, err := stostate.ParseStateKey(doc.st.Key(), stostate.STOPrefix)
	m["contract"] = parsedKey[1]
	m["holder"] = parsedKey[2]
	m["partition"] = parsedKey[3]
	m["height"] = doc.st.Height()
	m["balance"] = doc.am

	return bsonenc.Marshal(m)
}

type HolderPartitionOperatorsDoc struct {
	mongodbstorage.BaseDoc
	st   base.State
	oprs []base.Address
}

func NewHolderPartitionOperatorsDoc(st base.State, enc encoder.Encoder) (HolderPartitionOperatorsDoc, error) {
	oprs, err := stostate.StateTokenHolderPartitionOperatorsValue(st)
	if err != nil {
		return HolderPartitionOperatorsDoc{}, err
	}
	b, err := mongodbstorage.NewBaseDoc(nil, st, enc)
	if err != nil {
		return HolderPartitionOperatorsDoc{}, err
	}

	return HolderPartitionOperatorsDoc{
		BaseDoc: b,
		st:      st,
		oprs:    oprs,
	}, nil
}

func (doc HolderPartitionOperatorsDoc) MarshalBSON() ([]byte, error) {
	m, err := doc.BaseDoc.M()
	if err != nil {
		return nil, err
	}

	parsedKey, err := stostate.ParseStateKey(doc.st.Key(), stostate.STOPrefix)
	m["contract"] = parsedKey[1]
	m["holder"] = parsedKey[2]
	m["partition"] = parsedKey[3]
	m["height"] = doc.st.Height()
	m["operators"] = doc.oprs

	return bsonenc.Marshal(m)
}

type PartitionBalanceDoc struct {
	mongodbstorage.BaseDoc
	st base.State
	am common.Big
}

func NewPartitionBalanceDoc(st base.State, enc encoder.Encoder) (PartitionBalanceDoc, error) {
	am, err := stostate.StatePartitionBalanceValue(st)
	if err != nil {
		return PartitionBalanceDoc{}, err
	}
	b, err := mongodbstorage.NewBaseDoc(nil, st, enc)
	if err != nil {
		return PartitionBalanceDoc{}, err
	}

	return PartitionBalanceDoc{
		BaseDoc: b,
		st:      st,
		am:      am,
	}, nil
}

func (doc PartitionBalanceDoc) MarshalBSON() ([]byte, error) {
	m, err := doc.BaseDoc.M()
	if err != nil {
		return nil, err
	}

	parsedKey, err := stostate.ParseStateKey(doc.st.Key(), stostate.STOPrefix)
	m["contract"] = parsedKey[1]
	m["partition"] = parsedKey[2]
	m["height"] = doc.st.Height()
	m["balance"] = doc.am

	return bsonenc.Marshal(m)
}

type OperatorHoldersDoc struct {
	mongodbstorage.BaseDoc
	st  base.State
	hds []base.Address
}

func NewOperatorHoldersDoc(st base.State, enc encoder.Encoder) (OperatorHoldersDoc, error) {
	hds, err := stostate.StateOperatorTokenHoldersValue(st)
	if err != nil {
		return OperatorHoldersDoc{}, err
	}
	b, err := mongodbstorage.NewBaseDoc(nil, st, enc)
	if err != nil {
		return OperatorHoldersDoc{}, err
	}

	return OperatorHoldersDoc{
		BaseDoc: b,
		st:      st,
		hds:     hds,
	}, nil
}

func (doc OperatorHoldersDoc) MarshalBSON() ([]byte, error) {
	m, err := doc.BaseDoc.M()
	if err != nil {
		return nil, err
	}

	parsedKey, err := stostate.ParseStateKey(doc.st.Key(), stostate.STOPrefix)
	m["contract"] = parsedKey[1]
	m["operator"] = parsedKey[2]
	m["height"] = doc.st.Height()
	m["operators"] = doc.hds

	return bsonenc.Marshal(m)
}
