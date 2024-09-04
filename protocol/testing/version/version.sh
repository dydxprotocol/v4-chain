#!/bin/bash
CURRENT_VERSION=$(<VERSION_CURRENT)
PREUPGRADE_VERSION=$(<VERSION_PREUPGRADE)
PREUPGRADE_VERSION_FULL_NAME=$(<VERSION_FULL_NAME_PREUPGRADE)

TESTNET_CURRENT_VERSION="v6.0.0_testnet_fix"
TESTNET_PREUPGRADE_VERSION="v6.0.0"
TESTNET_PREUPGRADE_VERSION_FULL_NAME="v6.0.2"

# Define the mapping from version to amd64 URL
declare -A version_to_url
version_to_url["$PREUPGRADE_VERSION"]="https://github.com/dydxprotocol/v4-chain/releases/download/protocol%2F$PREUPGRADE_VERSION_FULL_NAME/dydxprotocold-$PREUPGRADE_VERSION_FULL_NAME-linux-amd64.tar.gz"
declare -A testnet_version_to_url
testnet_version_to_url["$TESTNET_PREUPGRADE_VERSION"]="https://github.com/dydxprotocol/v4-chain/releases/download/protocol%2F$TESTNET_PREUPGRADE_VERSION_FULL_NAME/dydxprotocold-$TESTNET_PREUPGRADE_VERSION_FULL_NAME-linux-amd64.tar.gz"