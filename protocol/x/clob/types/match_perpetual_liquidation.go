package types

import (
	gometrics "github.com/armon/go-metrics"
	"github.com/dydxprotocol/v4-chain/protocol/lib/metrics"
)

func (m MatchPerpetualLiquidation) GetMetricLabels(callback string) []gometrics.Label {
	return []gometrics.Label{
		metrics.GetLabelForStringValue(
			metrics.Callback,
			callback,
		),
		metrics.GetLabelForStringValue(
			metrics.SubaccountOwner,
			m.Liquidated.Owner,
		),
		metrics.GetLabelForIntValue(
			metrics.PerpetualId,
			int(m.PerpetualId),
		),
	}
}
