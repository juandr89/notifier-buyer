package server

import (
	"bytes"
	"fmt"
	"net/http"
	"time"
)

type RequestOptions struct {
	Method         string
	URL            string
	Body           []byte
	Headers        map[string]string
	MaxRetries     int
	RetryDelay     time.Duration
	RequestTimeout time.Duration
}

func DoRequestWithRetry(opts RequestOptions) (*http.Response, error) {
	client := &http.Client{
		Timeout: opts.RequestTimeout,
	}
	var resp *http.Response
	var err error

	for i := 0; i < opts.MaxRetries; i++ {
		req, err := http.NewRequest(opts.Method, opts.URL, bytes.NewBuffer(opts.Body))
		if err != nil {
			return nil, fmt.Errorf("error creating request: %w", err)
		}
		for key, value := range opts.Headers {
			req.Header.Set(key, value)
		}
		resp, err = client.Do(req)
		if err == nil && resp.StatusCode != int(500) {
			return resp, nil
		}
		fmt.Printf("Request failed (attempt %d/%d): %v\n", i+1, opts.MaxRetries, err)
		time.Sleep(opts.RetryDelay)
	}

	return nil, fmt.Errorf("request failed after %d attempts: %w", opts.MaxRetries, err)
}
