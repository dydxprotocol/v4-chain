package client_test

import (
	"errors"
	"fmt"
	"math/big"
	"testing"

	"github.com/cometbft/cometbft/libs/log"
	appflags "github.com/dydxprotocol/v4-chain/protocol/app/flags"
	"github.com/dydxprotocol/v4-chain/protocol/daemons/bridge/client"
	d_constants "github.com/dydxprotocol/v4-chain/protocol/daemons/constants"
	daemonflags "github.com/dydxprotocol/v4-chain/protocol/daemons/flags"
	"github.com/dydxprotocol/v4-chain/protocol/mocks"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/appoptions"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/grpc"
	bridgetypes "github.com/dydxprotocol/v4-chain/protocol/x/bridge/types"
	ethcoretypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestStart_EthRpcEndpointNotSet(t *testing.T) {
	errorMsg := "flag bridge-daemon-eth-rpc-endpoint is not set"

	require.EqualError(
		t,
		client.Start(
			grpc.Ctx,
			daemonflags.GetDefaultDaemonFlags(),
			appflags.GetFlagValuesFromOptions(appoptions.GetDefaultTestAppOptions("", nil)),
			log.NewNopLogger(),
			&mocks.GrpcClient{},
		),
		errorMsg,
	)
}

func TestStart_TcpConnectionFails(t *testing.T) {
	errorMsg := "Failed to create connection"

	// Mock the gRPC client to return an error when creating a TCP connection.
	mockGrpcClient := &mocks.GrpcClient{}
	mockGrpcClient.On("NewTcpConnection", grpc.Ctx, d_constants.DefaultGrpcEndpoint).Return(nil, errors.New(errorMsg))

	// Override default daemon flags with a non-empty EthRpcEndpoint.
	daemonFlagsWithEthRpcEndpoint := daemonflags.GetDefaultDaemonFlags()
	daemonFlagsWithEthRpcEndpoint.Bridge.EthRpcEndpoint = "http://localhost:8545"

	require.EqualError(
		t,
		client.Start(
			grpc.Ctx,
			daemonFlagsWithEthRpcEndpoint,
			appflags.GetFlagValuesFromOptions(appoptions.GetDefaultTestAppOptions("", nil)),
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

	// Mock the gRPC client to
	// - return a successful TCP connection.
	// - return an error when creating a gRPC connection.
	// - successfully close the TCP connection.
	mockGrpcClient := &mocks.GrpcClient{}
	mockGrpcClient.On("NewTcpConnection", grpc.Ctx, d_constants.DefaultGrpcEndpoint).Return(grpc.GrpcConn, nil)
	mockGrpcClient.On("NewGrpcConnection", grpc.Ctx, grpc.SocketPath).Return(nil, errors.New(errorMsg))
	mockGrpcClient.On("CloseConnection", grpc.GrpcConn).Return(nil)

	// Override default daemon flags with a non-empty EthRpcEndpoint.
	daemonFlagsWithEthRpcEndpoint := daemonflags.GetDefaultDaemonFlags()
	daemonFlagsWithEthRpcEndpoint.Bridge.EthRpcEndpoint = "http://localhost:8545"

	require.EqualError(
		t,
		client.Start(
			grpc.Ctx,
			daemonFlagsWithEthRpcEndpoint,
			appflags.GetFlagValuesFromOptions(appoptions.GetDefaultTestAppOptions("", nil)),
			log.NewNopLogger(),
			mockGrpcClient,
		),
		errorMsg,
	)
	mockGrpcClient.AssertCalled(t, "NewTcpConnection", grpc.Ctx, d_constants.DefaultGrpcEndpoint)
	mockGrpcClient.AssertCalled(t, "NewGrpcConnection", grpc.Ctx, grpc.SocketPath)

	// Assert that the connection from NewTcpConnection is closed.
	mockGrpcClient.AssertNumberOfCalls(t, "CloseConnection", 1)
	mockGrpcClient.AssertCalled(t, "CloseConnection", grpc.GrpcConn)
}

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

			err := client.RunBridgeDaemonTaskLoop(
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
