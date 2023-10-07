# Changelog

## [Unreleased]

### Features

* [#190](https://github.com/dydxprotocol/v4-chain/pull/190) Added `MsgCreateClobPair` to x/clob to allow creation of new clob pairs.

### Improvements

* [#479](https://github.com/dydxprotocol/v4-chain/pull/479) Ensure that rate limiting of short term order placements/cancellations is guarded against replay attacks.

* [#532](https://github.com/dydxprotocol/v4-chain/pull/532) Use updated `data-api.binance.vision` endpoint for Binance price data.

### Bug Fixes
* [#449](https://github.com/dydxprotocol/v4-chain/pull/449) Removes possible undesirable panics in x/rewards `EndBlocker`.

### API Breaking Changes
* [#477](https://github.com/dydxprotocol/v4-chain/pull/477) Changes the fields of `DelayedCompleteBridgeMessage`, `DelayedMessage`, `QueryBlockMessageIdsRequest`.

### State Breaking Changes
* [#477](https://github.com/dydxprotocol/v4-chain/pull/477) Changes the name of various keys in state.
