#!/bin/bash

# Script to show the orderbook levels for a ticker

# Load NVM and use the correct Node.js version
export NVM_DIR="$HOME/.nvm"
[ -s "$NVM_DIR/nvm.sh" ] && \. "$NVM_DIR/nvm.sh"
[ -s "$NVM_DIR/bash_completion" ] && \. "$NVM_DIR/bash_completion"

# Get the directory where this script is located
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )"

# Change to the script directory
cd "$SCRIPT_DIR"

# Load the correct Node.js version
cd ../.. && nvm use && cd services/scripts || exit 1

# Parse arguments
TICKER=""
PRICE=""
ALL=""
MIN_SIZE=""
SHOW_ZEROS=""
REDIS_URL="redis://master.redis-debug-backup.g88kxp.apne1.cache.amazonaws.com:6379"

# Process arguments
while [[ $# -gt 0 ]]; do
  case "$1" in
    --price|-p)
      PRICE="$2"
      shift 2
      ;;
    --all|-a)
      ALL="--all"
      shift
      ;;
    --min-size|-m)
      MIN_SIZE="--minSize $2"
      shift 2
      ;;
    --show-zeros|-z)
      SHOW_ZEROS="--showZeros"
      shift
      ;;
    --redis-url|-r)
      REDIS_URL="$2"
      shift 2
      ;;
    *)
      # First non-flag argument is the ticker
      if [ -z "$TICKER" ]; then
        TICKER="$1"
      fi
      shift
      ;;
  esac
done

# Check if ticker is provided
if [ -z "$TICKER" ]; then
  echo "Usage: ./show-orderbook.sh <ticker> [--price <price>] [--all] [--min-size <size>] [--show-zeros] [--redis-url <url>]"
  echo "Example: ./show-orderbook.sh SOL-USD"
  echo "Example: ./show-orderbook.sh SOL-USD --price 159.29 --all"
  echo "Example: ./show-orderbook.sh SOL-USD --min-size 1"
  echo "Example: ./show-orderbook.sh SOL-USD --show-zeros"
  exit 1
fi

# Add TLS prefix if not present and it's using the default Redis URL
if [[ "$REDIS_URL" == *"master.redis-debug-backup.g88kxp.apne1.cache.amazonaws.com"* ]] && [[ "$REDIS_URL" != *"rediss://"* ]]; then
  REDIS_URL="rediss://${REDIS_URL#redis://}"
fi

# Build price argument if provided
PRICE_ARG=""
if [ -n "$PRICE" ]; then
  PRICE_ARG="--price $PRICE"
fi

# Display info
echo "Showing orderbook levels for $TICKER from Redis at $REDIS_URL"

# Run the script with all parameters
NODE_ENV=development npx ts-node src/show-orderbook-levels.ts \
  --ticker "$TICKER" \
  $PRICE_ARG \
  $ALL \
  $MIN_SIZE \
  $SHOW_ZEROS \
  --redisUrl "$REDIS_URL" 