#!/usr/bin/env python3
"""
Submit a software upgrade proposal and auto-vote with test validators.
"""

import json
import os
import subprocess
import sys
import time

# Test validators
VALIDATORS = ["alice", "bob", "carl", "dave", "emily", "fiona", "greg", "henry", "ian", "jeff"]

def run_cmd(cmd):
    """Run command and return stdout."""
    result = subprocess.run(cmd, capture_output=True, text=True)
    if result.returncode != 0:
        print(f"Error: {result.stderr}")
        return None
    return result.stdout

def main():
    if len(sys.argv) < 2:
        print("Usage: python submit_upgrade_proposal.py <upgrade_name> [blocks_to_wait]")
        print("Example: python submit_upgrade_proposal.py v5.0.0 300")
        print("Default: 300 blocks from now for upgrade and voting period")
        sys.exit(1)
    
    upgrade_name = sys.argv[1]
    wait_blocks = int(sys.argv[2]) if len(sys.argv) > 2 else 300
    
    # Get current block height
    result = run_cmd(["dydxprotocold", "status"])
    if result:
        try:
            import json as j
            status = j.loads(result)
            current_height = int(status['sync_info']['latest_block_height'])
            upgrade_height = current_height + wait_blocks
            print(f"Current height: {current_height}, upgrade at: {upgrade_height}")
        except:
            # Fallback if we can't get current height
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
        "deposit": "10000000000000000000000adv4tnt",
        "title": f"Software Upgrade to {upgrade_name}",
        "summary": f"Upgrade the chain to {upgrade_name}"
    }
    
    with open("proposal.json", "w") as f:
        json.dump(proposal, f, indent=2)
    
    print(f"Submitting upgrade proposal for {upgrade_name} at height {upgrade_height}...")
    
    # Submit proposal
    cmd = [
        "dydxprotocold", "tx", "gov", "submit-proposal", "proposal.json",
        "--from", "alice",
        "--chain-id", "dydxprotocol-testnet",
        "--yes",
        "--broadcast-mode", "sync",
        "--gas", "auto",
        "--fees", "5000000000000000adv4tnt",
        "--keyring-backend", "test"
    ]
    
    result = run_cmd(cmd)
    if not result:
        os.remove("proposal.json")
        sys.exit(1)
    
    # Extract txhash
    for line in result.split('\n'):
        if 'txhash:' in line:
            print(f"Submitted: {line.split('txhash:')[1].strip()}")
    
    time.sleep(5)
    
    # Get proposal ID
    result = run_cmd(["dydxprotocold", "query", "gov", "proposals", "--output", "json"])
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
            "--chain-id", "dydxprotocol-testnet",
            "--yes",
            "--gas", "auto",
            "--fees", "5000000000000000adv4tnt",
            "--keyring-backend", "test"
        ]
        result = run_cmd(cmd)
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
    result = run_cmd(["dydxprotocold", "query", "gov", "proposal", str(proposal_id), "--output", "json"])
    if result:
        data = json.loads(result)
        status = data['proposal']['status']
        print(f"\nFinal status: {status}")
        if "PASSED" in status:
            print(f"✅ Upgrade to {upgrade_name} approved for height {upgrade_height}")

if __name__ == "__main__":
    main()