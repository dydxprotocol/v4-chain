package slinky_test

import (
	"fmt"
	"testing"

	"github.com/dydxprotocol/slinky/pkg/types"
	"github.com/stretchr/testify/require"

	"github.com/dydxprotocol/v4-chain/protocol/lib/slinky"
)

func TestMarketPairToCurrencyPair(t *testing.T) {
	testCases := []struct {
		mp  string
		cp  types.CurrencyPair
		err error
	}{
		{mp: "FOO-BAR", cp: types.CurrencyPair{Base: "FOO", Quote: "BAR"}, err: nil},
		{mp: "FOOBAR", cp: types.CurrencyPair{}, err: fmt.Errorf("incorrectly formatted CurrencyPair: FOOBAR")},
		{mp: "FOO/BAR", cp: types.CurrencyPair{}, err: fmt.Errorf("incorrectly formatted CurrencyPair: FOOBAR")},
	}
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("TestMarketPair %s", tc.mp), func(t *testing.T) {
			cp, err := slinky.MarketPairToCurrencyPair(tc.mp)
			if tc.err != nil {
				require.Error(t, err, tc.err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.cp, cp)
			}
		})
	}
}
