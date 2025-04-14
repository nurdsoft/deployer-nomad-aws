import os

def handler(event, context):
    """
    Lambda authorizer that checks x-api-key against a comma-separated list of
    valid keys passed via environment variable API_KEYS.
    """
    print("Received event:", event)

    allowed_keys = os.environ.get("API_KEYS", "").split(",")
    request_key = event.get("headers", {}).get("x-api-key")

    if request_key and request_key in allowed_keys:
        return { "isAuthorized": True }

    return { "isAuthorized": False }