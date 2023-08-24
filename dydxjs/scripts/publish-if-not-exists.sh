#!/bin/bash
set -euxo pipefail

VERSION=$(cat package.json | jq -r '.version')
NAME=$(cat package.json | jq -r '.name')

test -z "$(npm info $NAME@$VERSION)"
if [ $? -eq 0 ]; then
	set -e
	npm publish
else
	echo "skipping publish, package $NAME@$VERSION already published"
fi
