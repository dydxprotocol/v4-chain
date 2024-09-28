package client_test

import (
	"errors"
	"fmt"
	"math/big"
	"testing"

	"github.com/StreamFinance-Protocol/stream-chain/protocol/daemons/sdaioracle/client"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/mocks"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/constants"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/grpc"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestRunsDAIDaemonTaskLoop(t *testing.T) {
	errChainId := errors.New("error getting chain id")
	errQuerysDAIRate := errors.New("failed to fetch chain ID")
	errAddSDAIEvents := errors.New("failed to add sDAI events")

	tests := map[string]struct {
		chainId          int
		chainIdError     error
		daiRate          string
		queryDaiErr      error
		addsDAIEventsErr error

		expectedErrorString string
		expectedError       error
	}{
		"Success": {
			chainId: constants.EthChainId,
			daiRate: constants.SDAIRate,
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
		"Error getting dai conversion rare": {
			chainId:       constants.EthChainId,
			queryDaiErr:   errQuerysDAIRate,
			expectedError: errQuerysDAIRate,
		},
		"Error adding sDAI events": {
			chainId:          constants.EthChainId,
			daiRate:          constants.SDAIRate,
			addsDAIEventsErr: errAddSDAIEvents,
			expectedError:    errAddSDAIEvents,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			ctx := grpc.Ctx
			mockLogger := mocks.Logger{}
			mockEthClient := &ethclient.Client{}
			mockQueryClient := mocks.EthQueryClient{}
			mockServiceClient := mocks.SDAIServiceClient{}

			mockQueryClient.On("ChainID", ctx, mock.Anything).Return(big.NewInt(int64(tc.chainId)), tc.chainIdError)
			mockQueryClient.On("QueryDaiConversionRate", mock.Anything).Return(tc.daiRate, tc.queryDaiErr)
			mockServiceClient.On("AddsDAIEvent", ctx, mock.Anything).Return(nil, tc.addsDAIEventsErr)

			subTaskRunner := &client.SubTaskRunnerImpl{}
			err := subTaskRunner.RunsDAIDaemonTaskLoop(
				grpc.Ctx,
				&mockLogger,
				mockEthClient,
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
