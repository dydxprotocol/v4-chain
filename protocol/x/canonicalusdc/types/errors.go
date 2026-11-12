package types

import errorsmod "cosmossdk.io/errors"

var (
	ErrPaused                = errorsmod.Register(ModuleName, 2, "canonical USDC routing is paused")
	ErrUnsupportedChannel    = errorsmod.Register(ModuleName, 3, "unsupported canonical USDC channel")
	ErrInsufficientBacking   = errorsmod.Register(ModuleName, 4, "insufficient route backing")
	ErrPendingLimit          = errorsmod.Register(ModuleName, 5, "pending settlement limit reached")
	ErrInvalidControls       = errorsmod.Register(ModuleName, 6, "invalid canonical USDC controls")
	ErrInvalidLedger         = errorsmod.Register(ModuleName, 7, "invalid canonical USDC ledger")
	ErrUnauthorized          = errorsmod.Register(ModuleName, 9, "unauthorized canonical USDC action")
	ErrInvalidAmount         = errorsmod.Register(ModuleName, 10, "invalid canonical USDC amount")
	ErrSettlementNotFound    = errorsmod.Register(ModuleName, 11, "pending settlement not found")
	ErrDuplicateSettlement   = errorsmod.Register(ModuleName, 12, "pending settlement already exists")
	ErrNobleDepositsDisabled = errorsmod.Register(ModuleName, 13, "new Noble USDC deposits are disabled")
	ErrBypass                = errorsmod.Register(ModuleName, 16, "canonical USDC transfer bypass rejected")
	ErrNobleBackingRemains   = errorsmod.Register(
		ModuleName, 17, "noble withdrawals required while Noble backing remains",
	)
	ErrModuleAccount = errorsmod.Register(
		ModuleName, 18, "canonical USDC module account cannot be established",
	)
	ErrInjectiveDepositsDisabled = errorsmod.Register(ModuleName, 19, "Injective USDC deposits are disabled")
	ErrInvalidReceiver           = errorsmod.Register(ModuleName, 20, "invalid canonical USDC receiver")
	ErrMemoNotAllowed            = errorsmod.Register(ModuleName, 21, "memo not allowed on canonical USDC withdrawal")
	ErrNobleWithdrawalsDisabled  = errorsmod.Register(ModuleName, 22, "Noble USDC withdrawals are disabled")
	ErrSettlementMismatch        = errorsmod.Register(ModuleName, 23, "pending settlement does not match packet")
	ErrInvalidSettlement         = errorsmod.Register(ModuleName, 24, "invalid pending settlement")
	ErrRateLimitRequired         = errorsmod.Register(ModuleName, 25, "no rate limit configured for Injective USDC")
)
