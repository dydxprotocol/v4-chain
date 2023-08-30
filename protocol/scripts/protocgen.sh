#!/usr/bin/env bash

set -eo pipefail

echo "Generating gogo proto code"
cd ./proto
proto_dirs=$(find ./dydxprotocol -path -prune -o -name '*.proto' -print0 | xargs -0 -n1 dirname | sort | uniq)
for dir in $proto_dirs; do
	for file in $(find "${dir}" -maxdepth 1 -name '*.proto'); do
		if grep "option go_package" $file &>/dev/null; then
			buf generate --template buf.gen.gogo.yaml $file
		fi
	done
done

cd ..

# move proto files to the right places
find . -name "*.pb.go" -o -name "*.pb.gw.go" -type f -not -path "./proto/*" -delete
cp -r proto/.gen/github.com/dydxprotocol/v4-chain/protocol/* ./protocol/
rm -rf proto/.gen/github.com/

cd protocol && go mod tidy
