package types

var (
	DefaultAffiliateTiers = AffiliateTiers{
		Tiers: []AffiliateTiers_Tier{
			{
				ReqReferredVolumeQuoteQuantums: 0,
				ReqStakedWholeCoins:            0,
				TakerFeeSharePpm:               50_000,
			},
			{
				ReqReferredVolumeQuoteQuantums: 1_000_000_000_000,
				ReqStakedWholeCoins:            200,
				TakerFeeSharePpm:               100_000,
			},
			{
				ReqReferredVolumeQuoteQuantums: 5_000_000_000_000,
				ReqStakedWholeCoins:            1_000,
				TakerFeeSharePpm:               125_000,
			},
			{
				ReqReferredVolumeQuoteQuantums: 25_000_000_000_000,
				ReqStakedWholeCoins:            5_000,
				TakerFeeSharePpm:               150_000,
			},
		},
	}
)
