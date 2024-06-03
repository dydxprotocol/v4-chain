import requests
import sys
import time

node_url = sys.argv[1]

def fetch_unconfirmed_txs(node_url):
    endpoint = f"http://{node_url}/unconfirmed_txs"  # Ensure proper URL format
    try:
        response = requests.get(endpoint)
        print(f"Fetching from {endpoint}, Status Code: {response.status_code}")  # Debug print
        if response.status_code == 200:
            print(f"Response Content: {response.text}")  # Debug print
            return response.json()
        else:
            print(f"Failed to fetch data from {endpoint}, Status Code: {response.status_code}")
            return None
    except Exception as e:
        print(f"Error fetching data: {e}")
        return None

def decode_and_log_transactions(unconfirmed_txs_data):
    if "result" in unconfirmed_txs_data and "txs" in unconfirmed_txs_data["result"]:
        transactions = unconfirmed_txs_data["result"]["txs"]
        print(f"Number of unconfirmed transactions: {len(transactions)}")
        for tx in transactions:
            # Here, tx is a transaction string; we will print the transaction itself
            print(f"Transaction: {tx}")

def main():
    print("Welcome to the mempool logger ༼ つ ◕_◕ ༽つ!")
    while True:
        unconfirmed_txs_data = fetch_unconfirmed_txs(node_url)
        if unconfirmed_txs_data:
            decode_and_log_transactions(unconfirmed_txs_data)
        else:
            print("No data received")
        time.sleep(10)

if __name__ == "__main__":
    print("Script is starting")
    main()
