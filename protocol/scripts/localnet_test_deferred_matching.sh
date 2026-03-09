#!/bin/bash
# Test script for deferred matching on localnet.
#
# Prerequisites:
#   1. Localnet is running: `make localnet-start` (from protocol/)
#   2. Keys are imported (run the setup section below first)
#   3. Wait for ~2 blocks after localnet starts before running tests
#
# Usage:
#   ./scripts/localnet_test_deferred_matching.sh setup-keys
#   ./scripts/localnet_test_deferred_matching.sh test-match
#   ./scripts/localnet_test_deferred_matching.sh test-cancel-before-match
#   ./scripts/localnet_test_deferred_matching.sh query-height

set -euo pipefail

BINARY="./build/dydxprotocold"
CHAIN_ID="localdydxprotocol"
NODE="tcp://localhost:26657"

ALICE_ADDR="dydx199tqg4wdlnu4qjlxchpd7seg454937hjrknju4"
BOB_ADDR="dydx10fx7sy6ywd5senxae9dwytf8jxek3t2gcen2vs"
USDC_DENOM="ibc/8E27BA2D5493AF5636760E354E46004562C46AB7EC0CC4C1CA14E9E20E2545B5"
TX_FEES="--fees 25000$USDC_DENOM"

# place-order args: owner subaccount_number clientId clobPairId side quantums subticks goodTilBlock
# side: 1 = BUY, 2 = SELL
# cancel-order args: owner subaccount_number clientId clobPairId goodTilBlock

get_height() {
    curl -s "localhost:26657/status" | jq -r '.result.sync_info.latest_block_height'
}

wait_for_block() {
    local target=$1
    echo "Waiting for block $target..."
    while true; do
        current=$(get_height)
        if [ "$current" -ge "$target" ] 2>/dev/null; then
            echo "Reached block $current"
            return
        fi
        sleep 1
    done
}

setup_keys() {
    echo "=== Setting up keys ==="
    echo "Importing alice..."
    echo "merge panther lobster crazy road hollow amused security before critic about cliff exhibit cause coyote talent happy where lion river tobacco option coconut small" | \
        $BINARY keys add alice --recover --keyring-backend=test 2>/dev/null || echo "alice key may already exist"

    echo "Importing bob..."
    echo "color habit donor nurse dinosaur stable wonder process post perfect raven gold census inside worth inquiry mammal panic olive toss shadow strong name drum" | \
        $BINARY keys add bob --recover --keyring-backend=test 2>/dev/null || echo "bob key may already exist"

    echo ""
    echo "Keys imported. Verify:"
    $BINARY keys list --keyring-backend=test

    echo ""
    echo "=== Depositing USDC to subaccounts ==="
    # Deposit 1M USDC (1_000_000 * 1e6 = 1_000_000_000_000 quantums) to each subaccount 0.
    DEPOSIT_QUANTUMS=1000000000000

    echo "Depositing to alice's subaccount 0..."
    $BINARY tx sending deposit-to-subaccount dydx199tqg4wdlnu4qjlxchpd7seg454937hjrknju4 dydx199tqg4wdlnu4qjlxchpd7seg454937hjrknju4 0 1000000000000 \
        --from alice --keyring-backend=test \
        -y --broadcast-mode sync --fees 5000ibc/8E27BA2D5493AF5636760E354E46004562C46AB7EC0CC4C1CA14E9E20E2545B5 2>&1 | grep -E "code:|txhash:" || true

    echo "Waiting for alice's deposit to land..."
    sleep 5

    echo "Depositing to bob's subaccount 0..."
    $BINARY tx sending deposit-to-subaccount "$BOB_ADDR" "$BOB_ADDR" 0 $DEPOSIT_QUANTUMS \
        --from bob --chain-id $CHAIN_ID --node $NODE --keyring-backend=test \
        -y --broadcast-mode sync --fees 5000ibc/8E27BA2D5493AF5636760E354E46004562C46AB7EC0CC4C1CA14E9E20E2545B5 2>&1 | grep -E "code:|txhash:" || true

    echo "Waiting for bob's deposit to land..."
    sleep 5

    echo ""
    echo "Verifying subaccount balances:"
    echo "Alice subaccount 0:"
    curl -s "localhost:1317/dydxprotocol/subaccounts/subaccount/dydx199tqg4wdlnu4qjlxchpd7seg454937hjrknju4/0" | jq '.subaccount.asset_positions' 2>/dev/null || echo "  (query failed)"
    echo "Bob subaccount 0:"
    curl -s "localhost:1317/dydxprotocol/subaccounts/subaccount/dydx10fx7sy6ywd5senxae9dwytf8jxek3t2gcen2vs/0" | jq '.subaccount.asset_positions' 2>/dev/null || echo "  (query failed)"
}

test_basic_match() {
    echo "=== Test: Basic Order Matching ==="
    echo "This test places a buy from alice and a crossing sell from bob."
    echo "With deferred matching, both should go to the book during CheckTx,"
    echo "then match during PrepareProposal of the next block."
    echo ""

    HEIGHT=$(get_height)
    GTB=$((HEIGHT + 20))
    echo "Current block height: $HEIGHT"
    echo "Using GoodTilBlock: $GTB"
    echo ""

    # Alice: BUY 10 quantums of clobPair 0 at subticks 10000
    echo "1. Alice places BUY order (clientId=100)..."
    $BINARY tx clob place-order "$ALICE_ADDR" 0 100 0 1 1000000 1000000 $GTB \
        --from alice --chain-id $CHAIN_ID --node $NODE --keyring-backend=test \
        -y --broadcast-mode sync $TX_FEES 2>&1 | grep -E "code:|codespace:|raw_log:|txhash:" || true
    echo ""

    sleep 2

    # Bob: SELL 10 quantums of clobPair 0 at subticks 10000 (crosses alice's buy)
    echo "2. Bob places SELL order (clientId=200) - should cross alice's buy..."
    $BINARY tx clob place-order "$BOB_ADDR" 0 200 0 2 1000000 1000000 $GTB \
        --from bob --chain-id $CHAIN_ID --node $NODE --keyring-backend=test \
        -y --broadcast-mode sync $TX_FEES 2>&1 | grep -E "code:|codespace:|raw_log:|txhash:" || true
    echo ""

    echo "Orders placed. Wait for 1-2 blocks, then check the subaccount fills:"
    echo "  curl -s 'localhost:1317/dydxprotocol/subaccounts/subaccount/dydx199tqg4wdlnu4qjlxchpd7seg454937hjrknju4/0' | jq"
    echo "  curl -s 'localhost:1317/dydxprotocol/subaccounts/subaccount/dydx10fx7sy6ywd5senxae9dwytf8jxek3t2gcen2vs/0' | jq"
    echo ""
    echo "If the match occurred, both subaccounts will show perpetual positions."
}

test_cancel_before_match() {
    echo "=== Test: Cancel Before Match (Deferred Matching) ==="
    echo "This test places a buy from alice, then bob places a crossing sell"
    echo "and IMMEDIATELY cancels it in the same block."
    echo ""
    echo "With deferred matching: bob's cancel removes his order before matching,"
    echo "  so NO match occurs and alice's buy remains on the book."
    echo "With eager matching (old): bob's sell would match alice's buy instantly"
    echo "  during CheckTx, and the cancel would be too late."
    echo ""

    HEIGHT=$(get_height)
    GTB=$((HEIGHT + 20))
    echo "Current block height: $HEIGHT"
    echo "Using GoodTilBlock: $GTB"
    echo ""

    # Alice: BUY 10 quantums of clobPair 0 at subticks 10000
    # Use a different clientId to avoid conflicts with previous test
    echo "1. Alice places BUY order (clientId=300)..."
    $BINARY tx clob place-order "$ALICE_ADDR" 0 300 0 1 10 10000 $GTB \
        --from alice --chain-id $CHAIN_ID --node $NODE --keyring-backend=test \
        -y --broadcast-mode sync $TX_FEES 2>&1 | grep -E "code:|txhash:" || true
    echo ""

    # Wait for alice's order to land in a block
    echo "Waiting for alice's order to land..."
    sleep 2
    NEXT_HEIGHT=$(($(get_height) + 1))
    echo "Will target next block $NEXT_HEIGHT for bob's place+cancel"
    echo ""

    # Bob: SELL 10 quantums crossing alice's buy, then IMMEDIATELY cancel
    echo "2. Bob places SELL order (clientId=400) - would cross alice's buy..."
    $BINARY tx clob place-order "$BOB_ADDR" 0 400 0 2 10 10000 $GTB \
        --from bob --chain-id $CHAIN_ID --node $NODE --keyring-backend=test \
        -y --broadcast-mode sync $TX_FEES 2>&1 | grep -E "code:|txhash:" || true

    echo "3. Bob IMMEDIATELY cancels his order (clientId=400)..."
    $BINARY tx clob cancel-order "$BOB_ADDR" 0 400 0 $GTB \
        --from bob --chain-id $CHAIN_ID --node $NODE --keyring-backend=test \
        -y --broadcast-mode sync $TX_FEES 2>&1 | grep -E "code:|txhash:" || true
    echo ""

    echo "Both txs sent. Wait for 1-2 blocks, then check subaccount state."
    echo ""
    echo "EXPECTED RESULT (deferred matching):"
    echo "  - Alice's subaccount should NOT have a perpetual position (no match occurred)"
    echo "  - Bob's subaccount should NOT have a perpetual position"
    echo ""
    echo "Query subaccounts:"
    echo "  curl -s 'localhost:1317/dydxprotocol/subaccounts/subaccount/dydx199tqg4wdlnu4qjlxchpd7seg454937hjrknju4/0' | jq '.subaccount.perpetual_positions'"
    echo "  curl -s 'localhost:1317/dydxprotocol/subaccounts/subaccount/dydx10fx7sy6ywd5senxae9dwytf8jxek3t2gcen2vs/0' | jq '.subaccount.perpetual_positions'"
}

query_height() {
    HEIGHT=$(get_height)
    echo "Current block height: $HEIGHT"
}

query_subaccounts() {
    echo "=== Alice's Subaccount (0) ==="
    curl -s "localhost:1317/dydxprotocol/subaccounts/subaccount/$ALICE_ADDR/0" | jq '.subaccount.perpetual_positions'
    echo ""
    echo "=== Bob's Subaccount (0) ==="
    curl -s "localhost:1317/dydxprotocol/subaccounts/subaccount/$BOB_ADDR/0" | jq '.subaccount.perpetual_positions'
}

case "${1:-help}" in
    setup-keys)
        setup_keys
        ;;
    test-match)
        test_basic_match
        ;;
    test-cancel-before-match)
        test_cancel_before_match
        ;;
    query-height)
        query_height
        ;;
    query-subaccounts)
        query_subaccounts
        ;;
    help|*)
        echo "Usage: $0 {setup-keys|test-match|test-cancel-before-match|query-height|query-subaccounts}"
        echo ""
        echo "Commands:"
        echo "  setup-keys              Import alice and bob keys to local keyring"
        echo "  test-match              Place crossing orders (alice buy + bob sell) to test basic matching"
        echo "  test-cancel-before-match  Place alice buy, then bob sell+cancel to test deferred cancel"
        echo "  query-height            Show current block height"
        echo "  query-subaccounts       Show perpetual positions for alice and bob"
        ;;
esac
