package client_test

import (
	"errors"
	"fmt"
	"github.com/dydxprotocol/v4-chain/protocol/daemons/bridge/client"
	"github.com/dydxprotocol/v4-chain/protocol/mocks"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/grpc"
	bridgetypes "github.com/dydxprotocol/v4-chain/protocol/x/bridge/types"
	ethcoretypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"math/big"
	"testing"
)

func TestRunBridgeDaemonTaskLoop(t *testing.T) {
	errParams := errors.New("error getting event params")
	errPropose := errors.New("error getting propose params")
	errRecognizedEventInfo := errors.New("error getting recognized event info")
	errChainId := errors.New("error getting chain id")
	errEthereumLogs := errors.New("error getting Ethereum logs")
	errAddBridgeEvents := errors.New("error adding bridge events")

	tests := map[string]struct {
		eventParams            bridgetypes.EventParams
		eventParamsErr         error
		proposeParams          bridgetypes.ProposeParams
		proposeParamsErr       error
		recognizedEventInfo    bridgetypes.BridgeEventInfo
		recognizedEventInfoErr error
		chainId                int
		chainIdError           error
		filterLogs             []ethcoretypes.Log
		filterLogsErr          error
		addBridgeEventsErr     error

		expectedErrorString string
		expectedError       error
	}{
		"Success": {
			eventParams:         constants.EventParams,
			proposeParams:       constants.ProposeParams,
			recognizedEventInfo: constants.RecognizedEventInfo_Id2_Height0,
			chainId:             constants.EthChainId,
			filterLogs: []ethcoretypes.Log{
				constants.EthLog_Event0,
				constants.EthLog_Event1,
			},
		},
		"Error getting event params": {
			eventParamsErr: errParams,
			expectedError:  errParams,
		},
		"Error getting propose params": {
			eventParams:      constants.EventParams,
			proposeParamsErr: errPropose,
			expectedError:    errPropose,
		},
		"Error getting recognized event info": {
			eventParams:            constants.EventParams,
			proposeParams:          constants.ProposeParams,
			recognizedEventInfoErr: errRecognizedEventInfo,
			expectedError:          errRecognizedEventInfo,
		},
		"Error getting chain id": {
			eventParams:         constants.EventParams,
			proposeParams:       constants.ProposeParams,
			recognizedEventInfo: constants.RecognizedEventInfo_Id2_Height0,
			chainIdError:        errChainId,
			expectedError:       errChainId,
		},
		"Error chain ID not as expected": {
			eventParams:         constants.EventParams,
			proposeParams:       constants.ProposeParams,
			recognizedEventInfo: constants.RecognizedEventInfo_Id2_Height0,
			chainId:             constants.EthChainId + 1,
			expectedErrorString: fmt.Sprintf(
				"expected chain ID %d but node has chain ID %d",
				constants.EthChainId,
				constants.EthChainId+1,
			),
		},
		"Error getting Ethereum logs": {
			eventParams:         constants.EventParams,
			proposeParams:       constants.ProposeParams,
			recognizedEventInfo: constants.RecognizedEventInfo_Id2_Height0,
			chainId:             constants.EthChainId,
			filterLogsErr:       errEthereumLogs,
			expectedError:       errEthereumLogs,
		},
		"Error adding bridge events": {
			eventParams:         constants.EventParams,
			proposeParams:       constants.ProposeParams,
			recognizedEventInfo: constants.RecognizedEventInfo_Id2_Height0,
			chainId:             constants.EthChainId,
			filterLogs: []ethcoretypes.Log{
				constants.EthLog_Event0,
			},
			addBridgeEventsErr: errAddBridgeEvents,
			expectedError:      errAddBridgeEvents,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			ctx := grpc.Ctx
			mockLogger := mocks.Logger{}
			mockEthClient := mocks.EthClient{}
			mockQueryClient := mocks.BridgeQueryClient{}
			mockServiceClient := mocks.BridgeServiceClient{}

			mockQueryClient.On("EventParams", ctx, mock.Anything).Return(
				&bridgetypes.QueryEventParamsResponse{
					Params: tc.eventParams,
				},
				tc.eventParamsErr,
			)
			mockQueryClient.On("ProposeParams", ctx, mock.Anything).Return(
				&bridgetypes.QueryProposeParamsResponse{
					Params: tc.proposeParams,
				},
				tc.proposeParamsErr,
			)
			mockQueryClient.On("RecognizedEventInfo", ctx, mock.Anything).Return(
				&bridgetypes.QueryRecognizedEventInfoResponse{
					Info: tc.recognizedEventInfo,
				},
				tc.recognizedEventInfoErr,
			)
			mockEthClient.On("ChainID", ctx).Return(big.NewInt(int64(tc.chainId)), tc.chainIdError)
			mockEthClient.On("FilterLogs", ctx, mock.Anything).Return(tc.filterLogs, tc.filterLogsErr)
			mockServiceClient.On("AddBridgeEvents", ctx, mock.Anything).Return(nil, tc.addBridgeEventsErr)

			subTaskRunner := &client.SubTaskRunnerImpl{}
			err := subTaskRunner.RunBridgeDaemonTaskLoop(
				grpc.Ctx,
				&mockLogger,
				&mockEthClient,
				&mockQueryClient,
				&mockServiceClient,
			)
			if tc.expectedErrorString != "" {
				require.Error(t, err)
				require.ErrorContains(t, err, tc.expectedErrorString)
			}
			if tc.expectedError != nil {
				require.Error(t, err)
				require.ErrorIs(t, err, tc.expectedError)
			}
		})
	}
}
