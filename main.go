package main

import (
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/hashicorp/nomad/api"

	"deployer-nomad-aws/apikeys"
	"deployer-nomad-aws/nomad"
)

func Entrypoint(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	if !apikeys.Have(request.Headers["Authorization"]) {
		return events.APIGatewayProxyResponse{
			StatusCode: 403,
			Body:       "Forbidden",
		}, nil
	}

	var (
		reqContentType = request.Headers["Content-Type"]
		client         = nomad.GetInstance().Jobs()
		job            *api.Job
		resp           events.APIGatewayProxyResponse
		err            error
	)

	switch reqContentType {
	case "application/hcl":
		job, err = client.ParseHCL(string(resp.Body), true)
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

	if validateResp, _, err := client.Validate(job, nil); err != nil {
		b, _ := json.Marshal(validateResp)
		resp.StatusCode = 400
		resp.Headers = map[string]string{
			"Content-Type": "application/json",
		}
		resp.Body = string(b)
		return resp, nil
	}

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
