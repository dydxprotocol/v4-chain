//go:build all || integration_test

package cli_test

import (
	"strconv"
)

// Prevent strconv unused error
var _ = strconv.IntSize

// func TestGetSDAIPriceQuery(t *testing.T) {
// 	cfg := network.DefaultConfig(nil)

// 	chi := "1006681181716810314385961731"

// 	time.Sleep(15 * time.Second)

// 	rateQuery := "docker exec interchain-security-instance interchain-security-cd" +
// 		" query ratelimit get-sdai-price "
// 	data, _, err := network.QueryCustomNetwork(rateQuery)

// 	require.NoError(t, err)
// 	var resp types.GetSDAIPriceQueryResponse
// 	require.NoError(t, cfg.Codec.UnmarshalJSON(data, &resp))

// 	chiFloat, success := new(big.Float).SetString(chi)
// 	require.True(t, success, "Failed to parse chi as big.Float")

// 	priceFloat, success := new(big.Float).SetString(resp.Price)
// 	require.True(t, success, "Failed to parse price as big.Float")

// 	// Compare the big.Float values directly
// 	comparison := new(big.Float).Quo(priceFloat, chiFloat)

// 	minThreshold := big.NewFloat(0.99)
// 	maxThreshold := big.NewFloat(1.16)

// 	require.True(t, comparison.Cmp(minThreshold) >= 0, "Price should be at least 99% of chi")
// 	require.True(t, comparison.Cmp(maxThreshold) <= 0, "Price should be at most 116% of chi")
// }
