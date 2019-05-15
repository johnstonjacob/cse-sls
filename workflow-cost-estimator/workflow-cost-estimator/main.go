package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var creditPrice = "0.0006"
var resourceClasses = map[string]map[string]int{
	"docker": {
		"small":    5,
		"medium":   10,
		"medium+":  15,
		"large":    20,
		"xlarge":   40,
		"2xlarge":  80,
		"2xlarge+": 100,
		"3xlarge":  160,
		"4xlarge":  320,
	},
	"machine": {
		"small":   5,
		"medium":  10,
		"large":   20,
		"xlarge":  40,
		"2xlarge": 80,
		"3xlarge": 120,
	},
	"macOS":   {},
	"GPU":     {},
	"windows": {},
}

type circleURLs struct {
	circleURL string
	v1URL     string
	v2URL     string
}

type queryParameters struct {
	projectName string
	projectUser string
	circleToken string
	workflowID  string
	circleURL   string
}

// Response is of type APIGatewayProxyResponse since we're leveraging the
// AWS Lambda Proxy Request functionality (default behavior)
//
// https://serverless.com/framework/docs/providers/aws/events/apigateway/#lambda-proxy-integration
type Response events.APIGatewayProxyResponse

// Handler is our lambda handler invoked by the `lambda.Start` function call
func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (Response, error) {
	var buf bytes.Buffer

	params, urls, ok, reason := paramSetup(request.QueryStringParameters)

	if !ok {
		body, err := json.Marshal(map[string]interface{}{
			"message": reason,
		})

		if err != nil {
			return Response{StatusCode: 404}, err
		}

		json.HTMLEscape(&buf, body)
		return Response{StatusCode: 401, Body: buf.String()}, nil
	}

	getWorkflowStatus(urls, params)

	body, err := json.Marshal(map[string]interface{}{
		"message": "Go Serverless v1.0! Your function executed successfully!",
	})

	if err != nil {
		return Response{StatusCode: 404}, err
	}

	json.HTMLEscape(&buf, body)

	resp := Response{
		StatusCode:      200,
		IsBase64Encoded: false,
		Body:            buf.String(),
		Headers: map[string]string{
			"Content-Type":           "application/json",
			"X-MyCompany-Func-Reply": "hello-handler",
		},
	}

	return resp, nil
}

func getWorkflowStatus(urls circleURLs, params queryParameters) bool {
	workflowURL := fmt.Sprintf("%sworkflow/%s", urls.v2URL, params.workflowID)

	makeBasicAuthRequest(workflowURL, params.circleToken)

	return true
}

func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

func makeBasicAuthRequest(url string, token string) (string, bool) {
	client := &http.Client{}
	fmt.Println(url, token)

	req, err := http.NewRequest("GET", url, nil)

	req.Header.Add("Authorization", "Basic "+basicAuth(token, ""))

	if err != nil {
	}

	resp, err := client.Do(req)
	fmt.Println(resp)

	return "", true
}

func paramSetup(request map[string]string) (queryParameters, circleURLs, bool, string) {
	var params queryParameters
	var urls circleURLs

	if request == nil || request["circle_token"] == "" || request["workflow_id"] == "" || request["project_name"] == "" || request["project_user"] == "" {
		return params, urls, false, "Please provide query parameters: circle_token, workflow_id, project_name, project_user"
	}

	if request["circle_url"] == "" {
		urls.circleURL = "https://circleci.com"
	} else {
		urls.circleURL = request["circle_url"]
	}

	urls.v1URL = fmt.Sprintf("%s/api/v1/", urls.circleURL)
	urls.v2URL = fmt.Sprintf("%s/api/v2/", urls.circleURL)

	params.circleToken = request["circle_token"]
	params.workflowID = request["workflow_id"]
	params.projectName = request["project_name"]
	params.projectUser = request["project_user"]

	return params, urls, true, ""
}

func main() {
	lambda.Start(Handler)
}
