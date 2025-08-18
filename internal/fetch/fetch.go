package fetch

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// FetchData is a generic API client for any service
type FetchData struct {
	Scheme string
	Host   string
	Port   int
}

// NewFetchData is the constructor for a generic API client
func NewFetchData(scheme, host string, port int) *FetchData {
	return &FetchData{
		Scheme: scheme,
		Host:   host,
		Port:   port,
	}
}

func (f *FetchData) Fetch(endpoint, apiKey, qParam string, v interface{}) error {
	url := fmt.Sprintf("%s://%s%s?key=%s&q=%s", f.Scheme, f.Host, endpoint, apiKey, qParam)
	safeURL := fmt.Sprintf("%s://%s%s?q=%s", f.Scheme, f.Host, endpoint, qParam)
	fmt.Println("Fetching:", safeURL)

	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("HTTP request error: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading response body: %w", err)
	}

	if resp.StatusCode > 299 {
		return fmt.Errorf("response failed with status %d: %s", resp.StatusCode, string(body))
	}

	if err := json.Unmarshal(body, v); err != nil {
		return fmt.Errorf("JSON decode error: %w\nBody: %s", err, string(body))
	}

	return nil
}
