package types

import (
	"math/big"

	satypes "github.com/dydxprotocol/v4/x/subaccounts/types"
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

// MEVMatch represents all necessary data to calculate MEV for a regular match.
// If the MEV telemetry service is enabled, the validator's local matches will be encoded into
// this struct and sent to the MEV telemetry service.
type MEVMatch struct {
	TakerOrderSubaccountId satypes.SubaccountId `json:"taker_order_subaccount_id"`
	TakerFeePpm            uint32               `json:"taker_fee_ppm"`

	MakerOrderSubaccountId satypes.SubaccountId `json:"maker_order_subaccount_id"`
	MakerOrderSubticks     Subticks             `json:"maker_order_subticks"`
	MakerOrderIsBuy        bool                 `json:"maker_order_is_buy"`
	MakerFeePpm            uint32               `json:"maker_fee_ppm"`

	ClobPairId ClobPairId           `json:"clob_pair_id"`
	FillAmount satypes.BaseQuantums `json:"fill_amount"`
}

// MEVLiquidationMatch represents all necessary data to calculate MEV for a liquidation.
// If the MEV telemetry service is enabled, the validator's local matches will be encoded into
// this struct and sent to the MEV telemetry service.
type MEVLiquidationMatch struct {
	LiquidatedSubaccountId          satypes.SubaccountId `json:"liquidated_subaccount_id"`
	InsuranceFundDeltaQuoteQuantums int64                `json:"insurance_fund_delta_quote_quantums"`

	MakerOrderSubaccountId satypes.SubaccountId `json:"maker_order_subaccount_id"`
	MakerOrderSubticks     Subticks             `json:"maker_order_subticks"`
	MakerOrderIsBuy        bool                 `json:"maker_order_is_buy"`
	MakerFeePpm            uint32               `json:"maker_fee_ppm"`

	ClobPairId ClobPairId           `json:"clob_pair_id"`
	FillAmount satypes.BaseQuantums `json:"fill_amount"`
}

// MevNodeToNodeMetrics represents a data structure that will be sent to the MEV telemetry service.
type MevNodeToNodeMetrics struct {
	ValidatorMevMatches            []MEVMatch              `json:"validator_mev_matches"`
	ValidatorMevLiquidationMatches []MEVLiquidationMatch   `json:"validator_mev_liquidation_matches"`
	ClobMidPrices                  map[ClobPairId]Subticks `json:"clob_mid_prices"`
	ClobPairs                      map[ClobPairId]ClobPair `json:"clob_pairs"`
}

// MevMetrics represents all MEV metrics to send to the MEV telemetry service.
type MevMetrics struct {
	MevNodeToNode MevNodeToNodeMetrics `json:"mev_node_to_node"`
	MevDatapoint  MEVDatapoint         `json:"mev_datapoint"`
}
