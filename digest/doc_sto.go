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

type STODesignDoc struct {
	mongodbstorage.BaseDoc
	st base.State
	de stotypes.Design
}

func NewSTODesignDoc(st base.State, enc encoder.Encoder) (STODesignDoc, error) {
	de, err := stostate.StateDesignValue(st)
	if err != nil {
		return STODesignDoc{}, err
	}
	b, err := mongodbstorage.NewBaseDoc(nil, st, enc)
	if err != nil {
		return STODesignDoc{}, err
	}

	return STODesignDoc{
		BaseDoc: b,
		st:      st,
		de:      de,
	}, nil
}

func (doc STODesignDoc) MarshalBSON() ([]byte, error) {
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

type STOHolderPartitionsDoc struct {
	mongodbstorage.BaseDoc
	st  base.State
	pts []stotypes.Partition
}

func NewSTOHolderPartitionsDoc(st base.State, enc encoder.Encoder) (STOHolderPartitionsDoc, error) {
	pts, err := stostate.StateTokenHolderPartitionsValue(st)
	if err != nil {
		return STOHolderPartitionsDoc{}, err
	}
	b, err := mongodbstorage.NewBaseDoc(nil, st, enc)
	if err != nil {
		return STOHolderPartitionsDoc{}, err
	}

	return STOHolderPartitionsDoc{
		BaseDoc: b,
		st:      st,
		pts:     pts,
	}, nil
}

func (doc STOHolderPartitionsDoc) MarshalBSON() ([]byte, error) {
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

type STOHolderPartitionBalanceDoc struct {
	mongodbstorage.BaseDoc
	st base.State
	am common.Big
}

func NewSTOHolderPartitionBalanceDoc(st base.State, enc encoder.Encoder) (STOHolderPartitionBalanceDoc, error) {
	am, err := stostate.StateTokenHolderPartitionBalanceValue(st)
	if err != nil {
		return STOHolderPartitionBalanceDoc{}, err
	}
	b, err := mongodbstorage.NewBaseDoc(nil, st, enc)
	if err != nil {
		return STOHolderPartitionBalanceDoc{}, err
	}

	return STOHolderPartitionBalanceDoc{
		BaseDoc: b,
		st:      st,
		am:      am,
	}, nil
}

func (doc STOHolderPartitionBalanceDoc) MarshalBSON() ([]byte, error) {
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

type STOHolderPartitionOperatorsDoc struct {
	mongodbstorage.BaseDoc
	st   base.State
	oprs []base.Address
}

func NewSTOHolderPartitionOperatorsDoc(st base.State, enc encoder.Encoder) (STOHolderPartitionOperatorsDoc, error) {
	oprs, err := stostate.StateTokenHolderPartitionOperatorsValue(st)
	if err != nil {
		return STOHolderPartitionOperatorsDoc{}, err
	}
	b, err := mongodbstorage.NewBaseDoc(nil, st, enc)
	if err != nil {
		return STOHolderPartitionOperatorsDoc{}, err
	}

	return STOHolderPartitionOperatorsDoc{
		BaseDoc: b,
		st:      st,
		oprs:    oprs,
	}, nil
}

func (doc STOHolderPartitionOperatorsDoc) MarshalBSON() ([]byte, error) {
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

type STOPartitionBalanceDoc struct {
	mongodbstorage.BaseDoc
	st base.State
	am common.Big
}

func NewSTOPartitionBalanceDoc(st base.State, enc encoder.Encoder) (STOPartitionBalanceDoc, error) {
	am, err := stostate.StatePartitionBalanceValue(st)
	if err != nil {
		return STOPartitionBalanceDoc{}, err
	}
	b, err := mongodbstorage.NewBaseDoc(nil, st, enc)
	if err != nil {
		return STOPartitionBalanceDoc{}, err
	}

	return STOPartitionBalanceDoc{
		BaseDoc: b,
		st:      st,
		am:      am,
	}, nil
}

func (doc STOPartitionBalanceDoc) MarshalBSON() ([]byte, error) {
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

type STOOperatorHoldersDoc struct {
	mongodbstorage.BaseDoc
	st  base.State
	hds []base.Address
}

func NewSTOOperatorHoldersDoc(st base.State, enc encoder.Encoder) (STOOperatorHoldersDoc, error) {
	hds, err := stostate.StateOperatorTokenHoldersValue(st)
	if err != nil {
		return STOOperatorHoldersDoc{}, err
	}
	b, err := mongodbstorage.NewBaseDoc(nil, st, enc)
	if err != nil {
		return STOOperatorHoldersDoc{}, err
	}

	return STOOperatorHoldersDoc{
		BaseDoc: b,
		st:      st,
		hds:     hds,
	}, nil
}

func (doc STOOperatorHoldersDoc) MarshalBSON() ([]byte, error) {
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
