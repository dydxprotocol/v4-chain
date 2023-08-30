package types_test

import (
	"testing"
	"time"

	"github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/types"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	"github.com/stretchr/testify/require"
)

func TestNewPriceTimestamp_IsEmpty(t *testing.T) {
	pt := types.NewPriceTimestamp()

	require.Equal(t, time.Time{}, pt.LastUpdateTime)
	require.Equal(t, uint64(0), pt.Price)
}

func TestUpdatePrice_NoPreviousPriceSuccess(t *testing.T) {
	pt := types.NewPriceTimestamp()

	// No previous Price exists
	ok := pt.UpdatePrice(constants.Price1, &constants.TimeT)
	require.True(t, ok)

	require.Equal(t, constants.TimeT, pt.LastUpdateTime)
	require.Equal(t, constants.Price1, pt.Price)
}

func TestUpdatePrice_GreaterUpdateTimeSuccess(t *testing.T) {
	pt := types.NewPriceTimestamp()

	// Last update @ timeT
	ok := pt.UpdatePrice(constants.Price1, &constants.TimeT)
	require.True(t, ok)

	// New update @ timeT + threshold
	ok = pt.UpdatePrice(constants.Price2, &constants.TimeTPlusThreshold)
	require.True(t, ok)

	require.Equal(t, constants.TimeTPlusThreshold, pt.LastUpdateTime)
	require.Equal(t, constants.Price2, pt.Price)
}

func TestUpdatePrice_EqualUpdateTimeFail(t *testing.T) {
	pt := types.NewPriceTimestamp()

	// Last update @ timeT
	ok := pt.UpdatePrice(constants.Price1, &constants.TimeT)
	require.True(t, ok)

	// New update @ timeT
	ok = pt.UpdatePrice(constants.Price2, &constants.TimeT)
	require.False(t, ok)

	// No update should be made because the new update time is not greater.
	require.Equal(t, constants.TimeT, pt.LastUpdateTime)
	require.Equal(t, constants.Price1, pt.Price)
}

func TestUpdatePrice_SmallerUpdateTimeFail(t *testing.T) {
	pt := types.NewPriceTimestamp()

	// Last update @ timeT
	ok := pt.UpdatePrice(constants.Price1, &constants.TimeT)
	require.True(t, ok)

	// New update @ timeT - threshold
	ok = pt.UpdatePrice(constants.Price2, &constants.TimeTMinusThreshold)
	require.False(t, ok)

	// No update should be made because the new update time is not greater.
	require.Equal(t, constants.TimeT, pt.LastUpdateTime)
	require.Equal(t, constants.Price1, pt.Price)
}

func TestGetValidPrice_ValidLastUpdateTimeSuccess(t *testing.T) {
	pt := types.NewPriceTimestamp()

	// Last update @ timeT
	ok := pt.UpdatePrice(constants.Price1, &constants.TimeT)
	require.True(t, ok)

	r, ok := pt.GetValidPrice(constants.TimeT)
	require.True(t, ok)
	require.Equal(t, constants.Price1, r)
}

func TestGetValidPrice_InvalidLastUpdateTimeFail(t *testing.T) {
	pt := types.NewPriceTimestamp()

	// Last update @ timeT
	ok := pt.UpdatePrice(constants.Price1, &constants.TimeT)
	require.True(t, ok)

	// Updates @ timeT are no longer valid at this cutoff time.
	r, ok := pt.GetValidPrice(constants.TimeTPlus1)
	require.False(t, ok)
	require.Equal(t, uint64(0), r)
}
