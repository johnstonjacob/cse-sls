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

func TestNormalizeVCS(t *testing.T) {
	tables := []struct {
		input    string
		expected string
		errRes   string
	}{
		{"gh", "github", ""},
		{"bb", "bitbucket", ""},
		{"github", "github", ""},
		{"bitbucket", "bitbucket", ""},
		{"gl", "", "VCS gl is not valid."},
	}

	for _, table := range tables {
		if output, err := normalizeVCS(table.input); output != table.expected {
			t.Errorf("Got %s, expected %s", output, table.expected)
		} else if err != nil {
			if err.Error() != table.errRes {
				t.Errorf("Got %s, expected %s", err.Error(), table.errRes)
			}
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
func TestCostOfWorkflow(t *testing.T) {
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
		if output := costOfWorkflow(table.input); toFixed(output, 4) != toFixed(table.expected, 4) {
			t.Errorf("Got %f, expected %f", output, table.expected)
		}
	}
}
