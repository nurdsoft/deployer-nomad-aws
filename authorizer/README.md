# Lambda Authorizer for Nomad Deploy API

This Lambda function is used as a [Lambda Authorizer](https://docs.aws.amazon.com/apigateway/latest/developerguide/http-api-lambda-authorizer.html) to protect the `/deploy` route of the Nomad deploy API.

## How It Works

- Clients must include an `x-api-key` header.
- The authorizer compares the value against `VALID_API_KEY`, set via environment variable.
- If it matches, access is granted. Otherwise, API Gateway returns `403 Forbidden`.

## Deployment Notes

- The file is zipped as `build/authorizer.zip` via the Makefile.
- The Lambda uses Python 3.11 and is defined in `lambda_authorizer.py`.

## Example Header

```bash
curl -X POST https://api-id.execute-api.us-west-2.amazonaws.com/v1/deploy \
  -H "x-api-key: your-key-here" \
  -H "Content-Type: application/json" \
  -d '{"app":"nomad"}'