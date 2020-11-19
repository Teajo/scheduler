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

// Scheduling represents a scheduling object
type Scheduling struct {
	ID          string            `json:"id"`
	Date        time.Time         `json:"date"`
	Publishers  []string          `json:"publishers"`
	Settings    map[string]string `json:"settings"`
	Done        bool              `json:"done"`
	*RetryStrat `json:"retryStrategy"`
}

// NewScheduling creates a new scheduling struct
func NewScheduling(date time.Time, publishers []string, settings map[string]string) *Scheduling {
	return &Scheduling{
		ID:         uuid.New().String(),
		Date:       date,
		Publishers: publishers,
		Settings:   settings,
		RetryStrat: &RetryStrat{
			Timeout:     25 * time.Millisecond,
			Exponential: true,
			Limit:       5,
		},
	}
}

// NewSchedulingWithID creates a new scheduling struct
func NewSchedulingWithID(id string, date time.Time, publishers []string, settings map[string]string, done bool) *Scheduling {
	return &Scheduling{
		ID:         id,
		Date:       date,
		Publishers: publishers,
		Settings:   settings,
		Done:       done,
		RetryStrat: &RetryStrat{
			Timeout:     25 * time.Millisecond,
			Exponential: true,
			Limit:       5,
		},
	}
}

// LastDate last possible date
var LastDate time.Time = time.Unix(1<<63-62135596801, 999999999)

// FirstDate first possible date
var FirstDate time.Time = time.Time{}
