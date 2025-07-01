#!/usr/bin/env python3
"""
Script to parse the market-map.yml file and extract all markets into a Python dictionary.
"""

import argparse
import os
import shutil
import subprocess
import sys
import tempfile
import yaml
import json
from typing import Dict, Any, List
import time

alice_address = "dydx199tqg4wdlnu4qjlxchpd7seg454937hjrknju4"
bob_address = "dydx10fx7sy6ywd5senxae9dwytf8jxek3t2gcen2vs"
mainnet_node = "https://dydx-ops-rpc.kingnodes.com:443"
mainnet_chain = "dydx-mainnet-1"
PROPOSAL_STATUS_PASSED = 3

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

def make_alice_authority(node, chain):
    print("Making alice an authority")
    # Check first to see if alice is already an authority
    cmd = [
        "dydxprotocold",
        "query",
        "marketmap",
        "params",
        "--node=" + node,
        "--chain-id=" + chain
    ]
    with tempfile.NamedTemporaryFile(mode='w', suffix='.yml') as tmp_file:
        subprocess.run(cmd, stdout=tmp_file)
        result = load_yml(tmp_file.name)
        if alice_address in result['market_authorities']:
            # alice is authority, we're done
            print("alice is already authority")
            return True
    # alice is not an authority, we need to propose her as an authority
    # Create a temporary file to store the authority proposal
    with tempfile.NamedTemporaryFile(mode='w', suffix='.json', delete=False) as tmp_file:
        authority_proposal = {
            "title": "Make alice an authority for launching new markets",
            "deposit": "10000000000000000000000adv4tnt", 
            "summary": "Make alice an authority for launching new markets",
            "messages": [
                {
                    "@type": "/slinky.marketmap.v1.MsgParams",
                    "authority": "dydx10d07y265gmmuvt4z0w9aw880jnsr700jnmapky",
                    "params": {
                        "admin": "dydx10d07y265gmmuvt4z0w9aw880jnsr700jnmapky",
                        "market_authorities": [
                            "dydx10d07y265gmmuvt4z0w9aw880jnsr700jnmapky",
                            alice_address
                        ]
                    }
                }
            ]
        }
        json.dump(authority_proposal, tmp_file, indent=2)
        tmp_file_path = tmp_file.name
    print("submitting proposal for alice to be an authority")
    cmd = [
        "dydxprotocold",
        "tx",
        "gov",
        "submit-proposal",
        tmp_file_path,
        "--from=alice",
        "--gas=auto", 
        "--fees=10000000000000000000000adv4tnt",
        "--node=" + node,
        "--chain-id=" + chain,
        "--keyring-backend=test", 
        "--yes"
    ]
    
    result = subprocess.run(cmd, capture_output=True, text=True)
    if result.returncode != 0:
        raise Exception(f"Failed to submit proposal: {result.stderr}")
    # delete the temporary file
    os.remove(tmp_file_path)
    print("voting for alice to be an authority")
    time.sleep(5)
    # vote for alice
    voters = ["alice", "bob", "carl", "dave", "emily", "fiona", "greg", "henry", "ian", "jeff"]
    proposal_id = get_proposal_id(node, chain)
    for voter in voters:
        vote_for(node, chain, proposal_id, voter)
        
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
        "--node=" + node,
        "--chain-id=" + chain
    ]
    with tempfile.NamedTemporaryFile(mode='w', suffix='.json') as tmp_file:
        subprocess.run(cmd, stdout=tmp_file)
        result = load_yml(tmp_file.name)
        if result['proposal']['status'] == PROPOSAL_STATUS_PASSED:
            print("proposal passed, alice is now an authority")
            return True
        else:
            raise Exception("Failed to make alice an authority")

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

def sync_market_map(temp_dir: str, node, chain):
    # Query mainnet market map and testing env market map.
    cmd = [
        "dydxprotocold",
        "query",
        "marketmap",
        "market-map",
        "--node=" + mainnet_node,
        "--chain-id=" + mainnet_chain
    ]
    
    with open(f"{temp_dir}/mainnet-marketmap.yml", "w") as f:
        subprocess.run(cmd, stdout=f)
        
    cmd = [
        "dydxprotocold",
        "query",
        "marketmap", 
        "market-map",
        "--node=" + node,
        "--chain-id=" + chain
    ]
    
    with open(f"{temp_dir}/staging-marketmap.yml", "w") as f:
        subprocess.run(cmd, stdout=f)  

def markets_to_add(temp_dir: str):
    print(f"temp_dir: {temp_dir}")
    mainnet_marketmap = load_yml(f"{temp_dir}/mainnet-marketmap.yml")
    mainnet_marketmap = mainnet_marketmap['market_map']['markets']
    testnet_marketmap = load_yml(f"{temp_dir}/staging-marketmap.yml")
    testnet_marketmap = testnet_marketmap['market_map']['markets']
    toAdd = {}
    for market_name, market_data in mainnet_marketmap.items():
        if market_name not in testnet_marketmap and market_data["ticker"]["metadata_JSON"] != "":
            toAdd[market_name] = {
                "marketmap": market_data
            }

    return toAdd

def main():
    parser = argparse.ArgumentParser(description='Parse market map and sync markets')
    parser.add_argument('--chain-id', default='dydxprotocol-testnet', help='Chain ID, default is dydxprotocol-testnet')
    parser.add_argument('--node', default='https://validator.v4staging.dydx.exchange:443', help='Node URL, default is https://validator.v4staging.dydx.exchange:443')
    parser.add_argument('--number-markets', type=int, required=True, help='Number of markets to add')
    args = parser.parse_args()
    # alice needs to be an authority to launch markets into market map.
    make_alice_authority(args.node, args.chain_id)
    # determine what markets we need
    temp_dir = tempfile.mkdtemp()
    try: 
        sync_market_map(temp_dir, args.node, args.chain_id)
        toAdd = markets_to_add(temp_dir)
        counter = 0
        for name, market_data in toAdd.items():
            if market_data['marketmap'] is not None:
                # 3 retries for each market.
                for i in range(3):
                    try:
                        # Execute the marketmap create-markets command
                        print(f"Adding {name} to market map")
                        cmd = [
                            "dydxprotocold",
                            "tx",
                            "marketmap",
                            "create-markets",
                            f"--create-markets={json.dumps(market_data['marketmap'])}",
                            "--node=https://validator.v4staging.dydx.exchange:443",
                            "--chain-id=dydxprotocol-testnet", 
                            "--from=alice",
                            "--keyring-backend=test",
                            "--fees=5000000000000000adv4tnt", 
                            "--yes"
                        ]
                        
                        result = subprocess.run(cmd, capture_output=True, text=True)
                        if result.returncode != 0:
                            raise Exception(f"Failed to add {name} to market map")
                        print("stdout:", result.stdout)
                        print("stderr:", result.stderr)
                        print("returncode:", result.returncode)
                        print("waiting 5 seconds for market to be created")
                        time.sleep(5)
                        print("market created, launching market")
                        # no use checking the error log here, it will not fail here.
                        print("Launching market: ", market_data['marketmap']['ticker']['currency_pair']['Base'] + '-' + market_data['marketmap']['ticker']['currency_pair']['Quote'])
                        
                        cmd = [
                            "dydxprotocold",
                            "tx", 
                            "listing",
                            "create-market",
                            market_data['marketmap']['ticker']['currency_pair']['Base'] + '-' + market_data['marketmap']['ticker']['currency_pair']['Quote'],
                            alice_address,
                            "--from=alice", # bob has a decent amount of native tokens in subaccount 0 for gas. 
                            "--node=https://validator.v4staging.dydx.exchange:443",
                            "--chain-id=dydxprotocol-testnet",
                            "--keyring-backend=test",
                            "--fees=11270250000000000adv4tnt",
                            "--gas=auto",
                            "--yes"
                        ]
                        
                        result = subprocess.run(cmd, capture_output=True, text=True)
                        if result.returncode != 0:
                            raise Exception(f"Failed to launch {name}")
                        print("stdout:", result.stdout)
                        print("stderr:", result.stderr)
                        print("returncode:", result.returncode)
                        time.sleep(2)
                        counter += 1
                        break
                    except Exception as e:
                        print(e)    
                        print(f"got exception, retrying {i+1} time(s)")
            if counter >= args.number_markets:
                break
    finally:
        shutil.rmtree(temp_dir)

if __name__ == "__main__":
    main()
