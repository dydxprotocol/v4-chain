#!/bin/bash

# Name of the output zip file
OUTPUT_ZIP="gateway-demo.zip"

# Remove existing zip if it exists
if [ -f "$OUTPUT_ZIP" ]; then
    rm "$OUTPUT_ZIP"
    echo "Removed existing $OUTPUT_ZIP"
fi

echo "Zipping project files..."

# Create the zip file excluding unnecessary directories and files
# -r: recursive
# -x: exclude patterns
zip -r "$OUTPUT_ZIP" main.go PerpEngine.go

echo "Successfully created $OUTPUT_ZIP"
ls -lh "$OUTPUT_ZIP"
