package authenticator

import (
	"strconv"
	"strings"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/accountplus/types"
	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
)

var _ types.Authenticator = &SubaccountFilter{}

// SubaccountFilter filters incoming messages based on a whitelist of subaccount numbers.
// It ensures that only messages with whitelisted subaccount numbers are allowed.
type SubaccountFilter struct {
	whitelist map[uint32]struct{}
}

// NewSubaccountFilter creates a new SubaccountFilter with the provided EncodingConfig.
func NewSubaccountFilter() SubaccountFilter {
	return SubaccountFilter{}
}

// Type returns the type of the authenticator.
func (m SubaccountFilter) Type() string {
	return "SubaccountFilter"
}

// StaticGas returns the static gas amount for the authenticator. Currently, it's set to zero.
func (m SubaccountFilter) StaticGas() uint64 {
	return 0
}

// Initialize sets up the authenticator with the given data,
// which should be a string of subaccount numbers separated by the specified separator.
func (m SubaccountFilter) Initialize(config []byte) (types.Authenticator, error) {
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
func (m SubaccountFilter) Track(ctx sdk.Context, request types.AuthenticationRequest) error {
	return nil
}

// Authenticate checks if the message's subaccount numbers are in the whitelist.
func (m SubaccountFilter) Authenticate(ctx sdk.Context, request types.AuthenticationRequest) error {
	// Collect the clob pair ids from the request.
	requestSubaccountNums := make([]uint32, 0)
	switch msg := request.Msg.(type) {
	case *clobtypes.MsgPlaceOrder:
		requestSubaccountNums = append(requestSubaccountNums, msg.Order.OrderId.SubaccountId.Number)
	case *clobtypes.MsgCancelOrder:
		requestSubaccountNums = append(requestSubaccountNums, msg.OrderId.SubaccountId.Number)
	case *clobtypes.MsgBatchCancel:
		requestSubaccountNums = append(requestSubaccountNums, msg.SubaccountId.Number)
	default:
		// Skip other messages.
		return nil
	}

	// Make sure all the subaccount numbers are in the whitelist.
	for _, subaccountNum := range requestSubaccountNums {
		if _, ok := m.whitelist[subaccountNum]; !ok {
			return errorsmod.Wrapf(
				types.ErrSubaccountVerification,
				"subaccount number %d not in whitelist %v",
				subaccountNum,
				m.whitelist,
			)
		}
	}
	return nil
}

// ConfirmExecution confirms the execution of a message. Currently, it always confirms.
func (m SubaccountFilter) ConfirmExecution(ctx sdk.Context, request types.AuthenticationRequest) error {
	return nil
}

// OnAuthenticatorAdded is currently a no-op but can be extended for additional logic when an authenticator is added.
func (m SubaccountFilter) OnAuthenticatorAdded(
	ctx sdk.Context,
	account sdk.AccAddress,
	config []byte,
	authenticatorId string,
) (requireSigVerification bool, err error) {
	return false, nil
}

// OnAuthenticatorRemoved is a no-op in this implementation but can be used when an authenticator is removed.
func (m SubaccountFilter) OnAuthenticatorRemoved(
	ctx sdk.Context,
	account sdk.AccAddress,
	config []byte,
	authenticatorId string,
) error {
	return nil
}
