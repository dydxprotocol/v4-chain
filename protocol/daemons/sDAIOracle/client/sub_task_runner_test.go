package client_test

import (
	"errors"
	"fmt"
	"math/big"
	"testing"

	"github.com/StreamFinance-Protocol/stream-chain/protocol/daemons/sDAIOracle/client"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/mocks"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/constants"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/grpc"
	ethcoretypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestRunsDAIDaemonTaskLoop(t *testing.T) {
	errParams := errors.New("error getting event params")
	errPropose := errors.New("error getting propose params")
	errRecognizedEventInfo := errors.New("error getting recognized event info")
	errChainId := errors.New("error getting chain id")
	errEthereumLogs := errors.New("error getting Ethereum logs")
	errAddBridgeEvents := errors.New("error adding bridge events")

	tests := map[string]struct {
		chainId            int
		chainIdError       error
		filterLogs         []ethcoretypes.Log
		filterLogsErr      error
		addBridgeEventsErr error

		expectedErrorString string
		expectedError       error
	}{
		"Success": {
			chainId: constants.EthChainId,
			filterLogs: []ethcoretypes.Log{
				constants.EthLog_Event0,
				constants.EthLog_Event1,
			},
		},
		"Error getting event params": {
			expectedError: errParams,
		},
		"Error getting propose params": {
			expectedError: errPropose,
		},
		"Error getting recognized event info": {
			expectedError: errRecognizedEventInfo,
		},
		"Error getting chain id": {
			chainIdError:  errChainId,
			expectedError: errChainId,
		},
		"Error chain ID not as expected": {
			chainId: constants.EthChainId + 1,
			expectedErrorString: fmt.Sprintf(
				"expected chain ID %d but node has chain ID %d",
				constants.EthChainId,
				constants.EthChainId+1,
			),
		},
		"Error getting Ethereum logs": {
			chainId:       constants.EthChainId,
			filterLogsErr: errEthereumLogs,
			expectedError: errEthereumLogs,
		},
		"Error adding bridge events": {
			chainId: constants.EthChainId,
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
