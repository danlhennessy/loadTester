package main

import (
	"context"
	"testing"
	"time"
)

// func TestAddition(t *testing.T) {
// 	startNum := 1

// 	want := 2
// 	result := addOne(startNum)
// 	if result != want {
// 		t.Fatalf("Result: %d does not match expected value: %d", result, want)
// 	}
// }

func TestHitUrl(t *testing.T) {
	testUrl := "https://gesgergregerexample.com"

	want := traceResult{
		url:           testUrl,
		statusCode:    200,
		totalDuration: time.Millisecond * 1,
	}

	result, err := hitUrl(context.Background(), testUrl)

	if result.url != want.url || result.statusCode != want.statusCode || result.totalDuration < want.totalDuration || err != nil {
		t.Fatalf("Url hit failed, or totalDuration less than minimum expected: %v", result)
	}
}
