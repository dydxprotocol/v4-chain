# Canonical USDC: dydx-testnet-4 upgrade rehearsal

## Already in place (from this machine)
- Injective stand-in: ibc-go simd, chain id `injective-1`, mints `erc20:0xa00C59fF5a080D2b954d0c75e46E22a0c371235a`.
  State `~/.ibc-staging`, rpc 127.0.0.1:29657. Restart with `protocol/testing/ibc-local/staging/standin.sh` (idempotent).
- Hermes relaying injective-1 <-> dydx-testnet-4 (our node 3.22.142.107, rpc 26657 / grpc 9090), key store `~/.ibc-staging/hermes/keys`.
  Relayer on testnet: dydx18uny9ygfu6wzal7ljqvs2sq56t5at60jh0qxdp, funded 100 DV4TNT from dydx-1.
- Channel: dYdX side channel-94 / connection-58 / 07-tendermint-62 (open, ics20-1), stand-in side channel-0. 1 USDC relayed end to end.
  Voucher on testnet: ibc/26148908E8C7EB6D652F73103F9BDF39B3C602F9B5C73EEC6F1F20ACFB5B7466.
- Handler pins dydx-testnet-4 to Noble channel-0/connection-0/07-tendermint-0 and Injective channel-94/connection-58/07-tendermint-62 (commit 4132ee9c2).
- Baseline snapshot: `protocol/testing/ibc-local/staging/baseline-20261112-1442` (1,149,412 subaccounts with positions, 37 MegaVault owners).
  Retake right before the proposal: `protocol/testing/ibc-local/staging/snapshot.sh take <dir>`.

## 1. Build and stage the binary
- Branch tip must include 4132ee9c2. Build the testnet image the usual way (rc tag on the branch, e.g. protocol/v9.9.0-rc0) or a one-off image.
- Cosmovisor layout on both validators and the full nodes: `upgrades/v9.9/bin/dydxprotocold` (plan name is `v9.9`).

## 2. Pre-upgrade user state on testnet (optional but recommended, dydx-1 holds 484 USDC)
- Deposit to subaccount, fill a small position against dydx-2, rest a long-term bid, deposit into MegaVault — same steps as upgrade.sh's open_position / place_long_term_bid, against --node tcp://3.22.142.107:26657 --chain-id dydx-testnet-4.
- Noble withdrawals left in flight need a relayer on channel-0 (grand-1); skip unless one is running.

## 3. Proposal (expedited: 10 minutes, 75%, both validators vote)
```
height=$(( $(curl -s http://3.22.142.107:26657/status | jq -r .result.sync_info.latest_block_height) + 1200 ))
cat > upgrade.json <<JSON
{"messages":[{"@type":"/cosmos.upgrade.v1beta1.MsgSoftwareUpgrade","authority":"dydx10d07y265gmmuvt4z0w9aw880jnsr700jnmapky",
 "plan":{"name":"v9.9","time":"0001-01-01T00:00:00Z","height":"$height","info":"","upgraded_client_state":null}}],
 "deposit":"10000000000000000000adv4tnt","title":"v9.9 canonical USDC rehearsal","summary":"enter GRADUAL at upgrade","expedited":true}
JSON
dydxprotocold tx gov submit-proposal upgrade.json --from dydx-1 --keyring-backend test --chain-id dydx-testnet-4 --node tcp://3.22.142.107:26657 --gas auto --gas-adjustment 1.5 --gas-prices 25000000000adv4tnt -y
dydxprotocold tx gov vote <id> yes --from dydx-1 ...   ;  dydxprotocold tx gov vote <id> yes --from dydx-2 ...
```
Check the gov min deposit first (`q gov params`). Expect the halt at the height; testnet stops until both validators are on v9.9.

## 4. After the swap
```
dydxprotocold q canonicalusdc state --node tcp://3.22.142.107:26657 -o json
```
Expect: mode MODE_GRADUAL, noble channel-0, injective channel-94, injective_denom ibc/26148908E8C7EB6D652F73103F9BDF39B3C602F9B5C73EEC6F1F20ACFB5B7466,
noble_withdrawals_enabled true, noble_deposits_disabled true (Noble deposits bounce from the upgrade height), injective_deposits_disabled false, max_pending_settlements 10000, ledger.noble_backing == USDC supply, other ledger fields 0.
If the node log says "canonicalusdc stays DISABLED after the upgrade", the bindings did not validate; the proposal path remains.
Then `snapshot.sh take after` and `snapshot.sh diff before after`: subaccounts, vault, supply and perpetuals identical apart from accounts that traded meanwhile.

Outcome 2026-11-12: proposal 412, plan height 92344645 reached at 23:32:58Z; the six live hosts relaunched on the staged binary within the same minute (the indexer snapshot node was mid-tar and swaps when it catches up); state matched the expectations above with noble_backing == supply.

## 5. GRADUAL walk (loop.sh)
`loop.sh state|round [uusdc]|drain|reject` drives the Injective leg through the stand-in (standin.sh must be up). It sends with `--packet-timeout-height 0-0`: the stand-in's dYdX client only moves when packets flow, so the CLI's default relative timeout height is already in the past on dYdX and the packet times out and refunds.
Observed 2026-11-12/13 (three 25 USDC rounds plus drain and rejection, ~10 s per leg): a deposit mints uusdc to the receiver, raises injective_backing and parks the voucher in the module account; a withdrawal moves the amount from injective_backing to pending_injective at send and the ack clears it; supply == noble_backing + injective_backing + pending_injective at every step; a withdrawal above injective_backing fails in DeliverTx with "insufficient route backing" (fee charged, nothing moves); refillable_capacity accrues on Injective withdrawals and is consumed by deposits; the uusdc limiter capacity only followed the TVL baseline (~476M USDC per hour on testnet, so a shortfall regenerates within a block).
Noble route: never on testnet (expired clients). Finalize (disable Noble withdrawals) is reversible by proposal on testnet.

## 0. Cutting the testnet rc (rc0 = 1bc6a8009, superseded; rc1 = c6a7d60c3 with Noble deposits closed at the upgrade)
Testnet images are built by `protocol-build-and-push-testnet.yml` on pushes to `main` or `release/protocol/vX[.Y].x`
in dydxprotocol/v4-chain, tagged by 7-char commit hash in ECR `testnet-validator`. Tags `protocol/vX.Y.Z-rcN` trigger
`protocol-release.yml`, which builds reproducible linux binaries into a draft GitHub release. `release/protocol/v9.x`
carries the 9.8 line and has diverged from main, so 9.9 gets its own branch from the PR tip (main lineage).
```
tip=<green tip>                                   # after the last ibc-local upgrade run passes
# (branch is tq/canonicalusdc-v9.x on release/protocol/v9.x; the image for rc1 is testnet-validator:f9d4939)
git tag -s protocol/v9.9.0-rcN $tip -m "protocol/v9.9.0-rcN"   # signed, like the 9.8 tags; never re-tag, cut the next rc
git push origin tq/canonicalusdc-v9.x protocol/v9.9.0-rcN     # PR #3466 into release/protocol/v9.x builds testnet-validator:<merge-commit hash>
```
Then: approve nothing, CI runs on the branch push; publish the draft release (needed by cosmovisor flows and by the
staging image's PREUPGRADE_VERSION); note the ECR image tag (`git rev-parse --short=7 $tip`).
STAGED 2026-11-12 22:07 UTC on all seven testnet nodes (validators useast2 + apne1, full nodes useast2-0/1/2, indexer
full node and indexer snapshot node in apne1): the rc1 binary (sha256 b85ab68c...) sits at
`<chain home>/cosmovisor/upgrades/v9.9/bin/dydxprotocold` on each node's volume, `current` still points at v9.8, so
cosmovisor swaps by itself at the plan height. No task-definition change is needed for the upgrade. AFTER the upgrade,
move the task definitions to the rc1 image (`testnet-validator:f9d4939` or later) in a maintenance window: start.sh
copies the image's cosmovisor tree on every container start and would point `current` back at v9.8 on a restart.
The staging script is `protocol/testing/ibc-local/staging/stage-cosmovisor.sh`, run over SSM RunShellScript.

Optional cosmovisor dress rehearsal on the staging chain: run `protocol-build-and-push-staging-manual.yml` with
upgrade_name=v9.9, preupgrade_version=9.8.2. That chain (dydxprotocol-testnet) has no IBC channels, so expect the
store migration and the logged DISABLED fallback, not GRADUAL.

### Release-note draft (rc0)
- x/canonicalusdc: canonical USDC backing migration from Noble to Injective (#3465): GRADUAL/PAUSED controls,
  Noble and Injective routes, backing ledger, settlement tracking, Informal Systems audit fixes, audit report under audits/.
- v9.9 upgrade handler adds the module store and enters GRADUAL bound to Noble channel-0 and Injective channel-62
  on dydx-mainnet-1 (pinned testnet routes for dydx-testnet-4); falls back to DISABLED for governance if a binding fails.
- Injective deposits refill uusdc rate-limit capacity that Injective withdrawals spent (bounded, shared limiter).
- Since 9.8.2 (main): Go 1.26.10 and grpc/x deps (#3458), cometbft fork bump (#3461), x/clob withdrawal gating skips
  negative-TNC subaccounts with no positions (#3452), streaming liquidation taker orders (#3451).
