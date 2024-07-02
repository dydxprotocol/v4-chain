<p align="center"><img src="https://dydx.exchange/icon.svg?" width="256" /></p>

<h1 align="center">dYdX Chain Proto</h1>

<div align="center">
  <a href="https://github.com/dydxprotocol/v4-chain/actions/workflows/proto.yml?query=branch%3Amain" style="text-decoration:none;">
    <img src="https://github.com/dydxprotocol/v4-chain/actions/workflows/proto.yml/badge.svg?branch=main" />
  </a>
</div>

This directory defines all protos for `v4-chain`. We follow the Cosmos-SDK convention of using a tool called
[buf](https://github.com/bufbuild/buf) to manage proto dependencies. You can think of `buf` as being like `npm` for
protocol buffers. See the `buf` [documentation](https://docs.buf.build/how-to/iterate-on-modules#update-dependencies)
for further details.

## Building protos
After making changes to any .proto file(s), you will also need to build the protos. To do this, run `make proto-gen` under the root `/v4-chain` directory. If the changes to the protos is also used by Indexer, you will also need to run `pnpm build:proto` after `cd`-ing into `/v4-chain/indexer/packages/v4-protos`
