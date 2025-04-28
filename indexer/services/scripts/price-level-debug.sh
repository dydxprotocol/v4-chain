#!/bin/bash

# Script to analyze a price level in the orderbook

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

# Parse command line arguments
TICKER=""
SIDE=""
PRICE=""
ANALYZE_ALL=false
# REDIS_URL="redis://master.redis-debug-backup.g88kxp.apne1.cache.amazonaws.com:6379"
#REDIS_URL="redis://master.redis-debug-0426.g88kxp.apne1.cache.amazonaws.com:6379"
#REDIS_URL="redis://master.redis-debug-0427-1059pm.g88kxp.apne1.cache.amazonaws.com:6379"
REDIS_URL="redis://master.redis-debug-0428.g88kxp.apne1.cache.amazonaws.com:6379"

# Process arguments
while [[ $# -gt 0 ]]; do
  case "$1" in
    --analyze-all|-a)
      ANALYZE_ALL=true
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
        shift
      # Second non-flag argument is the side
      elif [ -z "$SIDE" ]; then
        SIDE="$1"
        shift
      # Third non-flag argument is the price (if not using --analyze-all)
      elif [ -z "$PRICE" ] && [ "$ANALYZE_ALL" = false ]; then
        PRICE="$1"
        shift
      else
        shift
      fi
      ;;
  esac
done

# Check if required parameters are provided
if [ -z "$TICKER" ] || [ -z "$SIDE" ]; then
  echo "Usage: ./price-level-debug.sh <ticker> <side> [price] [--analyze-all] [--redis-url <url>]"
  echo "Example: ./price-level-debug.sh SOL-USD SELL 159.29"
  echo "Example: ./price-level-debug.sh SOL-USD SELL --analyze-all"
  exit 1
fi

# Check if price is provided when not using analyze-all
if [ "$ANALYZE_ALL" = false ] && [ -z "$PRICE" ]; then
  echo "Error: Price must be provided when not using --analyze-all"
  echo "Example: ./price-level-debug.sh SOL-USD SELL 159.29"
  echo "Or use: ./price-level-debug.sh SOL-USD SELL --analyze-all"
  exit 1
fi

# Add TLS prefix if not present and it's using the default Redis URL
if [[ "$REDIS_URL" == *"master.redis-debug-backup.g88kxp.apne1.cache.amazonaws.com"* ]] && [[ "$REDIS_URL" != *"rediss://"* ]]; then
  REDIS_URL="rediss://${REDIS_URL#redis://}"
fi

# Display info
if [ "$ANALYZE_ALL" = true ]; then
  echo "Analyzing all price levels for $TICKER $SIDE from Redis at $REDIS_URL"
else
  echo "Analyzing price level $PRICE for $TICKER $SIDE from Redis at $REDIS_URL"
fi

# Build the command arguments
ARGS=""
if [ "$ANALYZE_ALL" = true ]; then
  ARGS="--analyzeAll"
else
  ARGS="--price $PRICE"
fi

# Run the script
NODE_ENV=development npx ts-node src/print-price-level-detail.ts \
  --ticker "$TICKER" \
  --side "$SIDE" \
  $ARGS \
  --redisUrl "$REDIS_URL"

# After analysis, offer usage hints
echo ""
if [ "$ANALYZE_ALL" = true ]; then
  echo "To view the entire orderbook, run: ./show-orderbook.sh $TICKER"
else
  echo "To view the entire orderbook, run: ./show-orderbook.sh $TICKER"
  echo "To view a specific price range, run: ./show-orderbook.sh $TICKER --price $PRICE --all"
  echo "To view only levels with orders, run: ./show-orderbook.sh $TICKER --min-size 1"
  echo "To analyze all price levels, run: ./price-level-debug.sh $TICKER $SIDE --analyze-all"
fi 