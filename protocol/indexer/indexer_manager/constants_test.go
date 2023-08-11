package indexer_manager_test

import (
	"time"

	indexerevents "github.com/dydxprotocol/v4/indexer/events"
	"github.com/dydxprotocol/v4/indexer/indexer_manager"
	"github.com/dydxprotocol/v4/testutil/constants"
	perptypes "github.com/dydxprotocol/v4/x/perpetuals/types"
	satypes "github.com/dydxprotocol/v4/x/subaccounts/types"
)

const (
	TxHash      = "txHash"
	TxHash1     = "txHash1"
	Data        = "data"
	Data2       = "data2"
	Data3       = "data3"
	BlockHeight = int64(5)
	ConsumedGas = uint64(100)
)

var BlockTime = time.Unix(1650000000, 0).UTC()

var OrderFillTendermintEvent = indexer_manager.IndexerTendermintEvent{
	Subtype: indexerevents.SubtypeOrderFill,
	Data:    Data3,
	OrderingWithinBlock: &indexer_manager.IndexerTendermintEvent_TransactionIndex{
		TransactionIndex: 0,
	},
	EventIndex: 0,
}

var TransferTendermintEvent = indexer_manager.IndexerTendermintEvent{
	Subtype: indexerevents.SubtypeTransfer,
	Data:    Data,
	OrderingWithinBlock: &indexer_manager.IndexerTendermintEvent_TransactionIndex{
		TransactionIndex: 0,
	},
	EventIndex: 1,
}

var SubaccountTendermintEvent = indexer_manager.IndexerTendermintEvent{
	Subtype: indexerevents.SubtypeSubaccountUpdate,
	Data:    Data2,
	OrderingWithinBlock: &indexer_manager.IndexerTendermintEvent_TransactionIndex{
		TransactionIndex: 1,
	},
	EventIndex: 0,
}

var OrderFillEvent = indexerevents.OrderFillEvent{
	MakerOrder: constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15,
	TakerOrder: &indexerevents.OrderFillEvent_Order{
		Order: &constants.Order_Alice_Num0_Id2_Clob1_Sell5_Price10_GTB15,
	},
	FillAmount: 5,
}

var FundingRateEvent = indexerevents.FundingEvent{
	Type: indexerevents.FundingEvent_TYPE_FUNDING_RATE,
	Values: []perptypes.FundingPremium{
		{
			PerpetualId: 0,
			PremiumPpm:  -1000,
		},
		{
			PerpetualId: 1,
			PremiumPpm:  0,
		},
		{
			PerpetualId: 2,
			PremiumPpm:  5000,
		},
	},
}

var FundingPremiumSampleEvent = indexerevents.FundingEvent{
	Type: indexerevents.FundingEvent_TYPE_PREMIUM_SAMPLE,
	Values: []perptypes.FundingPremium{
		{
			PerpetualId: 0,
			PremiumPpm:  1000,
		},
		{
			PerpetualId: 1,
			PremiumPpm:  0,
		},
	},
}

var SubaccountEvent = indexerevents.SubaccountUpdateEvent{
	SubaccountId: &constants.Alice_Num0,
	UpdatedPerpetualPositions: []*satypes.PerpetualPosition{
		&constants.Long_Perp_1BTC_PositiveFunding,
		&constants.Short_Perp_1ETH_NegativeFunding,
	},
	UpdatedAssetPositions: []*satypes.AssetPosition{
		&constants.Short_Asset_1BTC,
		&constants.Long_Asset_1ETH,
	},
}

var TransferEvent = indexerevents.TransferEvent{

	SenderSubaccountId:    constants.Alice_Num0,
	RecipientSubaccountId: constants.Alice_Num1,
	Amount:                uint64(5),
	AssetId:               uint32(0),
}
