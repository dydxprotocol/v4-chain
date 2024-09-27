package authenticator

import (
	"strconv"
	"strings"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
)

var _ Authenticator = &SubaccountFilter{}

// SubaccountFilter filters incoming messages based on a predefined JSON pattern.
// It allows for complex pattern matching to support advanced authentication flows.
type SubaccountFilter struct {
	whitelist []uint32
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

// Initialize sets up the authenticator with the given data, which should be a valid JSON pattern for message filtering.
func (m SubaccountFilter) Initialize(config []byte) (Authenticator, error) {
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
func (m SubaccountFilter) Track(ctx sdk.Context, request AuthenticationRequest) error {
	return nil
}

// Authenticate checks if the provided message conforms to the set JSON pattern.
// It returns an AuthenticationResult based on the evaluation.
func (m SubaccountFilter) Authenticate(ctx sdk.Context, request AuthenticationRequest) error {
	// Collect the clob pair ids from the request.
	requestSubaccountNums := make([]uint32, 0)
	switch msg := request.Msg.(type) {
	case *clobtypes.MsgPlaceOrder:
		requestSubaccountNums = append(requestSubaccountNums, msg.Order.OrderId.SubaccountId.Number)
	case *clobtypes.MsgCancelOrder:
		requestSubaccountNums = append(requestSubaccountNums, msg.OrderId.SubaccountId.Number)
	case *clobtypes.MsgBatchCancel:
		requestSubaccountNums = append(requestSubaccountNums, msg.SubaccountId.Number)
	}

	// Make sure all the subaccount numbers are in the whitelist.
	for _, subaccountNum := range requestSubaccountNums {
		whitelisted := false
		for _, whitelistId := range m.whitelist {
			if subaccountNum == whitelistId {
				whitelisted = true
				break
			}
		}

		if !whitelisted {
			return errorsmod.Wrapf(
				sdkerrors.ErrUnauthorized,
				"subaccount number %d not in whitelist %v",
				subaccountNum,
				m.whitelist,
			)
		}
	}
	return nil
}

// ConfirmExecution confirms the execution of a message. Currently, it always confirms.
func (m SubaccountFilter) ConfirmExecution(ctx sdk.Context, request AuthenticationRequest) error {
	return nil
}

// OnAuthenticatorAdded performs additional checks when an authenticator is added.
// Specifically, it ensures numbers in JSON are encoded as strings.
func (m SubaccountFilter) OnAuthenticatorAdded(
	ctx sdk.Context,
	account sdk.AccAddress,
	config []byte,
	authenticatorId string,
) error {
	return nil
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
