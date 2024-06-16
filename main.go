package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/httptrace"
	"os"
	"sync"
	"time"

	flag "github.com/spf13/pflag"
)

type traceResult struct {
	url           string
	statusCode    int
	totalDuration time.Duration
}

func main() {
	flag.Parse()
	urls := flag.Args()

	allResults, err := LoadTest(urls)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	for _, result := range allResults {
		log.Printf("Result: %+v\n", result)
	}
}

func LoadTest(urls []string) ([]traceResult, error) {
	ctx := context.Background()
	var results []traceResult
	resultChan := make(chan traceResult)
	var wg sync.WaitGroup

	wg.Add(len(urls))
	for i := range urls {
		go func(i int) {
			ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
			defer cancel()
			defer wg.Done()

			result, err := hitUrl(ctx, urls[i])
			resultChan <- result
			if err != nil {
				log.Fatal(err)
			}
		}(i)
	}
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	for result := range resultChan {
		results = append(results, result)
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

	client.CloseIdleConnections()
	return result, nil
}
