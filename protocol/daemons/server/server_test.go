package server_test

import (
	"cosmossdk.io/log"
	"errors"
	"fmt"
	pricefeedconstants "github.com/dydxprotocol/v4-chain/protocol/daemons/constants"
	"github.com/dydxprotocol/v4-chain/protocol/daemons/server"
	"github.com/dydxprotocol/v4-chain/protocol/mocks"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/grpc"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"net"
	"os"
	"testing"
)

const (
	RemoveAllError = "RemovalAll ERROR!"
	ServeError     = "Serve ERROR!"
)

func TestStartServer_ListenFailsWhenInUse(t *testing.T) {
	defer os.RemoveAll(grpc.SocketPath)

	// Delete path to socket if it exists.
	err := os.RemoveAll(grpc.SocketPath)
	require.NoError(t, err, "No error should occur on clearing socket")

	// Start Listening to socket.
	_, err = net.Listen(pricefeedconstants.UnixProtocol, grpc.SocketPath)
	require.NoError(t, err, "No error should occur on listening to socket")

	mockGrpcServer := &mocks.GrpcServer{}
	mockFileHandler := &mocks.FileHandler{}

	mockFileHandler.On("RemoveAll", grpc.SocketPath).
		Return(func(key string) error {
			return nil
		})

	s := createServerWithMocks(t, mockGrpcServer, mockFileHandler)

	errorString := fmt.Sprintf(
		"listen %s %s: bind: address already in use",
		pricefeedconstants.UnixProtocol,
		grpc.SocketPath,
	)
	require.PanicsWithError(t, errorString, s.Start)

	verifyGrpcServerMocks(
		t,
		mockGrpcServer,
		mockFileHandler,
		true,
		false,
	)
}

func TestStart_Valid(t *testing.T) {
	// Remove filepath in case net.Listen is reached.
	defer os.RemoveAll(grpc.SocketPath)

	mockGrpcServer := &mocks.GrpcServer{}
	mockFileHandler := &mocks.FileHandler{}

	s := createServerWithMocks(
		t,
		mockGrpcServer,
		mockFileHandler,
	)

	mockFileHandler.On("RemoveAll", grpc.SocketPath).
		Return(nil)

	mockGrpcServer.On("Serve", mock.Anything).
		Return(nil)
	mockGrpcServer.On("RegisterService", mock.Anything, mock.Anything).
		Return()

	require.NotPanics(t, s.Start)

	// Reset with Umask.
	verifyFilePermissions(t, grpc.SocketPath, grpc.DefaultPermissions)

	verifyGrpcServerMocks(
		t,
		mockGrpcServer,
		mockFileHandler,
		true,
		true,
	)
}

func TestStart_MixedInvalid(t *testing.T) {
	tests := map[string]struct {
		FileHandlerror error
		serveError     error
	}{
		"Remove Socket Fails": {
			FileHandlerror: errors.New(RemoveAllError),
		},
		"Serve Fails": {
			serveError: errors.New(ServeError),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Remove filepath in case net.Listen is reached.
			defer os.RemoveAll(grpc.SocketPath)

			mockGrpcServer := &mocks.GrpcServer{}
			mockFileHandler := &mocks.FileHandler{}

			s := createServerWithMocks(
				t,
				mockGrpcServer,
				mockFileHandler,
			)

			var expectedError error

			expectedSocketPermissions := grpc.UserReadWriteOnlyPermissions

			if tc.FileHandlerror != nil {
				expectedError = tc.FileHandlerror
			}
			mockFileHandler.On("RemoveAll", grpc.SocketPath).
				Return(tc.FileHandlerror)

			if tc.serveError != nil {
				expectedSocketPermissions = grpc.DefaultPermissions

				expectedError = tc.serveError
			}

			mockGrpcServer.On("RegisterService", mock.Anything, mock.Anything).
				Return()
			mockGrpcServer.On("Serve", mock.Anything).
				Return(tc.serveError)

			require.PanicsWithError(t, expectedError.Error(), s.Start)

			if tc.serveError != nil {
				// Failed to reset with Umask before panicking.
				verifyFilePermissions(t, grpc.SocketPath, expectedSocketPermissions)
			}

			verifyGrpcServerMocks(
				t,
				mockGrpcServer,
				mockFileHandler,
				true,
				tc.serveError != nil,
			)
		})
	}
}

func createServerWithMocks(
	t testing.TB,
	mockGrpcServer *mocks.GrpcServer,
	mockFileHandler *mocks.FileHandler,
) *server.Server {
	server := server.NewServer(
		log.NewNopLogger(),
		mockGrpcServer,
		mockFileHandler,
		grpc.SocketPath,
	)
	mockGrpcServer.On("Stop").Return().Once()
	t.Cleanup(server.Stop)
	return server
}

func verifyFilePermissions(
	t *testing.T,
	path string,
	expectedPermissions os.FileMode,
) {
	fileStats, err := os.Stat(path)
	require.NoError(t, err, "Stats should exist for file at path")

	permissions := fileStats.Mode().Perm()
	require.Equal(t, expectedPermissions, permissions)
}

func verifyGrpcServerMocks(
	t *testing.T,
	mockGrpcServer *mocks.GrpcServer,
	mockFileHandler *mocks.FileHandler,
	removeAllIsCalled bool,
	serveIsCalled bool,
) {
	if removeAllIsCalled {
		mockFileHandler.AssertCalled(t, "RemoveAll", grpc.SocketPath)
	} else {
		mockFileHandler.AssertNotCalled(t, "RemoveAll", grpc.SocketPath)
	}

	if serveIsCalled {
		mockGrpcServer.AssertCalled(t, "Serve", mock.Anything)
		mockGrpcServer.AssertCalled(t, "RegisterService", mock.Anything, mock.Anything)
	} else {
		mockGrpcServer.AssertNotCalled(t, "Serve", mock.Anything)
		mockGrpcServer.AssertNotCalled(t, "RegisterService", mock.Anything, mock.Anything)
	}
}
