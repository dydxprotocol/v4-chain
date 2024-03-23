package constants

import (
	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
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
		QuantumConversionExponent: -8,
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
		QuantumConversionExponent: -9,
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
		QuantumConversionExponent: 10,
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
		QuantumConversionExponent: -8,
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
		QuantumConversionExponent: -8,
		Status:                    clobtypes.ClobPair_STATUS_ACTIVE,
	}
	ClobPair_Btc_Initializing = clobtypes.ClobPair{
		Id: 0,
		Metadata: &clobtypes.ClobPair_PerpetualClobMetadata{
			PerpetualClobMetadata: &clobtypes.PerpetualClobMetadata{
				PerpetualId: 0,
			},
		},
		StepBaseQuantums:          5,
		SubticksPerTick:           5,
		QuantumConversionExponent: -8,
		Status:                    clobtypes.ClobPair_STATUS_INITIALIZING,
	}
	ClobPair_Btc_Final_Settlement = clobtypes.ClobPair{
		Id: 0,
		Metadata: &clobtypes.ClobPair_PerpetualClobMetadata{
			PerpetualClobMetadata: &clobtypes.PerpetualClobMetadata{
				PerpetualId: 0,
			},
		},
		StepBaseQuantums:          5,
		SubticksPerTick:           5,
		QuantumConversionExponent: -8,
		Status:                    clobtypes.ClobPair_STATUS_FINAL_SETTLEMENT,
	}
	ClobPair_Btc_Paused = clobtypes.ClobPair{
		Id: 0,
		Metadata: &clobtypes.ClobPair_PerpetualClobMetadata{
			PerpetualClobMetadata: &clobtypes.PerpetualClobMetadata{
				PerpetualId: 0,
			},
		},
		StepBaseQuantums:          5,
		SubticksPerTick:           5,
		QuantumConversionExponent: -8,
		Status:                    clobtypes.ClobPair_STATUS_PAUSED,
	}
	ClobPair_3_Iso = clobtypes.ClobPair{
		Id: 3,
		Metadata: &clobtypes.ClobPair_PerpetualClobMetadata{
			PerpetualClobMetadata: &clobtypes.PerpetualClobMetadata{
				PerpetualId: 3,
			},
		},
		StepBaseQuantums:          5,
		SubticksPerTick:           5,
		QuantumConversionExponent: -8,
		Status:                    clobtypes.ClobPair_STATUS_ACTIVE,
	}
	ClobPair_3_Iso_Final_Settlement = clobtypes.ClobPair{
		Id: 3,
		Metadata: &clobtypes.ClobPair_PerpetualClobMetadata{
			PerpetualClobMetadata: &clobtypes.PerpetualClobMetadata{
				PerpetualId: 3,
			},
		},
		StepBaseQuantums:          5,
		SubticksPerTick:           5,
		QuantumConversionExponent: -8,
		Status:                    clobtypes.ClobPair_STATUS_FINAL_SETTLEMENT,
	}
	ClobPair_4_Iso2 = clobtypes.ClobPair{
		Id: 4,
		Metadata: &clobtypes.ClobPair_PerpetualClobMetadata{
			PerpetualClobMetadata: &clobtypes.PerpetualClobMetadata{
				PerpetualId: 4,
			},
		},
		StepBaseQuantums:          5,
		SubticksPerTick:           5,
		QuantumConversionExponent: -8,
		Status:                    clobtypes.ClobPair_STATUS_ACTIVE,
	}
)
