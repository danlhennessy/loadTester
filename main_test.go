package main

import (
	"context"
	"testing"
	"time"
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
