import base64
import hashlib
import os
import sys

def main() -> None:
    """
    Read a base64-encoded tx from argv[1], env TX, or stdin, and print its
    SHA-256 hash in uppercase hex (same as shasum -a 256 | awk '{print toupper($1)}').
    """
    tx_b64 = sys.argv[1] if len(sys.argv) > 1 else os.environ.get("TX")
    if not tx_b64:
        tx_b64 = sys.stdin.read()
    tx_b64 = tx_b64.strip()
    if not tx_b64:
        raise SystemExit("no tx input provided")

    tx_bytes = base64.b64decode(tx_b64)
    digest = hashlib.sha256(tx_bytes).hexdigest().upper()

    os.system(f"dydxprotocold q tx {digest} -o json")


if __name__ == "__main__":
    main()