#!/bin/bash
set -xeo pipefail

# This script spins up a non-validating full node that periodically is halted
# and uploads snapshots of the data directory (contains all blocks) to S3.
# Snapshots can be used to start a new full-node from the height the snapshot
# was taken at.

# Example usage: ./snapshot.sh --s3_snapshot_bucket dev4fullnodesnapshots \
# --genesis_file_rpc_address http://18.178.88.89:26657 \
# --p2p_seeds dfa67970296bbecce14daba6cb0da516ed60458a@3.129.102.24:26656 \
# --upload_period 300

# Display usage information
function display_usage() {
    echo "Usage: ./snapshot.sh [options]"
    echo "Options:"
    echo "  --genesis_file_rpc_address      RPC address of a validator node to retrieve the genesis file from, e.g. http://18.178.88.89:26657"
    echo "  --p2p_seeds                     List of seed nodes to peer with, e.g. dfa67970296bbecce14daba6cb0da516ed60458a@3.129.102.24:26656"
    echo "  --upload_period                 Upload frequency in seconds, e.g. 300"
    echo "  --s3_snapshot_bucket            Name of the S3 bucket to upload snapshots to, e.g. dev4fullnodesnapshots"
    echo "  --dd_agent_host                 Datadog agent host"
}

if [[ "$1" == "-h" || "$1" == "--help" ]]; then
    display_usage
    exit 0
fi

# Copy the correct binaries into the full node home directories.
source "./vars.sh"

# Set up CosmosVisor for full-nodes.
for i in $(seq 0 $LAST_FULL_NODE_INDEX); do
	FULL_NODE_HOME_DIR="$HOME/chain/.full-node-$i"
	# Copy binaries for `cosmovisor` from the docker image into the home directory.
	# Work-around to ensure docker volume contains the same binaries as the git repo.
	cp -r "$HOME/cosmovisor" "$FULL_NODE_HOME_DIR/"
done

install_prerequisites() {
    apk add dasel jq curl
}

install_prerequisites


# local path to temporary snapshots. snapshots are deleted after uploading to S3.
SNAP_PATH="/dydxprotocol/chain/.full-node-0/snapshots/dydxprotocol/"
# logfile containing snapshot timestamps
LOG_PATH="/dydxprotocol/chain/.full-node-0/snapshots/dydxprotocol/dydxprotocol_log.txt"
# data directory to snapshot. this contains the blockchain state.
DATA_PATH="/dydxprotocol/chain/.full-node-0/data/"
RPC_ADDRESS="http://127.0.0.1:26657"

while [ $# -gt 0 ]; do

   if [[ $1 == *"--"* ]]; then
        v="${1/--/}"
        export $v="$2"
   fi

  shift
done


# initialize snapshot path and genesis file
mkdir -p $SNAP_PATH
touch $LOG_PATH
sleep 10
CHAIN_ID="dydx-testnet-1"

# Prune snapshots to prevent them from getting too big. We make 3 changes:
# Prune all app state except last 2 blocks
sed -i 's/pruning = "default"/pruning = "everything"/' /dydxprotocol/chain/.full-node-0/config/app.toml
# Tendermint pruning is decided by picking the most restrictive of multiple factors.
# Make the custom config setting as permissive as possible.
sed -i 's/min-retain-blocks = 0/min-retain-blocks = 2/' /dydxprotocol/chain/.full-node-0/config/app.toml
# Do not index tx_index.db
sed -i 's/indexer = "kv"/indexer = "null"/' /dydxprotocol/chain/.full-node-0/config/config.toml

# TODO: add metrics around snapshot upload latency/frequency/success rate
while true; do
  # p2p.seeds taken from --p2p.persistent_peers flag of full node
  cosmovisor run start --log_level info --home /dydxprotocol/chain/.full-node-0 --p2p.seeds "${p2p_seeds}" --non-validating-full-node=true --dd-agent-host=${dd_agent_host} &

  sleep ${upload_period}
  kill -TERM $(pidof cosmovisor)

  SNAP_NAME=$(echo "${CHAIN_ID}_$(date '+%Y-%m-%d-%H-%M').tar.gz")
  tar cvzf ${SNAP_PATH}/${SNAP_NAME} ${DATA_PATH}
  aws s3 cp ${SNAP_PATH}/${SNAP_NAME} s3://${s3_snapshot_bucket}/ --region ap-northeast-1 --debug || true
  rm -rf ${SNAP_PATH}/${SNAP_NAME}

done