package authenticator

import (
	"strconv"
	"strings"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/accountplus/types"
	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
)

var _ types.Authenticator = &ClobPairIdFilter{}

// ClobPairIdFilter filters incoming messages based on a whitelist of clob pair ids.
// It ensures that only messages with whitelisted clob pair ids are allowed.
type ClobPairIdFilter struct {
	whitelist map[uint32]struct{}
}

// NewClobPairIdFilter creates a new ClobPairIdFilter with the provided EncodingConfig.
func NewClobPairIdFilter() ClobPairIdFilter {
	return ClobPairIdFilter{}
}

// Type returns the type of the authenticator.
func (m ClobPairIdFilter) Type() string {
	return "ClobPairIdFilter"
}

// StaticGas returns the static gas amount for the authenticator. Currently, it's set to zero.
func (m ClobPairIdFilter) StaticGas() uint64 {
	return 0
}

// Initialize sets up the authenticator with the given configuration,
// which should be a list of clob pair ids separated by a predefined separator.
func (m ClobPairIdFilter) Initialize(config []byte) (types.Authenticator, error) {
	strSlice := strings.Split(string(config), SEPARATOR)

	m.whitelist = make(map[uint32]struct{})
	for _, str := range strSlice {
		num, err := strconv.ParseUint(str, 10, 32)
		if err != nil {
			return nil, err
		}
		m.whitelist[uint32(num)] = struct{}{}
	}
	return m, nil
}

// Track is a no-op in this implementation but can be used to track message handling.
func (m ClobPairIdFilter) Track(ctx sdk.Context, request types.AuthenticationRequest) error {
	return nil
}

// Authenticate checks if the message's clob pair ids are in the whitelist.
func (m ClobPairIdFilter) Authenticate(ctx sdk.Context, request types.AuthenticationRequest) error {
	// Collect the clob pair ids from the request.
	requestOrderIds := make([]uint32, 0)
	switch msg := request.Msg.(type) {
	case *clobtypes.MsgPlaceOrder:
		requestOrderIds = append(requestOrderIds, msg.Order.OrderId.ClobPairId)
	case *clobtypes.MsgCancelOrder:
		requestOrderIds = append(requestOrderIds, msg.OrderId.ClobPairId)
	case *clobtypes.MsgBatchCancel:
		for _, batch := range msg.ShortTermCancels {
			requestOrderIds = append(requestOrderIds, batch.ClobPairId)
		}
	default:
		// Skip other messages.
		return nil
	}

	// Make sure all the clob pair ids are in the whitelist.
	for _, clobPairId := range requestOrderIds {
		if _, ok := m.whitelist[clobPairId]; !ok {
			return errorsmod.Wrapf(
				types.ErrClobPairIdVerification,
				"order id %d not in whitelist %v",
				clobPairId,
				m.whitelist,
			)
		}
	}
	return nil
}

// ConfirmExecution confirms the execution of a message. Currently, it always confirms.
func (m ClobPairIdFilter) ConfirmExecution(ctx sdk.Context, request types.AuthenticationRequest) error {
	return nil
}

// OnAuthenticatorAdded is currently a no-op but can be extended for additional logic when an authenticator is added.
func (m ClobPairIdFilter) OnAuthenticatorAdded(
	ctx sdk.Context,
	account sdk.AccAddress,
	config []byte,
	authenticatorId string,
) (requireSigVerification bool, err error) {
	return false, nil
}

// OnAuthenticatorRemoved is a no-op in this implementation but can be used when an authenticator is removed.
func (m ClobPairIdFilter) OnAuthenticatorRemoved(
	ctx sdk.Context,
	account sdk.AccAddress,
	config []byte,
	authenticatorId string,
) error {
	return nil
}
