package price_function

import (
	"fmt"

	"github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/client/types"
)

// ExchangeError describes an error that is specific to a particular exchange. These errors are emitted
// by an exchange's price function whenever the exchange's API returns a valid http response.
type ExchangeError interface {
	error
	GetExchangeId() types.ExchangeId
}

// Ensure ExcahngeErrorImpl implements ExchangeError at compile time.
var _ ExchangeError = &ExchangeErrorImpl{}

// ExchangeErrorImpl implements ExchangeError.
type ExchangeErrorImpl struct {
	exchangeId types.ExchangeId
	err        error
}

// Error returns a string representation of the error.
func (e *ExchangeErrorImpl) Error() string {
	return fmt.Sprintf("%v exchange error: %v", e.exchangeId, e.err)
}

// GetExchangeId returns the exchange id associated with the error.
func (e *ExchangeErrorImpl) GetExchangeId() types.ExchangeId {
	return e.exchangeId
}

// NewExchangeError returns a new ExchangeError.
func NewExchangeError(exchangeId types.ExchangeId, msg string) ExchangeError {
	return &ExchangeErrorImpl{
		exchangeId: exchangeId,
		err:        fmt.Errorf("%s", msg),
	}
}
