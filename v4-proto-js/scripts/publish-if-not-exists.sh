#!/bin/bash
set -euxo pipefail

VERSION=$(cat package.json | jq -r '.version')
NAME=$(cat package.json | jq -r '.name')

# Attempt to extract the pre-release tag
if [[ $VERSION =~ ^[[:digit:]]+.[[:digit:]]+.[[:digit:]]+-([[:alpha:]]+).[[:digit:]]+$ ]]; then
	TAG=${BASH_REMATCH[1]}
else
	TAG=""
fi

test -z "$(npm info $NAME@$VERSION)"
if [ $? -eq 0 ]; then
	set -e
	# Publish with pre-release tag if it exists. Otherwise publish as latest (default tag)
	if [ -z "$TAG" ]; then
		npm publish
	else
		npm publish --tag $TAG
	fi
else
	echo "skipping publish, package $NAME@$VERSION already published"
fi
