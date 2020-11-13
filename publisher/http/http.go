package main

import (
	"bytes"
	"errors"
	"jpb/scheduler/publisher"
	"net/http"
)

// HTTPError is http error
type HTTPError int

const (
	resourceError HTTPError = iota
	networkError
	noError
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
func (p *HTTPPublisher) Publish(cfg map[string]string) *publisher.PublishError {
	method := cfg["method"]
	url := cfg["url"]
	json := cfg["json"]

	req, err := http.NewRequest(method, url, bytes.NewBuffer([]byte(json)))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return publisher.NewPublishError(err, true)
	}
	resp.Body.Close()

	httpError := getHTTPError(resp.StatusCode)

	switch httpError {
	case networkError:
		return publisher.NewPublishError(err, true)
	case resourceError:
		return publisher.NewPublishError(err, false)
	default:
		return nil
	}
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

func getHTTPError(statusCode int) HTTPError {
	if statusCode >= 500 {
		return networkError
	}

	if statusCode >= 400 {
		return resourceError
	}

	return noError
}
