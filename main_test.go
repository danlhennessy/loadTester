package main

import "testing"

func TestAddition(t *testing.T) {
	startNum := 1

	want := 2
	result := addOne(startNum)
	if result != want {
		t.Fatalf("Result: %d does not match expected value: %d", result, want)
	}
}
