#!/bin/bash
# Builds or downloads every binary the local IBC environment needs into
# build/ibc-local/bin. dydxprotocold is always rebuilt from the working tree;
# the other tools are fetched once and reused.
set -euo pipefail
source "$(dirname "${BASH_SOURCE[0]}")/lib.sh"

require_cmd go curl tar jq bc shasum

mkdir -p "$BIN_DIR"

log "building dydxprotocold from $PROTOCOL_DIR"
(cd "$PROTOCOL_DIR" && go build -o "$DYDX_BIN" ./cmd/dydxprotocold)

if [ ! -x "$SIMD_BIN" ]; then
	# ibc-go's simd cannot be built inside the protocol module because of the
	# cometbft fork replace, and `go install pkg@version` rejects modules with
	# replace directives, so build it from the downloaded module source.
	log "building ibc-go $IBC_GO_VERSION simd"
	src=$(cd "$PROTOCOL_DIR" && go mod download -json "github.com/cosmos/ibc-go/v8@$IBC_GO_VERSION" | jq -r '.Dir')
	[ -d "$src" ] || die "could not download ibc-go $IBC_GO_VERSION"
	(cd "$src" && go build -o "$SIMD_BIN" ./testing/simapp/simd)
fi

if [ ! -x "$DASEL_BIN" ]; then
	log "installing dasel $DASEL_VERSION"
	GOBIN="$BIN_DIR" go install "github.com/tomwright/dasel/v2/cmd/dasel@$DASEL_VERSION"
fi

if [ ! -x "$HERMES_BIN" ]; then
	case "$(uname -s)-$(uname -m)" in
	Darwin-arm64) target="aarch64-apple-darwin" ;;
	Darwin-x86_64) target="x86_64-apple-darwin" ;;
	Linux-aarch64 | Linux-arm64) target="aarch64-unknown-linux-gnu" ;;
	Linux-x86_64) target="x86_64-unknown-linux-gnu" ;;
	*) die "no hermes release for $(uname -s)-$(uname -m)" ;;
	esac
	log "downloading hermes $HERMES_VERSION ($target)"
	tmp=$(mktemp -d)
	curl -fsSL -o "$tmp/hermes.tar.gz" \
		"https://github.com/informalsystems/hermes/releases/download/$HERMES_VERSION/hermes-$HERMES_VERSION-$target.tar.gz"
	tar -xzf "$tmp/hermes.tar.gz" -C "$tmp"
	mv "$tmp/hermes" "$HERMES_BIN"
	command rm -rf "$tmp"
fi

"$HERMES_BIN" version >/dev/null || die "hermes binary at $HERMES_BIN does not run"
log "binaries ready in $BIN_DIR"
