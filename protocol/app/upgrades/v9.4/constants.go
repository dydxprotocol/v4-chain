package v_9_4

import (
	store "cosmossdk.io/store/types"
	"github.com/dydxprotocol/v4-chain/protocol/app/upgrades"
	affiliatetypes "github.com/dydxprotocol/v4-chain/protocol/x/affiliates/types"
)

const (
	UpgradeName = "v9.4"
)

var (
	Upgrade = upgrades.Upgrade{
		UpgradeName:   UpgradeName,
		StoreUpgrades: store.StoreUpgrades{},
	}

	PreviousAffilliateTiers = affiliatetypes.AffiliateTiers{
		Tiers: []affiliatetypes.AffiliateTiers_Tier{
			{
				ReqReferredVolumeQuoteQuantums: 0,
				ReqStakedWholeCoins:            0,
				TakerFeeSharePpm:               200_000,
			},
			{
				ReqReferredVolumeQuoteQuantums: 2_000_000_000_000, // 1M volume
				ReqStakedWholeCoins:            0,
				TakerFeeSharePpm:               300_000,
			},
		},
	}

	PreviousAffiliateParameters = affiliatetypes.AffiliateParameters{
		Maximum_30DAffiliateRevenuePerReferredUserQuoteQuantums:   1_000_000_000_000, // 10M volume
		RefereeMinimumFeeTierIdx:                                  1,
		Maximum_30DAttributableVolumePerReferredUserQuoteQuantums: 1_000_000_000, // 10k volume
	}

	PreviousAffiliateWhitelist = affiliatetypes.AffiliateWhitelist{
		Tiers: []affiliatetypes.AffiliateWhitelist_Tier{
			{
				Addresses: []string{
					"dydx1affiliate1",
					"dydx1affiliate2",
				},
				TakerFeeSharePpm: 200_000,
			},
		},
	}

	DefaultAffiliateTiers = affiliatetypes.AffiliateTiers{
		Tiers: []affiliatetypes.AffiliateTiers_Tier{
			{
				ReqReferredVolumeQuoteQuantums: 0,
				ReqStakedWholeCoins:            0,
				TakerFeeSharePpm:               400_000,
			},
			{
				ReqReferredVolumeQuoteQuantums: 1_000_000_000_000, // 1M volume
				ReqStakedWholeCoins:            0,
				TakerFeeSharePpm:               500_000,
			},
		},
	}

	DefaultAffiliateParameters = affiliatetypes.AffiliateParameters{
		Maximum_30DAffiliateRevenuePerReferredUserQuoteQuantums:   10_000_000_000_000, // 10M volume
		RefereeMinimumFeeTierIdx:                                  2,
		Maximum_30DAttributableVolumePerReferredUserQuoteQuantums: 10_000_000_000, // 10k volume
	}
)
