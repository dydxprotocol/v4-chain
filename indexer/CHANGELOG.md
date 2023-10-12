# Changelog

## [Unreleased]

### Features

* [#143](https://github.com/dydxprotocol/v4-chain/pull/143) Add websocket events for perpetual markets.

### Improvements
* [#577](https://github.com/dydxprotocol/v4-chain/pull/577) Correctly set order status for all order types

* [#552](https://github.com/dydxprotocol/v4-chain/pull/552) Updated Elliptic compliance client block addresses with risk scores equal to the threshold in addition to risk score greater than the threshold.

* [#469](https://github.com/dydxprotocol/v4-chain/pull/469) Added a reason field to `/screen` endpoint to display a reason for blocking an address.
  
### Bug Fixes
* [#579](https://github.com/dydxprotocol/v4-chain/pull/579) Fixed bug where open short-term orders were not being returned in the initial payload when subscribing to the v4_subaccounts channel.

* [#552](https://github.com/dydxprotocol/v4-chain/pull/552) Fixed bug with Elliptic compliance client where the API key was incorrectly used instead of the API secret to generate the auth headers for the Elliptic request.

* [#528](https://github.com/dydxprotocol/v4-chain/pull/528) Fixed bug with bulk SQL queries with nullable numeric / string / boolean values.

* [#496](https://github.com/dydxprotocol/v4-chain/pull/496) Don't geo-block requests made to `comlink` if it's from an internal ip.

* [#460](https://github.com/dydxprotocol/v4-chain/pull/460/files) Fixed track lag timing to output positive numbers.

### API Breaking Changes
