package authenticator

import (
	"strconv"
	"strings"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
)

var _ Authenticator = &ClobPairIdFilter{}

// ClobPairIdFilter filters incoming messages based on a predefined JSON pattern.
// It allows for complex pattern matching to support advanced authentication flows.
type ClobPairIdFilter struct {
	whitelist []uint32
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

// Initialize sets up the authenticator with the given data, which should be a valid JSON pattern for message filtering.
func (m ClobPairIdFilter) Initialize(config []byte) (Authenticator, error) {
	strSlice := strings.Split(string(config), SEPARATOR)

	m.whitelist = make([]uint32, len(strSlice))
	for i, str := range strSlice {
		num, err := strconv.ParseUint(str, 10, 32)
		if err != nil {
			return nil, err
		}
		m.whitelist[i] = uint32(num)
	}
	return m, nil
}

// Track is a no-op in this implementation but can be used to track message handling.
func (m ClobPairIdFilter) Track(ctx sdk.Context, request AuthenticationRequest) error {
	return nil
}

// Authenticate checks if the provided message conforms to the set JSON pattern.
// It returns an AuthenticationResult based on the evaluation.
func (m ClobPairIdFilter) Authenticate(ctx sdk.Context, request AuthenticationRequest) error {
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
	}

	// Make sure all the clob pair ids are in the whitelist.
	for _, clobPairId := range requestOrderIds {
		whitelisted := false
		for _, whitelistId := range m.whitelist {
			if clobPairId == whitelistId {
				whitelisted = true
				break
			}
		}

		if !whitelisted {
			return errorsmod.Wrapf(
				sdkerrors.ErrUnauthorized,
				"order id %d not in whitelist %v",
				clobPairId,
				m.whitelist,
			)
		}
	}
	return nil
}

// ConfirmExecution confirms the execution of a message. Currently, it always confirms.
func (m ClobPairIdFilter) ConfirmExecution(ctx sdk.Context, request AuthenticationRequest) error {
	return nil
}

// OnAuthenticatorAdded performs additional checks when an authenticator is added.
// Specifically, it ensures numbers in JSON are encoded as strings.
func (m ClobPairIdFilter) OnAuthenticatorAdded(
	ctx sdk.Context,
	account sdk.AccAddress,
	config []byte,
	authenticatorId string,
) error {
	return nil
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
