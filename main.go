package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"deployer-nomad-aws/apikeys"
	"deployer-nomad-aws/nomad"
)

func Entrypoint(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	fmt.Println("Lambda invoked.")
	fmt.Printf("Headers received: %+v\n", request.Headers)
	fmt.Printf("x-api-key received: %q\n", request.Headers["x-api-key"])
	fmt.Printf("API_KEYS env var: %s\n", os.Getenv("API_KEYS"))
	fmt.Printf("IsBase64Encoded: %v\n", request.IsBase64Encoded)

	// Authenticate API key
	if !apikeys.Have(request.Headers["x-api-key"]) {
		fmt.Println("Unauthorized request: missing or invalid API key.")
		return events.APIGatewayProxyResponse{
			StatusCode: 403,
			Body:       "Forbidden",
		}, nil
	}

	// Decode the body (base64 if needed)
	var rawBody string
	if request.IsBase64Encoded {
		decoded, err := base64.StdEncoding.DecodeString(request.Body)
		if err != nil {
			return events.APIGatewayProxyResponse{
				StatusCode: 400,
				Body:       "Failed to decode base64 body: " + err.Error(),
			}, nil
		}
		rawBody = string(decoded)
	} else {
		rawBody = request.Body
	}

	fmt.Println("Attempting to parse job...")

	// Parse HCL into Nomad job
	client := nomad.GetInstance().Jobs()
	job, err := client.ParseHCL(rawBody, true)
	if err != nil {
		fmt.Println("HCL parsing error:", err)
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "Invalid HCL: " + err.Error(),
		}, nil
	}

	// Validate Nomad job
	validateResp, _, err := client.Validate(job, nil)
	if err != nil {
		b, _ := json.Marshal(validateResp)
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Headers:    map[string]string{"Content-Type": "application/json"},
			Body:       string(b),
		}, nil
	}

	// Register Nomad job
	rr, _, err := client.Register(job, nil)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "Failed to register Nomad job: " + err.Error(),
		}, nil
	}

	respBytes, err := json.Marshal(rr)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "Failed to marshal response: " + err.Error(),
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Headers:    map[string]string{"Content-Type": "application/json"},
		Body:       string(respBytes),
	}, nil
}

func main() {
	lambda.Start(Entrypoint)
}