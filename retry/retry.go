package retry

import (
	"errors"
	"time"
)

// Do task with retry strategy
func Do(f func() error, limit int, timeout time.Duration, exponential bool) error {
	if limit > 0 {
		err := f()
		if err != nil {
			<-time.After(timeout)
			if exponential {
				timeout *= 2
			}
			return Do(f, limit-1, timeout, exponential)
		}
		return nil
	}
	return errors.New("Unable do execute task despite all retries")
}
