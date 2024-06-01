#!/bin/bash

jsonFilePath='./genesis.json' # Replace with the path to your JSON file


if [ -f "$jsonFilePath" ]; then
  jsonString=$(jq -c . "$jsonFilePath")

  # Function to escape JSON values
  escape_value() {
    echo "$1" | sed 's/\\/\\\\/g' | sed 's/"/\\"/g'
  }

  # Function to recursively process JSON and generate paths with values
  process_json() {
    local json=$1
    local prefix=$2

    for key in $(echo "$json" | jq -r 'keys_unsorted[]'); do
      local value=$(echo "$json" | jq -c --arg key "$key" '.[$key]')
      local escapedKey=$(echo "$key" | sed 's/\\/\\\\/g' | sed 's/"/\\"/g')

      if [[ "$value" == \{* ]]; then
        # Recursively process nested objects
        process_json "$value" "$prefix.$escapedKey"
      elif [[ "$value" == \[* ]]; then
        # Process arrays
        array_index=0
        for item in $(echo "$value" | jq -c '.[]'); do
          local escapedItem=$(escape_value "$item")
          if [[ "$item" == \{* ]]; then
            process_json "$item" "$prefix.$escapedKey[$array_index]"
          else
            if [[ "$item" != "true" ]] && [[ "$item" != "false" ]]; then
              escapedItem="\\\\\\\"${escapedItem}\\\\\\\""
            fi
            result="${result}${prefix}.$escapedKey[$array_index] = ${escapedItem} | "
          fi
          array_index=$((array_index + 1))
        done
      else
        # Quote non-boolean values
        if [[ "$value" != "true" ]] && [[ "$value" != "false" ]]; then
          escapedValue=$(escape_value "$value")
          escapedValue="\\\\\\\"${escapedValue}\\\\\\\""
        else
          escapedValue=$(escape_value "$value")
        fi
        # Append to result
        result="${result}${prefix}.${escapedKey} = ${escapedValue} | "
      fi
    done
  }

  # Initialize result string
  result=""

  # Start processing from the root object
  root_object=$(echo "$jsonString" | jq -c '.app_state')
  process_json "$root_object" ".app_state"

  # Remove the trailing " | "
  result=${result% | }

  echo "\"${result}\""
else
  echo "JSON file not found!"
fi