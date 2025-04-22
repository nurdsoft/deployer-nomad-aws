package main

import (
	"encoding/base64"
	"encoding/json"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/hashicorp/nomad/api"

	"deployer-nomad-aws/apikeys"
	"deployer-nomad-aws/nomad"
)

func Entrypoint(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Handle auth header or x-api-key
	authHeader := strings.TrimSpace(request.Headers["Authorization"])
	if strings.HasPrefix(authHeader, "Bearer ") {
		authHeader = strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))
	} else if authHeader == "" {
		authHeader = strings.TrimSpace(request.Headers["x-api-key"])
	}

	if !apikeys.Have(authHeader) {
		return events.APIGatewayProxyResponse{
			StatusCode: 403,
			Body:       "Forbidden",
		}, nil
	}

	var (
		client = nomad.GetInstance().Jobs()
		job    *api.Job
		resp   events.APIGatewayProxyResponse
		err    error
	)

	// Normalize content type
	contentType := strings.ToLower(strings.TrimSpace(request.Headers["Content-Type"]))
	if contentType == "" {
		contentType = strings.ToLower(strings.TrimSpace(request.Headers["content-type"]))
	}

	// Decode body if needed
	var body string
	if request.IsBase64Encoded {
		decoded, err := base64.StdEncoding.DecodeString(request.Body)
		if err != nil {
			return events.APIGatewayProxyResponse{
				StatusCode: 400,
				Body:       "Invalid base64-encoded body",
			}, nil
		}
		body = string(decoded)
	} else {
		body = request.Body
	}

	switch contentType {
	case "application/hcl":
		job, err = client.ParseHCL(strings.TrimSpace(body), true)
		if err != nil {
			resp.StatusCode = 400
			resp.Body = err.Error()
			return resp, nil
		}
	default:
		resp.StatusCode = 400
		resp.Body = "unsupported content type"
		return resp, nil
	}

	// Validate job
	if validateResp, _, err := client.Validate(job, nil); err != nil {
		b, _ := json.Marshal(validateResp)
		resp.StatusCode = 400
		resp.Headers = map[string]string{
			"Content-Type": "application/json",
		}
		resp.Body = string(b)
		return resp, nil
	}

	// Register job
	rr, _, err := client.Register(job, nil)
	if err != nil {
		resp.StatusCode = 500
		resp.Body = err.Error()
		return resp, nil
	}

	respBytes, err := json.Marshal(rr)
	if err != nil {
		resp.StatusCode = 500
		resp.Body = err.Error()
		return resp, nil
	}

	resp.Headers = map[string]string{
		"Content-Type": "application/json",
	}
	resp.Body = string(respBytes)
	resp.StatusCode = 200
	return resp, nil
}

func main() {
	lambda.Start(Entrypoint)
}