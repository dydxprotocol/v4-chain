#!/bin/bash

# Define variables
GOGOPROTO_REPO="https://github.com/gogo/protobuf.git"
GOGOPROTO_DIR="temp_gogoproto"
REPO_URL="https://github.com/circlefin/noble-cctp.git"
REPO_DIR="temp_proto_repo"
COSMOS_SDK_REPO="https://github.com/cosmos/cosmos-sdk.git"
COSMOS_SDK_DIR="temp_cosmos_sdk"
COSMOS_PROTO_REPO="https://github.com/cosmos/cosmos-proto.git"
COSMOS_PROTO_DIR="temp_cosmos_proto"
TX_PROTO_PATH="$REPO_DIR/proto/circle/cctp/v1/tx.proto"
OUT_DIR="temp_generated" 
TARGET_FILE="src/clients/lib/cctpProto.ts"

# Create directories
mkdir -p $OUT_DIR
mkdir -p $(dirname $TARGET_FILE)

# Clone required repos
if [ ! -d "$GOGOPROTO_DIR" ]; then
    git clone $GOGOPROTO_REPO $GOGOPROTO_DIR
fi

if [ ! -d "$REPO_DIR" ]; then
    git clone $REPO_URL $REPO_DIR
fi

if [ ! -d "$COSMOS_SDK_DIR" ]; then
    git clone --depth 1 $COSMOS_SDK_REPO $COSMOS_SDK_DIR
fi

if [ ! -d "$COSMOS_PROTO_DIR" ]; then
    git clone --depth 1 $COSMOS_PROTO_REPO $COSMOS_PROTO_DIR
fi

# Generate TypeScript code using ts-proto for tx.proto
./node_modules/.bin/grpc_tools_node_protoc \
    --plugin="./node_modules/.bin/protoc-gen-ts_proto" \
    --ts_proto_out="$OUT_DIR" \
    --proto_path=./$REPO_DIR/proto/circle/cctp/v1 \
    --proto_path=./$GOGOPROTO_DIR \
    --proto_path=./$COSMOS_SDK_DIR/proto \
    --proto_path=./$COSMOS_SDK_DIR/third_party/proto \
    --proto_path=./$COSMOS_PROTO_DIR/proto \
    --ts_proto_opt="esModuleInterop=true,forceLong=long,useOptionals=messages,env=browser,globalThisPolyfill=true" \
    $TX_PROTO_PATH

# Move the generated tx.ts file to the target location
[ -f $TARGET_FILE ] && rm $TARGET_FILE
mv "$OUT_DIR/tx.ts" $TARGET_FILE

# Cleanup 
rm -rf $GOGOPROTO_DIR
rm -rf $REPO_DIR
rm -rf $COSMOS_SDK_DIR
rm -rf $COSMOS_PROTO_DIR
rm -rf $OUT_DIR