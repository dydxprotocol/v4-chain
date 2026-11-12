package keeper

import (
	"fmt"

	errorsmod "cosmossdk.io/errors"
	transfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	connectiontypes "github.com/cosmos/ibc-go/v8/modules/core/03-connection/types"
	channeltypes "github.com/cosmos/ibc-go/v8/modules/core/04-channel/types"
	ibctm "github.com/cosmos/ibc-go/v8/modules/light-clients/07-tendermint"
	"github.com/dydxprotocol/v4-chain/protocol/x/canonicalusdc/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// RouteBinding is an open ICS20 transfer channel with its connection and
// light client.
type RouteBinding struct {
	Channel    string
	Connection string
	Client     string
}

// ResolveRoute finds the one open ICS20 transfer channel whose Tendermint
// light client tracks chainID. Zero or several such channels are an error, so
// a caller never binds a route by guesswork.
func (k Keeper) ResolveRoute(ctx sdk.Context, chainID string) (RouteBinding, error) {
	if k.channelKeeper == nil || k.connectionKeeper == nil || k.clientKeeper == nil {
		return RouteBinding{}, errorsmod.Wrap(types.ErrInvalidControls, "IBC keepers unavailable")
	}
	var matches []RouteBinding
	for _, channel := range k.channelKeeper.GetAllChannelsWithPortPrefix(ctx, transfertypes.PortID) {
		if channel.PortId != transfertypes.PortID || channel.State != channeltypes.OPEN ||
			channel.Ordering != channeltypes.UNORDERED || channel.Version != transfertypes.Version ||
			channel.Counterparty.PortId != transfertypes.PortID || len(channel.ConnectionHops) != 1 {
			continue
		}
		connection, found := k.connectionKeeper.GetConnection(ctx, channel.ConnectionHops[0])
		if !found || connection.State != connectiontypes.OPEN {
			continue
		}
		clientState, found := k.clientKeeper.GetClientState(ctx, connection.ClientId)
		if !found {
			continue
		}
		tendermint, ok := clientState.(*ibctm.ClientState)
		if !ok || tendermint.ChainId != chainID {
			continue
		}
		matches = append(matches, RouteBinding{
			Channel: channel.ChannelId, Connection: channel.ConnectionHops[0], Client: connection.ClientId,
		})
	}
	if len(matches) != 1 {
		return RouteBinding{}, errorsmod.Wrap(
			types.ErrInvalidControls, fmt.Sprintf("%d open transfer channels to %s", len(matches), chainID),
		)
	}
	return matches[0], nil
}
