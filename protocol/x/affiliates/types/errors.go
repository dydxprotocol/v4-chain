package types

import errorsmod "cosmossdk.io/errors"

var (
	ErrAffiliateAlreadyExistsForReferee = errorsmod.Register(ModuleName, 1, "Affiliate already exists for referee")
	ErrAffiliateTiersNotInitialized     = errorsmod.Register(ModuleName, 2, "Affiliate tier data not found")
	ErrInvalidAffiliateTiers            = errorsmod.Register(ModuleName, 3, "Invalid affiliate tier data")
	ErrUpdatingAffiliateReferredVolume  = errorsmod.Register(
		ModuleName, 4, "Error updating affiliate referred volume")
	ErrInvalidAddress          = errorsmod.Register(ModuleName, 5, "Invalid address")
	ErrAffiliateNotFound       = errorsmod.Register(ModuleName, 6, "Affiliate not found")
	ErrRevShareSafetyViolation = errorsmod.Register(
		ModuleName, 7, "Rev share safety violation")
	ErrDuplicateAffiliateAddressForWhitelist = errorsmod.Register(
		ModuleName, 8, "Duplicate affiliate address for whitelist")
	ErrAffiliateTiersNotSet = errorsmod.Register(ModuleName, 9,
		"Affiliate tiers not set (affiliate program is not active)")
	ErrSelfReferral                        = errorsmod.Register(ModuleName, 10, "Self referral not allowed")
	ErrUpdatingAffiliateReferredCommission = errorsmod.Register(
		ModuleName, 11, "Error updating affiliate referred commission")
	ErrUpdatingAttributedVolume = errorsmod.Register(
		ModuleName, 12, "Error updating attributed volume")
)
