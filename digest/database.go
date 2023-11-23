package digest

import (
	"github.com/ProtoconNet/mitum-currency/v3/common"
	currencydigest "github.com/ProtoconNet/mitum-currency/v3/digest"
	"github.com/ProtoconNet/mitum-currency/v3/digest/util"
	stostate "github.com/ProtoconNet/mitum-sto/state/sto"
	stotypes "github.com/ProtoconNet/mitum-sto/types/sto"
	mitumbase "github.com/ProtoconNet/mitum2/base"
	mitumutil "github.com/ProtoconNet/mitum2/util"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	defaultColNameAccount                     = "digest_ac"
	defaultColNameContractAccount             = "digest_ca"
	defaultColNameBalance                     = "digest_bl"
	defaultColNameCurrency                    = "digest_cr"
	defaultColNameOperation                   = "digest_op"
	defaultColNameBlock                       = "digest_bm"
	defaultColNameSTO                         = "digest_sto_de"
	defaultColNameSTOHolderPartitions         = "digest_sto_hac_pt"
	defaultColNameSTOHolderPartitionBalance   = "digest_sto_hac_pt_bl"
	defaultColNameSTOHolderPartitionOperators = "digest_sto_hac_pt_oac"
	defaultColNameSTOPartitionBalance         = "digest_sto_pt_bl"
	defaultColNameSTOPartitionControllers     = "digest_sto_pt_cac"
	defaultColNameSTOOperatorHolders          = "digest_sto_oac_hac"
)

func STOService(
	st *currencydigest.Database,
	contract string,
) (*stotypes.Design, error) {
	filter := util.NewBSONFilter("contract", contract)

	var design stotypes.Design
	var sta mitumbase.State
	var err error
	if err := st.DatabaseClient().GetByFilter(
		defaultColNameSTO,
		filter.D(),
		func(res *mongo.SingleResult) error {
			sta, err = currencydigest.LoadState(res.Decode, st.DatabaseEncoders())
			if err != nil {
				return err
			}

			design, err = stostate.StateDesignValue(sta)
			if err != nil {
				return err
			}

			return nil
		},
		options.FindOne().SetSort(util.NewBSONFilter("height", -1).D()),
	); err != nil {
		return nil, err
	}

	return &design, nil
}

func HolderPartitions(
	st *currencydigest.Database,
	contract,
	holder string,
) ([]stotypes.Partition, error) {
	filter := util.NewBSONFilter("contract", contract)
	filter = filter.Add("holder", holder)

	var partitions []stotypes.Partition
	var sta mitumbase.State
	var err error
	if err = st.DatabaseClient().GetByFilter(
		defaultColNameSTOHolderPartitions,
		filter.D(),
		func(res *mongo.SingleResult) error {
			sta, err = currencydigest.LoadState(res.Decode, st.DatabaseEncoders())
			if err != nil {
				return err
			}
			partitions, err = stostate.StateTokenHolderPartitionsValue(sta)
			if err != nil {
				return err
			}

			return nil
		},
		options.FindOne().SetSort(util.NewBSONFilter("height", -1).D()),
	); err != nil {
		return nil, err
	}

	return partitions, nil
}

func HolderPartitionBalance(
	st *currencydigest.Database,
	contract,
	holder,
	partition string,
) (common.Big, error) {
	filter := util.NewBSONFilter("contract", contract)
	filter = filter.Add("holder", holder)
	filter = filter.Add("partition", partition)

	var amount common.Big
	var sta mitumbase.State
	var err error
	if err := st.DatabaseClient().GetByFilter(
		defaultColNameSTOHolderPartitionBalance,
		filter.D(),
		func(res *mongo.SingleResult) error {
			sta, err = currencydigest.LoadState(res.Decode, st.DatabaseEncoders())
			if err != nil {
				return err
			}

			amount, err = stostate.StateTokenHolderPartitionBalanceValue(sta)
			if err != nil {
				return err
			}

			return nil
		},
		options.FindOne().SetSort(util.NewBSONFilter("height", -1).D()),
	); err != nil {
		return common.NilBig, mitumutil.ErrNotFound.Errorf(
			"sto holder partition balance by contract %s, account %s",
			contract,
			holder,
		)
	}

	return amount, nil
}

func HolderPartitionOperators(
	st *currencydigest.Database,
	contract,
	holder,
	partition string,
) ([]mitumbase.Address, error) {
	filter := util.NewBSONFilter("contract", contract)
	filter = filter.Add("holder", holder)
	filter = filter.Add("partition", partition)

	var operators []mitumbase.Address
	var sta mitumbase.State
	var err error
	if err = st.DatabaseClient().GetByFilter(
		defaultColNameSTOHolderPartitionOperators,
		filter.D(),
		func(res *mongo.SingleResult) error {
			sta, err = currencydigest.LoadState(res.Decode, st.DatabaseEncoders())
			if err != nil {
				return err
			}
			operators, err = stostate.StateTokenHolderPartitionOperatorsValue(sta)
			if err != nil {
				return err
			}

			return nil
		},
		options.FindOne().SetSort(util.NewBSONFilter("height", -1).D()),
	); err != nil {
		return nil, err
	}

	return operators, nil
}

func PartitionBalance(
	st *currencydigest.Database,
	contract,
	partition string,
) (common.Big, error) {
	filter := util.NewBSONFilter("contract", contract)
	filter = filter.Add("partition", partition)

	var amount common.Big
	var sta mitumbase.State
	var err error
	if err := st.DatabaseClient().GetByFilter(
		defaultColNameSTOPartitionBalance,
		filter.D(),
		func(res *mongo.SingleResult) error {
			sta, err = currencydigest.LoadState(res.Decode, st.DatabaseEncoders())
			if err != nil {
				return err
			}

			amount, err = stostate.StatePartitionBalanceValue(sta)
			if err != nil {
				return err
			}

			return nil
		},
		options.FindOne().SetSort(util.NewBSONFilter("height", -1).D()),
	); err != nil {
		return common.NilBig, mitumutil.ErrNotFound.Errorf(
			"sto partition balance by contract %s, account %s",
			contract,
			partition,
		)
	}

	return amount, nil
}

func OperatorHolders(
	st *currencydigest.Database,
	contract,
	operator string,
) ([]mitumbase.Address, error) {
	filter := util.NewBSONFilter("contract", contract)
	filter = filter.Add("operator", operator)

	var holders []mitumbase.Address
	var sta mitumbase.State
	var err error
	if err = st.DatabaseClient().GetByFilter(
		defaultColNameSTOOperatorHolders,
		filter.D(),
		func(res *mongo.SingleResult) error {
			sta, err = currencydigest.LoadState(res.Decode, st.DatabaseEncoders())
			if err != nil {
				return err
			}
			holders, err = stostate.StateOperatorTokenHoldersValue(sta)
			if err != nil {
				return err
			}

			return nil
		},
		options.FindOne().SetSort(util.NewBSONFilter("height", -1).D()),
	); err != nil {
		return nil, err
	}

	return holders, nil
}
