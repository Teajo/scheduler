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
func (p *HTTPPublisher) Publish(cfg map[string]interface{}) *publisher.PublishError {
	method, ok := cfg["method"].(string)
	if !ok {
		panic(errors.New("field 'method' has incorrect type"))
	}

	url, ok := cfg["url"].(string)
	if !ok {
		panic(errors.New("field 'url' has incorrect type"))
	}

	json, ok := cfg["json"].(string)
	if !ok {
		panic(errors.New("field 'json' has incorrect type"))
	}

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

// GetConfigDef returns needed config for this publisher
func (p *HTTPPublisher) GetConfigDef() map[string]*publisher.ConfigValueDef {
	m := make(map[string]*publisher.ConfigValueDef)
	m["method"] = &publisher.ConfigValueDef{Required: true, Possible: []string{"POST", "PUT"}, Default: "", Type: publisher.STRING}
	m["url"] = &publisher.ConfigValueDef{Required: true, Possible: nil, Default: "", Type: publisher.STRING, Placeholder: "http://127.0.0.1:8080"}
	m["json"] = &publisher.ConfigValueDef{Required: true, Possible: nil, Default: "", Type: publisher.JSON_STRING, Placeholder: "{}"}
	return m
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
