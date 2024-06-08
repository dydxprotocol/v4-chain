#!/bin/bash

# Define variables
GOGOPROTO_REPO="https://github.com/gogo/protobuf.git"
GOGOPROTO_DIR="temp_gogoproto"
REPO_URL="https://github.com/circlefin/noble-cctp.git"
REPO_DIR="temp_proto_repo"
TX_PROTO_PATH="$REPO_DIR/proto/circle/cctp/v1/tx.proto"
OUT_DIR="temp_enerated" 
TARGET_FILE="src/clients/lib/cctpProto.ts"

# Create the output directory if it does not exist
mkdir -p $OUT_DIR

# Clone gogoproto if it doesn't exist
if [ ! -d "$GOGOPROTO_DIR" ]; then
    git clone $GOGOPROTO_REPO $GOGOPROTO_DIR
fi

# Clone the repository if it doesn't exist
if [ ! -d "$REPO_DIR" ]; then
    git clone $REPO_URL $REPO_DIR
fi

# Generate TypeScript code using ts-proto for tx.proto
./node_modules/.bin/grpc_tools_node_protoc \
    --plugin="./node_modules/.bin/protoc-gen-ts_proto" \
    --ts_proto_out="$OUT_DIR" \
    --proto_path=./$REPO_DIR/proto/circle/cctp/v1 \
    --proto_path=./$GOGOPROTO_DIR \
    --ts_proto_opt="esModuleInterop=true,forceLong=long,useOptionals=messages,env=browser,globalThisPolyfill=true" \
    $TX_PROTO_PATH

# Move the generated tx.ts file to the target location
rm $TARGET_FILE
mv "$OUT_DIR/tx.ts" $TARGET_FILE

# Cleanup 
rm -rf $GOGOPROTO_DIR
rm -rf $REPO_DIR
rm -rf $OUT_DIR
