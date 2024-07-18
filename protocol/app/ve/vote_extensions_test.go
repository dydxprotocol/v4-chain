package ve_test

import (
	"cosmossdk.io/log"

	"testing"

	ve "github.com/StreamFinance-Protocol/stream-chain/protocol/app/ve"
	vecodec "github.com/StreamFinance-Protocol/stream-chain/protocol/app/ve/codec"
	vetypes "github.com/StreamFinance-Protocol/stream-chain/protocol/app/ve/types"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/mocks"
	constants "github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/constants"
	keepertest "github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/keeper"
	vetestutils "github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/ve"
	pricestypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/prices/types"
	cometabci "github.com/cometbft/cometbft/abci/types"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
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
		"price daemon returns no prices": {
			pricesKeeper: func() *mocks.ExtendVotePricesKeeper {
				mPricesKeeper := &mocks.ExtendVotePricesKeeper{}
				mPricesKeeper.On("GetValidMarketPriceUpdates", mock.Anything).Return(
					&pricestypes.MarketPriceUpdates{
						MarketPriceUpdates: []*pricestypes.MarketPriceUpdates_MarketPriceUpdate{},
					},
				)

				mPricesKeeper.On("GetMarketParam", mock.Anything, mock.Anything).Return(
					&pricestypes.MarketParam{},
					false,
				)

				return mPricesKeeper
			},
			expectedResponse: &vetypes.DaemonVoteExtension{
				Prices: nil,
			},
		},
		"oracle service returns single price": {
			pricesKeeper: func() *mocks.ExtendVotePricesKeeper {
				mpricesKeeper := &mocks.ExtendVotePricesKeeper{}
				mpricesKeeper.On("GetValidMarketPriceUpdates", mock.Anything).Return(
					constants.ValidSingleMarketPriceUpdate,
				)
				mpricesKeeper.On("GetMarketParam", mock.Anything, mock.Anything).Return(
					constants.TestSingleMarketParam,
					true,
				)

				return mpricesKeeper
			},
			expectedResponse: &vetypes.DaemonVoteExtension{
				Prices: map[uint32][]byte{
					constants.MarketId0: constants.Price5Bytes,
				},
			},
		},
		"oracle service returns multiple prices": {
			pricesKeeper: func() *mocks.ExtendVotePricesKeeper {
				mPricesKeeper := &mocks.ExtendVotePricesKeeper{}
				mPricesKeeper.On("GetValidMarketPriceUpdates", mock.Anything).Return(
					constants.ValidMarketPriceUpdates,
				)
				mPricesKeeper.On("GetMarketParam", mock.Anything, constants.MarketId0).Return(
					constants.TestMarketParams[0],
					true,
				)
				mPricesKeeper.On("GetMarketParam", mock.Anything, constants.MarketId1).Return(
					constants.TestMarketParams[1],
					true,
				)
				mPricesKeeper.On("GetMarketParam", mock.Anything, constants.MarketId2).Return(
					constants.TestMarketParams[2],
					true,
				)
				return mPricesKeeper
			},
			expectedResponse: &vetypes.DaemonVoteExtension{
				Prices: map[uint32][]byte{
					constants.MarketId0: constants.Price5Bytes,
					constants.MarketId1: constants.Price6Bytes,
					constants.MarketId2: constants.Price7Bytes,
				},
			},
		},
		"getting prices panics": {
			pricesKeeper: func() *mocks.ExtendVotePricesKeeper {
				mPricesKeeper := &mocks.ExtendVotePricesKeeper{}
				mPricesKeeper.On("GetValidMarketPriceUpdates", mock.Anything).Panic("panic")
				return mPricesKeeper
			},
			expectedResponse: &vetypes.DaemonVoteExtension{
				Prices: nil,
			},
			expectedError: true,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			ctx, _, _, _, _, _ := keepertest.PricesKeepers(t)
			ctx = vetestutils.GetVeEnabledCtx(ctx, 3)
			votecodec := vecodec.NewDefaultVoteExtensionCodec()

			mPriceApplier := &mocks.VEPriceApplier{}

			h := ve.NewVoteExtensionHandler(
				log.NewTestLogger(t),
				votecodec,
				tc.pricesKeeper(),
				mPriceApplier,
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
				mPriceApplier.On("ApplyPricesFromVE", ctx, finalizeBlockReq).Return(nil, nil)
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
	votecodec := vecodec.NewDefaultVoteExtensionCodec()
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
		"empty vote extension": {
			pricesKeeper: func() *mocks.ExtendVotePricesKeeper {
				mPricesKeeper := &mocks.ExtendVotePricesKeeper{}
				return mPricesKeeper
			},
			getReq: func() *cometabci.RequestVerifyVoteExtension {
				return &cometabci.RequestVerifyVoteExtension{}
			},
			expectedResponse: &cometabci.ResponseVerifyVoteExtension{
				Status: cometabci.ResponseVerifyVoteExtension_ACCEPT,
			},
			expectedError: false,
		},
		"empty vote extension - 1 valid price": {
			pricesKeeper: func() *mocks.ExtendVotePricesKeeper {
				mPricesKeeper := &mocks.ExtendVotePricesKeeper{}
				mPricesKeeper.On("GetMaxPairs", mock.Anything).Return(1)
				return mPricesKeeper
			},
			getReq: func() *cometabci.RequestVerifyVoteExtension {
				return &cometabci.RequestVerifyVoteExtension{}
			},
			expectedResponse: &cometabci.ResponseVerifyVoteExtension{
				Status: cometabci.ResponseVerifyVoteExtension_ACCEPT,
			},
			expectedError: false,
		},
		"malformed bytes returns error": {
			pricesKeeper: func() *mocks.ExtendVotePricesKeeper {
				mPricesKeeper := &mocks.ExtendVotePricesKeeper{}
				return mPricesKeeper
			},
			getReq: func() *cometabci.RequestVerifyVoteExtension {
				return &cometabci.RequestVerifyVoteExtension{
					VoteExtension: []byte("malformed"),
				}
			},
			expectedResponse: &cometabci.ResponseVerifyVoteExtension{
				Status: cometabci.ResponseVerifyVoteExtension_REJECT,
			},
			expectedError: true,
		},
		"valid vote extension - multple valid prices": {
			pricesKeeper: func() *mocks.ExtendVotePricesKeeper {
				mPricesKeeper := &mocks.ExtendVotePricesKeeper{}
				mPricesKeeper.On("GetAllMarketParams", mock.Anything).Return(
					constants.TestMarketParams,
				)
				return mPricesKeeper
			},
			getReq: func() *cometabci.RequestVerifyVoteExtension {
				extBz, err := vetestutils.CreateVoteExtensionBytes(
					constants.ValidVEPrice,
				)
				require.NoError(t, err)
				return &cometabci.RequestVerifyVoteExtension{
					VoteExtension: extBz,
					Height:        3,
				}
			},
			expectedResponse: &cometabci.ResponseVerifyVoteExtension{
				Status: cometabci.ResponseVerifyVoteExtension_ACCEPT,
			},
			expectedError: false,
		},
		"invalid vote extension - multple valid prices - should fail": {
			pricesKeeper: func() *mocks.ExtendVotePricesKeeper {
				mPricesKeeper := &mocks.ExtendVotePricesKeeper{}
				mPricesKeeper.On("GetAllMarketParams", mock.Anything).Return(
					constants.TestMarketParams[1:], // two prices
				)
				return mPricesKeeper
			},
			getReq: func() *cometabci.RequestVerifyVoteExtension {
				extBz, err := vetestutils.CreateVoteExtensionBytes(
					constants.ValidVEPrice,
				)
				require.NoError(t, err)
				return &cometabci.RequestVerifyVoteExtension{
					VoteExtension: extBz,
					Height:        3,
				}
			},
			expectedResponse: &cometabci.ResponseVerifyVoteExtension{
				Status: cometabci.ResponseVerifyVoteExtension_REJECT,
			},
			expectedError: true,
		},
		"vote extension with no prices": {
			pricesKeeper: func() *mocks.ExtendVotePricesKeeper {
				mPricesKeeper := &mocks.ExtendVotePricesKeeper{}
				mPricesKeeper.On("GetAllMarketParams", mock.Anything).Return(
					[]pricestypes.MarketParam{}, // two prices
				)
				return mPricesKeeper
			},
			getReq: func() *cometabci.RequestVerifyVoteExtension {
				prices := map[uint32][]byte{}

				extBz, err := vetestutils.CreateVoteExtensionBytes(
					prices,
				)
				require.NoError(t, err)

				return &cometabci.RequestVerifyVoteExtension{
					VoteExtension: extBz,
					Height:        3,
				}
			},

			expectedResponse: &cometabci.ResponseVerifyVoteExtension{
				Status: cometabci.ResponseVerifyVoteExtension_ACCEPT,
			},
			expectedError: false,
		},
		"vote extension with malformed prices": {
			pricesKeeper: func() *mocks.ExtendVotePricesKeeper {
				mPricesKeeper := &mocks.ExtendVotePricesKeeper{}
				mPricesKeeper.On("GetAllMarketParams", mock.Anything).Return(
					[]pricestypes.MarketParam{constants.TestMarketParams[0]},
				)
				return mPricesKeeper
			},
			getReq: func() *cometabci.RequestVerifyVoteExtension {
				prices := map[uint32][]byte{
					constants.MarketId0: make([]byte, 34),
				}

				extBz, err := vetestutils.CreateVoteExtensionBytes(
					prices,
				)
				require.NoError(t, err)

				return &cometabci.RequestVerifyVoteExtension{
					VoteExtension: extBz,
					Height:        3,
				}
			},
			expectedResponse: &cometabci.ResponseVerifyVoteExtension{
				Status: cometabci.ResponseVerifyVoteExtension_REJECT,
			},
			expectedError: true,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			ctx, _, _, _, _, _ := keepertest.PricesKeepers(t)
			ctx = vetestutils.GetVeEnabledCtx(ctx, 3)
			mPriceApplier := &mocks.VEPriceApplier{}
			mPricesKeeper := tc.pricesKeeper()

			handler := ve.NewVoteExtensionHandler(
				log.NewTestLogger(t),
				votecodec,
				mPricesKeeper,
				mPriceApplier,
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
