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
				TakerFeeSharePpm:               50_000,
			},
			{
				ReqReferredVolumeQuoteQuantums: 1_000_000_000_000, // 1M volume
				ReqStakedWholeCoins:            200,
				TakerFeeSharePpm:               100_000,
			},
			{
				ReqReferredVolumeQuoteQuantums: 5_000_000_000_000, // 5M volume
				ReqStakedWholeCoins:            1_000,
				TakerFeeSharePpm:               125_000,
			},
			{
				ReqReferredVolumeQuoteQuantums: 25_000_000_000_000, // 25M volume
				ReqStakedWholeCoins:            5_000,
				TakerFeeSharePpm:               150_000,
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
					"dydx199tqg4wdlnu4qjlxchpd7seg454937hjrknju4", // Carl
					"dydx10fx7sy6ywd5senxae9dwytf8jxek3t2gcen2vs", // Dave
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
				ReqStakedWholeCoins:            200,
				TakerFeeSharePpm:               400_000,
			},
		},
	}

	DefaultAffiliateParameters = affiliatetypes.AffiliateParameters{
		Maximum_30DAffiliateRevenuePerReferredUserQuoteQuantums:   10_000_000_000_000, // 10M volume
		RefereeMinimumFeeTierIdx:                                  2,
		Maximum_30DAttributableVolumePerReferredUserQuoteQuantums: 10_000_000_000, // 10k volume
	}
)
