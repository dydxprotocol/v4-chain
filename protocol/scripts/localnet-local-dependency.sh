#!/bin/bash

set -e

echo "Starting localnet with local cometbft and cosmos-sdk..."

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
NC='\033[0m' # No Color

# Check if local dependencies exist
if [ ! -d "../../cometbft" ]; then
    echo -e "${RED}Error: ../cometbft not found${NC}"
    exit 1
fi

if [ ! -d "../../cosmos-sdk" ]; then
    echo -e "${RED}Error: ../cosmos-sdk not found${NC}"
    exit 1
fi

# Get version and commit info BEFORE changing directories
VERSION=$(git describe --tags --always --match "protocol/v*" | sed 's/^protocol\/v//')
COMMIT=$(git rev-parse HEAD)
echo "Version: $VERSION"
echo "Commit: $COMMIT"

# Create a temporary build context
BUILD_DIR=$(mktemp -d)
trap "rm -rf $BUILD_DIR" EXIT

echo "Creating build context in $BUILD_DIR..."

# Copy everything to build directory (including .git for version info)
cp -r . "$BUILD_DIR/"
cp -r ../../cometbft "$BUILD_DIR/cometbft"
cp -r ../../cosmos-sdk "$BUILD_DIR/cosmos-sdk"

# Update go.mod in build directory
cd "$BUILD_DIR"
go mod edit -replace github.com/cometbft/cometbft=./cometbft
go mod edit -replace github.com/cosmos/cosmos-sdk=./cosmos-sdk

# Create a modified Dockerfile that handles local deps
cat > Dockerfile.local << 'EOF'
ARG GOLANG_1_23_ALPINE_DIGEST="ac67716dd016429be8d4c2c53a248d7bcdf06d34127d3dc451bda6aa5a87bc06"

FROM golang@sha256:${GOLANG_1_23_ALPINE_DIGEST} as builder
ARG VERSION
ARG COMMIT

RUN set -eux; apk add --no-cache ca-certificates build-base git linux-headers bash binutils-gold

# Copy local dependencies
WORKDIR /
COPY cometbft /cometbft
COPY cosmos-sdk /cosmos-sdk

# Set up the main project
WORKDIR /dydxprotocol
COPY go.mod go.sum ./

# Update replace directives to container paths
RUN go mod edit -replace github.com/cometbft/cometbft=/cometbft && \
    go mod edit -replace github.com/cosmos/cosmos-sdk=/cosmos-sdk

# Download other dependencies
RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/root/go/pkg/mod \
    go mod download

# Copy the remaining files
COPY . .

# Build dydxprotocold binary
RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/root/go/pkg/mod \
    go build \
      -mod=readonly \
      -tags "netgo,ledger,muslc" \
      -ldflags "-X github.com/cosmos/cosmos-sdk/version.Name=dydxprotocol \
              -X github.com/cosmos/cosmos-sdk/version.AppName=dydxprotocold \
              -X github.com/cosmos/cosmos-sdk/version.Version=$VERSION \
              -X github.com/cosmos/cosmos-sdk/version.Commit=$COMMIT \
              -X github.com/cosmos/cosmos-sdk/version.BuildTags='netgo,ledger,muslc' \
              -w -s -linkmode=external -extldflags '-Wl,-z,muldefs -static'" \
      -trimpath \
      -o /dydxprotocol/build/ \
      ./...

# Runner stage
FROM golang@sha256:${GOLANG_1_23_ALPINE_DIGEST}
RUN apk add --no-cache bash
COPY --from=builder /dydxprotocol/build/dydxprotocold /bin/dydxprotocold
ENV HOME /dydxprotocol
WORKDIR $HOME
EXPOSE 26656 26657 1317 9090
ENTRYPOINT ["dydxprotocold"]
EOF

# Reset chain if needed (do this from the original directory)
echo "Resetting chain data..."
cd - > /dev/null  # Go back to original directory
make reset-chain
cd "$BUILD_DIR"  # Go back to build directory

# Build the base image with local dependencies
echo -e "${GREEN}Building Docker image with local dependencies...${NC}"
docker build \
  --build-arg VERSION=$VERSION \
  --build-arg COMMIT=$COMMIT \
  -f Dockerfile.local \
  -t dydxprotocol-base \
  --no-cache \
  .

# Build the localnet image
echo -e "${GREEN}Building localnet image...${NC}"
docker build . -t local:dydxprotocol -f testing/testnet-local/Dockerfile --no-cache

# Start the localnet (from original directory for docker-compose files)
cd - > /dev/null
echo -e "${GREEN}Starting localnet...${NC}"
docker compose -f docker-compose.yml -f docker-compose.local.yml up --force-recreate