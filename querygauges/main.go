package main

import (
	"flag"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"time"
)

func main() {
	promURL := flag.String("url", "http://localhost:9090", "Prometheus server URL")
	query := flag.String("query", "", "PromQL query to execute")
	flag.Parse()

	if *query == "" {
		fmt.Fprintln(os.Stderr, "Error: -query flag is required")
		os.Exit(1)
	}

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	queryURL := fmt.Sprintf("%s/api/v1/query", *promURL)

	for {
		params := url.Values{}
		params.Set("query", *query)

		resp, err := client.Get(queryURL + "?" + params.Encode())
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error executing query: %v\n", err)
			os.Exit(1)
		}

		// Read and discard the response body
		_, err = io.Copy(io.Discard, resp.Body)
		resp.Body.Close()

		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading response: %v\n", err)
			os.Exit(1)
		}

		if resp.StatusCode != http.StatusOK {
			fmt.Fprintf(os.Stderr, "Error: received status code %d\n", resp.StatusCode)
			os.Exit(1)
		}

		// Random delay between 0-0.1 seconds
		delay := time.Duration(rand.Float64() * float64(100*time.Millisecond))
		time.Sleep(delay)
	}
}
