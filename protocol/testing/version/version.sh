CURRENT_VERSION=$(<VERSION_CURRENT)
PREUPGRADE_VERSION=$(<VERSION_PREUPGRADE)

# Define the mapping from version to amd64 URL
declare -A version_to_url
version_to_url["$PREUPGRADE_VERSION"]="https://github.com/dydxprotocol/v4-chain/releases/download/protocol%2F$PREUPGRADE_VERSION/dydxprotocold-$PREUPGRADE_VERSION-linux-amd64.tar.gz"
