#!/bin/bash

# This script is used to create a non validating full node for the dydx protocol. 
# It will install all the necessary dependencies, download the dydxprotocold binary, 
# initialize the node, create a service for the node, and download a snapshot to speed up the syncing process.


# Function which reads in y or exits if anything else is entered
function read_yes_or_exit {
  read -p "$1" answer
  if [ "$answer" != "y" ]; then
    echo "Exiting"
    exit 1
  fi
}

if [ $(nproc) -lt 8 ]; then
  echo "This device has less than 8 cpus, recommended to use a device with at least 8 cpus"
  # check if user wants to continue using y/n
  echo "Do you want to continue? (y)"
  read_yes_or_exit
fi

if [ $(free -g | grep Mem | awk '{print $2}') -lt 64 ]; then
  echo "This device has lass than 64 gigs of ram, recommended to use a device with at least 64 gigs of ram"
  echo "Do you want to continue? (y)"
  read_yes_or_exit
fi

VERSION="v4.1.0"
WORKDIR=$HOME/.dydxprotocol
CHAIN_ID="dydx-mainnet-1"
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

# Update system and install dependencies

sudo apt-get -y update
sudo apt-get install -y curl jq lz4


# Get architecture
if [ $(uname -m) = "aarch64" ]; then
ARCH=arm64
elif [ $(uname -m) = "x86_64" ]; then
ARCH=amd64
else
echo "This device is not arm64 or amd64, please use a device with arm64 or amd64 architecture"
exit 1
fi

# check if go is installed and install if not
if ! command -v go &> /dev/null; then
  echo "Go is not installed, installing go"
  wget https://golang.org/dl/go1.22.2.linux-$ARCH.tar.gz 
  sudo tar -C /usr/local -xzf go1.22.2.linux-$ARCH.tar.gz
  rm go1.22.2.linux-$ARCH.tar.gz
  echo 'export PATH=$PATH:/usr/local/go/bin:$HOME/go/bin' >> $HOME/.bashrc
  eval "$(cat $HOME/.bashrc | tail -n +10)"
else
  echo "Go is installed"
fi

# Install Cosmovisor
go install cosmossdk.io/tools/cosmovisor/cmd/cosmovisor@v1.5.0

# Download dydxprotocold binaries and initialize node
mkdir -p $WORKDIR
cd $WORKDIR
FILE_REGEX="dydxprotocold-v[0-9]+\.[0-9]+\.[0-9]+-linux-$ARCH.tar.gz"
curl -s "https://api.github.com/repos/dydxprotocol/v4-chain/releases/tags/protocol/$VERSION" \
| grep -E $FILE_REGEX \
| cut -d : -f 2,3 \
| tr -d \" \
| wget -qi -
tar -xvf dydxprotocold-v*-linux-$ARCH.tar.gz && rm dydxprotocold-v*-linux-$ARCH.tar.gz
cd $WORKDIR/build
mv dydxprotocold-*-linux-$ARCH dydxprotocold
chmod +x dydxprotocold
./dydxprotocold init --chain-id=$CHAIN_ID $NODE_NAME

# Create cosmovisor directories and move binaries
mkdir -p $HOME/.dydxprotocol/cosmovisor/genesis/bin
mkdir -p $HOME/.dydxprotocol/cosmovisor/upgrades
mv dydxprotocold $HOME/.dydxprotocol/cosmovisor/genesis/bin/

# Update config
curl https://dydx-rpc.lavenderfive.com/genesis | python3 -c 'import json,sys;print(json.dumps(json.load(sys.stdin)["result"]["genesis"], indent=2))' > $WORKDIR/config/genesis.json
sed -i 's/seeds = ""/seeds = "'"${SEED_NODES[*]}"'"/' $WORKDIR/config/config.toml

# Create Service
sudo tee /etc/systemd/system/dydxprotocold.service > /dev/null << EOF
[Unit]
Description=dydx node service
After=network-online.target

[Service]
User=$USER
ExecStart=/$HOME/go/bin/cosmovisor run start --non-validating-full-node=true
WorkingDirectory=$HOME/.dydxprotocol
Restart=always
RestartSec=5
LimitNOFILE=4096
Environment="DAEMON_HOME=$HOME/.dydxprotocol"
Environment="DAEMON_NAME=dydxprotocold"
Environment="DAEMON_ALLOW_DOWNLOAD_BINARIES=false"
Environment="DAEMON_RESTART_AFTER_UPGRADE=true"
Environment="UNSAFE_SKIP_BACKUP=true"

[Install]
WantedBy=multi-user.target
EOF

sudo systemctl daemon-reload
sudo systemctl enable dydxprotocold


# Get snapshot
cd $WORKDIR
xml_content=$(curl -s "${BASE_SNAPSHOT_URL}")
file_key=$(echo "$xml_content" | grep -oP '<Key>dydx/dydx_\d+\.tar\.lz4</Key>' | head -1 | sed 's/<[^>]*>//g')
file_url="${BASE_SNAPSHOT_URL}${file_key}"
cp $WORKDIR/data/priv_validator_state.json $WORKDIR/priv_validator_state.json.backup
rm -rf $WORKDIR/data
wget -O "snapshot.tar.lz4" "${file_url}"
lz4 -c -d snapshot.tar.lz4 | tar -x -C $WORKDIR
mv $WORKDIR/priv_validator_state.json.backup $WORKDIR/data/priv_validator_state.json
rm snapshot.tar.lz4


echo "Full node setup complete"
echo "To start the node run 'sudo systemctl start dydxprotocold'"
echo "To stop the node run 'sudo systemctl stop dydxprotocold'"