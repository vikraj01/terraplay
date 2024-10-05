#!/bin/bash

if [ -z "$1" ]; then
  echo "Usage: $0 <pem-file-path>"
  exit 1
fi

PEM_FILE=$1
ENV_FILE=".env"

if [ ! -f "$PEM_FILE" ]; then
  echo "Error: PEM file not found!"
  exit 1
fi

PEM_CONTENT=$(awk '{printf "%s\\n", $0}' "$PEM_FILE")

echo "PRIVATE_KEY=\"$PEM_CONTENT\"" >> $ENV_FILE

echo "Private key from $PEM_FILE has been added to $ENV_FILE"
