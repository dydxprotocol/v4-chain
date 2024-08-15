#!/bin/bash

# Note: Don't use -x to avoid logging secrets.
set -euo pipefail

if [[ ! -z "${SECRET_ID:-}" ]]
then
  # Check if the secret exists
  secrets=$(aws secretsmanager list-secrets --filters Key="name",Values="$SECRET_ID" | jq -r ".SecretList[].Name")
  # If secrets are found, set environment variables, otherwise use defaults
  if [[ ! -z "${secrets}" ]]
  then
    # Retrieve the secret from secrets manager if it exists
    values=$(aws secretsmanager get-secret-value --secret-id $SECRET_ID | jq -r ".SecretString")
      # Parse the JSON secret into key-value pairs
    env_vars=$(echo "$values" | jq -r 'to_entries | map(
      if .key == "FIREBASE_PRIVATE_KEY" then
        "\(.key)=\(.value|tostring)"
      else
        "\(.key)=\(.value|tostring|gsub("\n"; " "))"
      end
    ) | .[]')

    # Log the environment variables to the console
    echo "Environment Variables:"
    echo "$env_vars"

    while read line; do
      export "$line"
    done <<< "$env_vars"
  fi
fi

exec pnpm run start
