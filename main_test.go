package main

import (
	"testing"

	"github.com/aws/aws-lambda-go/events"
)

func Test_Entrypoint(t *testing.T) {
	// Mock the request
	request := events.APIGatewayProxyRequest{
		Headers: map[string]string{
			"Content-Type": "application/hcl",
		},
		Body: "job \"example\" {}",
	}

	// Call the Entrypoint function
	response, err := Entrypoint(request)

	// Check for errors
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Check the response status code
	if response.StatusCode != 200 {
		t.Fatalf("expected status code 200, got %d: %s", response.StatusCode, response.Body)
	}

	// Check the response body
	if response.Body == "" {
		t.Fatal("expected non-empty response body")
	}
}
