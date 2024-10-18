import argparse
import requests
import subprocess

"""
Instructions:

1. Install `requests` library

    `pip3 install requests`

2. Build `dydxprotocold` by running `make build` in the root directory of the protocol repository

    `cd v4-chain/protocol && make clean && make build`

3. Run the script using the correct endpoint and binary path. Example:

    `python3 scripts/markets/get_isolated_market_insurance_fund.py \
        --endpoint=https://dydx-ops-rest.kingnodes.com \
        --binary_path=/Users/taehoonlee/v4-chain/protocol/build/dydxprotocold`
"""

def get_id_to_pair(base_endpoint_url):
    endpoint_url = base_endpoint_url + "/dydxprotocol/prices/params/market"

    id_to_pair = {}

    # Pagination key
    pagination_key = None

    while True:
        # Build the URL with pagination key if it exists
        url = f"{endpoint_url}"
        if pagination_key:
            url += f"?pagination.key={pagination_key}"

        try:
            # Make the GET request
            response = requests.get(url)
            response.raise_for_status()  # Raise an exception for HTTP errors

            # Parse the JSON response
            data = response.json()

            # Extract
            for market in data.get('market_params', []):
                id_to_pair[market.get('id')] = market.get('pair')
        
            # Check if there is a next pagination key
            pagination_key = data.get('pagination', {}).get('next_key')

            # If no more pages, break the loop
            if not pagination_key:
                break

        except requests.RequestException as e:
            print(f"Request failed: {e}")
            break

    return id_to_pair


def get_isolated_market_ids(base_endpoint_url):
    endpoint_url = base_endpoint_url + "/dydxprotocol/perpetuals/perpetual"

    # Initialize the result list for matching market_ids
    matching_market_ids = []

    # Pagination key
    pagination_key = None

    while True:
        # Build the URL with pagination key if it exists
        url = f"{endpoint_url}"
        if pagination_key:
            url += f"?pagination.key={pagination_key}"

        try:
            # Make the GET request
            response = requests.get(url)
            response.raise_for_status()  # Raise an exception for HTTP errors

            # Parse the JSON response
            data = response.json()

            # Extract and filter market_ids where market_type matches the criteria
            for market in data.get('perpetual', []):
                market_params = market.get('params')
                if market_params.get('market_type') == 'PERPETUAL_MARKET_TYPE_ISOLATED':
                    matching_market_ids.append(market_params.get('market_id'))

            # Check if there is a next pagination key
            pagination_key = data.get('pagination', {}).get('next_key')

            # If no more pages, break the loop
            if not pagination_key:
                break

        except requests.RequestException as e:
            print(f"Request failed: {e}")
            break

    return matching_market_ids

def run_dydxprotocold(command):
    try:
        # Run the command and capture the output
        result = subprocess.run(command, capture_output=True, text=True, check=True)
        return result.stdout

    except subprocess.CalledProcessError as e:
        # Handle errors during command execution
        print(f"Error running dydxprotocold \n{e.stderr}")

    except FileNotFoundError:
        # Handle the case where the binary is not found
        print("Error: dydxprotocold binary not found. Ensure you ran `make build`.")

def get_insurance_fund_address_for_markets(binary_path, market_ids):
    """
    This function takes a list of market_ids and runs the `dydxprotocold` binary for each.
    It prints the output and errors, if any.
    """

    market_id_to_address = []

    for market_id in market_ids:
        # Construct the command
        command = [binary_path, "q", "module-name-to-address", "insurance_fund:" + str(market_id)]
        address = run_dydxprotocold(command)
        if not address:
            raise Exception(f"Failed to get insurance fund address for market_id: {market_id}")
        market_id_to_address.append((market_id, address))

    return market_id_to_address

def get_bank_balance(market_id_to_address, base_endpoint_url):
    url = base_endpoint_url + "/cosmos/bank/v1beta1/balances/"

    address_to_balance = {}

    for _, address in market_id_to_address:
        try:
            # Make the GET request
            response = requests.get(url + address)
            response.raise_for_status()  # Raise an exception for HTTP errors

            # Parse the JSON response
            data = response.json()
            balances = data.get('balances', [])
            bank_balance = 0

            if len(balances) > 1:
                raise ValueError("Expected at most one balance entry.")
            if len(balances) == 1:
                if balances[0].get('denom') != 'ibc/8E27BA2D5493AF5636760E354E46004562C46AB7EC0CC4C1CA14E9E20E2545B5':
                    raise ValueError("Unexpected 'denom' value in balances.")
                bank_balance = int(balances[0].get('amount'))
                bank_balance = bank_balance / (10 ** 6)
            address_to_balance[address] = bank_balance

        except requests.RequestException as e:
            print(f"Request failed: {e}")
            break

    return address_to_balance


def print_market_info(market_id_to_address, address_to_balance, id_to_pair):
    """
    Prints market information with aligned columns for better readability.
    """
    header = f"{'Market':<15} {'ID':<5} {'Insurance Addr':<58} {'Bank Balance':>15}"
    print(header)
    print("-" * len(header))  # Separator line

    for market_id, address in market_id_to_address:
        pair = id_to_pair.get(market_id)
        balance = address_to_balance.get(address)

        print(
            f"{pair:<15} "
            f"{market_id:<5} "
            f"{address:<58} "
            f"{balance:>15.6f}"
        )


if __name__ == "__main__":
    parser = argparse.ArgumentParser(description='Get insurance fund balances for all isolated markets.')
    parser.add_argument('--endpoint', required=True, help='The endpoint URL to fetch the markets.')
    parser.add_argument('--binary_path', required=True, help='The local path to the `dydxprotocold` binary.')
    args = parser.parse_args()
    endpoint = args.endpoint
    binary_path = args.binary_path

    id_to_pair = get_id_to_pair(endpoint)
    market_ids = get_isolated_market_ids(endpoint)
    market_id_to_address = get_insurance_fund_address_for_markets(binary_path, market_ids)
    address_to_balance = get_bank_balance(market_id_to_address, endpoint)
    print_market_info(market_id_to_address, address_to_balance, id_to_pair)
