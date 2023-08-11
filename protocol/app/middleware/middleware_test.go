package middleware_test

import (
	"fmt"
	"testing"

	"github.com/dydxprotocol/v4/app/middleware"
	"github.com/dydxprotocol/v4/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestRunTxPanicLoggingMiddleware(t *testing.T) {
	logger := &mocks.Logger{}
	func() {
		defer func() {
			if r := recover(); r != nil {
				handler := middleware.NewRunTxPanicLoggingMiddleware(logger)
				err := handler(r)
				require.Nil(t, err)
			}
		}()
		// Do something that does not panic.
	}()
	logger.AssertExpectations(t)

	logger.On(
		"Error",
		"runTx panic'ed with a problem",
		"stack trace",
		mock.Anything,
	).Return(nil).Once()

	func() {
		defer func() {
			if r := recover(); r != nil {
				handler := middleware.NewRunTxPanicLoggingMiddleware(logger)
				err := handler(r)
				require.Nil(t, err)
			}
		}()
		// Panic with a string.
		panic("a problem")
	}()
	logger.AssertExpectations(t)

	logger.On(
		"Error",
		"runTx panic'ed with a problem",
		"stack trace",
		mock.Anything,
	).Return(nil).Once()

	func() {
		defer func() {
			if r := recover(); r != nil {
				handler := middleware.NewRunTxPanicLoggingMiddleware(logger)
				err := handler(r)
				require.Nil(t, err)
			}
		}()
		// Panic with an error.
		panic(fmt.Errorf("a problem"))
	}()
	logger.AssertExpectations(t)
}
