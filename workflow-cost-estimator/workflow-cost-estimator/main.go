package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

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

	errorMessage, statusCode, ok := getWorkflowStatus(urls, params)

	if !ok {
		body, err := json.Marshal(map[string]interface{}{
			"message": errorMessage,
		})

		if err != nil {
			return Response{StatusCode: 500}, err
		}

		json.HTMLEscape(&buf, body)

		resp := Response{
			StatusCode:      statusCode,
			IsBase64Encoded: false,
			Body:            buf.String(),
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		}

		return resp, nil
	}
	errorMessage, statusCode, ok, jobs := getWorkflowJobs(urls, params)

	errorMessage, totalCost, ok := tallyJobCost(jobs, urls, params)
	_ = totalCost

	if !ok {
		return Response{StatusCode: 500}, errors.New(errorMessage)
	}

	body, err := json.Marshal(map[string]interface{}{
		"message": "Go Serverless v1.0! Your function executed successfully!",
	})

	if err != nil {
		return Response{StatusCode: 500}, err
	}

	json.HTMLEscape(&buf, body)

	resp := Response{
		StatusCode:      200,
		IsBase64Encoded: false,
		Body:            buf.String(),
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}

	return resp, nil
}

func tallyJobCost(jobs workflowJobsResponse, urls circleURLs, params queryParameters) (string, float64, bool) {
	var totalCredits float64

	for _, job := range jobs.Jobs {
		jobURL := fmt.Sprintf("%sproject/%s/%s/%s/%d", urls.v1URL, params.projectVCS, params.projectUser, params.projectName, job.JobNumber)
		errorMessage, cost, ok := getJobDetails(jobURL, params)

		if !ok {
			return errorMessage, 0, false
		}
		totalCredits += cost
	}
	fmt.Println(totalCredits, float64(totalCredits))
	fmt.Println("total cost:", float64(totalCredits)*creditPrice)
	fmt.Println(totalCredits, "total credits")
	return "", totalCredits, true
}

func getJobDetails(url string, params queryParameters) (string, float64, bool) {
	var response jobDetailResponse
	var buildTime int

	resp, errorMessage, ok := makeBasicAuthRequest(url, params.circleToken)
	defer resp.Body.Close()

	if !ok {
		return errorMessage, 500, false
	}

	errorMessage, ok = unmarshalAPIResp(resp, &response)

	resourceClass := response.Picard.ResourceClass.Class
	exeuctor := response.Picard.Executor
	creditPerMin := resourceClasses[exeuctor][resourceClass]

	if !ok {
		return errorMessage, 500, false
	}

	for _, step := range response.Steps {
		for _, action := range step.Actions {
			if action.Background {
				continue
			}

			buildTime += action.RunTimeMillis
		}
	}

	seconds := float64((buildTime / 1000) % 60)
	minutes := float64((buildTime / (1000 * 60)) % 60)
	fmt.Printf("Minutes: %f, Seconds: %f\nCredit per min: %f\n", minutes, seconds, creditPerMin)

	cost := (minutes + seconds/60) * creditPerMin

	fmt.Printf("build credits for %s: %f\n", response.Workflows.JobName, cost)

	return "", cost, true

}

func getWorkflowJobs(urls circleURLs, params queryParameters) (string, int, bool, workflowJobsResponse) {
	var response workflowJobsResponse
	workflowJobsURL := fmt.Sprintf("%sworkflow/%s/jobs", urls.v2URL, params.workflowID)

	resp, errorMessage, ok := makeBasicAuthRequest(workflowJobsURL, params.circleToken)
	defer resp.Body.Close()

	if !ok {
		return errorMessage, 500, false, response
	}

	errorMessage, ok = unmarshalAPIResp(resp, &response)

	if !ok {
		return errorMessage, 500, false, response
	}

	return "", 0, true, response
}

func getWorkflowStatus(urls circleURLs, params queryParameters) (string, int, bool) {
	var response workflowResponse
	workflowURL := fmt.Sprintf("%sworkflow/%s", urls.v2URL, params.workflowID)

	resp, errorMessage, ok := makeBasicAuthRequest(workflowURL, params.circleToken)
	defer resp.Body.Close()

	if !ok {
		return errorMessage, 500, false
	}

	errorMessage, ok = unmarshalAPIResp(resp, &response)

	if !ok {
		return errorMessage, 500, false
	}

	if response.Status != "success" && response.Status != "failed" {
		return fmt.Sprintf("Workflow status is %s. Status must be 'success' or 'failed' to estimate cost", response.Status), 202, false
	}

	return "", 0, true
}

func unmarshalAPIResp(resp *http.Response, f interface{}) (string, bool) {
	var bodyString string

	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return fmt.Sprintf("Error reading API response. Error: %s", err), false
		}
		bodyString = string(bodyBytes)
	} else {
		return fmt.Sprintf("Bad status code from API response. Status code: %d", resp.StatusCode), false
	}

	if err := json.Unmarshal([]byte(bodyString), &f); err != nil {
		return fmt.Sprintf("Error unmarshalling JSON resposnse. Error: %s", err), false
	}

	return "", true
}

func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

func makeBasicAuthRequest(url string, token string) (*http.Response, string, bool) {
	client := &http.Client{
		Timeout: time.Second * 10,
	}

	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return nil, fmt.Sprintf("Error creating HTTP client with provided URL. URL: %s. Error: %s", url, err), false
	}

	req.Header.Add("Authorization", "Basic "+basicAuth(token, ""))
	req.Header.Add("Accept", "application/json")

	resp, err := client.Do(req)

	if err != nil {
		return nil, fmt.Sprintf("Error getting requested URL. URL: %s. Error: %s", url, err), false
	}

	return resp, "", true
}

func paramSetup(request map[string]string) (queryParameters, circleURLs, bool, string) {
	var params queryParameters
	var urls circleURLs

	if request == nil || request["circle_token"] == "" || request["workflow_id"] == "" || request["project_name"] == "" || request["project_user"] == "" || request["project_vcs"] == "" {
		return params, urls, false, "Please provide query parameters: circle_token, workflow_id, project_name, project_user, project_vcs"
	}

	if request["circle_url"] == "" {
		urls.circleURL = "https://circleci.com"
	} else {
		urls.circleURL = request["circle_url"]
	}

	urls.v1URL = fmt.Sprintf("%s/api/v1.1/", urls.circleURL)
	urls.v2URL = fmt.Sprintf("%s/api/v2/", urls.circleURL)

	params.circleToken = request["circle_token"]
	params.workflowID = request["workflow_id"]
	params.projectName = request["project_name"]
	params.projectUser = request["project_user"]
	params.projectVCS = request["project_vcs"]

	return params, urls, true, ""
}

func main() {
	lambda.Start(Handler)
}
