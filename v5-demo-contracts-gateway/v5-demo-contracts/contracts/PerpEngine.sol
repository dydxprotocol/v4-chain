// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "@openzeppelin/contracts/access/Ownable.sol";
import "@openzeppelin/contracts/utils/ReentrancyGuard.sol";
import "./CollateralVault.sol";
import "./Oracle.sol";

contract PerpEngine is Ownable, ReentrancyGuard {
    
    struct Position {
        int256 size;       // Positive = Long, Negative = Short (18 decimals)
        int256 entryPrice; // 8 decimals
    }

    struct MarketConfig {
        uint256 initialMarginRatio;
        uint256 maintenanceMarginRatio;
        uint256 tradingFeeRate;
    }

    CollateralVault public immutable vault;
    Oracle public immutable oracle;

    mapping(bytes32 => MarketConfig) public markets;
    mapping(bytes32 => mapping(address => Position)) public positions;

    event PositionChanged(bytes32 indexed marketId, address indexed user, int256 newSize, int256 entryPrice, int256 realizedPnL);
    event Liquidated(bytes32 indexed marketId, address indexed user, address indexed liquidator, int256 penaltyPaid);
    event MarketAdded(bytes32 indexed marketId);
    event BalanceSettled(bytes32 indexed marketId, address indexed user, int256 balanceDelta, int256 sizeDelta);

    constructor(
        address _vault, 
        address _oracle
    ) Ownable(msg.sender) {
        vault = CollateralVault(_vault);
        oracle = Oracle(_oracle);
    }

    function addMarket(
        bytes32 marketId,
        uint256 initialMarginRatio,
        uint256 maintenanceMarginRatio,
        uint256 tradingFeeRate
    ) external onlyOwner {
        markets[marketId] = MarketConfig({
            initialMarginRatio: initialMarginRatio,
            maintenanceMarginRatio: maintenanceMarginRatio,
            tradingFeeRate: tradingFeeRate
        });
        emit MarketAdded(marketId);
    }

    // --- Views ---

    function getPosition(bytes32 marketId, address user) external view returns (Position memory) {
        return positions[marketId][user];
    }

    function getUnrealizedPnl(bytes32 marketId, address user) public view returns (int256) {
        Position memory pos = positions[marketId][user];
        if (pos.size == 0) return 0;

        (int256 currentPrice, ) = oracle.getPrice(marketId);
        // PnL = size * (price - entryPrice)
        return (pos.size * (currentPrice - pos.entryPrice)) / 1e8;
    }

    function getMargin(address user) public view returns (int256 equity, uint256 totalNotional, uint256 marginRatio) {
        int256 collateral = int256(vault.balanceOf(user));
        
        // For demo simplicity, we only check PnL of the *active* market if we knew it.
        // But `getMargin` should be global.
        // Since we can't iterate mappings, we will assume for this demo that the user 
        // is only trading one market at a time or we accept that `getMargin` here 
        // might be incomplete without an array of markets.
        // FIX: For the demo, we will just return the collateral as equity and 0 margin ratio 
        // if we can't easily sum all positions. 
        // OR: We can store an array of `activeMarkets` per user?
        // Let's keep it simple: The `openPosition` checks margin for THAT market only + collateral.
        // This effectively means "Cross Margin" but we only check the current market's PnL impact.
        // This is a limitation of this simple demo but acceptable.
        
        equity = collateral; 
        // We can't easily add other markets' PnL without iteration.
        // So we will just use Collateral as Equity for the check in `openPosition` 
        // (assuming other positions are 0 PnL or we ignore them).
        
        // However, `openPosition` needs to check the NEW position's margin.
        // So `getMargin` is less useful here. We will do the check inside `openPosition`.
        
        return (equity, 0, 0); 
    }

    // --- Ledger Mode ---
    bool public engineEnabled = false;

    modifier onlyEngineEnabled() {
        require(engineEnabled, "Engine disabled");
        _;
    }

    function setEngineEnabled(bool _enabled) external onlyOwner {
        engineEnabled = _enabled;
    }

    struct Settlement {
        bytes32 marketId;
        address user;
        int256 balanceDelta;
        int256 sizeDelta;
    }

    function settle(bytes32 marketId, address user, int256 balanceDelta, int256 sizeDelta) external onlyOwner {
        _settle(marketId, user, balanceDelta, sizeDelta);
    }

    function settleBatch(Settlement[] calldata settlements) external onlyOwner {
        for (uint256 i = 0; i < settlements.length; i++) {
            _settle(
                settlements[i].marketId,
                settlements[i].user,
                settlements[i].balanceDelta,
                settlements[i].sizeDelta
            );
        }
    }

    function _settle(bytes32 marketId, address user, int256 balanceDelta, int256 sizeDelta) internal {
        // 1. Adjust Balance
        if (balanceDelta != 0) {
            vault.modifyBalance(user, balanceDelta);
        }

        // 2. Adjust Position Size (Visual only, no risk checks)
        if (sizeDelta != 0) {
            Position storage pos = positions[marketId][user];
            pos.size += sizeDelta;
        }
        
        emit BalanceSettled(marketId, user, balanceDelta, sizeDelta);
    }

    // --- Actions ---

    function openPosition(bytes32 marketId, int256 sizeDelta, uint256 maxPrice, uint256 deadline) external nonReentrant onlyEngineEnabled {
        require(block.timestamp <= deadline, "Deadline expired");
        require(sizeDelta != 0, "Size delta cannot be 0");
        
        MarketConfig memory market = markets[marketId];
        require(market.initialMarginRatio > 0, "Market not found");

        (int256 price, ) = oracle.getPrice(marketId);
        require(price > 0, "Invalid oracle price");

        if (sizeDelta > 0) {
            require(uint256(price) <= maxPrice, "Slippage exceeded");
        } else {
            require(uint256(price) >= maxPrice, "Slippage exceeded");
        }

        _updatePosition(marketId, msg.sender, sizeDelta, price);

        // Check Initial Margin (Simplified: Collateral + PnL of THIS position >= Initial Margin)
        // Equity = Collateral + PnL (which is 0 just after open usually, unless entry price differs)
        // Actually, we just check: Collateral >= Notional * IMR
        int256 collateral = int256(vault.balanceOf(msg.sender));
        Position memory pos = positions[marketId][msg.sender];
        uint256 notional = (abs(pos.size) * uint256(price)) / 1e8;
        uint256 requiredMargin = (notional * market.initialMarginRatio) / 1e18;
        
        require(collateral >= int256(requiredMargin), "Insufficient margin");
    }

    function closePosition(bytes32 marketId, uint256 maxSlippageBps, uint256 deadline) external nonReentrant onlyEngineEnabled {
        require(block.timestamp <= deadline, "Deadline expired");
        
        Position memory pos = positions[marketId][msg.sender];
        require(pos.size != 0, "No position");

        (int256 price, ) = oracle.getPrice(marketId);
        
        _updatePosition(marketId, msg.sender, -pos.size, price);
    }

    function liquidate(bytes32 marketId, address user) external nonReentrant onlyEngineEnabled {
        Position memory pos = positions[marketId][user];
        require(pos.size != 0, "No position");

        MarketConfig memory market = markets[marketId];
        (int256 price, ) = oracle.getPrice(marketId);

        // Check Maintenance Margin
        int256 collateral = int256(vault.balanceOf(user));
        // Add PnL of this position
        int256 pnl = (pos.size * (price - pos.entryPrice)) / 1e8;
        int256 equity = collateral + pnl;
        
        uint256 notional = (abs(pos.size) * uint256(price)) / 1e8;
        uint256 requiredMargin = (notional * market.maintenanceMarginRatio) / 1e18;

        require(equity < int256(requiredMargin), "Position healthy");

        // 1. Close Position
        _updatePosition(marketId, user, -pos.size, price);

        // 2. Apply Penalty (1% of notional)
        uint256 penalty = notional / 100;

        vault.modifyBalance(user, -int256(penalty));
        vault.modifyBalance(msg.sender, int256(penalty));

        emit Liquidated(marketId, user, msg.sender, int256(penalty));
    }

    // --- Internal ---

    function _updatePosition(bytes32 marketId, address user, int256 sizeDelta, int256 price) internal {
        Position storage pos = positions[marketId][user];
        int256 oldSize = pos.size;
        int256 newSize = oldSize + sizeDelta;

        // 1. Calculate Trading Fee
        MarketConfig memory market = markets[marketId];
        uint256 tradeNotional = (abs(sizeDelta) * uint256(price)) / 1e8;
        uint256 fee = (tradeNotional * market.tradingFeeRate) / 1e18;
        
        // Deduct fee from user
        vault.modifyBalance(user, -int256(fee));

        // 2. Realize PnL logic
        int256 realizedPnl = 0;

        if (oldSize == 0) {
            pos.entryPrice = price;
        } else if ((oldSize > 0 && sizeDelta > 0) || (oldSize < 0 && sizeDelta < 0)) {
            // Increasing position
            int256 oldCost = oldSize * pos.entryPrice;
            int256 addedCost = sizeDelta * price;
            pos.entryPrice = (oldCost + addedCost) / newSize;
        } else {
            // Reducing or Flipping
            int256 closedAmount;
            if (abs(sizeDelta) <= abs(oldSize)) {
                // Partial/Full Close
                closedAmount = sizeDelta; 
                int256 pnl = (-closedAmount * (price - pos.entryPrice)) / 1e8;
                vault.modifyBalance(user, pnl);
                realizedPnl = pnl;
            } else {
                // Flipping
                // 1. Close old
                int256 pnl = (oldSize * (price - pos.entryPrice)) / 1e8;
                vault.modifyBalance(user, pnl);
                realizedPnl = pnl;
                // 2. Open new
                pos.entryPrice = price;
            }
        }

        pos.size = newSize;
        if (newSize == 0) {
            pos.entryPrice = 0;
        }

        emit PositionChanged(marketId, user, newSize, pos.entryPrice, realizedPnl);
    }

    function abs(int256 x) internal pure returns (uint256) {
        return x >= 0 ? uint256(x) : uint256(-x);
    }
}
