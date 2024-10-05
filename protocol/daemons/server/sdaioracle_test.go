package server_test

import (
	"errors"
	"testing"

	"github.com/StreamFinance-Protocol/stream-chain/protocol/daemons/sdaioracle/api"
	sdaitypes "github.com/StreamFinance-Protocol/stream-chain/protocol/daemons/server/types/sdaioracle"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/mocks"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/grpc"
	"github.com/stretchr/testify/mock"
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

	resp, err := s.AddsDAIEvent(grpc.Ctx, &api.AddsDAIEventRequest{})
	require.NoError(t, err)
	require.Empty(t, resp)
	require.Empty(t, sDAIEventManager.GetSDaiPrice())
}

func TestAddsDAIEvents_Error(t *testing.T) {
	mockGrpcServer := &mocks.GrpcServer{}
	mockFileHandler := &mocks.FileHandler{}
	mockSDaiEventManager := &mocks.SDAIEventManager{}

	mockSDaiEventManager.On("AddsDAIEvent", mock.Anything).Return(errors.New("error"))
	s := createServerWithMocks(
		t,
		mockGrpcServer,
		mockFileHandler,
	).WithsDAIEventManager(
		mockSDaiEventManager,
	)

	resp, err := s.AddsDAIEvent(grpc.Ctx, &api.AddsDAIEventRequest{
		ConversionRate: sdaitypes.TestSDAIEventRequest.ConversionRate,
	})

	require.Error(t, err)
	require.Nil(t, resp)
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

	expectedEvent := &api.AddsDAIEventRequest{
		ConversionRate: sdaitypes.TestSDAIEventRequest.ConversionRate,
	}

	resp, err := s.AddsDAIEvent(grpc.Ctx, &api.AddsDAIEventRequest{
		ConversionRate: sdaitypes.TestSDAIEventRequest.ConversionRate,
	})
	require.NoError(t, err)
	require.Empty(t, resp)

	event := sDAIEventManager.GetSDaiPrice()
	require.Equal(t, expectedEvent.ConversionRate, event.ConversionRate)

	secondEvent := &api.AddsDAIEventRequest{
		ConversionRate: "1",
	}

	resp, err = s.AddsDAIEvent(grpc.Ctx, secondEvent)
	require.NoError(t, err)
	require.Empty(t, resp)

	event = sDAIEventManager.GetSDaiPrice()
	require.Equal(t, secondEvent.ConversionRate, event.ConversionRate)
}
