FROM dydxprotocol-base

COPY ./testing/testnet-dev/dev.sh /dydxprotocol/dev.sh
COPY ./testing/genesis.sh /dydxprotocol/genesis.sh
COPY ./testing/start.sh /dydxprotocol/start.sh
COPY ./daemons/pricefeed/client/constants/testdata /dydxprotocol/exchange_config
COPY ./testing/delaymsg_config /dydxprotocol/delaymsg_config

# Install cosmovisor, then purge the Go module cache in the SAME layer so the
# downloaded module source (which includes npm lockfiles like cosmos-sdk's
# docs/package-lock.json) is never committed into the shipped image. Must be
# same-layer: image scanners inspect per-layer, so a later cleanup wouldn't help.
RUN go install cosmossdk.io/tools/cosmovisor/cmd/cosmovisor@v1.5.0 && go clean -modcache

RUN /dydxprotocol/dev.sh

ENV DAEMON_NAME=dydxprotocold
ENTRYPOINT ["cosmovisor", "run"]
