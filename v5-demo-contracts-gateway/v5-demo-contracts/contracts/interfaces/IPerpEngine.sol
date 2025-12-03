// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

/**
 * @title IPerpEngine
 * @notice Core interface for the Perpetual Engine supporting multi-market trading
 */
interface IPerpEngine {
    struct Position {
        int256 size;       // Positive = Long, Negative = Short (18 decimals)
        int256 entryPrice; // 8 decimals
    }

    struct MarketConfig {
        uint256 initialMarginRatio;     // e.g. 0.1e18 = 10%
        uint256 maintenanceMarginRatio; // e.g. 0.05e18 = 5%
        uint256 tradingFeeRate;         // e.g. 0.001e18 = 0.1%
    }

    // --- Admin Functions ---
    function addMarket(
        bytes32 marketId,
        uint256 initialMarginRatio,
        uint256 maintenanceMarginRatio,
        uint256 tradingFeeRate
    ) external;

    function setEngineEnabled(bool _enabled) external;

    struct Settlement {
        bytes32 marketId;
        address user;
        int256 balanceDelta;
        int256 sizeDelta;
    }

    // --- Ledger Mode (Off-chain Settlement) ---
    function settle(
        bytes32 marketId,
        address user,
        int256 balanceDelta,
        int256 sizeDelta
    ) external;

    function settleBatch(Settlement[] calldata settlements) external;

    // --- Standalone Mode (On-chain Trading) ---
    function openPosition(
        bytes32 marketId,
        int256 sizeDelta,
        uint256 maxPrice,
        uint256 deadline
    ) external;

    function closePosition(
        bytes32 marketId,
        uint256 maxSlippageBps,
        uint256 deadline
    ) external;

    function liquidate(bytes32 marketId, address user) external;

    // --- Views ---
    function getPosition(bytes32 marketId, address user) external view returns (Position memory);
    function getUnrealizedPnl(bytes32 marketId, address user) external view returns (int256);
    function getMargin(address user) external view returns (int256 equity, uint256 totalNotional, uint256 marginRatio);

    // --- Events ---
    event PositionChanged(bytes32 indexed marketId, address indexed user, int256 newSize, int256 entryPrice, int256 realizedPnL);
    event Liquidated(bytes32 indexed marketId, address indexed user, address indexed liquidator, int256 penaltyPaid);
    event MarketAdded(bytes32 indexed marketId);
    event BalanceSettled(bytes32 indexed marketId, address indexed user, int256 balanceDelta, int256 sizeDelta);
}

/**
 * @title ICollateralVault
 * @notice Vault for managing user collateral (USDC)
 */
interface ICollateralVault {
    function deposit(uint256 amount) external;
    function withdraw(uint256 amount) external;
    function balanceOf(address user) external view returns (uint256);
    function totalCollateral() external view returns (uint256);
    
    // Engine-only
    function modifyBalance(address user, int256 amountDelta) external;
    function setPerpEngine(address _perpEngine) external;

    event Deposit(address indexed user, uint256 amount);
    event Withdraw(address indexed user, uint256 amount);
}

/**
 * @title IOracle
 * @notice Multi-market price oracle
 */
interface IOracle {
    function setPrice(bytes32 marketId, int256 price) external;
    function getPrice(bytes32 marketId) external view returns (int256 price, uint256 timestamp);

    event PriceUpdated(bytes32 indexed marketId, int256 price, uint256 timestamp);
}
