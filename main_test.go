package main

import (
	"context"
	"fmt"
	"testing"
	"time"

	"go.uber.org/goleak"
)

func TestHitUrl(t *testing.T) {
	testUrl := "https://example.com"

	want := traceResult{
		url:           testUrl,
		statusCode:    200,
		totalDuration: time.Millisecond * 1,
	}

	result, err := hitUrl(context.Background(), testUrl)

	if result.url != want.url || result.statusCode != want.statusCode || result.totalDuration < want.totalDuration || err != nil {
		t.Fatalf("Url hit failed, or totalDuration less than minimum expected: %v, Error: %s\n", result, err)
	}
}

func TestConcurrentLoadTest(t *testing.T) {
	defer goleak.VerifyNone(t)

	testUrls := []string{
		"https://example.com",
		"https://grrrrrrqwddwqdwoogle.com",
		"https://bbc.co.uk",
	}
	testGoroutines := len(testUrls)
	want := make([]traceResult, testGoroutines)

	allResults, err := LoadTest(testUrls, &testGoroutines)

	fmt.Println("Check")

	if err != nil {
		t.Fatalf("Concurrent load test failed, error: %s\n", err)
	}
	if len(allResults) < len(want) {
		t.Fatalf("Not enough results returned\n")
	}
}
