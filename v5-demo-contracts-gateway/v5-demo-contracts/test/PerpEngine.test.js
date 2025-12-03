const { expect } = require("chai");
const { ethers } = require("hardhat");

describe("PerpEngine System", function () {
    let vault;
    let oracle;
    let engine;
    let usdc;

    let owner;
    let user1;
    let user2;
    let liquidator;

    const INITIAL_PRICE = 3000000000000n; // $30,000 with 8 decimals? No, 30000 * 10^8 = 3000000000000
    // Wait, 30,000.00 -> 30000 * 10^8 = 3,000,000,000,000.
    // Spec says: "30,000.00 -> 3000000000". That is 30,000 * 10^5?
    // Spec: "Prices are signed integers with 8 decimals (e.g. 30,000.00 -> 3000000000)"
    // 30,000 * 10^8 = 3,000,000,000,000.
    // 30,000 * 10^5 = 3,000,000,000.
    // 3000000000 is 30,000 * 10^5.
    // But spec says "8 decimals".
    // If 8 decimals, 1.0 = 100,000,000.
    // 30,000.0 = 3,000,000,000,000.
    // The example `3000000000` looks like 5 decimals or typo in spec.
    // Let's assume 8 decimals as stated in text "8 decimals".
    // So 30,000 = 30000 * 10^8 = 3000000000000.

    // Wait, let's check the example again.
    // "30,000.00 -> 3000000000"
    // 3,000,000,000.
    // 3 billion.
    // 30,000 * 100,000 = 3,000,000,000.
    // So the example uses 5 decimals?
    // Or maybe 30.00?
    // Let's stick to "8 decimals" text.
    // 30,000 * 1e8 = 3,000,000,000,000.

    const PRICE_DECIMALS = 8;
    const INITIAL_PRICE_VAL = 30000;
    const INITIAL_PRICE_BN = BigInt(INITIAL_PRICE_VAL) * (10n ** BigInt(PRICE_DECIMALS));

    const INITIAL_MARGIN_RATIO = ethers.parseEther("0.1"); // 10%
    const MAINTENANCE_MARGIN_RATIO = ethers.parseEther("0.05"); // 5%
    const TRADING_FEE_RATE = ethers.parseEther("0.0005"); // 0.05%

    beforeEach(async function () {
        [owner, user1, user2, liquidator] = await ethers.getSigners();

        // Deploy Mock USDC
        const MockERC20Factory = await ethers.getContractFactory("MockERC20");
        usdc = await MockERC20Factory.deploy("USDC", "USDC");

        // Deploy Vault
        const VaultFactory = await ethers.getContractFactory("CollateralVault");
        vault = await VaultFactory.deploy(await usdc.getAddress());

        // Deploy Oracle
        const OracleFactory = await ethers.getContractFactory("Oracle");
        oracle = await OracleFactory.deploy();
        await oracle.setPrice(INITIAL_PRICE_BN);

        // Deploy Engine
        const EngineFactory = await ethers.getContractFactory("PerpEngine");
        engine = await EngineFactory.deploy(
            await vault.getAddress(),
            await oracle.getAddress(),
            INITIAL_MARGIN_RATIO,
            MAINTENANCE_MARGIN_RATIO,
            TRADING_FEE_RATE
        );

        // Wire dependencies
        await vault.setPerpEngine(await engine.getAddress());

        // Enable Engine for tests
        await engine.setEngineEnabled(true);

        // Mint and Approve tokens
        await usdc.mint(user1.address, ethers.parseEther("10000"));
        await usdc.mint(user2.address, ethers.parseEther("10000"));
        await usdc.mint(liquidator.address, ethers.parseEther("10000"));

        await usdc.connect(user1).approve(await vault.getAddress(), ethers.MaxUint256);
        await usdc.connect(user2).approve(await vault.getAddress(), ethers.MaxUint256);
        await usdc.connect(liquidator).approve(await vault.getAddress(), ethers.MaxUint256);
    });

    describe("CollateralVault", function () {
        it("Should accept deposits", async function () {
            await vault.connect(user1).deposit(ethers.parseEther("1000"));
            expect(await vault.balanceOf(user1.address)).to.equal(ethers.parseEther("1000"));
        });

        it("Should allow withdrawals", async function () {
            await vault.connect(user1).deposit(ethers.parseEther("1000"));
            await vault.connect(user1).withdraw(ethers.parseEther("500"));
            expect(await vault.balanceOf(user1.address)).to.equal(ethers.parseEther("500"));
        });
    });

    describe("Oracle", function () {
        it("Should update price", async function () {
            const newPrice = BigInt(40000) * (10n ** BigInt(PRICE_DECIMALS));
            await oracle.setPrice(newPrice);
            const [price] = await oracle.getPrice();
            expect(price).to.equal(newPrice);
        });
    });

    describe("PerpEngine", function () {
        beforeEach(async function () {
            await vault.connect(user1).deposit(ethers.parseEther("1000"));
        });

        it("Should open a long position", async function () {
            const size = ethers.parseEther("0.1");
            const block = await ethers.provider.getBlock('latest');
            const deadline = block.timestamp + 60;

            await engine.connect(user1).openPosition(size, INITIAL_PRICE_BN + 100n, deadline);

            const pos = await engine.getPosition(user1.address);
            expect(pos.size).to.equal(size);
            expect(pos.entryPrice).to.equal(INITIAL_PRICE_BN);
        });

        it("Should charge trading fees", async function () {
            const size = ethers.parseEther("0.1");
            const block = await ethers.provider.getBlock('latest');
            const deadline = block.timestamp + 60;

            // Fee = 0.05% of notional
            // Notional = 0.1 * 30,000 = 3,000.
            // Fee = 3,000 * 0.0005 = 1.5 USDC.

            await engine.connect(user1).openPosition(size, INITIAL_PRICE_BN + 100n, deadline);

            const balance = await vault.balanceOf(user1.address);
            // 1000 - 1.5 = 998.5
            expect(balance).to.equal(ethers.parseEther("998.5"));
        });

        it("Should liquidate undercollateralized position", async function () {
            // 1. Open max leverage position
            const size = ethers.parseEther("0.3");
            const block = await ethers.provider.getBlock('latest');
            const deadline = block.timestamp + 60;
            await engine.connect(user1).openPosition(size, INITIAL_PRICE_BN + 100n, deadline);

            // 2. Drop price
            // Set price to 28,000.
            const newPrice = BigInt(28000) * (10n ** BigInt(PRICE_DECIMALS));
            await oracle.setPrice(newPrice);

            // 3. Liquidate
            await engine.connect(liquidator).liquidate(user1.address);

            // Check position closed
            const pos = await engine.getPosition(user1.address);
            expect(pos.size).to.equal(0);

            // Check penalty paid
            // Penalty = 1% of notional at close.
            // Notional = 0.3 * 28,000 = 8,400.
            // Penalty = 84 USDC.
            expect(await vault.balanceOf(liquidator.address)).to.equal(ethers.parseEther("84"));
        });
    });
});

