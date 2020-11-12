package utils

import (
	"time"

	"github.com/google/uuid"
)

// Scheduling represents a scheduling object
type Scheduling struct {
	ID        string
	Date      time.Time
	Publisher string
	Settings  map[string]string
}

// NewScheduling creates a new scheduling struct
func NewScheduling(date time.Time, publisher string, settings map[string]string) *Scheduling {
	return &Scheduling{
		ID:        uuid.New().String(),
		Date:      date,
		Publisher: publisher,
		Settings:  settings,
	}
}
