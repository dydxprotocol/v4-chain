package authenticator

import (
	"encoding/json"
	"strconv"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/dydxprotocol/v4-chain/protocol/x/accountplus/types"
)

func subTrack(
	ctx sdk.Context,
	request types.AuthenticationRequest,
	subAuthenticators []types.Authenticator,
) error {
	baseId := request.AuthenticatorId
	for id, auth := range subAuthenticators {
		request.AuthenticatorId = compositeId(baseId, id)
		err := auth.Track(ctx, request)
		if err != nil {
			return errorsmod.Wrapf(err, "sub-authenticator track failed (sub-authenticator id = %s)", request.AuthenticatorId)
		}
	}
	return nil
}

func splitSignatures(signature []byte, total int) ([][]byte, error) {
	var signatures [][]byte
	err := json.Unmarshal(signature, &signatures)
	if err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "failed to parse signatures")
	}
	if len(signatures) != total {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "invalid number of signatures")
	}
	return signatures, nil
}

func onSubAuthenticatorsAdded(
	ctx sdk.Context,
	account sdk.AccAddress,
	data []byte,
	authenticatorId string,
	am *AuthenticatorManager,
	isAnyOf bool, // If true, the composite is AnyOf, otherwise AllOf
) (bool, error) {
	var initDatas []types.SubAuthenticatorInitData
	if err := json.Unmarshal(data, &initDatas); err != nil {
		return false, errorsmod.Wrapf(err, "failed to unmarshal sub-authenticator init data")
	}

	if len(initDatas) <= 1 {
		return false, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "no sub-authenticators provided")
	}

	// If authenticator type is AllOf, we just need to check if ANY of the sub-authenticators require
	// signature verification.
	// Use `false` as identify value for OR operation.
	requireSigVerification := false
	if isAnyOf {
		// For `AnyOf`, we need to check if ALL of the sub-authenticators require signature verification.
		// Use `true` as identify value for AND operation.
		requireSigVerification = true
	}

	baseId := authenticatorId
	subAuthenticatorCount := 0
	for id, initData := range initDatas {
		authenticatorCode := am.GetAuthenticatorByType(initData.Type)
		if authenticatorCode == nil {
			return false, errorsmod.Wrapf(
				sdkerrors.ErrInvalidRequest,
				"sub-authenticator failed to be added in function `OnAuthenticatorAdded` as type is not registered in manager",
			)
		}
		subId := compositeId(baseId, id)
		subRequireSigVerification, err := authenticatorCode.OnAuthenticatorAdded(
			ctx,
			account,
			initData.Config,
			subId,
		)
		if err != nil {
			return false, errorsmod.Wrapf(
				err,
				"sub-authenticator `OnAuthenticatorAdded` failed (sub-authenticator id = %s)",
				subId,
			)
		}

		if isAnyOf {
			// For `AnyOf`, we require ALL sub-authenticators to require signature verification.
			requireSigVerification = requireSigVerification && subRequireSigVerification
		} else {
			// For `AllOf`, we just need ANY of sub-authenticators to require signature verification.
			requireSigVerification = requireSigVerification || subRequireSigVerification
		}

		subAuthenticatorCount++
	}

	// If not all sub-authenticators are registered, return an error
	if subAuthenticatorCount != len(initDatas) {
		return false, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "failed to initialize all sub-authenticators")
	}

	return requireSigVerification, nil
}

func onSubAuthenticatorsRemoved(
	ctx sdk.Context,
	account sdk.AccAddress,
	data []byte,
	authenticatorId string,
	am *AuthenticatorManager,
) error {
	var initDatas []types.SubAuthenticatorInitData
	if err := json.Unmarshal(data, &initDatas); err != nil {
		return err
	}

	baseId := authenticatorId
	for id, initData := range initDatas {
		authenticatorCode := am.GetAuthenticatorByType(initData.Type)
		if authenticatorCode == nil {
			return errorsmod.Wrapf(
				sdkerrors.ErrInvalidRequest,
				"sub-authenticator failed to be removed in function `OnAuthenticatorRemoved` as type is not registered in manager",
			)
		}
		subId := compositeId(baseId, id)
		err := authenticatorCode.OnAuthenticatorRemoved(ctx, account, initData.Config, subId)
		if err != nil {
			return errorsmod.Wrapf(err, "sub-authenticator `OnAuthenticatorRemoved` failed (sub-authenticator id = %s)", subId)
		}
	}
	return nil
}

func compositeId(baseId string, subId int) string {
	return baseId + "." + strconv.Itoa(subId)
}
