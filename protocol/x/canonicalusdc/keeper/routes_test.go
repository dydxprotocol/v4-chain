package keeper

import (
	"testing"

	transfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	channeltypes "github.com/cosmos/ibc-go/v8/modules/core/04-channel/types"
	canonicaltypes "github.com/dydxprotocol/v4-chain/protocol/x/canonicalusdc/types"
	"github.com/stretchr/testify/require"
)

func TestResolveRouteFindsTheOneChannelToAChain(t *testing.T) {
	ctx, k, _, _, _ := setupKeeper(t)
	controls := activeControls()
	bindOpenRoutes(k, controls)
	k.clientKeeper = clientMock{controls.NobleClient: "noble-1", controls.InjectiveClient: "injective-1"}

	noble, err := k.ResolveRoute(ctx, "noble-1")
	require.NoError(t, err)
	require.Equal(t, RouteBinding{
		Channel: controls.NobleChannel, Connection: controls.NobleConnection, Client: controls.NobleClient,
	}, noble)
	injective, err := k.ResolveRoute(ctx, "injective-1")
	require.NoError(t, err)
	require.Equal(t, controls.InjectiveChannel, injective.Channel)

	_, err = k.ResolveRoute(ctx, "osmosis-1")
	require.ErrorIs(t, err, canonicaltypes.ErrInvalidControls)
	require.ErrorContains(t, err, "0 open transfer channels to osmosis-1")
}

func TestResolveRouteRefusesAmbiguousAndClosedChannels(t *testing.T) {
	ctx, k, _, _, _ := setupKeeper(t)
	controls := activeControls()
	bindOpenRoutes(k, controls)
	k.clientKeeper = clientMock{controls.NobleClient: "noble-1", controls.InjectiveClient: "injective-1"}
	channels := k.channelKeeper.(channelMock)

	second := channels[controls.InjectiveChannel]
	channels["channel-9"] = second
	_, err := k.ResolveRoute(ctx, "injective-1")
	require.ErrorContains(t, err, "2 open transfer channels to injective-1")

	closed := channeltypes.NewChannel(
		channeltypes.CLOSED, channeltypes.UNORDERED,
		channeltypes.NewCounterparty(transfertypes.PortID, "channel-remote"),
		second.ConnectionHops, transfertypes.Version,
	)
	channels["channel-9"] = closed
	injective, err := k.ResolveRoute(ctx, "injective-1")
	require.NoError(t, err)
	require.Equal(t, controls.InjectiveChannel, injective.Channel)
}

func TestResolveRouteNeedsIBCKeepers(t *testing.T) {
	ctx, k, _, _, _ := setupKeeper(t)
	_, err := k.ResolveRoute(ctx, "noble-1")
	require.ErrorIs(t, err, canonicaltypes.ErrInvalidControls)
}
