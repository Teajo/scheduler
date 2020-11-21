package utils

import (
	"time"

	"github.com/google/uuid"
)

// RetryStrat represents a retry strategy
type RetryStrat struct {
	Timeout     time.Duration `json:"timeout"`
	Exponential bool          `json:"exponential"`
	Limit       int           `json:"limit"`
}

// Publisher represents a publisher
type Publisher struct {
	Publisher  string                 `json:"publisher"`
	Settings   map[string]interface{} `json:"settings"`
	RetryStrat *RetryStrat            `json:"retryStrategy"`
}

// Scheduling represents a scheduling object
type Scheduling struct {
	ID         string       `json:"id"`
	Date       time.Time    `json:"date"`
	Done       bool         `json:"done"`
	Publishers []*Publisher `json:"publishers"`
}

// NewScheduling creates a new scheduling struct
func NewScheduling(date time.Time, publishers []*Publisher) *Scheduling {
	return &Scheduling{
		ID:         uuid.New().String(),
		Date:       date,
		Publishers: publishers,
	}
}

// NewSchedulingWithID creates a new scheduling struct
func NewSchedulingWithID(id string, date time.Time, publishers []*Publisher, done bool) *Scheduling {
	return &Scheduling{
		ID:         id,
		Date:       date,
		Done:       done,
		Publishers: publishers,
	}
}

// LastDate last possible date
var LastDate time.Time = time.Unix(1<<63-62135596801, 999999999)

// FirstDate first possible date
var FirstDate time.Time = time.Time{}
