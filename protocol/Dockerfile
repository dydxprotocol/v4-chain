# NB: This is a digest for a multi-arch manifest list, you will want to get this by running
# `docker buildx imagetools inspect golang:1.22-alpine`
ARG GOLANG_1_22_ALPINE_DIGEST="6522f0ca555a7b14c46a2c9f50b86604a234cdc72452bf6a268cae6461d9000b"

# This Dockerfile is a stateless build of the `klyra` binary as a Docker container.
# It does not include any configuration, state, or genesis information.

# --------------------------------------------------------
# Builder
# --------------------------------------------------------

FROM golang@sha256:${GOLANG_1_22_ALPINE_DIGEST} as builder
ARG VERSION
ARG COMMIT

ENV DAEMON_NAME=klyra
ENV DAEMON_HOME=/klyraprotocol/chain/.alice

RUN set -eux; apk add --no-cache ca-certificates build-base; apk add git linux-headers bash binutils-gold

# Download go dependencies
WORKDIR /klyraprotocol
COPY go.* ./
RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/root/go/pkg/mod \
    go mod download

# Copy the remaining files
COPY . .

# Build klyra binary
RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/root/go/pkg/mod \
    go build \
    -mod=readonly \
    -tags "netgo,ledger,muslc" \
    -ldflags "-X github.com/cosmos/cosmos-sdk/version.Name="klyraprotocol" \
    -X github.com/cosmos/cosmos-sdk/version.AppName="klyra" \
    -X github.com/cosmos/cosmos-sdk/version.Version=$VERSION \
    -X github.com/cosmos/cosmos-sdk/version.Commit=$COMMIT \
    -X github.com/cosmos/cosmos-sdk/version.BuildTags='netgo,ledger,muslc' \
    -w -s -linkmode=external -extldflags '-Wl,-z,muldefs -static'" \
    -trimpath \
    -o /klyraprotocol/build/ \
    ./...

# --------------------------------------------------------
# Runner
# --------------------------------------------------------

FROM golang@sha256:${GOLANG_1_22_ALPINE_DIGEST}

RUN apk add --no-cache bash

COPY --from=builder /klyraprotocol/build/klyra /bin/klyra

ENV HOME /klyraprotocol
WORKDIR $HOME

# tendermint p2p
EXPOSE 26656
# tendermint rpc
EXPOSE 26657
# rest server
EXPOSE 1317
# grpc
EXPOSE 9090

ENTRYPOINT ["klyra"]
