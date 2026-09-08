package types

import errorsmod "cosmossdk.io/errors"

var (
	ErrDisabled              = errorsmod.Register(ModuleName, 1, "canonical USDC routing is disabled")
	ErrPaused                = errorsmod.Register(ModuleName, 2, "canonical USDC routing is paused")
	ErrUnsupportedChannel    = errorsmod.Register(ModuleName, 3, "unsupported canonical USDC channel")
	ErrInsufficientBacking   = errorsmod.Register(ModuleName, 4, "insufficient route backing")
	ErrPendingLimit          = errorsmod.Register(ModuleName, 5, "pending settlement limit reached")
	ErrInvalidControls       = errorsmod.Register(ModuleName, 6, "invalid canonical USDC controls")
	ErrInvalidLedger         = errorsmod.Register(ModuleName, 7, "invalid canonical USDC ledger")
	ErrInvalidParticipant    = errorsmod.Register(ModuleName, 8, "invalid backing-swap participant")
	ErrUnauthorized          = errorsmod.Register(ModuleName, 9, "unauthorized canonical USDC action")
	ErrInvalidAmount         = errorsmod.Register(ModuleName, 10, "invalid canonical USDC amount")
	ErrSettlementNotFound    = errorsmod.Register(ModuleName, 11, "pending settlement not found")
	ErrDuplicateSettlement   = errorsmod.Register(ModuleName, 12, "pending settlement already exists")
	ErrNobleDepositsDisabled = errorsmod.Register(ModuleName, 13, "new Noble USDC deposits are disabled")
	ErrInvalidFundingMemo    = errorsmod.Register(ModuleName, 14, "invalid backing-funding memo")
	ErrMigrationCeiling      = errorsmod.Register(ModuleName, 15, "migration ceiling exceeded")
	ErrBypass                = errorsmod.Register(ModuleName, 16, "canonical USDC transfer bypass rejected")
	ErrUnclassifiedReturn    = errorsmod.Register(ModuleName, 17, "unclassified legacy canonical USDC return")
)
