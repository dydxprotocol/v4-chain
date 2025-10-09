# /// script
# requires-python = ">=3.10"
# dependencies = [
#   "requests>=2.31.0",
#   "websockets>=12.0",
# ]
# ///
"""
Async script: fetch markets metadata, then subscribe to every market on the chosen channel.
If a limit error is returned, log how many subscriptions succeeded and exit.
"""

import asyncio
import json
import logging
import re
import sys

import requests
import websockets

MARKETS_URL = "https://indexer.v4mainnet.dydx.exchange/v4/perpetualMarkets"
WS_URL = "wss://indexer.v4mainnet.dydx.exchange/v4/ws"

logging.basicConfig(
    level=logging.INFO,
    format="%(asctime)s %(levelname)s %(message)s",
    datefmt="%H:%M:%S",
)
log = logging.getLogger("dydx-limit-test")


def fetch_tickers(markets_url: str) -> list[str]:
    """Fetch markets and return list of ticker strings.
    Handles pure JSON or HTML-with-<pre> wrapper.
    """
    resp = requests.get(markets_url, timeout=20)
    try:
        data = resp.json()
    except ValueError:
        # Fallback: extract JSON from <pre>...</pre> if server wrapped it
        m = re.search(r"<pre>(\{.*\})</pre>", resp.text, flags=re.S | re.I)
        if not m:
            raise
        data = json.loads(m.group(1))
    markets = data.get("markets") or {}
    # ticker is the id used for subscriptions
    tickers = [str(v["ticker"]) for v in markets.values() if "ticker" in v]
    # Deduplicate and sort for stable order
    tickers = sorted(set(tickers))
    if not tickers:
        raise RuntimeError("No tickers found from metadata")
    return tickers


async def subscribe_until_limit(
    ws_url: str,
    channel: str,
    tickers: list[str],
    number: int | None,
    subscribe_delay: float,
    gather_delay: float,
) -> None:
    """Subscribe to each ticker on the given channel until a limit error occurs.
    On limit error: log count and exit(0).
    """
    stop_event = asyncio.Event()
    successful: set[str] = set()

    async with websockets.connect(ws_url, max_queue=None, ping_interval=20) as ws:
        # Receiver: track success and detect limit errors
        async def recv_loop():
            try:
                while not stop_event.is_set():
                    raw = await ws.recv()
                    msg = json.loads(raw)
                    mtype = msg.get("type", "")
                    if mtype == "error":
                        # Detect any limit-related error
                        # Match common fields/messages without relying on exact code
                        text = json.dumps(msg).lower()
                        if "limit" in text or "too many" in text:
                            log.error("Limit error received: %s", msg)
                            stop_event.set()
                            return
                        else:
                            log.warning("Non-limit error: %s", msg)
                    else:
                        # Count a subscription as successful when we first see either a 'subscribed'
                        # ack or any first data message for that channel+id.
                        if msg.get("channel") == channel:
                            msg_id = msg.get("id")
                            if msg_id is not None:
                                mid = str(msg_id)
                                if mid not in successful:
                                    # Accept either explicit ack ('subscribed') or first data
                                    successful.add(mid)
            except websockets.ConnectionClosedOK:
                pass
            except websockets.ConnectionClosedError as e:
                log.error("WebSocket closed with error: %s", e)
                stop_event.set()

        # Sender: sequentially subscribe to each market id until a limit error is hit
        async def send_loop():
            try:
                for ticker in tickers[slice(0, number)]:
                    if stop_event.is_set():
                        break
                    sub = {"type": "subscribe", "channel": channel, "id": ticker}
                    await ws.send(json.dumps(sub))
                    # Gentle pacing to avoid instant flood; adjust if needed
                    await asyncio.sleep(subscribe_delay)
            finally:
                # Allow receiver to drain any final messages
                await asyncio.sleep(gather_delay)
                stop_event.set()

        recv_task = asyncio.create_task(recv_loop())
        send_task = asyncio.create_task(send_loop())

        await stop_event.wait()
        # Cancel outstanding tasks and close
        recv_task.cancel()
        send_task.cancel()
        # Report result and exit if this stop was caused by a limit error
        count = len(successful)
        log.info(
            "Channel '%s' subscriptions succeeded before limit: %d", channel, count
        )
        # Exit immediately per requirement
        sys.exit(0)


def parse_args():
    import argparse

    parser = argparse.ArgumentParser(description="dYdX v4 WS limit tester")
    parser.add_argument(
        "--channel",
        default="v4_orderbook",
        choices=["v4_orderbook", "v4_trades"],
        help="Channel to subscribe on (default: v4_orderbook)",
    )
    parser.add_argument(
        "--number",
        default=None,
        type=int,
        help="Number of markets to subscribe to (default: all)",
    )
    parser.add_argument(
        "--markets-url",
        default=MARKETS_URL,
        help=f"URL to fetch markets from (default: {MARKETS_URL})",
    )
    parser.add_argument(
        "--ws-url",
        default=WS_URL,
        help=f"Endpoint to connect to (default: {WS_URL})",
    )
    parser.add_argument(
        "--subscribe-delay",
        type=float,
        default=0.1,
        help="Delay between subscription requests (s) (default: 0.1s)",
    )
    parser.add_argument(
        "--gather-delay",
        type=float,
        default=1,
        help="Wait following subscription requests to collect responses (s) (default: 1s)",
    )
    return parser.parse_args()


def main():
    args = parse_args()
    tickers = fetch_tickers(args.markets_url)
    log.info("Fetched %d tickers", len(tickers))
    log.info(
        f"Subscribing to {len(tickers) if args.number is None else args.number} markets",
    )
    asyncio.run(
        subscribe_until_limit(
            args.ws_url,
            args.channel,
            tickers,
            args.number,
            args.subscribe_delay,
            args.gather_delay,
        )
    )


if __name__ == "__main__":
    main()
