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

## 5. GRADUAL walk
- Injective deposit: from the stand-in, `ibc-simd tx ibc-transfer transfer transfer channel-0 <dydx addr> <amt>erc20:0xa00C... --home ~/.ibc-staging/injective --node tcp://127.0.0.1:29657 --chain-id injective-1 --from user --keyring-backend test --fees 0stake`.
  Expect uusdc minted to the receiver, injective_backing up, the voucher held by the module.
- Injective withdrawal: `dydxprotocold tx ibc-transfer transfer transfer channel-94 <cosmos addr> <amt>ibc/8E27... --from dydx-1 ...`; settles on ack, burns backing.
- Rate limit: `q ratelimit capacity-by-denom ibc/8E27...` before/after; a round trip is neutral.
- Noble route only with a channel-0 relayer. Finalize (disable Noble withdrawals) is reversible by proposal on testnet.

## 0. Cutting protocol/v9.9.0-rc0 for testnet
Testnet images are built by `protocol-build-and-push-testnet.yml` on pushes to `main` or `release/protocol/vX[.Y].x`
in dydxprotocol/v4-chain, tagged by 7-char commit hash in ECR `testnet-validator`. Tags `protocol/vX.Y.Z-rcN` trigger
`protocol-release.yml`, which builds reproducible linux binaries into a draft GitHub release. `release/protocol/v9.x`
carries the 9.8 line and has diverged from main, so 9.9 gets its own branch from the PR tip (main lineage).
```
tip=<green tip>                                   # after the last ibc-local upgrade run passes
git branch release/protocol/v9.9.x $tip
git tag -s protocol/v9.9.0-rc0 $tip -m "protocol/v9.9.0-rc0"   # signed, like the 9.8 tags
git push origin release/protocol/v9.9.x protocol/v9.9.0-rc0    # one push: image gets version 9.9.0-rc0
```
Then: approve nothing, CI runs on the branch push; publish the draft release (needed by cosmovisor flows and by the
staging image's PREUPGRADE_VERSION); note the ECR image tag (`git rev-parse --short=7 $tip`).
Testnet validators on 9.8.2 crash-loop at the plan height with `UPGRADE "v9.9" NEEDED`; switch both validator task
definitions (useast2 first, apne1 ~8 min) and the full nodes to the rc0 image then. The new binary applies the upgrade
on start and enters GRADUAL with the pinned testnet routes.

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
