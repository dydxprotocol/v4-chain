# Protocol v5.0.0 Release Notes

## Highlights

This release includes changes in following areas:

- [Isolated Markets](#isolated-market)
- Batch order cancellation
- Protocol-enshrined LP Vault
- Signature Verification Parallelization
- Soft Open Interest Cap
- Full Node Streaming

See below for a detailed breakdown of the changes.

## What's New?

### Isolated Market

[Description]

- Move insurance fund into x/perpetuals ([#1106](https://github.com/dydxprotocol/v4-chain/pull/1106))
- Add market type to the CreatePerpetual API ([#1118](https://github.com/dydxprotocol/v4-chain/pull/1118))
- Use market specific insurance fund for cross or isolated markets ([#1132](https://github.com/dydxprotocol/v4-chain/pull/1132))
- Add function to retrieve collateral pool addr for a subaccount ([#1142](https://github.com/dydxprotocol/v4-chain/pull/1142))
- Move SA module address transfers to use perpetual based SA accounts ([#1146](https://github.com/dydxprotocol/v4-chain/pull/1146))
- Add state migrations for isolated markets ([#1155](https://github.com/dydxprotocol/v4-chain/pull/1155))
- Update upper limit of subaccount number constraint ([#1164](https://github.com/dydxprotocol/v4-chain/pull/1164))
- Check isolated market constraints in UpdateSubaccount. ([#1158](https://github.com/dydxprotocol/v4-chain/pull/1158))
- Update subaccount-to-subaccount transfers for collateral pools. ([#1178](https://github.com/dydxprotocol/v4-chain/pull/1178))
- Add market type to indexer events. ([#1201](https://github.com/dydxprotocol/v4-chain/pull/1201))
- Add logic to transfer collateral when open/close isolated perpetual position. ([#1200](https://github.com/dydxprotocol/v4-chain/pull/1200))
- Update withdrawal gating to be per collateral pool. ([#1208](https://github.com/dydxprotocol/v4-chain/pull/1208))
- v5.0.0 state migration negative tnc subaccount store. ([#1226](https://github.com/dydxprotocol/v4-chain/pull/1226))
- Add query function for collateral pool address. ([#1256](https://github.com/dydxprotocol/v4-chain/pull/1256))
- Add new error code to be returned when orders fail isolated subaccount checks. ([#1331](https://github.com/dydxprotocol/v4-chain/pull/1331))
- Add V2 for perpetual market create event + handlers. ([#1395](https://github.com/dydxprotocol/v4-chain/pull/1395))
- Fix IsIsolatedPerpetual function orientation ([#1403](https://github.com/dydxprotocol/v4-chain/pull/1403))
- Fix consensus failure from unhandled error. ([#1418](https://github.com/dydxprotocol/v4-chain/pull/1418))

### Protocol Vault

[Description]

- add x/vault protos ([#1169](https://github.com/dydxprotocol/v4-chain/pull/1169))
- add x/vault total shares store ([#1179](https://github.com/dydxprotocol/v4-chain/pull/1179))
- construct simple vault orders ([#1206](https://github.com/dydxprotocol/v4-chain/pull/1206))
- refresh all vault orders ([#1213](https://github.com/dydxprotocol/v4-chain/pull/1213))
- store x/vault parameters in state, add query, and update via governance ([#1238](https://github.com/dydxprotocol/v4-chain/pull/1238))
- Implement MsgDepositToVault ([#1220](https://github.com/dydxprotocol/v4-chain/pull/1220))
- change vault state key to not use serialization ([#1248](https://github.com/dydxprotocol/v4-chain/pull/1248))
- Track shares of owners in vaults ([#1253](https://github.com/dydxprotocol/v4-chain/pull/1253))
- decommission vaults at the beginning of a block ([#1264](https://github.com/dydxprotocol/v4-chain/pull/1264))
- add cli for deposit to vault ([#1276](https://github.com/dydxprotocol/v4-chain/pull/1276))
- rename order_size_ppm to order_size_pct_ppm in vault params ([#1271](https://github.com/dydxprotocol/v4-chain/pull/1271))
- implement full vault strategy ([#1262](https://github.com/dydxprotocol/v4-chain/pull/1262))
- keep track of vault shares as integers ([#1267](https://github.com/dydxprotocol/v4-chain/pull/1267))
- vaults do not quote if order size is 0 ([#1325](https://github.com/dydxprotocol/v4-chain/pull/1325))
- Skip various order placement logic for vault orders ([#1332](https://github.com/dydxprotocol/v4-chain/pull/1332))
- return error if shares to mint is rounded down to zero ([#1333](https://github.com/dydxprotocol/v4-chain/pull/1333))
- initialize vault module params in v5.0.0 upgrade handler ([#1359](https://github.com/dydxprotocol/v4-chain/pull/1359))
- add activation threshold on vaults ([#1344](https://github.com/dydxprotocol/v4-chain/pull/1344))
- add OwnerShares and AllVaults query and cli ([#1362](https://github.com/dydxprotocol/v4-chain/pull/1362))
- bounded quotes by oracle price ([#1390](https://github.com/dydxprotocol/v4-chain/pull/1390))
- updated default vault parameters ([#1450](https://github.com/dydxprotocol/v4-chain/pull/1450))

### BatchCancel

`MsgBatchCancel` added to the protocol. Allow for batches (max 100) of short term orders to be cancelled. Treated as 2 order cancels by rate limiter.

- MsgBatchCancel Protos ([#1094](https://github.com/dydxprotocol/v4-chain/pull/1094))
- Validatebasic for new MsgBatchCancel ([#1101](https://github.com/dydxprotocol/v4-chain/pull/1101))
- msgBatchCancel success + failure fields ([#1130](https://github.com/dydxprotocol/v4-chain/pull/1130))
- Clob `MsgBatchCancel` functionality ([#1110](https://github.com/dydxprotocol/v4-chain/pull/1110))
- Batch Cancel CLI + additional validation ([#1204](https://github.com/dydxprotocol/v4-chain/pull/1204))
- Batch Cancel Rate Limit Metrics ([#1233](https://github.com/dydxprotocol/v4-chain/pull/1233))

### Soft Open Interest Caps

[Description]

- OIMF protos and genesis values ([#1125](https://github.com/dydxprotocol/v4-chain/pull/1125))
- Modify `GetInitialMarginQuoteQuantums` to reflect OIMF ([#1159](https://github.com/dydxprotocol/v4-chain/pull/1159))
- Calculate current OI and pass to `GetMarginRequirements` calculation ([#1161](https://github.com/dydxprotocol/v4-chain/pull/1161))
- Evaluate subaccounts update end state with OIMF ([#1172](https://github.com/dydxprotocol/v4-chain/pull/1172))
- OIMF: Simplify interface for `GetDeltaOpenInterestFromUpdates` ([#1227](https://github.com/dydxprotocol/v4-chain/pull/1227))
- Add upgrade handler for OIMF caps ([#1232](https://github.com/dydxprotocol/v4-chain/pull/1232))
- Update OI after fill and added unit tests ([#1231](https://github.com/dydxprotocol/v4-chain/pull/1231))
- Update LiquidityTierUpsertEvent to add OI caps ([#1242](https://github.com/dydxprotocol/v4-chain/pull/1242))
- Add upgrade handler to initialize OI during upgrade handler ([#1302](https://github.com/dydxprotocol/v4-chain/pull/1302))
- Update OI caps to new values ([#1351](https://github.com/dydxprotocol/v4-chain/pull/1351))
- update caps in genesis file ([#1355](https://github.com/dydxprotocol/v4-chain/pull/1355))
- update DB and APIs to add OI caps info ([#1305](https://github.com/dydxprotocol/v4-chain/pull/1305))
- Add open interest update event and handlers ([#1352](https://github.com/dydxprotocol/v4-chain/pull/1352))
- Add `QueryAllLiquidityTiers` to CLI ([#1434](https://github.com/dydxprotocol/v4-chain/pull/1434))
- Refactor perpetual upgrade handlers in `v5.0.0` ([#1442](https://github.com/dydxprotocol/v4-chain/pull/1442))
- Update iavl with WorkingHash fix ([#1484](https://github.com/dydxprotocol/v4-chain/pull/1484))

### Order Rate Limits

Rate Limits for short term orders and cancels were combined. Previously, rate limits were set to be 200 short term place orders over 1 block and 200 short term cancel orders over 1 block. To relax rate limits and allow for burst placing of orders, place and cancel order rates were combined and moved to a 5-block window. Max is now `(200 + 200) * 5 = 2000` combined short term place and cancel orders over 5 blocks.

- Combine place, cancel, batch cancel rate limiters ([#1165](https://github.com/dydxprotocol/v4-chain/pull/1165))
- Fix upgrade rate limit config logic, add helpful logs ([#1212](https://github.com/dydxprotocol/v4-chain/pull/1212))

### Full Node Streaming

[Description]

- stream offchain updates through stream manager ([#1138](https://github.com/dydxprotocol/v4-chain/pull/1138))
- add command line flag for full node streaming ([#1145](https://github.com/dydxprotocol/v4-chain/pull/1145))
- construct the initial orderbook snapshot ([#1147](https://github.com/dydxprotocol/v4-chain/pull/1147))
- separate indexer and grpc streaming events ([#1209](https://github.com/dydxprotocol/v4-chain/pull/1209))
- only send response when there is at least one update ([#1216](https://github.com/dydxprotocol/v4-chain/pull/1216))
- send order update when short term order state fill amounts a… ([#1241](https://github.com/dydxprotocol/v4-chain/pull/1241))
- send fill amount updates for reverted operations in prepare check state ([#1240](https://github.com/dydxprotocol/v4-chain/pull/1240))
- add block number + stage to grpc updates ([#1252](https://github.com/dydxprotocol/v4-chain/pull/1252))
- avoid state reads when sending updates ([#1261](https://github.com/dydxprotocol/v4-chain/pull/1261))
- send updates for both normal order matches and liquidation ([#1280](https://github.com/dydxprotocol/v4-chain/pull/1280))
- decouple grpc and non validating full node flags ([#1400](https://github.com/dydxprotocol/v4-chain/pull/1400))

### Performance Improvements

[Description]

- Migrate over to cosmos version now that performance improvement had been released. ([#1007](https://github.com/dydxprotocol/v4-chain/pull/1007))
- Push CheckTx concurrency into protocol ([#1030](https://github.com/dydxprotocol/v4-chain/pull/1030))
- Cache current collateralization and margin requirements per subaccount ([#1064](https://github.com/dydxprotocol/v4-chain/pull/1064))
- Increase maximum number of connections ([#1080](https://github.com/dydxprotocol/v4-chain/pull/1080))
- Reduce the cost of adding an indexer event. ([#1078](https://github.com/dydxprotocol/v4-chain/pull/1078))
- Update Cosmos SDK and add parallel query tests. ([#1086](https://github.com/dydxprotocol/v4-chain/pull/1086))
- Optimize string building loop to write to a single buffer instead of to multiple buffers that are concatenated.
([#1119](https://github.com/dydxprotocol/v4-chain/pull/1119))
- Create ante handlers which parallelize processing per account until a global lock is necessary. ([#1108](https://github.com/dydxprotocol/v4-chain/pull/1108))
- Upgrade cosmos fork to use libsecp256k1 from go-ethereum ([#1210](https://github.com/dydxprotocol/v4-chain/pull/1210))
- Store pruneable orders in a key-per-order format rather than key-per-height ([#1230](https://github.com/dydxprotocol/v4-chain/pull/1230))
- Fix bugs in pruneable order migration ([#1250](https://github.com/dydxprotocol/v4-chain/pull/1250))
- Skip add to book collat check upon placement for conditional orders ([#1369](https://github.com/dydxprotocol/v4-chain/pull/1369))

### Slinky/Vote Extension

[Description] 

- feat: Slinky sidecar integration ([#1109](https://github.com/dydxprotocol/v4-chain/pull/1109))
- feat: VoteExtension Slinky logic ([#1139](https://github.com/dydxprotocol/v4-chain/pull/1139))
- feat: Proposal logic for Slinky ([#1135](https://github.com/dydxprotocol/v4-chain/pull/1135))
- feat: Enforce MarketParam.Pair uniqueness constraint ([#1193](https://github.com/dydxprotocol/v4-chain/pull/1193))
- feat: Slinky full integration PR ([#1141](https://github.com/dydxprotocol/v4-chain/pull/1141))
- bug: Handle paginated querying to prices-keeper in slinky + pricefeed daemon ([#1177](https://github.com/dydxprotocol/v4-chain/pull/1177))
- Remove smoothed prices ([#1215](https://github.com/dydxprotocol/v4-chain/pull/1215))
- perf: optimize `slinky_adapter` via in-memory cache [SKI-13] ([#1182](https://github.com/dydxprotocol/v4-chain/pull/1182))
- Bump slinky version to v0.3.1 ([#1275](https://github.com/dydxprotocol/v4-chain/pull/1275))
- Add v5.0.0 enable vote extensions upgrade handler ([#1308](https://github.com/dydxprotocol/v4-chain/pull/1308))
- fix: Fix v5 upgrade VE handler ([#1468](https://github.com/dydxprotocol/v4-chain/pull/1468))

### Bug fixes

- [bugfix][protocol] fix settled funding event emission ([#1051](https://github.com/dydxprotocol/v4-chain/pull/1051))
- Fix --home flag being ignored in cosmos v0.50 ([#1053](https://github.com/dydxprotocol/v4-chain/pull/1053))
- Fix unnecessary home directory creation issues ([#1073](https://github.com/dydxprotocol/v4-chain/pull/1073))
- Upgrade iavl to v1.0.1 to prune orphan nodes async ([#1081](https://github.com/dydxprotocol/v4-chain/pull/1081))
- Upgrade client/v2 to fix autocli signing issues ([#1090](https://github.com/dydxprotocol/v4-chain/pull/1090))
- Don't attempt to write the store during simulated transactions. ([#1122](https://github.com/dydxprotocol/v4-chain/pull/1122))
- Disable fast nodes and update iavl to address bugs ([#1173](https://github.com/dydxprotocol/v4-chain/pull/1173))
- Fix lib.ErrorLogWithError ([#1306](https://github.com/dydxprotocol/v4-chain/pull/1306))
- Upgrade IAVL to fix issues saving to disk ([#1357](https://github.com/dydxprotocol/v4-chain/pull/1357))
- Fix stateful order removal event ordering. ([#1456](https://github.com/dydxprotocol/v4-chain/pull/1456))
- Stop passing in app PreBlocker to slinky VE handler ([#1491](https://github.com/dydxprotocol/v4-chain/pull/1491))

### Miscellaneous

- [feature] Add request_id to all code paths using lib.UnwrapSDKContext ([#1058](https://github.com/dydxprotocol/v4-chain/pull/1058))
- Metric for when subaccounts are created ([#1071](https://github.com/dydxprotocol/v4-chain/pull/1071))
- deprecate pessimistic add-to-book collateralization check ([#1107](https://github.com/dydxprotocol/v4-chain/pull/1107))
- Update deleveraging handler for DeleveragingEventV1 ([#1121](https://github.com/dydxprotocol/v4-chain/pull/1121))
- Add `v5` upgrade handler and set up container upgrade test ([#1153](https://github.com/dydxprotocol/v4-chain/pull/1153))
- Emit metrics gated through execution mode ([#1157](https://github.com/dydxprotocol/v4-chain/pull/1157))
- Update go version to 1.22 in [README.md](http://readme.md/) ([#1196](https://github.com/dydxprotocol/v4-chain/pull/1196))
- Add query for PendingSendPacket ([#1176](https://github.com/dydxprotocol/v4-chain/pull/1176))
- Upgrade cosmos fork to v0.50.5 ([#1246](https://github.com/dydxprotocol/v4-chain/pull/1246))
- Move price cache population to Preblocker ([#1417](https://github.com/dydxprotocol/v4-chain/pull/1417))
- Move clob hydration to preblocker ([#1412](https://github.com/dydxprotocol/v4-chain/pull/1412))
- Upgrade to newly rebased cometbft fork ([#1438](https://github.com/dydxprotocol/v4-chain/pull/1438))
- Increase conditional order trigger multiplier to 25 ([#1460](https://github.com/dydxprotocol/v4-chain/pull/1460))
- Add a cli to convert module name to address ([#1462](https://github.com/dydxprotocol/v4-chain/pull/1462))