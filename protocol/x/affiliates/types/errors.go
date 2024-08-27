package types

import errorsmod "cosmossdk.io/errors"

var (
	ErrAffiliateAlreadyExistsForReferee = errorsmod.Register(ModuleName, 1, "Affiliate already exists for referee")
	ErrAffiliateTiersNotInitialized     = errorsmod.Register(ModuleName, 2, "Affiliate tier data not found")
	ErrInvalidAffiliateTiers            = errorsmod.Register(ModuleName, 3, "Invalid affiliate tier data")
	ErrUpdatingAffiliateReferredVolume  = errorsmod.Register(ModuleName, 4, "Error updating affiliate referred volume")
	ErrInvalidAddress                   = errorsmod.Register(ModuleName, 5, "Invalid address")
)
