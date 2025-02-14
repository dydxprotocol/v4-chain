import json
import requests


MAINNET_URL_BASE = "https://dydx-ops-rest.kingnodes.com"
MARKETMAP_URL = f"{MAINNET_URL_BASE}/slinky/marketmap/v1/marketmap"
PERPETUALS_URL = f"{MAINNET_URL_BASE}/dydxprotocol/perpetuals/perpetual?pagination.limit=10000"
CLOB_URL = f"{MAINNET_URL_BASE}/dydxprotocol/clob/clob_pair?pagination.limit=10000"

DELISTED_CLOB_STATUS = "STATUS_FINAL_SETTLEMENT"

TICKERS_TO_DELIST = [] # should be in BASE-QUOTE format (ex. PAIN-USD)

def main():
    marketmap_data = requests.get(MARKETMAP_URL).json()
    perpetuals_data = requests.get(PERPETUALS_URL).json()
    clob_data = requests.get(CLOB_URL).json()

    proposal = {
        "title": "Delist " + ", ".join(TICKERS_TO_DELIST) + " and disable them in the marketmap",
        "deposit":"2000000000000000000000adydx"
    }
    proposal["summary"] = proposal["title"]
    proposal_messages = []
    update_markets_message = [] # for disabling in marketmap

    ticker_to_perpetual_id = get_ticker_to_perpetual_id(perpetuals_data)
    perpetual_id_to_clob_pair = get_perpetual_id_to_clob_pair(clob_data)
    for ticker in TICKERS_TO_DELIST:
        # Create clob delisting proposal message
        clob_delisting_message = {}
        clob_delisting_message["@type"] = "/dydxprotocol.clob.MsgUpdateClobPair"
        clob_delisting_message["authority"] = "dydx10d07y265gmmuvt4z0w9aw880jnsr700jnmapky"

        perpetual_id = ticker_to_perpetual_id[ticker]
        clob_pair = perpetual_id_to_clob_pair[perpetual_id]
        clob_pair["status"] = DELISTED_CLOB_STATUS

        clob_delisting_message["clob_pair"] = clob_pair
        proposal_messages.append(clob_delisting_message)

        # Append to marketmap disable message
        marketmap_ticker = "/".join(ticker.split("-"))
        market = marketmap_data["market_map"]["markets"][marketmap_ticker]
        market["ticker"]["enabled"] = False
        update_markets_message.append(market)

    # Create marketmap disable proposal message
    marketmap_disable_message = {}
    marketmap_disable_message["@type"] = "/slinky.marketmap.v1.MsgUpdateMarkets"
    marketmap_disable_message["authority"] = "dydx10d07y265gmmuvt4z0w9aw880jnsr700jnmapky"
    marketmap_disable_message["update_markets"] = update_markets_message
    proposal_messages.append(marketmap_disable_message)

    proposal["messages"] = proposal_messages
    print(json.dumps(proposal, indent=4))

def get_ticker_to_perpetual_id(perpetuals_data):
    ticker_to_perpetual_id = {}
    for data in perpetuals_data["perpetual"]:
        ticker = data["params"]["ticker"]
        perpetual_id = data["params"]["id"]
        ticker_to_perpetual_id[ticker] = perpetual_id
    return ticker_to_perpetual_id

def get_perpetual_id_to_clob_pair(clob_data):
    perpetual_id_to_clob_pair = {}
    for data in clob_data["clob_pair"]:
        perpetual_id = data["perpetual_clob_metadata"]["perpetual_id"]
        perpetual_id_to_clob_pair[perpetual_id] = data
    return perpetual_id_to_clob_pair

if __name__ == '__main__':
    main()