ARG GO_VERSION="1.19"

# This Dockerfile is a stateless build of the `dydxprotocold` binary as a Docker container.
# It does not include any configuration, state, or genesis information.

# --------------------------------------------------------
# Builder
# --------------------------------------------------------

FROM golang:${GO_VERSION}-alpine as builder

RUN set -eux; apk add --no-cache ca-certificates build-base; apk add git linux-headers bash

# Download go dependencies
WORKDIR /dydxprotocol
COPY go.* ./
RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/root/go/pkg/mod \
    go mod download

# Copy the remaining files
COPY . .

# Build dydxprotocold binary
RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/root/go/pkg/mod \
    VERSION=$(echo $(git describe --tags) | sed 's/^v//') && \
    COMMIT=$(git log -1 --format='%H') && \
    go build \
      -mod=readonly \
      -tags "netgo,ledger,muslc" \
      -ldflags "-X github.com/cosmos/cosmos-sdk/version.Name="dydxprotocol" \
              -X github.com/cosmos/cosmos-sdk/version.AppName="dydxprotocold" \
              -X github.com/cosmos/cosmos-sdk/version.Version=$VERSION \
              -X github.com/cosmos/cosmos-sdk/version.Commit=$COMMIT \
              -X github.com/cosmos/cosmos-sdk/version.BuildTags='netgo,ledger,muslc' \
              -w -s -linkmode=external -extldflags '-Wl,-z,muldefs -static'" \
      -trimpath \
      -o /dydxprotocol/build/ \
      ./...

# --------------------------------------------------------
# Runner
# --------------------------------------------------------

FROM golang:${GO_VERSION}-alpine

RUN apk add --no-cache bash

COPY --from=builder /dydxprotocol/build/dydxprotocold /bin/dydxprotocold

ENV HOME /dydxprotocol
WORKDIR $HOME

# tendermint p2p
EXPOSE 26656
# tendermint rpc
EXPOSE 26657
# rest server
EXPOSE 1317
# grpc
EXPOSE 9090

ENTRYPOINT ["dydxprotocold"]
