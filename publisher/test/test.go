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
func (p *TestPublisher) Publish(cfg map[string]interface{}) *publisher.PublishError {
	return nil
}

// GetConfigDef returns needed config for this publisher
func (p *TestPublisher) GetConfigDef() map[string]*publisher.ConfigValueDef {
	m := make(map[string]*publisher.ConfigValueDef)
	return m
}
