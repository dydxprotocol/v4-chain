package ante_test

import (
	"testing"

	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/cosmos/cosmos-sdk/testutil/testdata"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/mocks"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	sdktest "github.com/dydxprotocol/v4-chain/protocol/testutil/sdk"
	txtest "github.com/dydxprotocol/v4-chain/protocol/testutil/sdk/tx"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/ante"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// The rate limiters mutate unsynchronized maps and app/ante.go skips its CheckTx locks under simulate, so the
// decorator must not reach them when simulate is true. The mock has no Rate* expectations for that case, so
// any call fails the test.
func TestClobRateLimitDecorator_SkipsRateLimitingUnderSimulate(t *testing.T) {
	for name, isSimulate := range map[string]bool{
		"simulate=true: rate limiting is skipped":  true,
		"simulate=false: rate limiting still runs": false,
	} {
		t.Run(name, func(t *testing.T) {
			ctx, _, _ := sdktest.NewSdkContextWithMultistore()
			ctx = ctx.WithIsCheckTx(true)

			mockClobKeeper := &mocks.ClobKeeper{}
			if !isSimulate {
				mockClobKeeper.On("RateLimitPlaceOrder", ctx, constants.Msg_PlaceOrder).Return(nil)
			}

			decorator := ante.NewRateLimitDecorator(mockClobKeeper)
			antehandler := sdk.ChainAnteDecorators(decorator)

			priv, _, _ := testdata.KeyTestPubAddr()
			tx, err := txtest.CreateTestTx(
				[]cryptotypes.PrivKey{priv},
				[]uint64{0},
				[]uint64{0},
				"dydx",
				[]sdk.Msg{constants.Msg_PlaceOrder},
				0,
			)
			require.NoError(t, err)

			_, err = antehandler(ctx, tx, isSimulate)
			require.NoError(t, err)

			if !isSimulate {
				mockClobKeeper.AssertExpectations(t)
			}
			mockClobKeeper.AssertNotCalled(t, "RateLimitCancelOrder", mock.Anything, mock.Anything)
			mockClobKeeper.AssertNotCalled(t, "RateLimitBatchCancel", mock.Anything, mock.Anything)
			mockClobKeeper.AssertNotCalled(t, "RateLimitUpdateLeverage", mock.Anything, mock.Anything)
			if isSimulate {
				mockClobKeeper.AssertNotCalled(t, "RateLimitPlaceOrder", mock.Anything, mock.Anything)
			}
		})
	}
}
