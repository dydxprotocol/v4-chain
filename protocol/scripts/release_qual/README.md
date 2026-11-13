# Release Qualification Scripts

This directory contains scripts for testing and qualifying new releases on testnet/staging environments.

## Scripts

### submit_upgrade_proposal.py

A Python script to submit software upgrade proposals and automatically vote with test validators on testnet/staging environments.

**Features:**
- Submits upgrade proposals with a specified upgrade name and block height
- Automatically votes "yes" with all test validators
- Restricted to testnet/staging environments only (mainnet explicitly blocked)
- Configurable via command-line arguments or environment variables

**Usage:**
```bash
# Basic usage - upgrade in 300 blocks (default)
./submit_upgrade_proposal.py v5.0.0

# Specify number of blocks to wait
./submit_upgrade_proposal.py v5.0.0 500

# Use custom node and chain-id
./submit_upgrade_proposal.py v5.0.0 --node https://validator.custom.com:443 --chain-id dydxprotocol-testnet

# Using environment variables
export DYDX_NODE=https://validator.custom.com:443
export DYDX_CHAIN_ID=dydxprotocol-testnet
./submit_upgrade_proposal.py v5.0.0
```

**Configuration Priority:**
1. Command-line arguments (highest priority)
2. Environment variables (`DYDX_NODE`, `DYDX_CHAIN_ID`)
3. Default values (v4staging node and testnet chain-id)

**Safety Features:**
- Only works on allowed testnet/staging chain IDs
- Explicitly blocks mainnet chain IDs
- Requires user confirmation before submitting proposal
- Validates chain ID against allowed list

**Test Validators:**
The script automatically votes with these test validators:
- alice, bob, carl, dave, emily, fiona, greg, henry, ian, jeff

## Requirements

- Python 3
- `dydxprotocold` CLI installed and configured
- Test validator keys in the keyring (using `test` backend)
- Access to a testnet/staging node

## Security

These scripts are designed for **testnet/staging environments only**. The `submit_upgrade_proposal.py` script includes multiple safeguards to prevent accidental use on mainnet:
- Hardcoded list of allowed chain IDs (testnet/staging only)
- Explicit blocking of known mainnet chain IDs
- Chain ID validation before execution