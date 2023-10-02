#!/bin/sh
set -x

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
}

if [[ "$1" == "-h" || "$1" == "--help" ]]; then
    display_usage
    exit 0
fi


CHAIN_ID="dydxprotocol-testnet"
# local path to temporary snapshots. snapshots are deleted after uploading to S3.
SNAP_PATH="/dydxprotocol/chain/local_node/snapshots/dydxprotocol/"
# logfile containing snapshot timestamps
LOG_PATH="/dydxprotocol/chain/local_node/snapshots/dydxprotocol/dydxprotocol_log.txt"
# data directory to snapshot. this contains the blockchain state.
DATA_PATH="/dydxprotocol/chain/local_node/data/"
RPC_ADDRESS="http://127.0.0.1:26657"

while [ $# -gt 0 ]; do

   if [[ $1 == *"--"* ]]; then
        v="${1/--/}"
        export $v="$2"
   fi

  shift
done


now_date() {
    echo -n $(TZ="UTC" date '+%Y-%m-%d_%H:%M:%S')
}

install_prerequisites() {
    apk add dasel jq curl
    apk add --no-cache \
        python3 \
        py3-pip \
    && pip3 install --upgrade pip \
    && pip3 install --no-cache-dir \
        awscli \
    && rm -rf /var/cache/apk/*
}

setup_cosmovisor() {
    VAL_HOME_DIR="$HOME/chain/local_node"
    export DAEMON_NAME=dydxprotocold
    export DAEMON_HOME="$HOME/chain/local_node"

    cosmovisor init /bin/dydxprotocold
}

install_prerequisites

# log messages prefixed by highlighted timestamp
log_this() {
    YEL='\033[1;33m' # yellow
    NC='\033[0m'     # No Color
    local logging="$@"
    printf "|$(now_date)| $logging\n" | tee -a ${LOG_PATH}
}

# initialize snapshot path and genesis file
mkdir -p $SNAP_PATH
touch $LOG_PATH
sleep 10
dydxprotocold init --chain-id=${CHAIN_ID} --home /dydxprotocol/chain/local_node local_node
curl -X GET ${genesis_file_rpc_address}/genesis | jq '.result.genesis' > /dydxprotocol/chain/local_node/config/genesis.json

setup_cosmovisor

# TODO: add metrics around snapshot upload latency/frequency/success rate
while true; do
  # p2p.seeds taken from --p2p.persistent_peers flag of full node
  cosmovisor run start --log_level info --home /dydxprotocol/chain/local_node --p2p.seeds "${p2p_seeds}" --non-validating-full-node=true &

  sleep ${upload_period}
  kill -TERM $(pidof cosmovisor)

  log_this "Creating new snapshot"
  SNAP_NAME=$(echo "${CHAIN_ID}_$(date '+%Y-%m-%d-%M-%H').tar.gz")
  tar cvzf ${SNAP_PATH}/${SNAP_NAME} ${DATA_PATH}
  aws s3 cp ${SNAP_PATH}/${SNAP_NAME} s3://${s3_snapshot_bucket}/ --region ap-northeast-1 --debug || true
  rm -rf ${SNAP_PATH}/${SNAP_NAME}
  log_this "Done creating snapshot\n---------------------------\n"

done
