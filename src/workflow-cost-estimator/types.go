package main

import "time"

type circleURLs struct {
	circleURL string
	v1URL     string
	v2URL     string
}

type queryParameters map[string]string

type job struct {
	Name         string  `json:"job_name"`
	TotalCost    float64 `json:"total_cost"`
	TotalCredits float64 `json:"total_credits"`
	TotalRuntime string  `json:"total_runtime"`
}

type responseBody struct {
	TotalCost    float64 `json:"total_cost"`
	TotalCredits float64 `json:"total_credits"`
	TotalRuntime string  `json:"total_runtime"`
	Disclaimer   string  `json:"disclaimer"`
	Jobs         []job   `json:"jobs"`
}

func newResponseBody(totalCredits, totalCost float64, totalRuntime time.Duration, jobs []job) *responseBody {
	b := new(responseBody)
	disclaimer := "This is a cost estimate. This is not an official CircleCI endpoint. Please contact jacobjohnston@circleci.com for questions."
	b.TotalCost = totalCost
	b.TotalCredits = totalCredits
	b.TotalRuntime = totalRuntime.String()
	b.Jobs = jobs
	b.Disclaimer = disclaimer

	return b
}

type responseErr struct {
	err        string
	statusCode int
}

func (e responseErr) Error() string {
	return e.err
}

// Actions is a struct of job details actions
type Actions struct {
	Index              int         `json:"index"`
	Parallel           bool        `json:"parallel"`
	Failed             interface{} `json:"failed"`
	InfrastructureFail interface{} `json:"infrastructure_fail"`
	Name               string      `json:"name"`
	Status             string      `json:"status"`
	Background         bool        `json:"background"`
	ExitCode           interface{} `json:"exit_code"`
	Canceled           interface{} `json:"canceled"`
	Step               int         `json:"step"`
	RunTimeMillis      int         `json:"run_time_millis"`
	HasOutput          bool        `json:"has_output"`
}

// Steps is a struct of job actions
type Steps struct {
	Name    string    `json:"name"`
	Actions []Actions `json:"actions"`
}

//ResourceClass is
type ResourceClass struct {
	Class string `json:"class"`
}

// Picard is
type Picard struct {
	ResourceClass ResourceClass `json:"resource_class"`
	Executor      string        `json:"executor"`
}
type jobDetailResponse struct {
	Steps              []Steps `json:"steps"`
	Parallel           int     `json:"parallel"`
	InfrastructureFail bool    `json:"infrastructure_fail"`
	Status             string  `json:"status"`
	Lifecycle          string  `json:"lifecycle"`
	BuildTimeMillis    int     `json:"build_time_millis"`
	Picard             Picard  `json:"picard"`
	Workflows          struct {
		JobName string `json:"job_name"`
	} `json:"workflows"`
}

// Jobs is
type Jobs struct {
	ID           string        `json:"id"`
	Name         string        `json:"name"`
	Type         string        `json:"type"`
	Status       string        `json:"status"`
	StartTime    time.Time     `json:"start_time"`
	Dependencies []interface{} `json:"dependencies"`
	JobNumber    int           `json:"job_number"`
	ProjectSlug  string        `json:"project_slug"`
	StopTime     time.Time     `json:"stop_time"`
}
type workflowJobsResponse struct {
	NextPageToken interface{} `json:"next_page_token"`
	Jobs          []Jobs      `json:"items"`
}

type workflowResponse struct {
	CreatedAt      time.Time   `json:"created_at"`
	ID             string      `json:"id"`
	Name           string      `json:"name"`
	PipelineID     interface{} `json:"pipeline_id"`
	PipelineNumber interface{} `json:"pipeline_number"`
	Project        struct {
		ID string `json:"id"`
	} `json:"project"`
	Status    string    `json:"status"`
	StoppedAt time.Time `json:"stopped_at"`
}
