package constants

import (
	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
)

const (
	MakerFeePpm = uint32(200) // 0.02% fee.
	TakerFeePpm = uint32(500) // 0.05% fee.
)

var (
	ClobPair_Spot_Btc = clobtypes.ClobPair{
		Id: 1000,
		Metadata: &clobtypes.ClobPair_SpotClobMetadata{
			SpotClobMetadata: &clobtypes.SpotClobMetadata{
				BaseAssetId:  0,
				QuoteAssetId: 0,
			},
		},
		StepBaseQuantums:          10,
		SubticksPerTick:           100,
		MinOrderBaseQuantums:      10,
		QuantumConversionExponent: -8,
		MakerFeePpm:               MakerFeePpm,
		TakerFeePpm:               TakerFeePpm,
		Status:                    clobtypes.ClobPair_STATUS_ACTIVE,
	}
	ClobPair_Btc = clobtypes.ClobPair{
		Id: 0,
		Metadata: &clobtypes.ClobPair_PerpetualClobMetadata{
			PerpetualClobMetadata: &clobtypes.PerpetualClobMetadata{
				PerpetualId: 0,
			},
		},
		StepBaseQuantums:          5,
		SubticksPerTick:           5,
		MinOrderBaseQuantums:      5,
		QuantumConversionExponent: -8,
		MakerFeePpm:               MakerFeePpm,
		TakerFeePpm:               TakerFeePpm,
		Status:                    clobtypes.ClobPair_STATUS_ACTIVE,
	}
	ClobPair_Btc_No_Fee = clobtypes.ClobPair{
		Id: 0,
		Metadata: &clobtypes.ClobPair_PerpetualClobMetadata{
			PerpetualClobMetadata: &clobtypes.PerpetualClobMetadata{
				PerpetualId: 0,
			},
		},
		StepBaseQuantums:          5,
		SubticksPerTick:           5,
		MinOrderBaseQuantums:      5,
		QuantumConversionExponent: -8,
		Status:                    clobtypes.ClobPair_STATUS_ACTIVE,
	}
	ClobPair_Eth = clobtypes.ClobPair{
		Id: 1,
		Metadata: &clobtypes.ClobPair_PerpetualClobMetadata{
			PerpetualClobMetadata: &clobtypes.PerpetualClobMetadata{
				PerpetualId: 1,
			},
		},
		StepBaseQuantums:          1000,
		SubticksPerTick:           1000,
		MinOrderBaseQuantums:      1000,
		QuantumConversionExponent: -9,
		MakerFeePpm:               MakerFeePpm,
		TakerFeePpm:               TakerFeePpm,
		Status:                    clobtypes.ClobPair_STATUS_ACTIVE,
	}
	ClobPair_Eth_No_Fee = clobtypes.ClobPair{
		Id: 1,
		Metadata: &clobtypes.ClobPair_PerpetualClobMetadata{
			PerpetualClobMetadata: &clobtypes.PerpetualClobMetadata{
				PerpetualId: 1,
			},
		},
		StepBaseQuantums:          1000,
		SubticksPerTick:           1000,
		MinOrderBaseQuantums:      1000,
		QuantumConversionExponent: -9,
		Status:                    clobtypes.ClobPair_STATUS_ACTIVE,
	}
	ClobPair_Asset = clobtypes.ClobPair{
		Id: 100,
		Metadata: &clobtypes.ClobPair_SpotClobMetadata{
			SpotClobMetadata: &clobtypes.SpotClobMetadata{
				BaseAssetId:  0,
				QuoteAssetId: 1,
			},
		},
		StepBaseQuantums:          1000,
		SubticksPerTick:           100,
		MinOrderBaseQuantums:      10,
		QuantumConversionExponent: 10,
		MakerFeePpm:               MakerFeePpm,
		TakerFeePpm:               TakerFeePpm,
		Status:                    clobtypes.ClobPair_STATUS_ACTIVE,
	}
	ClobPair_Btc2 = clobtypes.ClobPair{
		Id: 101,
		Metadata: &clobtypes.ClobPair_PerpetualClobMetadata{
			PerpetualClobMetadata: &clobtypes.PerpetualClobMetadata{
				PerpetualId: 0,
			},
		},
		StepBaseQuantums:          100,
		SubticksPerTick:           1000,
		MinOrderBaseQuantums:      100,
		QuantumConversionExponent: -8,
		MakerFeePpm:               200,
		TakerFeePpm:               500,
		Status:                    clobtypes.ClobPair_STATUS_ACTIVE,
	}
	ClobPair_Btc3 = clobtypes.ClobPair{
		Id: 0,
		Metadata: &clobtypes.ClobPair_PerpetualClobMetadata{
			PerpetualClobMetadata: &clobtypes.PerpetualClobMetadata{
				PerpetualId: 0,
			},
		},
		StepBaseQuantums:          10,
		SubticksPerTick:           100,
		MinOrderBaseQuantums:      100,
		QuantumConversionExponent: -8,
		MakerFeePpm:               MakerFeePpm,
		TakerFeePpm:               TakerFeePpm,
		Status:                    clobtypes.ClobPair_STATUS_ACTIVE,
	}
)
