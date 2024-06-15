package main

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptrace"
	"os"
	"time"

	flag "github.com/spf13/pflag"
	"golang.org/x/sync/errgroup"
)

type traceResult struct {
	url           string
	statusCode    int
	totalDuration time.Duration
}

func main() {

	var maxGoroutines = flag.IntP("goroutines", "g", 0, "Maximum number of goroutines")
	flag.Parse()
	urls := flag.Args()

	allResults, err := LoadTest(urls, maxGoroutines)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	for result := range allResults {
		fmt.Println(result)
	}
}

func LoadTest(urls []string, maxGoroutines *int) ([]traceResult, error) {

	group, ctx := errgroup.WithContext(context.Background())
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	group.SetLimit(*maxGoroutines)

	errChan := make(chan error, len(urls))
	results := make([]traceResult, len(urls))

	for i := range urls {
		group.Go(func() error {
			result, err := hitUrl(ctx, urls[i])
			if err != nil {
				errChan <- err
			}
			results[i] = result
			return nil
		})
	}

	go func() {
		for err := range errChan {
			fmt.Fprintln(os.Stderr, "Error:", err)
		}
	}()

	if err := group.Wait(); err != nil {
		return nil, err
	}

	return results, nil
}

func hitUrl(ctx context.Context, url string) (traceResult, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return traceResult{url: url}, err
	}

	trace := &httptrace.ClientTrace{
		GotConn: func(connInfo httptrace.GotConnInfo) {
			fmt.Printf("Got Conn: %+v\n", connInfo)
		},
		DNSStart: func(dnsStartInfo httptrace.DNSStartInfo) {
			fmt.Printf("DNS lookup started for: %v at time: %s\n", dnsStartInfo, time.Now())
		},
		DNSDone: func(dnsDoneInfo httptrace.DNSDoneInfo) {
			fmt.Printf("DNS Info: %+v - Completed at time: %s\n", dnsDoneInfo, time.Now())
		},
	}

	req = req.WithContext(httptrace.WithClientTrace(req.Context(), trace))

	client := &http.Client{}

	start := time.Now()

	resp, err := client.Do(req)
	if err != nil {
		return traceResult{url: url}, err
	}
	defer resp.Body.Close()

	elapsed := time.Since(start)

	result := traceResult{
		url:           url,
		statusCode:    resp.StatusCode,
		totalDuration: elapsed,
	}

	return result, nil
}
