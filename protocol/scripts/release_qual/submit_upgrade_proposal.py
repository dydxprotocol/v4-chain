#!/usr/bin/env python3
"""
Submit a software upgrade proposal and auto-vote with test validators.
"""

import argparse
import json
import os
import subprocess
import sys
import time

# Test validators
VALIDATORS = ["alice", "bob", "carl", "dave", "emily", "fiona", "greg", "henry", "ian", "jeff"]

# Allowed chain IDs for testnet/staging only (blocking mainnet)
ALLOWED_CHAIN_IDS = [
    "dydxprotocol-testnet",  # Standard testnet/staging chain
]

# Explicitly blocked mainnet chain IDs
BLOCKED_CHAIN_IDS = [
    "dydx-mainnet-1",
    "dydxprotocol-mainnet",
    "mainnet"
]

def load_config(args):
    """Load configuration from command-line args or environment variables.
    Priority: command-line args > environment variables > defaults
    """
    # Default values
    default_node = "https://validator.v4staging.dydx.exchange:443"
    default_chain_id = "dydxprotocol-testnet"

    # First priority: command-line arguments
    if args.node:
        node = args.node
        print(f"Using node from command-line: {node}")
    else:
        # Second priority: environment variables, fallback to defaults
        node = os.environ.get("DYDX_NODE", default_node)

    if args.chain_id:
        chain_id = args.chain_id
        print(f"Using chain-id from command-line: {chain_id}")
    else:
        # Second priority: environment variables, fallback to defaults
        chain_id = os.environ.get("DYDX_CHAIN_ID", default_chain_id)

    # Check if explicitly blocked (mainnet)
    if chain_id in BLOCKED_CHAIN_IDS or "mainnet" in chain_id.lower():
        print(f"Error: This script cannot run on mainnet (chain_id: {chain_id})")
        print("This script is only for testnet/staging environments.")
        sys.exit(1)

    # Validate chain_id is allowed
    if chain_id not in ALLOWED_CHAIN_IDS:
        print(f"Error: Chain ID '{chain_id}' is not in the allowed list.")
        print(f"Allowed chains: {', '.join(ALLOWED_CHAIN_IDS)}")
        print("This script is restricted to testnet/staging environments only.")
        sys.exit(1)

    print(f"Using node: {node}")
    print(f"Using chain-id: {chain_id}")

    return node, chain_id

def run_cmd(cmd, node=None):
    """Run command and return stdout."""
    # Add node flag if provided
    if node and "--node" not in cmd:
        cmd.extend(["--node", node])
    try:
        result = subprocess.run(
            cmd, capture_output=True, text=True, check=True
        )
        return result.stdout
    except subprocess.CalledProcessError as e:
        print(f"Error: {e.stderr}")
        return None

def main():
    """Submit a software upgrade proposal and auto-vote with test validators."""
    parser = argparse.ArgumentParser(
        description='Submit a software upgrade proposal and auto-vote with test validators.'
    )
    parser.add_argument('upgrade_name', help='Name of the upgrade (e.g., v5.0.0)')
    parser.add_argument(
        'blocks_to_wait', nargs='?', type=int, default=300,
        help='Number of blocks to wait for an upgrade and voting period (default: 300)'
    )
    parser.add_argument(
        '--node',
        help='Node RPC endpoint (e.g., http://validator.v4staging.dydx.exchange:26657)'
    )
    parser.add_argument(
        '--chain-id', dest='chain_id', help='Chain ID (e.g., dydxprotocol-testnet)'
    )

    args = parser.parse_args()

    upgrade_name = args.upgrade_name
    wait_blocks = args.blocks_to_wait

    # Load configuration
    node, chain_id = load_config(args)

    # Display configuration and ask for confirmation
    print("\n" + "="*60)
    print("UPGRADE PROPOSAL CONFIGURATION")
    print("="*60)
    print(f"Chain ID:      {chain_id}")
    print(f"Node:          {node}")
    print(f"Upgrade Name:  {upgrade_name}")
    print(f"Block Wait:    {wait_blocks} blocks")
    print("="*60)

    response = input("\nDo you want to proceed with this upgrade proposal? (yes/no): ")
    if response.lower() not in ['yes', 'y']:
        print("Upgrade proposal cancelled.")
        sys.exit(0)

    print("\nProceeding with upgrade proposal...")

    # Get current block height
    result = run_cmd(["dydxprotocold", "status"], node=node)
    if result:
        try:
            status = json.loads(result)
            current_height = int(status['sync_info']['latest_block_height'])
            upgrade_height = current_height + wait_blocks
            print(f"Current height: {current_height}, upgrade at: {upgrade_height}")
        except (json.JSONDecodeError, KeyError) as e:
            # Fallback if we can't get current height
            print(f"Could not parse block height, using default. Error: {e}")
            upgrade_height = 1000000
            print(f"Using default upgrade height: {upgrade_height}")
    else:
        upgrade_height = 1000000

    # Create proposal.json
    proposal = {
        "messages": [{
            "@type": "/cosmos.upgrade.v1beta1.MsgSoftwareUpgrade",
            "authority": "dydx10d07y265gmmuvt4z0w9aw880jnsr700jnmapky",
            "plan": {
                "name": upgrade_name,
                "height": str(upgrade_height),
                "info": f"Upgrade to {upgrade_name}"
            }
        }],
        "metadata": "",
        "deposit": "20000000adv4tnt",
        "title": f"Software Upgrade to {upgrade_name}",
        "summary": f"Upgrade the chain to {upgrade_name}"
    }

    with open("proposal.json", "w", encoding="utf-8") as f:
        json.dump(proposal, f, indent=2)

    print(f"Submitting upgrade proposal for {upgrade_name} at height {upgrade_height}...")

    # Submit proposal
    cmd = [
        "dydxprotocold", "tx", "gov", "submit-proposal", "proposal.json",
        "--from", "alice",
        "--chain-id", chain_id,
        "--yes",
        "--broadcast-mode", "sync",
        "--gas", "auto",
        "--fees", "5000000000000000adv4tnt",
        "--keyring-backend", "test"
    ]

    result = run_cmd(cmd, node=node)
    if not result:
        os.remove("proposal.json")
        sys.exit(1)

    # Extract txhash
    for line in result.split('\n'):
        if 'txhash:' in line:
            print(f"Submitted: {line.split('txhash:')[1].strip()}")

    time.sleep(5)

    # Get proposal ID
    result = run_cmd(["dydxprotocold", "query", "gov", "proposals", "--output", "json"], node=node)
    if not result:
        os.remove("proposal.json")
        sys.exit(1)

    proposals = json.loads(result)
    proposal_id = proposals['proposals'][-1]['id']
    print(f"Proposal ID: {proposal_id}")

    # Vote
    print(f"Voting with {len(VALIDATORS)} validators...")
    for voter in VALIDATORS:
        cmd = [
            "dydxprotocold", "tx", "gov", "vote", str(proposal_id), "yes",
            "--from", voter,
            "--chain-id", chain_id,
            "--yes",
            "--gas", "auto",
            "--fees", "5000000000000000adv4tnt",
            "--keyring-backend", "test"
        ]
        result = run_cmd(cmd, node=node)
        if result and 'txhash:' in result:
            print(f"  {voter}: ✓")
        else:
            print(f"  {voter}: ✗")
        time.sleep(1)

    # Clean up
    os.remove("proposal.json")

    # Wait for voting period (same as blocks to upgrade)
    print(f"\nWaiting {wait_blocks} seconds (~{wait_blocks} blocks) for voting period...")
    for i in range(0, wait_blocks, 30):
        remaining = wait_blocks - i
        if remaining > 0:
            print(f"  {remaining}s remaining...")
            time.sleep(min(30, remaining))

    # Check status
    cmd = [
        "dydxprotocold", "query", "gov", "proposal",
        str(proposal_id), "--output", "json"
    ]
    result = run_cmd(cmd, node=node)
    if result:
        data = json.loads(result)
        status = data['proposal']['status']
        print(f"\nFinal status: {status}")
        if "PASSED" in status:
            print(f"✅ Upgrade to {upgrade_name} approved for height {upgrade_height}")

if __name__ == "__main__":
    main()
