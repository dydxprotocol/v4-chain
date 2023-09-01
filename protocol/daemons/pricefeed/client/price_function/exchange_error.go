package price_function

import (
	"fmt"
	"github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/client/types"
)

type ExchangeError interface {
	error
	GetExchangeId() types.ExchangeId
}

type ExchangeErrorImpl struct {
	exchangeId types.ExchangeId
	err        error
}

func (e *ExchangeErrorImpl) Error() string {
	return fmt.Sprintf("%v exchange error: %v", e.exchangeId, e.err)
}

func (e *ExchangeErrorImpl) GetExchangeId() types.ExchangeId {
	return e.exchangeId
}

func NewExchangeError(exchangeId types.ExchangeId, msg string) ExchangeError {
	return &ExchangeErrorImpl{
		exchangeId: exchangeId,
		err:        fmt.Errorf(msg),
	}
}
