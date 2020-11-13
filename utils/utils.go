package utils

import (
	"time"

	"github.com/google/uuid"
)

// RetryStrat represents a retry strategy
type RetryStrat struct {
	Timeout     time.Duration
	Exponential bool
	Limit       int
}

// Scheduling represents a scheduling object
type Scheduling struct {
	ID        string
	Date      time.Time
	Publisher string
	Settings  map[string]string
	*RetryStrat
}

// NewScheduling creates a new scheduling struct
func NewScheduling(date time.Time, publisher string, settings map[string]string) *Scheduling {
	return &Scheduling{
		ID:        uuid.New().String(),
		Date:      date,
		Publisher: publisher,
		Settings:  settings,
		RetryStrat: &RetryStrat{
			Timeout:     25 * time.Millisecond,
			Exponential: true,
			Limit:       5,
		},
	}
}

// NewSchedulingWithID creates a new scheduling struct
func NewSchedulingWithID(id string, date time.Time, publisher string, settings map[string]string) *Scheduling {
	return &Scheduling{
		ID:        id,
		Date:      date,
		Publisher: publisher,
		Settings:  settings,
		RetryStrat: &RetryStrat{
			Timeout:     25 * time.Millisecond,
			Exponential: true,
			Limit:       5,
		},
	}
}
