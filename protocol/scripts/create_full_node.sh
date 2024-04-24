#!/bin/bash

# version needs to be changed on new releases
VERSION="v4.0.5"
WORKDIR=$HOME/dydxprotocol
CHAIN_ID="dydx-mainnet-1"
DATA_DIR=$WORKDIR/data
SEED_NODES=("ade4d8bc8cbe014af6ebdf3cb7b1e9ad36f412c0@seeds.polkachu.com:23856", 
"65b740ee326c9260c30af1f044e9cda63c73f7c1@seeds.kingnodes.net:23856", 
"f04a77b92d0d86725cdb2d6b7a7eb0eda8c27089@dydx-mainnet-seed.bwarelabs.com:36656",
"20e1000e88125698264454a884812746c2eb4807@seeds.lavenderfive.com:23856",
"c2c2fcb5e6e4755e06b83b499aff93e97282f8e8@tenderseed.ccvalidators.com:26401",
"a9cae4047d5c34772442322b10ef5600d8e54900@dydx-mainnet-seednode.allthatnode.com:26656",
"802607c6db8148b0c68c8a9ec1a86fd3ba606af6@64.227.38.88:26656",
"ebc272824924ea1a27ea3183dd0b9ba713494f83@dydx-mainnet-seed.autostake.com:27366"
)
BASE_SNAPSHOT_URL="https://snapshots.polkachu.com/snapshots/"
# Add random number to the end of name to keep a unique name for the node
NODE_NAME="FullNodeMainnet-$RANDOM" 


if [ $(nproc) -lt 8 ]; then
  echo "This device has less than 8 cpus, recommended to use a device with at least 8 cpus"
fi

if [ $(free -g | grep Mem | awk '{print $2}') -lt 64 ]; then
  echo "This device has lass than 64 gigs of ram, recommended to use a device with at least 64 gigs of ram"
fi

# Check if dydxprotocold is not installed
if [ ! $(which dydxprotocold) ]; then
  echo "dydxprotocold is not installed, installing dydxprotocold"
    mkdir -p $WORKDIR
    cd $WORKDIR

    # check arch of the device
    if [ $(uname -m) = "aarch64" ]; then
    ARCH=arm64
    elif [ $(uname -m) = "x86_64" ]; then
    ARCH=amd64
    else
    echo "This device is not arm64 or amd64, please use a device with arm64 or amd64 architecture"
    fi
    FILE_REGEX="dydxprotocold-v[0-9]+\.[0-9]+\.[0-9]+-linux-$ARCH.tar.gz"

    #Download latest github release for repo
    curl -s "https://api.github.com/repos/dydxprotocol/v4-chain/releases/tags/protocol/$VERSION" \
    | grep -E $FILE_REGEX \
    | cut -d : -f 2,3 \
    | tr -d \" \
    | wget -qi -

    tar -xvf dydxprotocold-v*-linux-$ARCH.tar.gz
    rm dydxprotocold-v*-linux-$ARCH.tar.gz
    cd $WORKDIR/build
    mv dydxprotocold-*-linux-$ARCH dydxprotocold
    chmod +x dydxprotocold

    # check if is not in path
    if [ ! $(echo $PATH | grep "$HOME/.local/bin") ]; then
        # if .local/bin does not exist, create it
        if [ ! -d "$HOME/.local/bin" ]; then
        mkdir -p $HOME/.local/bin
        fi
        echo "export PATH=$HOME/.local/bin:$PATH" >> $HOME/.bashrc
        eval "$(cat $HOME/.bashrc | tail -n +10)"
    fi
    cp dydxprotocold $HOME/.local/bin
    rm -rf $WORKDIR/build

    mkdir $DATA_DIR 
    dydxprotocold init --chain-id=$CHAIN_ID --home=$DATA_DIR $NODE_NAME
fi

# Move the priv_validator_key.json to the workdir while getting the snapshot
mv $DATA_DIR/config/priv_validator_key.json $WORKDIR

# Get snapshot
cd $DATA_DIR
xml_content=$(curl -s "${BASE_SNAPSHOT_URL}")
file_key=$(echo "$xml_content" | grep -oP '<Key>dydx/dydx_\d+\.tar\.lz4</Key>' | head -1 | sed 's/<[^>]*>//g')
file_url="${BASE_SNAPSHOT_URL}${file_key}"
wget -O "snapshot.tar.lz4" "${file_url}"
lz4 -c -d snapshot.tar.lz4 | tar -x -C $DATA_DIR
rm snapshot.tar.lz4

# Move the priv_validator_key.json back to the config folder
mv $WORKDIR/priv_validator_key.json $DATA_DIR/config

# Get genesis file
curl https://dydx-rpc.lavenderfive.com/genesis | python3 -c 'import json,sys;print(json.dumps(json.load(sys.stdin)["result"]["genesis"], indent=2))' > $DATA_DIR/config/genesis.json

dydxprotocold start --p2p.seeds="${SEED_NODES[*]}" --home=$DATA_DIR --non-validating-full-node=true