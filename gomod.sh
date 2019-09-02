#!/bin/bash
set -eu

touch go.mod

PROJECT_NAME=EwanValentine
CURRENT_DIR=$(basename $(pwd))

CONTENT=$(cat <<-EOD
module github.com/${PROJECT_NAME}/${CURRENT_DIR}

require github.com/aws/aws-lambda-go v1.6.0
EOD
)

echo "$CONTENT" > go.mod
