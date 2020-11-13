package main

import (
	"jpb/scheduler/publisher"
)

// TestPublisher represents a http publisher
type TestPublisher struct {
}

func main() {
}

// New creates a new publisher
func New() publisher.Publisher {
	return &TestPublisher{}
}

// Publish publishes
func (p *TestPublisher) Publish(cfg map[string]string) *publisher.PublishError {
	return nil
}

// CheckConfig checks publisher config
func (p *TestPublisher) CheckConfig(cfg map[string]string) error {
	return nil
}
