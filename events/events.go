package events

// Events is events enum
type Events string

const (
	// TICK tick event
	TICK Events = "tick"
)

func (e Events) string() string {
	return string(e)
}
