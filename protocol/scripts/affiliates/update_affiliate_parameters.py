#!/usr/bin/env python3
"""
Script to add an order router rev share to the protocol.
"""

import argparse
import os
import shutil
import subprocess
import sys
import tempfile
import yaml
import json
import time
from typing import Dict, Any

# Mainnet configuration
mainnet_node = "https://dydx-ops-rpc.kingnodes.com:443"
mainnet_chain = "dydx-mainnet-1"

# Staging configuration
staging_node = "https://validator.v4staging.dydx.exchange:443"
staging_chain = "dydxprotocol-testnet"

# Testnet configuration
testnet_node = "https://validator.v4testnet.dydx.exchange:443"
testnet_chain = "dydxprotocol-testnet"

PROPOSAL_STATUS_PASSED = 3

def vote_for(node, chain, proposal_id, person):
    print("voting as " + person)
    cmd = [
        "dydxprotocold",
        "tx",
        "gov",
        "vote",
        proposal_id,
        "yes",
        "--from=" + person,
        "--node=" + node,
        "--chain-id=" + chain,
        "--keyring-backend=test",
        "--fees=5000000000000000adv4tnt",
        "--yes"
    ]
    result = subprocess.run(cmd, capture_output=True, text=True)
    if result.returncode != 0:
        raise Exception(f"Failed to vote: {result.stderr}")

def load_yml(file_path) -> Dict[str, Any]:
    """
    Loads any yml file and returns the data as a dictionary.
    
    Args:
        file_path: Path to the yml file
        
    Returns:
        Dictionary containing the parsed data
    """
    try:
        with open(file_path, 'r', encoding='utf-8') as file:
            data = yaml.safe_load(file)
        return data
    except FileNotFoundError:
        print(f"Error: File '{file_path}' not found.")
        return {}
    except yaml.YAMLError as e:
        print(f"Error parsing YAML file: {e}")
        return {}

def get_proposal_id(node, chain):
    cmd = [
        "dydxprotocold",
        "query",
        "gov",
        "proposals",
        "--node=" + node,
        "--chain-id=" + chain
    ]
    with tempfile.NamedTemporaryFile(mode='w', suffix='.json') as tmp_file:
        subprocess.run(cmd, stdout=tmp_file)
        result = load_yml(tmp_file.name)
        return result['proposals'][-1]['id']

def main():
    parser = argparse.ArgumentParser(description='Parse market map and sync markets')
    parser.add_argument('--chain-id', default=staging_chain, help='Chain ID, default is dydxprotocol-testnet')
    parser.add_argument('--node', default=staging_node, help='Node URL, default is https://validator.v4staging.dydx.exchange:443')
    parser.add_argument('--max-30d-commission', type=int, required=True, help='Maximum 30d commission per referred')
    parser.add_argument('--referee-min-fee-tier', type=int, required=True, help='Referee minimum fee tier idx')
    parser.add_argument('--max-30d-revenue', type=int, required=True, help='Maximum 30d affiliate revenue per affiliate')
    args = parser.parse_args()

    counter = 0
    # 3 retries for the process.
    for i in range(3):
        try:
            with tempfile.NamedTemporaryFile(mode='w', suffix='.json', delete=False) as tmp_file:
                affiliate_parameters_msg = {
                    "messages": [
                        {
                            "@type": "/dydxprotocol.affiliates.MsgUpdateAffiliateParameters",
                            "authority": "dydx10d07y265gmmuvt4z0w9aw880jnsr700jnmapky",
                            "affiliate_parameters": {
                                "maximum_30d_commission_per_referred_quote_quantums": int(args.max_30d_commission),
                                "referee_minimum_fee_tier_idx": int(args.referee_min_fee_tier),
                                "maximum_30d_affiliate_revenue_per_affiliate_quote_quantums": int(args.max_30d_revenue),
                            }
			            }
                    ],
                    "deposit": "10000000000000000000000adv4tnt",
                    "metadata": "",
                    "title": "Update affiliate parameters",
                    "summary": f"Update affiliate parameters: max_30d_commission={args.max_30d_commission}, referee_min_fee_tier={args.referee_min_fee_tier}, max_30d_revenue={args.max_30d_revenue}"
                }
                json.dump(affiliate_parameters_msg, tmp_file, indent=2)
                print(affiliate_parameters_msg)
                tmp_file_path = tmp_file.name
            print("submitting proposal for affiliate parameters update")
            cmd = [
                "dydxprotocold",
                "tx",
                "gov",
                "submit-proposal",
                tmp_file_path,
                "--from=alice",
                "--gas=auto", 
                "--fees=10000000000000000000000adv4tnt",
                "--node=" + args.node,
                "--chain-id=" + args.chain_id,
                "--keyring-backend=test", 
                "--yes"
            ]
            
            result = subprocess.run(cmd, capture_output=True, text=True)
            if result.returncode != 0:
                raise Exception(f"Failed to submit proposal: {result.stderr}")
            # delete the temporary file
            os.remove(tmp_file_path)
            print("voting for affiliate parameters update")
            time.sleep(5)
            # vote for alice
            voters = ["alice", "bob", "carl", "dave", "emily", "fiona", "greg", "henry", "ian", "jeff"]
            proposal_id = get_proposal_id(args.node, args.chain_id)
            for voter in voters:
                vote_for(args.node, args.chain_id, proposal_id, voter)
                
            # wait for the proposal to pass
            print("Waiting 2 minutes for proposal to pass")
            time.sleep(120)
            # check if the proposal passed
            cmd = [
                "dydxprotocold",
                "query",
                "gov",
                "proposal",
                proposal_id,
                "--node=" + args.node,
                "--chain-id=" + args.chain_id
            ]
            with tempfile.NamedTemporaryFile(mode='w', suffix='.json') as tmp_file:
                subprocess.run(cmd, stdout=tmp_file)
                result = load_yml(tmp_file.name)
                if result['proposal']['status'] == PROPOSAL_STATUS_PASSED:
                    print("proposal passed, affiliate parameters updated")
                    return True
                else:
                    raise Exception("Failed to update affiliate parameters")
            break
        except Exception as e:
            print(e)    
            print(f"got exception, retrying {i+1} time(s)")


if __name__ == "__main__":
    main()
