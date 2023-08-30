package client_test

import (
	"errors"
	"fmt"
	"math/big"
	"testing"

	"github.com/cometbft/cometbft/libs/log"
	"github.com/dydxprotocol/v4-chain/protocol/daemons/bridge/client"
	d_constants "github.com/dydxprotocol/v4-chain/protocol/daemons/constants"
	"github.com/dydxprotocol/v4-chain/protocol/daemons/flags"
	"github.com/dydxprotocol/v4-chain/protocol/mocks"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/grpc"
	bridgetypes "github.com/dydxprotocol/v4-chain/protocol/x/bridge/types"
	ethcoretypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestStart_TcpConnectionFails(t *testing.T) {
	errorMsg := "Failed to create connection"

	mockGrpcClient := &mocks.GrpcClient{}
	mockGrpcClient.On("NewTcpConnection", grpc.Ctx, d_constants.DefaultGrpcEndpoint).Return(nil, errors.New(errorMsg))

	require.EqualError(
		t,
		client.Start(
			grpc.Ctx,
			flags.GetDefaultDaemonFlags(),
			log.NewNopLogger(),
			mockGrpcClient,
		),
		errorMsg,
	)
	mockGrpcClient.AssertCalled(t, "NewTcpConnection", grpc.Ctx, d_constants.DefaultGrpcEndpoint)
	mockGrpcClient.AssertNotCalled(t, "NewGrpcConnection", grpc.Ctx, grpc.SocketPath)
	mockGrpcClient.AssertNotCalled(t, "CloseConnection", grpc.GrpcConn)
}

func TestStart_UnixSocketConnectionFails(t *testing.T) {
	errorMsg := "Failed to create connection"

	mockGrpcClient := &mocks.GrpcClient{}
	mockGrpcClient.On("NewTcpConnection", grpc.Ctx, d_constants.DefaultGrpcEndpoint).Return(grpc.GrpcConn, nil)
	mockGrpcClient.On("NewGrpcConnection", grpc.Ctx, grpc.SocketPath).Return(nil, errors.New(errorMsg))
	mockGrpcClient.On("CloseConnection", grpc.GrpcConn).Return(nil)

	require.EqualError(
		t,
		client.Start(
			grpc.Ctx,
			flags.GetDefaultDaemonFlags(),
			log.NewNopLogger(),
			mockGrpcClient,
		),
		errorMsg,
	)
	mockGrpcClient.AssertCalled(t, "NewTcpConnection", grpc.Ctx, d_constants.DefaultGrpcEndpoint)
	mockGrpcClient.AssertCalled(t, "NewGrpcConnection", grpc.Ctx, grpc.SocketPath)
	mockGrpcClient.AssertNumberOfCalls(t, "CloseConnection", 1)
}

func TestRunBridgeDaemonTaskLoop(t *testing.T) {
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

		expectedResponse error
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
			eventParamsErr:   errors.New("error getting event params"),
			expectedResponse: errors.New("error getting event params"),
		},
		"Error getting propose params": {
			eventParams:      constants.EventParams,
			proposeParamsErr: errors.New("error getting propose params"),
			expectedResponse: errors.New("error getting propose params"),
		},
		"Error getting recognized event info": {
			eventParams:            constants.EventParams,
			proposeParams:          constants.ProposeParams,
			recognizedEventInfoErr: errors.New("error getting recognized event info"),
			expectedResponse:       errors.New("error getting recognized event info"),
		},
		"Error getting chain id": {
			eventParams:         constants.EventParams,
			proposeParams:       constants.ProposeParams,
			recognizedEventInfo: constants.RecognizedEventInfo_Id2_Height0,
			chainIdError:        errors.New("error getting chain id"),
			expectedResponse:    errors.New("error getting chain id"),
		},
		"Error chain ID not as expected": {
			eventParams:         constants.EventParams,
			proposeParams:       constants.ProposeParams,
			recognizedEventInfo: constants.RecognizedEventInfo_Id2_Height0,
			chainId:             constants.EthChainId + 1,
			expectedResponse: fmt.Errorf(
				"Expected chain ID %d but node has chain ID %d",
				constants.EthChainId,
				constants.EthChainId+1,
			),
		},
		"Error getting Ethereum logs": {
			eventParams:         constants.EventParams,
			proposeParams:       constants.ProposeParams,
			recognizedEventInfo: constants.RecognizedEventInfo_Id2_Height0,
			chainId:             constants.EthChainId,
			filterLogsErr:       errors.New("error getting Ethereum logs"),
			expectedResponse:    errors.New("error getting Ethereum logs"),
		},
		"Error adding bridge events": {
			eventParams:         constants.EventParams,
			proposeParams:       constants.ProposeParams,
			recognizedEventInfo: constants.RecognizedEventInfo_Id2_Height0,
			chainId:             constants.EthChainId,
			filterLogs: []ethcoretypes.Log{
				constants.EthLog_Event0,
			},
			addBridgeEventsErr: errors.New("error adding bridge events"),
			expectedResponse:   errors.New("error adding bridge events"),
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

			err := client.RunBridgeDaemonTaskLoop(
				grpc.Ctx,
				&mockLogger,
				&mockEthClient,
				&mockQueryClient,
				&mockServiceClient,
			)
			require.Equal(t, tc.expectedResponse, err)
		})
	}
}
