#!/bin/bash
# snapshot.sh take DIR | diff BEFORE_DIR AFTER_DIR
# Dumps the testnet state the canonical USDC upgrade must not disturb:
# subaccount collateral and positions, MegaVault shares, USDC supply, perpetual
# open interest, Noble channel packet commitments, and canonicalusdc state once
# the module exists.
set -euo pipefail
REST="${REST:-https://test-dydx-rest.kingnodes.com}"
USDC="ibc/8E27BA2D5493AF5636760E354E46004562C46AB7EC0CC4C1CA14E9E20E2545B5"
get() { curl -sf -m 60 "$REST/$1"; }
paginate() { # paginate PATH JQ_ARRAY_KEY
	local next="" page out
	while :; do
		page=$(get "$1&pagination.limit=1000${next:+&pagination.key=$(printf %s "$next" | jq -sRr @uri)}")
		echo "$page" | jq -c ".$2[]"
		next=$(echo "$page" | jq -r '.pagination.next_key // empty')
		[ -n "$next" ] || break
	done
}
take() {
	local dir=$1
	mkdir -p "$dir"
	get "cosmos/base/tendermint/v1beta1/blocks/latest" | jq '{height: .block.header.height, time: .block.header.time}' >"$dir/meta.json"
	paginate "dydxprotocol/subaccounts/subaccount?x=1" subaccount |
		jq -c 'select((.asset_positions | length) + (.perpetual_positions | length) > 0) |
			{id: "\(.id.owner)/\(.id.number)", assets: [.asset_positions[] | {asset_id, quantums}],
			 perps: [.perpetual_positions[] | {perpetual_id, quantums}]}' | jq -s 'sort_by(.id)' >"$dir/subaccounts.json"
	{ paginate "dydxprotocol/vault/megavault/all_owner_shares?x=1" owner_shares | jq -s 'sort_by(.owner)'
	  get "dydxprotocol/vault/megavault/total_shares"; } | jq -s '{owners: .[0], total: .[1].total_shares}' >"$dir/vault.json"
	get "cosmos/bank/v1beta1/supply/by_denom?denom=$USDC" | jq '.amount' >"$dir/supply.json"
	paginate "dydxprotocol/perpetuals/perpetual?x=1" perpetual | jq -c '{id: .params.id, open_interest}' | jq -s 'sort_by(.id)' >"$dir/perpetuals.json"
	get "ibc/core/channel/v1/channels/channel-0/ports/transfer/packet_commitments?pagination.limit=1000" |
		jq '[.commitments[].sequence]' >"$dir/noble_commitments.json"
	get "dydxprotocol/v4/canonicalusdc/state" >"$dir/canonicalusdc.json" 2>/dev/null || echo '{"absent": true}' >"$dir/canonicalusdc.json"
	echo "snapshot at height $(jq -r .height "$dir/meta.json"): $(jq length "$dir/subaccounts.json") subaccounts with positions, $(jq '.owners | length' "$dir/vault.json") MegaVault owners, USDC supply $(jq -r .amount "$dir/supply.json")"
}
diff_snapshots() {
	local a=$1 b=$2 f
	for f in subaccounts vault supply perpetuals noble_commitments; do
		if cmp -s "$a/$f.json" "$b/$f.json"; then
			echo "$f: identical"
		else
			echo "$f: DIFFERENT"
			diff <(jq -S . "$a/$f.json") <(jq -S . "$b/$f.json") | head -40
		fi
	done
	echo "canonicalusdc before: $(jq -c '.controls.mode // .absent' "$a/canonicalusdc.json")  after: $(jq -c '.controls.mode // .absent' "$b/canonicalusdc.json")"
}
case "${1:-}" in
take) take "$2" ;;
diff) diff_snapshots "$2" "$3" ;;
*) echo "usage: $0 take DIR | diff BEFORE AFTER" >&2; exit 1 ;;
esac
