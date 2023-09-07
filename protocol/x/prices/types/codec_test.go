package types_test

import (
	"bytes"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	"testing"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
	"github.com/stretchr/testify/require"
)

func TestRegisterCodec(t *testing.T) {
	cdc := codec.NewLegacyAmino()
	types.RegisterCodec(cdc)
	var buf bytes.Buffer
	err := cdc.Amino.PrintTypes(&buf)
	require.NoError(t, err)
	require.Contains(t, buf.String(), "MsgUpdateMarketPrices")
	require.Contains(t, buf.String(), "prices/UpdateMarketPrices")
}

func TestRegisterInterfaces(t *testing.T) {
	registry := cdctypes.NewInterfaceRegistry()
	types.RegisterInterfaces(registry)
	require.NoError(t, registry.EnsureRegistered(&types.MsgUpdateMarketPrices{}))
	require.NoError(t, registry.EnsureRegistered(&types.MsgUpdateMarketPricesResponse{}))
	require.NoError(t, registry.EnsureRegistered(&types.MsgCreateOracleMarket{}))
	require.NoError(t, registry.EnsureRegistered(&types.MsgCreateOracleMarketResponse{}))
}
