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
	totalDuration time.Duration
}

func main() {
	LoadTest := func(ctx context.Context) ([]traceResult, error) {
		var maxGoroutines = flag.IntP("goroutines", "g", 0, "Maximum number of goroutines")
		flag.Parse()
		urls := flag.Args()
		g, ctx := errgroup.WithContext(context.Background())
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
		g.SetLimit(*maxGoroutines)

		errChan := make(chan error, len(urls))
		results := make([]traceResult, len(urls))

		for i := range urls {
			g.Go(func() error {
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

		if err := g.Wait(); err != nil {
			return nil, err
		}

		return results, nil
	}

	allResults, err := LoadTest(context.Background())
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	for _, result := range allResults {
		fmt.Println(result)
	}
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
		DNSDone: func(dnsInfo httptrace.DNSDoneInfo) {
			fmt.Printf("DNS Info: %+v\n", dnsInfo)
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
		totalDuration: elapsed,
	}

	return result, nil
}
