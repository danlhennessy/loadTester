package main

import (
	"context"
	"fmt"
	"os"
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
		t.Fatalf("Url hit failed, or totalDuration less than minimum expected: %v", result)
	}
}

func TestLeaks(t *testing.T) {
	defer goleak.VerifyNone(t)

	allResults, err := LoadTest(context.Background())
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	for _, result := range allResults {
		fmt.Println(result)
	}
}
