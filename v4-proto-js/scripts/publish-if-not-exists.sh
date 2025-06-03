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

# Skip publishing if version already exists
set +e
if npm info "$NAME@$VERSION" --loglevel verbose >/dev/null 2>&1; then
  echo "Skipping publish as $NAME@$VERSION already exists on npm."
  exit 0
fi
set -e

if [ -z "$TAG" ]; then
  echo "Running: npm publish --loglevel verbose"
  npm publish --loglevel verbose
else
  echo "Running: npm publish --tag $TAG --loglevel verbose"
  npm publish --tag "$TAG" --loglevel verbose
fi

echo "Success: Published $NAME@$VERSION."
