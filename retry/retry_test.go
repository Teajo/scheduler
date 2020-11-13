package retry

import (
	"errors"
	"testing"
	"time"
)

func TestRetryWithFullErrors(t *testing.T) {
	i := 0

	err := Do(func() error {
		i++
		return errors.New("random error")
	}, 3, 100*time.Millisecond, true)

	if err == nil {
		t.Error("err should be returned")
	}

	if i != 3 {
		t.Error("func should be called 3 times")
	}
}

func TestRetryWithSuccess(t *testing.T) {
	err := Do(func() error {
		return nil
	}, 3, 100*time.Millisecond, false)

	if err != nil {
		t.Error("err should not be returned")
	}
}
