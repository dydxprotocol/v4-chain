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

	resp, err := s.AddsDAIEvent(grpc.Ctx, &api.AddsDAIEventsRequest{})
	require.NoError(t, err)
	require.Empty(t, resp)
	require.Empty(t, sDAIEventManager.GetSDaiPrice())
}

func TestAddsDAIEvents(t *testing.T) {
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
		ConversionRate: sdaitypes.TestSDAIEventRequest.ConversionRate,
	}

	resp, err := s.AddsDAIEvent(grpc.Ctx, &api.AddsDAIEventsRequest{
		ConversionRate: sdaitypes.TestSDAIEventRequest.ConversionRate,
	})
	require.NoError(t, err)
	require.Empty(t, resp)

	event := sDAIEventManager.GetSDaiPrice()
	require.Equal(t, expectedEvent.ConversionRate, event.ConversionRate)

	secondEvent := &api.AddsDAIEventsRequest{
		ConversionRate: "1",
	}

	resp, err = s.AddsDAIEvent(grpc.Ctx, secondEvent)
	require.NoError(t, err)
	require.Empty(t, resp)

	event = sDAIEventManager.GetSDaiPrice()
	require.Equal(t, secondEvent.ConversionRate, event.ConversionRate)
}
