package main

import (
	"bytes"
	"errors"
	"jpb/scheduler/publisher"
	"net/http"
)

// HTTPPublisher represents a http publisher
type HTTPPublisher struct {
}

func main() {
}

// New creates a new publisher
func New() publisher.Publisher {
	return &HTTPPublisher{}
}

// Publish publishes
func (p *HTTPPublisher) Publish(cfg map[string]string) error {
	method := cfg["method"]
	url := cfg["url"]
	json := cfg["json"]

	req, err := http.NewRequest(method, url, bytes.NewBuffer([]byte(json)))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	resp.Body.Close()

	if resp.StatusCode >= 400 {
		return errors.New("HTTP status error code")
	}
	return err
}

// CheckConfig checks publisher config
func (p *HTTPPublisher) CheckConfig(cfg map[string]string) error {
	method := cfg["method"]
	if method != http.MethodPost && method != http.MethodPut {
		return errors.New("Http method must be POST or PUT")
	}

	_, urlOk := cfg["url"]
	if !urlOk {
		return errors.New("Must provide an url in url field")
	}

	_, jsonOk := cfg["json"]
	if !jsonOk {
		return errors.New("Must provide a json field")
	}

	return nil
}
