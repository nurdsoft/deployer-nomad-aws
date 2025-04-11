package apikeys

import (
	"os"
	"strings"
)

const EnvVar = "API_KEYS"

var apiKeys []string

func parseKeys() []string {
	value := os.Getenv(EnvVar)
	if value == "" {
		return nil
	}
	strings.TrimSpace(value)
	list := strings.Split(value, ",")
	if len(list) == 0 {
		return nil
	}
	return list
}

// Have returns true if the key is in the list of API keys or authentication
// is disabled (i.e. no keys are set).
func Have(key string) bool {
	// Return true if no keys are set (i.e. authentication is disabled)
	// We use the length rather than nil becuase this will be passed in as an env.
	// variable and will be an empty list rather than nil.
	if len(apiKeys) == 0 {
		return true
	}

	for _, k := range apiKeys {
		if key == k {
			return true
		}
	}
	return false
}

func init() {
	apiKeys = parseKeys()
}
