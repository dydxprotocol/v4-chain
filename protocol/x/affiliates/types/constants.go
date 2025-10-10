package types

var (
	DefaultAffiliateTiers = AffiliateTiers{
		Tiers: []AffiliateTiers_Tier{
			{
				ReqReferredVolumeQuoteQuantums: 0,      // 0 USDC
				ReqStakedWholeCoins:            0,      // 0 coins staked
				TakerFeeSharePpm:               50_000, // 5%
			},
			{
				ReqReferredVolumeQuoteQuantums: 1_000_000_000_000, // 1 million USDC
				ReqStakedWholeCoins:            200,               // 200 whole coins
				TakerFeeSharePpm:               100_000,           // 10%
			},
			{
				ReqReferredVolumeQuoteQuantums: 5_000_000_000_000, // 5 million USDC
				ReqStakedWholeCoins:            1_000,             // 1000 whole coins
				TakerFeeSharePpm:               125_000,           // 12.5%
			},
			{
				ReqReferredVolumeQuoteQuantums: 25_000_000_000_000, // 25 million USDC
				ReqStakedWholeCoins:            5_000,              // 5000 whole coins
				TakerFeeSharePpm:               150_000,            // 15%
			},
			{
				ReqReferredVolumeQuoteQuantums: 50_000_000_000_000, // 50 million USDC
				ReqStakedWholeCoins:            100_000_000,        // 100m whole coins
				TakerFeeSharePpm:               250_000,            // 25%
			},
		},
	}

	DefaultAffiliateParameters = AffiliateParameters{
		Maximum_30DAttributableVolumePerReferredUserQuoteQuantums: 100_000_000_000_000, // 100 million USDC
		RefereeMinimumFeeTierIdx:                                  2,
		Maximum_30DAffiliateRevenuePerReferredUserQuoteQuantums:   10_000_000_000, // 10k commission
	}

	AffiliatesRevSharePpmCap = uint32(500_000) // 50%
)
