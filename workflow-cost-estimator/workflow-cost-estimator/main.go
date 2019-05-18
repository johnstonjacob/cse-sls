package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"sync"
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
	var body body

	params, urls, err := paramSetup(request.QueryStringParameters)
	if err != nil {
		return *generateResponse(body, err), nil
	}

	err = getWorkflowStatus(urls, params)
	if err != nil {
		return *generateResponse(body, err), nil
	}

	jobs, err := getWorkflowJobs(urls, params)
	if err != nil {
		return *generateResponse(body, err), nil
	}

	body, err = tallyJobCost(jobs, urls, params)

	if err != nil {
		return *generateResponse(body, err), nil
	}

	return *generateResponse(body, nil), nil
}

func generateResponse(body body, err error) *Response {
	var buf bytes.Buffer
	var bodyBytes []byte
	var statusCode int

	if err, ok := err.(responseErr); ok {
		statusCode = err.statusCode
		errBody, err := json.Marshal(map[string]interface{}{
			"error": err.Error(),
		})

		if err != nil {
			return generateResponse(body, responseErr{err: err.Error(), statusCode: 500})
		}
		bodyBytes = errBody
	} else {
		responseBytes, err := json.Marshal(body)
		bodyBytes = responseBytes
		statusCode = 200

		if err != nil {
			return generateResponse(body, responseErr{err: err.Error(), statusCode: 500})
		}
	}

	json.HTMLEscape(&buf, bodyBytes)

	return &Response{
		StatusCode:      statusCode,
		IsBase64Encoded: false,
		Body:            buf.String(),
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}
}

func tallyJobCost(jobs workflowJobsResponse, urls circleURLs, params queryParameters) (body, error) {
	var totalCredits float64
	var wg sync.WaitGroup
	var body body

	wg.Add(len(jobs.Jobs))

	c := make(chan float64, 4)
	//ec := make(chan error)

	for _, job := range jobs.Jobs {
		go func(job Jobs) {

			jobURL := fmt.Sprintf("%sproject/%s/%s/%s/%d", urls.v1URL, params.projectVCS, params.projectUser, params.projectName, job.JobNumber)
			cost, err := getJobDetails(jobURL, params)
			_ = err

			/*		if err != nil {
					ec <- err
					return
				}*/
			c <- cost
		}(job)
	}

	go func(totalCredits *float64) {
		for credits := range c {
			*totalCredits += credits
			wg.Done()
		}

	}(&totalCredits)

	// TODO: return all errors not just the first one
	/*for err := range ec {
		return body, err
	}*/

	wg.Wait()
	body.TotalCredits = math.Ceil(totalCredits)
	totalPrice := totalCredits * creditPrice
	body.TotalCost = math.Ceil(totalPrice*100) / 100
	body.Disclaimer = disclaimer

	return body, nil
}

func getJobDetails(url string, params queryParameters) (float64, error) {
	var response jobDetailResponse
	var buildTime time.Duration
	var creditPerMin float64
	var ok bool

	resp, err := makeBasicAuthRequest(url, params.circleToken)
	defer resp.Body.Close()

	if err != nil {
		return 0, err
	}

	err = unmarshalAPIResp(resp, &response)
	if err != nil {
		return 0, err
	}

	resourceClass := response.Picard.ResourceClass.Class
	executor := response.Picard.Executor
	if creditPerMin, ok = resourceClasses[executor][resourceClass]; !ok {
		return 0, responseErr{fmt.Sprintf("Missing resource class cost for %s:%s in job %s", executor, resourceClass, response.Workflows.JobName), 500}
	}

	for _, step := range response.Steps {
		for _, action := range step.Actions {
			if action.Background {
				continue
			}

			buildTime += time.Duration(action.RunTimeMillis) * time.Millisecond
		}
	}

	buildTime = buildTime.Round(time.Second)

	cost := buildTime.Minutes() * creditPerMin

	return cost, nil
}

func getWorkflowJobs(urls circleURLs, params queryParameters) (workflowJobsResponse, error) {
	var response workflowJobsResponse
	workflowJobsURL := fmt.Sprintf("%sworkflow/%s/jobs", urls.v2URL, params.workflowID)

	resp, err := makeBasicAuthRequest(workflowJobsURL, params.circleToken)
	defer resp.Body.Close()

	if err != nil {
		return response, err
	}

	err = unmarshalAPIResp(resp, &response)

	if err != nil {
		return response, err
	}

	return response, err
}

func getWorkflowStatus(urls circleURLs, params queryParameters) error {
	var response workflowResponse
	workflowURL := fmt.Sprintf("%sworkflow/%s", urls.v2URL, params.workflowID)

	resp, err := makeBasicAuthRequest(workflowURL, params.circleToken)
	defer resp.Body.Close()

	if err != nil {
		return err
	}

	err = unmarshalAPIResp(resp, &response)

	if err != nil {
		return err
	}

	if response.Status != "success" && response.Status != "failed" {
		return responseErr{fmt.Sprintf("Workflow status is %s. Status must be 'success' or 'failed' to estimate cost", response.Status), 201}
	}

	return nil
}

func unmarshalAPIResp(resp *http.Response, f interface{}) error {
	var bodyString string

	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return responseErr{fmt.Sprintf("Error reading API response. Error: %s", err), 500}
		}
		bodyString = string(bodyBytes)
	} else {
		return responseErr{fmt.Sprintf("Bad status code from CCI API response. Status code: %d", resp.StatusCode), 500}
	}

	if err := json.Unmarshal([]byte(bodyString), &f); err != nil {
		return responseErr{fmt.Sprintf("Error unmarshalling JSON resposnse. Error: %s", err), 500}
	}

	return nil
}

func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

func makeBasicAuthRequest(url string, token string) (*http.Response, error) {
	client := &http.Client{
		Timeout: time.Second * 10,
	}

	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return nil, responseErr{fmt.Sprintf("Error creating HTTP client with provided URL. URL: %s. Error: %s", url, err), 500}
	}

	req.Header.Add("Authorization", "Basic "+basicAuth(token, ""))
	req.Header.Add("Accept", "application/json")

	resp, err := client.Do(req)

	if err != nil {
		return nil, responseErr{fmt.Sprintf("Error getting requested URL. URL: %s. Error: %s", url, err), 500}
	}

	return resp, nil
}

func paramSetup(request map[string]string) (queryParameters, circleURLs, error) {
	var params queryParameters
	var urls circleURLs

	if request == nil || request["circle_token"] == "" || request["workflow_id"] == "" || request["project_name"] == "" || request["project_user"] == "" || request["project_vcs"] == "" {
		err := responseErr{"Please provide query parameters: circle_token, workflow_id, project_name, project_user, project_vcs", 400}
		return params, urls, err
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

	return params, urls, nil
}

func main() {
	lambda.Start(Handler)
}
