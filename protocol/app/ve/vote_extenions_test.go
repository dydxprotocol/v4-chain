package ve_test

import (
	"cosmossdk.io/log"

	"testing"

	"github.com/stretchr/testify/require"

	keepertest "github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/keeper"

	ve "github.com/StreamFinance-Protocol/stream-chain/protocol/app/ve"
	vecodec "github.com/StreamFinance-Protocol/stream-chain/protocol/app/ve/codec"
	vetypes "github.com/StreamFinance-Protocol/stream-chain/protocol/app/ve/types"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/mocks"
	cometabci "github.com/cometbft/cometbft/abci/types"
)

type TestExtendedVoteTC struct {
	expectedResponse  *vetypes.DaemonVoteExtension
	pricesKeeper      func() *mocks.ExtendVotePricesKeeper
	extendVoteRequest func() *cometabci.RequestExtendVote
	expectedError     bool
}

func TestExtendVoteHandler(t *testing.T) {
	tests := map[string]TestExtendedVoteTC{
		"nil request returns error": {
			pricesKeeper: func() *mocks.ExtendVotePricesKeeper {
				mPricesKeeper := &mocks.ExtendVotePricesKeeper{}
				return mPricesKeeper
			},
			extendVoteRequest: func() *cometabci.RequestExtendVote {
				return nil
			},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			ctx, _, _, indexPriceCache, _, _ := keepertest.PricesKeepers(t)

			votecodec := vecodec.NewDefaultVoteExtensionCodec()

			mPriceApplier := &mocks.PriceApplier{}

			h := ve.NewVoteExtensionHandler(
				log.NewTestLogger(t),
				indexPriceCache,
				votecodec,
				tc.pricesKeeper(),
			)

			req := &cometabci.RequestExtendVote{}
			if tc.extendVoteRequest != nil {
				req = tc.extendVoteRequest()
			}
			if req != nil {
				finalizeBlockReq := &cometabci.RequestFinalizeBlock{
					Txs:    req.Txs,
					Height: req.Height,
				}
				mPriceApplier.On("ApplyPricesFromVoteExtensions", ctx, finalizeBlockReq).Return(nil, nil)
			}
			resp, err := h.ExtendVoteHandler()(ctx, req)
			if !tc.expectedError {
				if resp == nil || len(resp.VoteExtension) == 0 {
					return
				}
				require.NoError(t, err)
				require.NotNil(t, resp)
				ext, err := votecodec.Decode(resp.VoteExtension)
				require.NoError(t, err)
				require.Equal(t, tc.expectedResponse.Prices, ext.Prices)
			} else {
				require.Error(t, err)
			}
		})
	}
}

type TestVerifyExtendedVoteTC struct {
	getReq           func() *cometabci.RequestVerifyVoteExtension
	pricesKeeper     func() *mocks.ExtendVotePricesKeeper
	expectedResponse *cometabci.ResponseVerifyVoteExtension
	expectedError    bool
}

func TestVerifyVoteHandler(t *testing.T) {
	tests := map[string]TestVerifyExtendedVoteTC{
		"nil request returns error": {
			pricesKeeper: func() *mocks.ExtendVotePricesKeeper {
				mPricesKeeper := &mocks.ExtendVotePricesKeeper{}
				return mPricesKeeper
			},
			getReq: func() *cometabci.RequestVerifyVoteExtension {
				return nil
			},
			expectedResponse: nil,
			expectedError:    true,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			ctx, _, _, indexPriceCache, _, _ := keepertest.PricesKeepers(t)

			votecodec := vecodec.NewDefaultVoteExtensionCodec()

			handler := ve.NewVoteExtensionHandler(
				log.NewTestLogger(t),
				indexPriceCache,
				votecodec,
				tc.pricesKeeper(),
			).VerifyVoteExtensionHandler()

			resp, err := handler(ctx, tc.getReq())
			require.Equal(t, tc.expectedResponse, resp)

			if tc.expectedError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}

}
