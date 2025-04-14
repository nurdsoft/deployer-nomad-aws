package apikeys

import (
	"fmt"
	"os"
	"strings"
)

const EnvVar = "API_KEYS"

var apiKeys []string

func parseKeys() []string {
	value := os.Getenv(EnvVar)
	fmt.Printf("Raw API_KEYS from env: %q\n", value)

	if value == "" {
		return nil
	}
	value = strings.TrimSpace(value)
	list := strings.Split(value, ",")
	for i := range list {
		list[i] = strings.TrimSpace(list[i])
	}
	fmt.Printf("Parsed API keys: %+v\n", list)
	return list
}

// Have returns true if the key is in the list of API keys or authentication is disabled.
func Have(key string) bool {
	fmt.Printf("Checking if key %q exists in %+v\n", key, apiKeys)

	if len(apiKeys) == 0 {
		fmt.Println("No API keys defined — allowing request")
		return true
	}

	for _, k := range apiKeys {
		if key == k {
			fmt.Println("API key matched")
			return true
		}
	}
	fmt.Println("API key did not match")
	return false
}

func init() {
	apiKeys = parseKeys()
}