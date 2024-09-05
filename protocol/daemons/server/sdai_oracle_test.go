package server_test

import (
	"testing"

	"github.com/StreamFinance-Protocol/stream-chain/protocol/daemons/sdaioracle/api"
	sdaitypes "github.com/StreamFinance-Protocol/stream-chain/protocol/daemons/server/types/sdaioracle"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/mocks"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/grpc"
	"github.com/stretchr/testify/require"
)

func TestAddsDAIEvents_EmptyRequest(t *testing.T) {
	mockGrpcServer := &mocks.GrpcServer{}
	mockFileHandler := &mocks.FileHandler{}
	sDAIEventManager := sdaitypes.SetupMockEventManagerWithNoEvents()

	s := createServerWithMocks(
		t,
		mockGrpcServer,
		mockFileHandler,
	).WithsDAIEventManager(
		sDAIEventManager,
	)

	resp, err := s.AddsDAIEvents(grpc.Ctx, &api.AddsDAIEventsRequest{})
	require.NoError(t, err)
	require.Empty(t, resp)
	require.Empty(t, sDAIEventManager.GetLastTensDAIEventsUnordered())
}

func TestAddsDAIEvents_SingleEvent(t *testing.T) {
	mockGrpcServer := &mocks.GrpcServer{}
	mockFileHandler := &mocks.FileHandler{}
	sDAIEventManager := sdaitypes.SetupMockEventManagerWithNoEvents()

	s := createServerWithMocks(
		t,
		mockGrpcServer,
		mockFileHandler,
	).WithsDAIEventManager(
		sDAIEventManager,
	)

	expectedEvent := &api.AddsDAIEventsRequest{
		ConversionRate: sdaitypes.TestSDAIEventRequests[0].ConversionRate,
	}

	resp, err := s.AddsDAIEvents(grpc.Ctx, &api.AddsDAIEventsRequest{
		ConversionRate: sdaitypes.TestSDAIEventRequests[0].ConversionRate,
	})
	require.NoError(t, err)
	require.Empty(t, resp)

	events := sDAIEventManager.GetLastTensDAIEventsUnordered()
	require.Equal(t, expectedEvent.ConversionRate, events[0].ConversionRate)
}

func TestAddsDAIEvents_MultipleEvents(t *testing.T) {
	mockGrpcServer := &mocks.GrpcServer{}
	mockFileHandler := &mocks.FileHandler{}
	sDAIEventManager := sdaitypes.SetupMockEventManagerWithNoEvents()

	s := createServerWithMocks(
		t,
		mockGrpcServer,
		mockFileHandler,
	).WithsDAIEventManager(
		sDAIEventManager,
	)

	expectedEvents := []*api.AddsDAIEventsRequest{
		{
			ConversionRate: sdaitypes.TestSDAIEventRequests[0].ConversionRate,
		},
		{
			ConversionRate: sdaitypes.TestSDAIEventRequests[1].ConversionRate,
		},
		{
			ConversionRate: sdaitypes.TestSDAIEventRequests[2].ConversionRate,
		},
		{
			ConversionRate: sdaitypes.TestSDAIEventRequests[1].ConversionRate,
		},
	}

	for _, event := range expectedEvents {
		resp, err := s.AddsDAIEvents(grpc.Ctx, event)
		require.NoError(t, err)
		require.Empty(t, resp)
	}

	actualEvents := sDAIEventManager.GetLastTensDAIEventsUnordered()
	for i, event := range expectedEvents {
		require.Equal(t, event.ConversionRate, actualEvents[i].ConversionRate)
	}
}
