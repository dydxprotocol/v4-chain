package authenticator

import (
	"encoding/json"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/dydxprotocol/v4-chain/protocol/x/accountplus/types"
)

type AllOf struct {
	SubAuthenticators   []types.Authenticator
	am                  *AuthenticatorManager
	signatureAssignment SignatureAssignment
}

var _ types.Authenticator = &AllOf{}

func NewAllOf(am *AuthenticatorManager) AllOf {
	return AllOf{
		am:                  am,
		SubAuthenticators:   []types.Authenticator{},
		signatureAssignment: Single,
	}
}

func NewPartitionedAllOf(am *AuthenticatorManager) AllOf {
	return AllOf{
		am:                  am,
		SubAuthenticators:   []types.Authenticator{},
		signatureAssignment: Partitioned,
	}
}

func (aoa AllOf) Type() string {
	if aoa.signatureAssignment == Single {
		return "AllOf"
	}
	return "PartitionedAllOf"
}

func (aoa AllOf) StaticGas() uint64 {
	var totalGas uint64
	for _, auth := range aoa.SubAuthenticators {
		totalGas += auth.StaticGas()
	}
	return totalGas
}

func (aoa AllOf) Initialize(config []byte) (types.Authenticator, error) {
	var initDatas []types.SubAuthenticatorInitData
	if err := json.Unmarshal(config, &initDatas); err != nil {
		return nil, errorsmod.Wrap(err, "failed to parse sub-authenticators initialization data")
	}

	if len(initDatas) <= 1 {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "allOf must have at least 2 sub-authenticators")
	}

	for _, initData := range initDatas {
		authenticatorCode := aoa.am.GetAuthenticatorByType(initData.Type)
		instance, err := authenticatorCode.Initialize(initData.Config)
		if err != nil {
			return nil, errorsmod.Wrapf(err, "failed to initialize sub-authenticator (type = %s)", initData.Type)
		}
		aoa.SubAuthenticators = append(aoa.SubAuthenticators, instance)
	}

	// If not all sub-authenticators are registered, return an error
	if len(aoa.SubAuthenticators) != len(initDatas) {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "failed to initialize all sub-authenticators")
	}

	return aoa, nil
}

func (aoa AllOf) Authenticate(ctx sdk.Context, request types.AuthenticationRequest) error {
	if len(aoa.SubAuthenticators) == 0 {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "no sub-authenticators provided")
	}

	var signatures [][]byte
	var err error
	if aoa.signatureAssignment == Partitioned {
		// Partitioned signatures are decoded and passed one by one as the signature of the sub-authenticator
		signatures, err = splitSignatures(request.Signature, len(aoa.SubAuthenticators))
		if err != nil {
			return err
		}
	}

	baseId := request.AuthenticatorId
	for i, auth := range aoa.SubAuthenticators {
		// update the authenticator id to include the sub-authenticator id
		request.AuthenticatorId = compositeId(baseId, i)
		// update the request to include the sub-authenticator signature
		if aoa.signatureAssignment == Partitioned {
			request.Signature = signatures[i]
		}
		if err := auth.Authenticate(ctx, request); err != nil {
			return errorsmod.Wrap(
				types.ErrAllOfVerification,
				err.Error(),
			)
		}
	}
	return nil
}

func (aoa AllOf) Track(ctx sdk.Context, request types.AuthenticationRequest) error {
	return subTrack(ctx, request, aoa.SubAuthenticators)
}

func (aoa AllOf) ConfirmExecution(ctx sdk.Context, request types.AuthenticationRequest) error {
	var signatures [][]byte
	var err error
	if aoa.signatureAssignment == Partitioned {
		// Partitioned signatures are decoded and passed one by one as the signature of the sub-authenticator
		signatures, err = splitSignatures(request.Signature, len(aoa.SubAuthenticators))
		if err != nil {
			return err
		}
	}

	baseId := request.AuthenticatorId
	for i, auth := range aoa.SubAuthenticators {
		// update the authenticator id to include the sub-authenticator id
		request.AuthenticatorId = compositeId(baseId, i)
		// update the request to include the sub-authenticator signature
		if aoa.signatureAssignment == Partitioned {
			request.Signature = signatures[i]
		}

		if err := auth.ConfirmExecution(ctx, request); err != nil {
			return err
		}
	}
	return nil
}

func (aoa AllOf) OnAuthenticatorAdded(
	ctx sdk.Context,
	account sdk.AccAddress,
	config []byte,
	authenticatorId string,
) (bool, error) {
	return onSubAuthenticatorsAdded(
		ctx,
		account,
		config,
		authenticatorId,
		aoa.am,
		false, // isAnyOf = false
	)
}

func (aoa AllOf) OnAuthenticatorRemoved(
	ctx sdk.Context,
	account sdk.AccAddress,
	config []byte,
	authenticatorId string,
) error {
	return onSubAuthenticatorsRemoved(ctx, account, config, authenticatorId, aoa.am)
}
