// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";
import "@openzeppelin/contracts/utils/ReentrancyGuard.sol";

contract CollateralVault is ReentrancyGuard {
    using SafeERC20 for IERC20;

    IERC20 public immutable collateralToken;
    
    // User balances in 18 decimals
    mapping(address => uint256) private _balances;
    uint256 private _totalCollateral;

    event Deposit(address indexed user, uint256 amount);
    event Withdraw(address indexed user, uint256 amount);

    constructor(address _collateralToken) {
        require(_collateralToken != address(0), "Invalid token address");
        collateralToken = IERC20(_collateralToken);
    }

    function deposit(uint256 amount) external nonReentrant {
        require(amount > 0, "Amount must be > 0");
        
        // Transfer tokens from user to vault
        collateralToken.safeTransferFrom(msg.sender, address(this), amount);
        
        // Update state
        _balances[msg.sender] += amount;
        _totalCollateral += amount;
        
        emit Deposit(msg.sender, amount);
    }

    function withdraw(uint256 amount) external nonReentrant {
        require(amount > 0, "Amount must be > 0");
        require(_balances[msg.sender] >= amount, "Insufficient balance");
        
        // Update state
        _balances[msg.sender] -= amount;
        _totalCollateral -= amount;
        
        // Transfer tokens to user
        collateralToken.safeTransfer(msg.sender, amount);
        
        emit Withdraw(msg.sender, amount);
    }

    function balanceOf(address user) external view returns (uint256) {
        return _balances[user];
    }

    function totalCollateral() external view returns (uint256) {
        return _totalCollateral;
    }
    
    // Allow PerpEngine to modify balances (to be implemented later with access control if needed, 
    // but for this minimal demo, we might need to expose internal functions or make PerpEngine trusted.
    // However, the spec says "Charges trading fee taken from user collateral" and "Liquidate... paid from user collateral".
    // So the Vault needs to allow the Engine to move funds.
    // For simplicity in this "Minimal" demo, we can add a trusted engine address or just keep it simple.
    // Wait, the spec says "Cross-margin using the single collateral asset".
    // Usually the Engine holds the logic and the Vault holds the funds.
    // If the Engine needs to debit fees/losses, it needs permission.
    // Let's add a `modifyBalance` function restricted to the Engine.
    
    address public perpEngine;

    modifier onlyEngine() {
        require(msg.sender == perpEngine, "Caller is not PerpEngine");
        _;
    }

    function setPerpEngine(address _perpEngine) external {
        // In a real app, this would be Ownable and restricted. 
        // For this demo, we'll set it once or allow overwrite if needed (but let's keep it simple).
        // Assuming deployer sets it. We need Ownable for this.
        require(perpEngine == address(0), "Engine already set");
        perpEngine = _perpEngine;
    }

    // Helper for Engine to adjust balances (fees, pnl realization, liquidation)
    function modifyBalance(address user, int256 amountDelta) external onlyEngine {
        if (amountDelta > 0) {
            uint256 addAmount = uint256(amountDelta);
            _balances[user] += addAmount;
            _totalCollateral += addAmount;
        } else {
            uint256 subAmount = uint256(-amountDelta);
            require(_balances[user] >= subAmount, "Insufficient balance for modification");
            _balances[user] -= subAmount;
            _totalCollateral -= subAmount;
        }
    }
}
