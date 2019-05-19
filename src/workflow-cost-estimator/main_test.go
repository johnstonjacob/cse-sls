package main

import "testing"

func TestSnakeCaseToCamelCase(t *testing.T) {
	input := "this_is_snake_case"
	expected := "thisIsSnakeCase"

	out := snakeCaseToCamelCase(input)

	if out != expected {
		t.Errorf("Got %s, expected %s", out, expected)
	}
}
