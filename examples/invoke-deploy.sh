#!/bin/bash

API_ID=""
API_KEY=""
FILE="${1:-examples/deploy-batch-job.hcl}"  # Pass file path or default to batch

curl -X POST "https://${API_ID}.execute-api.us-west-2.amazonaws.com/v1/deploy" \
  -H "x-api-key: ${API_KEY}" \
  --data-binary @"${FILE}"