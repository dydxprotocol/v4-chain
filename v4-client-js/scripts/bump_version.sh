#!/bin/sh


# search for the first line that starts with "version" in package.json
# get the value in the quotes
VERSION=$(cat package.json | jq -r '.version')

echo "Current version is $VERSION. Enter new version (or press enter to skip):"
read NEW_VERSION

#if NEW_VERSION is not empty, replace the version in package.json
if [ -n "$NEW_VERSION" ]; then
  sed -i '' "s/  \"version\": \"$VERSION\"/  \"version\": \"$NEW_VERSION\"/" package.json
  echo "Version bumped to $NEW_VERSION"
  npm i
fi
