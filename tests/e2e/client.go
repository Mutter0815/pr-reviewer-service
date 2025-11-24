package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"time"
)

type httpClient struct {
	baseURL string
	client  *http.Client
}

func newHTTPClient() *httpClient {
	base := os.Getenv("APP_URL")
	if base == "" {
		base = "http://app:8080"
	}

	return &httpClient{
		baseURL: base,
		client: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func (c *httpClient) do(method, path string, body any) (*http.Response, []byte, error) {
	var reader io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, nil, err
		}
		reader = bytes.NewReader(b)
	}

	req, err := http.NewRequest(method, c.baseURL+path, reader)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, nil, err
	}
	defer func() {
		io.Copy(io.Discard, resp.Body)
		resp.Body.Close()
	}()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}

	return resp, data, nil
}
