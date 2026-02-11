"""
Localnet harness to sanity-check cancel-first behavior.

Assumptions:
  - You already have signed txs (base64) for the scenarios you want to test.
    Provide them via env vars (see ENV section below) or edit SCENARIOS.
  - The node is running locally and reachable via --node (default http://localhost:26657).
  - The binary `dydxprotocold` is on PATH.

What it does:
  - Broadcasts provided tx sets (e.g., place + cancel) to localnet.
  - Waits for inclusion, fetches the block, verifies cancels are ordered
    before places in the block tx list, and prints the MsgProposedOperations.

Environment variables:
  NODE            RPC, default http://localhost:26657
  CHAIN_ID        e.g., localdydxprotocol
  PLACE_TX_B64    signed place tx (base64)
  CANCEL_TX_B64   signed cancel tx (base64)
  BATCH_CANCEL_TX_B64 optional signed MsgBatchCancel tx (base64)
  OTHER_TX_B64    optional non-clob tx (base64)
"""

import base64
import json
import os
import subprocess
import sys
import tempfile
from dataclasses import dataclass
from typing import List, Optional


def run(cmd: List[str]) -> str:
    out = subprocess.check_output(cmd, stderr=subprocess.STDOUT)
    return out.decode()


def broadcast_tx(tx_b64: str, node: str, chain_id: str) -> dict:
    if tx_b64 is None:
        raise ValueError("tx_b64 is None")
    # Ensure we have a trimmed string (env vars may inject whitespace/newlines)
    if not isinstance(tx_b64, str):
        tx_b64 = str(tx_b64)
    tx_b64 = tx_b64.strip()

    with tempfile.NamedTemporaryFile(delete=False) as f:
        f.write(base64.b64decode(tx_b64.encode("ascii")))
        f.flush()
        path = f.name
    try:
        raw = run(
            [
                "dydxprotocold",
                "tx",
                "broadcast",
                path,
                "--node",
                node,
                "--chain-id",
                chain_id,
                "--broadcast-mode",
                "block",
                "-y",
                "-o",
                "json",
            ]
        )
        return json.loads(raw)
    finally:
        os.remove(path)


def get_block(height: int, node: str) -> dict:
    raw = run(["dydxprotocold", "q", "block", str(height), "--node", node, "-o", "json"])
    return json.loads(raw)


def decode_tx(tx_b64: str) -> dict:
    with tempfile.NamedTemporaryFile(delete=False) as f:
        f.write(base64.b64decode(tx_b64))
        f.flush()
        path = f.name
    try:
        raw = run(["dydxprotocold", "tx", "decode", path, "-o", "json"])
        return json.loads(raw)
    finally:
        os.remove(path)


def find_ops_tx(block: dict) -> Optional[str]:
    txs = block.get("block", {}).get("data", {}).get("txs", [])
    if len(txs) < 1:
        return None
    # The ops tx is typically the first (after price/funding/bridge ordering may apply).
    # We’ll just return the first tx for manual inspection.
    return txs[0]


@dataclass
class Scenario:
    name: str
    txs: List[str]
    expect_cancel_first: bool = True


def verify_order(block: dict, scenario: Scenario) -> bool:
    txs = block.get("block", {}).get("data", {}).get("txs", [])
    if not txs:
        print(f"[{scenario.name}] no txs in block")
        return False
    # decode all txs to get type info
    decoded = [decode_tx(t) for t in txs]
    types = []
    for d in decoded:
        msgs = d.get("body", {}).get("messages", [])
        msg_types = [m.get("@type", "") for m in msgs]
        types.append(msg_types)

    # find first cancel and first place
    first_cancel = None
    first_place = None
    for idx, msg_types in enumerate(types):
        if first_cancel is None and any("MsgCancel" in t or "MsgBatchCancel" in t for t in msg_types):
            first_cancel = idx
        if first_place is None and any("MsgPlaceOrder" in t for t in msg_types):
            first_place = idx
    if first_cancel is None or first_place is None:
        print(f"[{scenario.name}] could not find both cancel and place in block txs")
        return False
    ok = first_cancel <= first_place if scenario.expect_cancel_first else True
    print(f"[{scenario.name}] cancel idx={first_cancel}, place idx={first_place}, ok={ok}")
    return ok


def main() -> None:
    node = os.environ.get("NODE", "http://localhost:26657")
    chain_id = os.environ.get("CHAIN_ID", "localdydxprotocol")

    place = 1
    cancel = 1
    batch_cancel = os.environ.get("BATCH_CANCEL_TX_B64")
    other = os.environ.get("OTHER_TX_B64")

    scenarios: List[Scenario] = []
    if place and cancel:
        scenarios.append(Scenario("place+cancel", [place, cancel], expect_cancel_first=True))
        scenarios.append(Scenario("cancel+place", [cancel, place], expect_cancel_first=True))
    if batch_cancel and place:
        scenarios.append(Scenario("batch-cancel+place", [batch_cancel, place], expect_cancel_first=True))
    if other and cancel:
        scenarios.append(Scenario("cancel+other", [cancel, other], expect_cancel_first=True))

    if not scenarios:
        print("No scenarios configured. Set PLACE_TX_B64/CANCEL_TX_B64/etc.")
        sys.exit(1)

    for scenario in scenarios:
        print(f"\n=== Running scenario: {scenario.name} ===")
        # broadcast all txs in given order; PrepareProposal should reorder cancels anyway
        heights = []
        for tx_b64 in scenario.txs:
            resp = broadcast_tx(tx_b64, node=node, chain_id=chain_id)
            height = int(resp.get("height", 0) or resp.get("tx_response", {}).get("height", 0))
            heights.append(height)
            print(f"  broadcast height={height}, code={resp.get('code') or resp.get('tx_response', {}).get('code')}")

        height = max(heights)
        block = get_block(height, node=node)
        verify_order(block, scenario)

        ops_b64 = find_ops_tx(block)
        if ops_b64:
            ops_decoded = decode_tx(ops_b64)
            print("  Ops tx messages:")
            print(json.dumps(ops_decoded.get("body", {}).get("messages", []), indent=2))
        else:
            print("  No ops tx found")


if __name__ == "__main__":
    main()
