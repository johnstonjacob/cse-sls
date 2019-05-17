package main

import "time"

var creditPrice = 0.0006
var resourceClasses = map[string]map[string]float64{
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
	projectVCS  string
	circleToken string
	workflowID  string
	circleURL   string
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

type workflowJobsResponse struct {
	NextPageToken interface{} `json:"next_page_token"`
	Jobs          []struct {
		ID           string        `json:"id"`
		Name         string        `json:"name"`
		Type         string        `json:"type"`
		Status       string        `json:"status"`
		StartTime    time.Time     `json:"start_time"`
		Dependencies []interface{} `json:"dependencies"`
		JobNumber    int           `json:"job_number"`
		StopTime     time.Time     `json:"stop_time"`
	} `json:"jobs"`
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
