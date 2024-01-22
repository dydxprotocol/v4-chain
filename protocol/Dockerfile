# NB: This is a digest for a multi-arch manifest list, you will want to get this by running
# `docker buildx imagetools inspect golang:1.21-alpine`
ARG GOLANG_1_21_ALPINE_DIGEST="926f7f7e1ab8509b4e91d5ec6d5916ebb45155b0c8920291ba9f361d65385806"

# This Dockerfile is a stateless build of the `dydxprotocold` binary as a Docker container.
# It does not include any configuration, state, or genesis information.

# --------------------------------------------------------
# Builder
# --------------------------------------------------------

FROM golang@sha256:${GOLANG_1_21_ALPINE_DIGEST} as builder
ARG VERSION
ARG COMMIT

RUN set -eux; apk add --no-cache ca-certificates build-base; apk add git linux-headers bash

# Download go dependencies
WORKDIR /dydxprotocol
COPY go.* ./
RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/root/go/pkg/mod \
    go mod download

# Copy the remaining files
COPY . .

# See https://github.com/CosmWasm/wasmvm/releases
#ADD https://github.com/CosmWasm/wasmvm/releases/download/v1.5.0/libwasmvm_muslc.aarch64.a /lib/libwasmvm_muslc.aarch64.a
#RUN sha256sum /lib/libwasmvm_muslc.aarch64.a | grep 2687afbdae1bc6c7c8b05ae20dfb8ffc7ddc5b4e056697d0f37853dfe294e913

ADD https://github.com/CosmWasm/wasmvm/releases/download/v1.5.0/libwasmvm_muslc.x86_64.a /lib/libwasmvm_muslc.x86_64.a

# Copy the library you want to the final location that will be found by the linker flag `-lwasmvm_muslc`
#RUN cp /lib/libwasmvm_muslc.aarch64.a /lib/libwasmvm_muslc.a
RUN cp /lib/libwasmvm_muslc.x86_64.a /lib/libwasmvm_muslc.a

# Build dydxprotocold binary
# force it to use static lib (from above) not standard libgo_cosmwasm.so file
RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/root/go/pkg/mod \
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

RUN echo "Ensuring binary is statically linked ..." \
  && (file /dydxprotocol/build/dydxprotocold | grep "statically linked")

# --------------------------------------------------------
# Runner
# --------------------------------------------------------

FROM golang@sha256:${GOLANG_1_21_ALPINE_DIGEST}

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
