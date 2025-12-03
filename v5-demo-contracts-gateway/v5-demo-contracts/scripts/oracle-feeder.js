/**
 * Oracle feeder script compatible with Foundry
 * Uses ethers.js with Foundry artifacts instead of Hardhat
 */

const { ethers } = require("ethers");
require("dotenv").config();

async function main() {
    const oracleAddress = process.env.ORACLE_ADDRESS;
    if (!oracleAddress) {
        throw new Error("ORACLE_ADDRESS env var not set");
    }

    const priceArg = process.argv[2];
    if (!priceArg) {
        throw new Error("Please provide price as argument (e.g. 30000)");
    }

    const marketArg = process.argv[3] || "BTC-USD"; // Default to BTC-USD
    const marketId = ethers.encodeBytes32String(marketArg);

    const priceVal = BigInt(priceArg);
    // Assume input is "30000" for $30k. We need 8 decimals.
    // 30000 * 10^8 = 3,000,000,000,000.
    const price = priceVal * 100000000n;

    console.log(`Updating Oracle at ${oracleAddress} for ${marketArg} to price ${priceVal} ($${priceVal}) -> raw ${price}`);

    // Connect to network
    const rpcUrl = process.env.RPC_URL || "https://sepolia.base.org";
    const privateKey = process.env.PRIVATE_KEY;
    
    if (!privateKey) {
        throw new Error("PRIVATE_KEY environment variable is required");
    }

    const provider = new ethers.JsonRpcProvider(rpcUrl);
    const signer = new ethers.Wallet(privateKey, provider);

    // Load Oracle artifact
    const OracleArtifact = require("../out/Oracle.sol/Oracle.json");
    const oracle = new ethers.Contract(oracleAddress, OracleArtifact.abi, signer);

    const tx = await oracle.setPrice(marketId, price);
    console.log("Tx sent:", tx.hash);
    await tx.wait();
    console.log("Price updated.");
}

main().catch((error) => {
    console.error(error);
    process.exitCode = 1;
});
