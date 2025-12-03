/**
 * Demo Batch Order Script
 * Places a single order via Gateway API and listens for confirmation
 */

const { ethers } = require("ethers");
require("dotenv").config();

async function main() {
    const userAddress = process.env.USER_ADDRESS;
    if (!userAddress) {
        throw new Error("USER_ADDRESS environment variable is required");
    }
    
    const gatewayUrl = process.env.GATEWAY_URL || "http://localhost:8080";
    const rpcUrl = process.env.RPC_URL || "https://sepolia.base.org";
    
    // Load deployment config
    const fs = require("fs");
    const path = require("path");
    const deploymentPath = path.resolve(__dirname, "../deployment.json");
    
    if (!fs.existsSync(deploymentPath)) {
        throw new Error("deployment.json not found. Deploy contracts first.");
    }
    
    const deploymentConfig = JSON.parse(fs.readFileSync(deploymentPath));
    const provider = new ethers.JsonRpcProvider(rpcUrl);
    const network = await provider.getNetwork();
    const chainId = network.chainId.toString();
    
    const chainConfig = deploymentConfig[chainId];
    if (!chainConfig) {
        throw new Error(`No config found for Chain ID ${chainId}`);
    }
    
    const engineAddress = chainConfig.addresses.engine;
    const BTC_MARKET = chainConfig.markets.find(m => m.name === "BTC-USD")?.id;
    
    if (!BTC_MARKET) {
        throw new Error("BTC-USD market not found in config");
    }
    
    console.log("Demo Batch Order");
    console.log("================");
    console.log("Chain ID:", chainId);
    console.log("User Address:", userAddress);
    console.log("Engine Address:", engineAddress);
    console.log("Gateway URL:", gatewayUrl);
    console.log("");
    
    // Setup event listener
    const PerpEngineArtifact = require("../out/PerpEngine.sol/PerpEngine.json");
    const engine = new ethers.Contract(engineAddress, PerpEngineArtifact.abi, provider);
    
    // Place order via Gateway
    // Using 0.01 USDC per order (allows 100 orders from 1 USDC deposit)
    const order = {
        marketId: "BTC-USD",
        userId: "demo-user-1",
        evmAddress: userAddress,
        balanceDelta: "-10000",      // -0.01 USDC (6 decimals)
        sizeDelta: "10000000000000000", // +0.01 BTC (18 decimals)
        reason: "demo-batch-order"
    };
    
    console.log("Placing order:", JSON.stringify(order, null, 2));
    console.log("");
    
    const response = await fetch(`${gatewayUrl}/settle-batch`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ settlements: [order] })
    });
    
    const result = await response.json();
    
    if (!response.ok || result.error) {
        throw new Error(result.error || `HTTP ${response.status}`);
    }
    
    console.log("✅ Order submitted");
    console.log("Tx Hash:", result.txHash);
    console.log("Waiting for confirmation...");
    console.log("");
    
    // Base Sepolia RPC doesn't support event filters, so we poll the transaction receipt
    const receipt = await provider.waitForTransaction(result.txHash, 1, 60000);
    
    // Parse events from receipt
    const iface = new ethers.Interface(PerpEngineArtifact.abi);
    let eventData = null;
    for (const log of receipt.logs) {
        try {
            const parsed = iface.parseLog(log);
            if (parsed && parsed.name === "BalanceSettled") {
                eventData = {
                    marketId: parsed.args.marketId,
                    user: parsed.args.user,
                    balanceDelta: parsed.args.balanceDelta.toString(),
                    sizeDelta: parsed.args.sizeDelta.toString(),
                    txHash: receipt.hash,
                    blockNumber: receipt.blockNumber
                };
                break;
            }
        } catch (e) {
            // Not our event, continue
        }
    }
    
    if (!eventData) {
        throw new Error("BalanceSettled event not found in transaction receipt");
    }
    
    console.log("✅ Order confirmed!");
    console.log("Event Data:");
    console.log("  Market ID:", eventData.marketId);
    console.log("  User:", eventData.user);
    console.log("  Balance Delta:", eventData.balanceDelta, "(USDC, 6 decimals)");
    console.log("  Size Delta:", eventData.sizeDelta, "(BTC, 18 decimals)");
    console.log("  Tx Hash:", eventData.txHash);
    console.log("  Block:", eventData.blockNumber);
    console.log("");
    
    // Verify on-chain state
    console.log("Verifying on-chain state...");
    const position = await engine.getPosition(BTC_MARKET, userAddress);
    const vault = new ethers.Contract(chainConfig.addresses.vault, require("../out/CollateralVault.sol/CollateralVault.json").abi, provider);
    const balance = await vault.balanceOf(userAddress);
    
    console.log("Position Size:", ethers.formatUnits(position.size, 18), "BTC");
    console.log("Entry Price:", ethers.formatUnits(position.entryPrice, 8), "USD");
    console.log("Vault Balance:", ethers.formatUnits(balance, 18), "USDC");
    console.log("");
    console.log("✅ Demo complete!");
}

main().catch((error) => {
    console.error(error);
    process.exitCode = 1;
});

