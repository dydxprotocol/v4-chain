// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "forge-std/Test.sol";
import "../contracts/MockERC20.sol";
import "../contracts/CollateralVault.sol";
import "../contracts/Oracle.sol";
import "../contracts/PerpEngine.sol";

contract PerpEngineTest is Test {
    CollateralVault public vault;
    Oracle public oracle;
    PerpEngine public engine;
    MockERC20 public usdc;

    address public owner;
    address public user1;
    address public user2;
    address public liquidator;

    uint256 constant INITIAL_PRICE_VAL = 30000;
    uint8 constant PRICE_DECIMALS = 8;
    uint256 constant INITIAL_PRICE = INITIAL_PRICE_VAL * 10**PRICE_DECIMALS;

    bytes32 constant BTC_MARKET = keccak256("BTC-USD");

    function setUp() public {
        owner = address(this);
        user1 = address(0x1);
        user2 = address(0x2);
        liquidator = address(0x3);

        // Deploy Mock USDC
        usdc = new MockERC20("USDC", "USDC");

        // Deploy Vault
        vault = new CollateralVault(address(usdc));

        // Deploy Oracle
        oracle = new Oracle();
        oracle.setPrice(BTC_MARKET, int256(INITIAL_PRICE));

        // Deploy Engine
        engine = new PerpEngine(address(vault), address(oracle));

        // Wire dependencies
        vault.setPerpEngine(address(engine));
        engine.setEngineEnabled(true);

        // Setup market
        engine.addMarket(
            BTC_MARKET,
            0.1e18,  // 10% initial margin
            0.05e18, // 5% maintenance margin
            0.001e18 // 0.1% trading fee
        );

        // Mint tokens
        usdc.mint(user1, 10000e18);
        usdc.mint(user2, 10000e18);
        usdc.mint(liquidator, 10000e18);

        // Approve vault
        vm.prank(user1);
        usdc.approve(address(vault), type(uint256).max);
        vm.prank(user2);
        usdc.approve(address(vault), type(uint256).max);
        vm.prank(liquidator);
        usdc.approve(address(vault), type(uint256).max);
    }

    function testVaultDeposit() public {
        vm.prank(user1);
        vault.deposit(1000e18);
        assertEq(vault.balanceOf(user1), 1000e18);
    }

    function testVaultWithdraw() public {
        vm.prank(user1);
        vault.deposit(1000e18);
        vm.prank(user1);
        vault.withdraw(500e18);
        assertEq(vault.balanceOf(user1), 500e18);
    }

    function testOracleSetPrice() public {
        uint256 newPrice = 40000 * 10**PRICE_DECIMALS;
        oracle.setPrice(BTC_MARKET, int256(newPrice));
        (int256 price,) = oracle.getPrice(BTC_MARKET);
        assertEq(uint256(price), newPrice);
    }

    function testOpenPosition() public {
        vm.prank(user1);
        vault.deposit(1000e18);

        uint256 size = 0.1e18;
        uint256 maxPrice = INITIAL_PRICE + 100;
        uint256 deadline = block.timestamp + 60;

        vm.prank(user1);
        engine.openPosition(BTC_MARKET, int256(size), maxPrice, deadline);

        PerpEngine.Position memory pos = engine.getPosition(BTC_MARKET, user1);
        assertEq(uint256(pos.size), size);
        assertEq(uint256(pos.entryPrice), INITIAL_PRICE);
    }

    function testTradingFee() public {
        vm.prank(user1);
        vault.deposit(1000e18);

        uint256 size = 0.1e18;
        uint256 maxPrice = INITIAL_PRICE + 100;
        uint256 deadline = block.timestamp + 60;

        vm.prank(user1);
        engine.openPosition(BTC_MARKET, int256(size), maxPrice, deadline);

        // Fee = 0.1% of notional = 0.1 * 30000 * 0.001 = 3 USDC
        uint256 balance = vault.balanceOf(user1);
        assertEq(balance, 1000e18 - 3e18);
    }

    function testLiquidation() public {
        vm.prank(user1);
        vault.deposit(1000e18);

        // Open max leverage position
        uint256 size = 0.3e18;
        uint256 maxPrice = INITIAL_PRICE + 100;
        uint256 deadline = block.timestamp + 60;

        vm.prank(user1);
        engine.openPosition(BTC_MARKET, int256(size), maxPrice, deadline);

        // Drop price to trigger liquidation
        uint256 newPrice = 28000 * 10**PRICE_DECIMALS;
        oracle.setPrice(BTC_MARKET, int256(newPrice));

        // Liquidate
        vm.prank(liquidator);
        engine.liquidate(BTC_MARKET, user1);

        // Check position closed
        PerpEngine.Position memory pos = engine.getPosition(BTC_MARKET, user1);
        assertEq(pos.size, 0);

        // Check penalty paid (1% of notional)
        // Notional = 0.3 * 28000 = 8400, Penalty = 84 USDC
        uint256 penalty = (size * newPrice / 1e8) / 100;
        uint256 liquidatorBalance = vault.balanceOf(liquidator);
        assertGe(liquidatorBalance, penalty); // Greater than or equal (liquidator starts with 10000)
    }
}

