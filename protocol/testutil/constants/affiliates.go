package constants

import (
	affiliates "github.com/dydxprotocol/v4-chain/protocol/x/affiliates/types"
)

var (
	DefaultAffiliateTiers = affiliates.AffiliateTiers{
		Tiers: []affiliates.AffiliateTiers_Tier{
			{
				ReqReferredVolume:   0,
				ReqStakedWholeCoins: 0,
				TakerFeeSharePpm:    50000,
			},
			{
				ReqReferredVolume:   1000000,
				ReqStakedWholeCoins: 200,
				TakerFeeSharePpm:    100000,
			},
			{
				ReqReferredVolume:   5000000,
				ReqStakedWholeCoins: 1000,
				TakerFeeSharePpm:    125000,
			},
			{
				ReqReferredVolume:   25000000,
				ReqStakedWholeCoins: 5000,
				TakerFeeSharePpm:    150000,
			},
		},
	}
)
