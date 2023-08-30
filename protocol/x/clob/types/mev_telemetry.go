package types

import (
	"math/big"
)

// MEVDatapoint contains the Volume (ValidatorVolumeQuoteQuantums) and MEV per market
// to be sent to the MEV telemetry service. Every datapoint contains a self-reported
// identifier and a block height for which the metric is reported.
type MEVDatapoint struct {
	Height              uint32                  `json:"block_height"`
	ChainID             string                  `json:"chain_id"`
	VolumeQuoteQuantums map[ClobPairId]*big.Int `json:"volume_quote_quantums"`
	MEV                 map[ClobPairId]float32  `json:"mev"`
	Identifier          string                  `json:"identifier"`
}

// MevMetrics represents all MEV metrics to send to the MEV telemetry service.
type MevMetrics struct {
	MevNodeToNode MevNodeToNodeMetrics `json:"mev_node_to_node"`
	MevDatapoint  MEVDatapoint         `json:"mev_datapoint"`
}
