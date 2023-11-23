package digest

import (
	stostate "github.com/ProtoconNet/mitum-sto/state/sto"
	mitumbase "github.com/ProtoconNet/mitum2/base"
	"go.mongodb.org/mongo-driver/mongo"
)

func (bs *BlockSession) prepareSTO() error {
	if len(bs.sts) < 1 {
		return nil
	}

	var stoDesignModels []mongo.WriteModel
	var stoHolderPartitionsModels []mongo.WriteModel
	var stoHolderPartitionBalanceModels []mongo.WriteModel
	var stoHolderPartitionOperatorsModels []mongo.WriteModel
	var stoPartitionBalanceModels []mongo.WriteModel
	//var stoPartitionControllersModels []mongo.WriteModel
	var stoOperatorHoldersModels []mongo.WriteModel

	for i := range bs.sts {
		st := bs.sts[i]
		switch {
		case stostate.IsStateDesignKey(st.Key()):
			j, err := bs.handleSTODesignState(st)
			if err != nil {
				return err
			}
			stoDesignModels = append(stoDesignModels, j...)
		case stostate.IsStateTokenHolderPartitionsKey(st.Key()):
			j, err := bs.handleSTOHolderPartitionsState(st)
			if err != nil {
				return err
			}
			stoHolderPartitionsModels = append(stoHolderPartitionsModels, j...)
		case stostate.IsStateTokenHolderPartitionBalanceKey(st.Key()):
			j, err := bs.handleSTOHolderPartitionBalanceState(st)
			if err != nil {
				return err
			}
			stoHolderPartitionBalanceModels = append(stoHolderPartitionBalanceModels, j...)
		case stostate.IsStateTokenHolderPartitionOperatorsKey(st.Key()):
			j, err := bs.handleSTOHolderPartitionOperatorsState(st)
			if err != nil {
				return err
			}
			stoHolderPartitionOperatorsModels = append(stoHolderPartitionOperatorsModels, j...)
		case stostate.IsStatePartitionBalanceKey(st.Key()):
			j, err := bs.handleSTOPartitionBalanceState(st)
			if err != nil {
				return err
			}
			stoPartitionBalanceModels = append(stoPartitionBalanceModels, j...)
		//case stostate.IsStatePartitionControllersKey(st.Key()):
		//	j, err := bs.handlePartitionControllersState(st)
		//	if err != nil {
		//		return err
		//	}
		//	stoPartitionControllersModels = append(stoPartitionControllersModels, j...)
		case stostate.IsStateOperatorTokenHoldersKey(st.Key()):
			j, err := bs.handleSTOperatorHoldersState(st)
			if err != nil {
				return err
			}
			stoOperatorHoldersModels = append(stoOperatorHoldersModels, j...)
		default:
			continue
		}
	}

	bs.stoDesignModels = stoDesignModels
	bs.stoHolderPartitionsModels = stoHolderPartitionsModels
	bs.stoHolderPartitionBalanceModels = stoHolderPartitionBalanceModels
	bs.stoHolderPartitionOperatorsModels = stoHolderPartitionOperatorsModels
	bs.stoPartitionBalanceModels = stoPartitionBalanceModels
	//bs.stoPartitionControllersModels = stoPartitionControllersModels
	bs.stoOperatorHoldersModels = stoOperatorHoldersModels

	return nil
}

func (bs *BlockSession) handleSTODesignState(st mitumbase.State) ([]mongo.WriteModel, error) {
	if doc, err := NewDesignDoc(st, bs.st.DatabaseEncoder()); err != nil {
		return nil, err
	} else {
		return []mongo.WriteModel{
			mongo.NewInsertOneModel().SetDocument(doc),
		}, nil
	}
}

func (bs *BlockSession) handleSTOHolderPartitionsState(st mitumbase.State) ([]mongo.WriteModel, error) {
	if doc, err := NewHolderPartitionsDoc(st, bs.st.DatabaseEncoder()); err != nil {
		return nil, err
	} else {
		return []mongo.WriteModel{
			mongo.NewInsertOneModel().SetDocument(doc),
		}, nil
	}
}

func (bs *BlockSession) handleSTOHolderPartitionBalanceState(st mitumbase.State) ([]mongo.WriteModel, error) {
	if doc, err := NewHolderPartitionBalanceDoc(st, bs.st.DatabaseEncoder()); err != nil {
		return nil, err
	} else {
		return []mongo.WriteModel{
			mongo.NewInsertOneModel().SetDocument(doc),
		}, nil
	}
}

func (bs *BlockSession) handleSTOHolderPartitionOperatorsState(st mitumbase.State) ([]mongo.WriteModel, error) {
	if doc, err := NewHolderPartitionOperatorsDoc(st, bs.st.DatabaseEncoder()); err != nil {
		return nil, err
	} else {
		return []mongo.WriteModel{
			mongo.NewInsertOneModel().SetDocument(doc),
		}, nil
	}
}

func (bs *BlockSession) handleSTOPartitionBalanceState(st mitumbase.State) ([]mongo.WriteModel, error) {
	if doc, err := NewPartitionBalanceDoc(st, bs.st.DatabaseEncoder()); err != nil {
		return nil, err
	} else {
		return []mongo.WriteModel{
			mongo.NewInsertOneModel().SetDocument(doc),
		}, nil
	}
}

func (bs *BlockSession) handleSTOperatorHoldersState(st mitumbase.State) ([]mongo.WriteModel, error) {
	if doc, err := NewOperatorHoldersDoc(st, bs.st.DatabaseEncoder()); err != nil {
		return nil, err
	} else {
		return []mongo.WriteModel{
			mongo.NewInsertOneModel().SetDocument(doc),
		}, nil
	}
}
