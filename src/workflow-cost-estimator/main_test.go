package main

import (
	"math"
	"testing"
)

func TestSnakeCaseToCamelCase(t *testing.T) {
	tables := []struct {
		input    string
		expected string
	}{
		{"this_is_snake_case", "thisIsSnakeCase"},
		{"workflow_ID", "workflowID"},
		{"hello____Snake______CASE", "helloSnakeCASE"},
	}

	for _, table := range tables {
		if output := snakeCaseToCamelCase(table.input); output != table.expected {
			t.Errorf("Got %s, expected %s", output, table.expected)
		}
	}

}

func round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

func toFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(round(num*output)) / output
}

func TestCreditCost(t *testing.T) {
	tables := []struct {
		input    float64
		expected float64
	}{
		{500, .3},
		{20659, 12.395400},
		{1378, .8268},
		{10, .006},
		{1000, .6},
		{100000, 60},
	}

	for _, table := range tables {
		if output := creditCost(table.input); toFixed(output, 4) != toFixed(table.expected, 4) {
			t.Errorf("Got %f, expected %f", output, table.expected)
		}
	}
}

func TestLookupCreditPerMin(t *testing.T) {
	tables := []struct {
		executor  string
		rc        string
		jobName   string
		expected  float64
		errString string
	}{
		{"docker", "2xlarge+", "job", 100, ""},
		{"machine", "2xlarge", "job", 80, ""},
		{"machine", "gpu.large", "job", 320, ""},
		// TODO probably a better way to do these
		{"watson", "qubit.large", "job", 0, "Missing Executor. Please contact jacobjohnston@circleci.com with this error message, your parameters, and executor type. Executor: watson"},
		{"macos", "mac128k", "job", 0, "Missing resource class cost for macos:mac128k in job job"},
	}

	for _, table := range tables {
		if output, err := lookupCreditPerMin(table.executor, table.rc, table.jobName); output != table.expected {
			t.Errorf("Got %f, expected %f", output, table.expected)
		} else if err != nil {
			if err.Error() != table.errString {
				t.Errorf("Got %s, expected %s", err.Error(), table.errString)
			}
		}
	}
}
