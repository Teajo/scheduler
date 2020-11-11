package publisher

import (
	"bytes"
	"errors"
	"net/http"
)

// HTTPPublisher represents a publisher
type HTTPPublisher struct {
	method string // PUT or POST
	url    string
	json   string
}

// NewHTTPPublisher creates a new publisher
func NewHTTPPublisher(method string, url string, json string) *HTTPPublisher {
	return &HTTPPublisher{
		method: method,
		url:    url,
		json:   json,
	}
}

// Publish publishes
func (p *HTTPPublisher) Publish() error {
	jsonStr := []byte(p.json)
	req, err := http.NewRequest(p.method, p.url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	resp.Body.Close()

	if resp.StatusCode >= 400 {
		return errors.New("HTTP status error code")
	}
	return err
}
