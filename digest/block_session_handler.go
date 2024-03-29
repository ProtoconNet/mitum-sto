package digest

import (
	crcydigest "github.com/ProtoconNet/mitum-currency/v3/digest"
	"github.com/ProtoconNet/mitum2/base"
	"go.mongodb.org/mongo-driver/mongo"
)

func (bs *BlockSession) handleAccountState(st base.State) ([]mongo.WriteModel, error) {
	if rs, err := crcydigest.NewAccountValue(st); err != nil {
		return nil, err
	} else if doc, err := crcydigest.NewAccountDoc(rs, bs.st.DatabaseEncoder()); err != nil {
		return nil, err
	} else {
		return []mongo.WriteModel{mongo.NewInsertOneModel().SetDocument(doc)}, nil
	}
}

func (bs *BlockSession) handleBalanceState(st base.State) ([]mongo.WriteModel, string, error) {
	doc, address, err := crcydigest.NewBalanceDoc(st, bs.st.DatabaseEncoder())
	if err != nil {
		return nil, "", err
	}
	return []mongo.WriteModel{mongo.NewInsertOneModel().SetDocument(doc)}, address, nil
}

func (bs *BlockSession) handleContractAccountState(st base.State) ([]mongo.WriteModel, error) {
	doc, err := crcydigest.NewContractAccountStatusDoc(st, bs.st.DatabaseEncoder())
	if err != nil {
		return nil, err
	}
	return []mongo.WriteModel{mongo.NewInsertOneModel().SetDocument(doc)}, nil
}

func (bs *BlockSession) handleCurrencyState(st base.State) ([]mongo.WriteModel, error) {
	doc, err := crcydigest.NewCurrencyDoc(st, bs.st.DatabaseEncoder())
	if err != nil {
		return nil, err
	}
	return []mongo.WriteModel{mongo.NewInsertOneModel().SetDocument(doc)}, nil
}
