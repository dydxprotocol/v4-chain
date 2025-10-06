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
import toml
from pathlib import Path
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

def load_client_config() -> Dict[str, Any]:
    """
    Loads configuration from ~/.dydxprotocol/config/client.toml if it exists.
    
    Returns:
        Dictionary containing chain-id and node from client.toml, or empty dict if not found
    """
    config_path = Path.home() / ".dydxprotocol" / "config" / "client.toml"
    if config_path.exists():
        try:
            with open(config_path, 'r') as f:
                config = toml.load(f)
                return config
        except Exception as e:
            print(f"Warning: Could not load client.toml: {e}")
            return {}
    return {}

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
    # Load configuration from client.toml if available
    client_config = load_client_config()
    default_chain_id = client_config.get('chain-id', staging_chain)
    default_node = client_config.get('node', staging_node)
    
    parser = argparse.ArgumentParser(description='Update affiliate parameters')
    parser.add_argument('--chain-id', default=default_chain_id, help=f'Chain ID, default from client.toml or {staging_chain}')
    parser.add_argument('--node', default=default_node, help=f'Node URL, default from client.toml or {staging_node}')
    parser.add_argument('--max-30d-commission', type=int, required=True, help='Maximum 30d commission per referred')
    parser.add_argument('--referee-min-fee-tier', type=int, required=True, help='Referee minimum fee tier idx')
    parser.add_argument('--max-30d-volume', type=int, required=True, help='Maximum 30d attributable volume per referred')
    args = parser.parse_args()
    
    # Print configuration source
    if client_config:
        print(f"Loaded configuration from ~/.dydxprotocol/config/client.toml")
    print(f"Using chain-id: {args.chain_id}")
    print(f"Using node: {args.node}")

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
                                "maximum_30d_attributable_volume_per_referred_user_notional": int(args.max_30d_volume),
                                "referee_minimum_fee_tier_idx": int(args.referee_min_fee_tier),
                                "maximum_30d_attributable_revenue_per_referred_user_quote_quantums": int(args.max_30d_commission),
                            }
			            }
                    ],
                    "deposit": "10000000000000000000000adv4tnt",
                    "metadata": "",
                    "title": "Update affiliate parameters",
                    "summary": f"Update affiliate parameters: max_30d_commission={args.max_30d_commission}, referee_min_fee_tier={args.referee_min_fee_tier}, max_30d_volume={args.max_30d_volume}"
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
            
            # Print the full command
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
                "/Users/justinbarnett/projects/v4-chain/protocol/build/dydxprotocold",
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
