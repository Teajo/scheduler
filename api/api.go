package api

import "jpb/scheduler/utils"

type api interface {
	Listen()
}

// Scheduling represents a scheduling object
type Scheduling struct {
	ID         string             `json:"id"`
	Date       string             `json:"date"`
	Done       bool               `json:"done"`
	Publishers []*utils.Publisher `json:"publishers"`
}
