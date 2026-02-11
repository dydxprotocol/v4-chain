#!/usr/bin/env python3
import json
import subprocess
import sys
import time

RPC = "http://localhost:26657"
CHAIN_ID = "localdydxprotocol"
KEYRING = "test"
MAKER = "alice"  # resting order signer
MAKER_ADDRESS = "dydx199tqg4wdlnu4qjlxchpd7seg454937hjrknju4"
TAKER = "bob"  # crossing order signer
TAKER_ADDRESS = "dydx10fx7sy6ywd5senxae9dwytf8jxek3t2gcen2vs"
CLOB_PAIR_ID_BTC = 0  # set to your BTC clob pair id
RESTING_PRICE = 30000_000000  # price in subticks (adjust to your tick size)
CROSS_PRICE   = 31000_000000  # higher price to cross
SIZE = 100_0000              # base quantums (adjust to your lot size)
SUBTICKS = RESTING_PRICE     # subticks == price for short-term orders

def run(cmd):
    out = subprocess.check_output(cmd, stderr=subprocess.STDOUT)
    return out.decode().strip()

def latest_height():
    raw = run(["dydxprotocold", "status", "--node", RPC])
    print(raw)
    return int(json.loads(raw)["sync_info"]["latest_block_height"])

def deposit_to_subaccount(owner_key, amount):
    cmd = [
        "/Users/kefan/go/src/github.com/dydxprotocol/v4-chain/protocol/build/dydxprotocold", "tx", "sending", "deposit-to-subaccount",
        run(["dydxprotocold", "keys", "show", owner_key, "-a", "--keyring-backend", KEYRING]),
        MAKER_ADDRESS if owner_key == MAKER else TAKER_ADDRESS,
        "0", str(amount),
        "--yes",
        "--broadcast-mode", "sync",
        "--from", owner_key,
        "--keyring-backend", KEYRING,
        "-o", "json",
        "--fees=5000ibc/8E27BA2D5493AF5636760E354E46004562C46AB7EC0CC4C1CA14E9E20E2545B5"
    ]
    print("Deposit to subaccount command:", " ".join(arg for arg in cmd if arg))
    resp = json.loads(run(cmd))


def place_order(owner_key, side, price, size, client_id, short=False):
    gtb = latest_height() + 40
    cmd = [
        "/Users/kefan/go/src/github.com/dydxprotocol/v4-chain/protocol/build/dydxprotocold", "tx", "clob", "place-order",
        run(["dydxprotocold", "keys", "show", owner_key, "-a", "--keyring-backend", KEYRING]),
        "0", str(client_id), str(CLOB_PAIR_ID_BTC), str(side),
        str(size), str(price), str(gtb),
        "--from", owner_key,
        "--keyring-backend", KEYRING,
        "--chain-id", CHAIN_ID,
        "--node", RPC,
        "--yes",
        "--broadcast-mode", "sync",
        "-o", "json",
        "--fees=5000ibc/8E27BA2D5493AF5636760E354E46004562C46AB7EC0CC4C1CA14E9E20E2545B5"
    ]
    if short: 
        cmd.append("--short=true")
    print("Place order command:", " ".join(arg for arg in cmd if arg))
    resp = json.loads(run(cmd))
    print("Place order response:", resp)
    return resp["txhash"], int(resp.get("height", 0))

def cancel_order(owner_key, client_id):
    gtb = latest_height() + 50
    cmd = [
        "/Users/kefan/go/src/github.com/dydxprotocol/v4-chain/protocol/build/dydxprotocold", "tx", "clob", "cancel-order",
    run(["dydxprotocold", "keys", "show", owner_key, "-a", "--keyring-backend", KEYRING]),
        "0", str(client_id), str(CLOB_PAIR_ID_BTC), str(gtb),
        "--from", owner_key,
        "--keyring-backend", KEYRING,
        "--chain-id", CHAIN_ID,
        "--node", RPC,
        "--yes",
        "--broadcast-mode", "sync",
        "-o", "json",
        "--fees=5000ibc/8E27BA2D5493AF5636760E354E46004562C46AB7EC0CC4C1CA14E9E20E2545B5"
    ]
    print("Cancel command:", " ".join(arg for arg in cmd if arg))
    resp = json.loads(run(cmd))
    print("Cancel order response:", resp)
    return resp["txhash"], int(resp.get("height", 0))

def main():
    deposit_to_subaccount(MAKER, 10000000000)
    deposit_to_subaccount(TAKER, 10000000000)
    # 1) Resting SELL from maker
    print("Placing resting SELL...")
    resting_hash, resting_h = place_order(MAKER, side=2, price=RESTING_PRICE, size=SIZE, client_id=1)
    print(f"Resting order tx: {resting_hash} at height {resting_h}")
    time.sleep(5)
    # 2) Crossing BUY from taker
    print("Placing crossing BUY...")
    cross_hash, cross_h = place_order(TAKER, side=1, price=CROSS_PRICE, size=SIZE, client_id=3, short=True)
    print(f"Crossing order tx: {cross_hash} at height {cross_h}")

    # 3) Immediate cancel for crossing order (same client_id/subaccount)
    print("Cancelling crossing BUY...")
    cancel_hash, cancel_h = cancel_order(MAKER, client_id=1)
    print(f"Cancel tx: {cancel_hash} at height {cancel_h}")

    print("\nExpectations:")
    print("- Cancel is reordered ahead of the crossing order in the block proposal.")
    print("- Resting order remains unmatched; no fill against the canceled order.")

if __name__ == "__main__":
    try:
        main()
    except subprocess.CalledProcessError as e:
        sys.stderr.write(e.output.decode())
        sys.exit(1)
