package events

// Events is events enum
type Events string

const (
	// TICK tick event
	TICK Events = "tick"

	// TASKDONE task is done event
	TASKDONE Events = "task_done"
)

func (e Events) string() string {
	return string(e)
}
