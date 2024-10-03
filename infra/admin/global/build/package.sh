#!/bin/bash

FUNCTION_NAME="lambda"
OUTPUT_FILE="${FUNCTION_NAME}.zip"

if [ -f $OUTPUT_FILE ]; then
    rm $OUTPUT_FILE
fi

zip -r $OUTPUT_FILE . -x "*.sh"

if [ -f $OUTPUT_FILE ]; then
    echo "$OUTPUT_FILE created successfully."
else
    echo "Failed to create $OUTPUT_FILE."
fi
