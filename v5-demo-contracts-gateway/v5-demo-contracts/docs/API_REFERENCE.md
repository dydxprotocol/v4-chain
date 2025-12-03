# Developer API Reference

## 1. Smart Contracts (EVM)
**Network**: Base Sepolia
**PerpEngine Address**: `0xFeBd8490C545466DC4b0d566c33ff02beDF6f303`

### Core Settlement APIs (Write)
These functions are restricted to the **Owner** (Operator). Services can call these directly to settle trades on the ledger.

*   **`settle(bytes32 marketId, address user, int256 balanceDelta, int256 sizeDelta)`**
    *   **Purpose**: Settle a single trade/update for a user.
    *   **Params**:
        *   `marketId`: e.g., `0x4254432d555344...` (BTC-USD).
        *   `user`: The user's EVM address.
        *   `balanceDelta`: Change in USDC balance (6 decimals). Positive = Credit, Negative = Debit.
        *   `sizeDelta`: Change in position size (18 decimals). Positive = Long, Negative = Short.

*   **`settleBatch(Settlement[] settlements)`**
    *   **Purpose**: Batch multiple settlements in a single transaction to save gas.
    *   **Struct**: `Settlement { bytes32 marketId; address user; int256 balanceDelta; int256 sizeDelta; }`

### View APIs (Read - Experimental For Now)
*   **`getPosition(bytes32 marketId, address user)`**
    *   **Returns**: `(int256 size, int256 entryPrice)`
    *   **Description**: Get a user's current position size and entry price.
*   **`getUnrealizedPnl(bytes32 marketId, address user)`**
    *   **Returns**: `int256 pnl`
    *   **Description**: Calculate estimated PnL based on the current Oracle price.

### CollateralVault (Custody- Experimental For Now)
*   **`deposit(uint256 amount)`**
    *   **User Action**: User calls this to deposit USDC.
*   **`withdraw(uint256 amount)`**
    *   **User Action**: User calls this to withdraw USDC.
*   **`balanceOf(address user)`**
    *   **Returns**: `uint256`
    *   **Description**: Get a user's deposited collateral balance.

---

## 2. Gateway Service (REST)
**Purpose**: A convenience service that manages nonces and signs transactions for the Operator.

### Trading Engine
*   **`POST /settle`**
    *   **Body**: `{ "marketId": "BTC-USD", "userId": "...", "evmAddress": "0x...", "balanceDelta": "...", "sizeDelta": "..." }`
*   **`POST /settle-batch`**
    *   **Body**: `{ "settlements": [ ... ] }`

### State Verification
*   **`GET /user-state?address=0x...&marketId=BTC-USD`**
    *   **Returns**: `{ "balance": "...", "position": "...", "entryPrice": "..." }`
