#!/bin/bash
set -euxo pipefail

VERSION=$(cat package.json | jq -r '.version')
NAME=$(cat package.json | jq -r '.name')

test -z "$(npm info $NAME@$VERSION)"
if [ $? -eq 0 ]; then
	set -e

	git config --global user.email "ci@dydx.exchange"
	git config --global user.name "github_actions"

	# Get version and tag
	git tag v4-client-js@${VERSION}
	git push --tags

	npm publish
else
	echo "skipping publish, package $NAME@$VERSION already published"
fi
