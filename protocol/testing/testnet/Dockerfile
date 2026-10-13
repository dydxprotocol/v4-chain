FROM dydxprotocol-base

RUN apk add --no-cache bash jq aws-cli
# Install cosmovisor, then purge the Go module cache in the SAME layer so the
# downloaded module source (which includes npm lockfiles like cosmos-sdk's
# docs/package-lock.json) is never committed into the shipped image. Must be
# same-layer: image scanners inspect per-layer, so a later cleanup wouldn't help.
RUN go install cosmossdk.io/tools/cosmovisor/cmd/cosmovisor@v1.5.0 && go clean -modcache

COPY ./testing/version/. /dydxprotocol/
COPY ./testing/testnet/. /dydxprotocol/

ENV HOME /dydxprotocol
WORKDIR $HOME

RUN /dydxprotocol/testnet.sh

ENTRYPOINT ["/dydxprotocol/start.sh"]
