package constants

import (
	"github.com/dydxprotocol/v4-chain/protocol/x/feetiers/types"
)

var PerpetualFeeParams = types.PerpetualFeeParams{
	Tiers: []*types.PerpetualFeeTier{
		{
			Name:        "1",
			MakerFeePpm: 200,
			TakerFeePpm: 500,
		},
	},
}

var PerpetualFeeParamsMakerRebate = types.PerpetualFeeParams{
	Tiers: []*types.PerpetualFeeTier{
		{
			Name:        "1",
			MakerFeePpm: -200,
			TakerFeePpm: 500,
		},
	},
}

var PerpetualFeeParamsNoFee = types.PerpetualFeeParams{
	Tiers: []*types.PerpetualFeeTier{
		{
			Name:        "1",
			MakerFeePpm: 0,
			TakerFeePpm: 0,
		},
	},
}
