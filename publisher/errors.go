package publisher

// PublishError represents publish error
type PublishError struct {
	err         error
	shouldRetry bool
}

// NewPublishError creates a publish error
func NewPublishError(err error, shouldRetry bool) *PublishError {
	return &PublishError{
		err:         err,
		shouldRetry: shouldRetry,
	}
}

func (pe *PublishError) Error() string {
	return pe.err.Error()
}

// ShouldRetry indicates that should retry publish
func (pe *PublishError) ShouldRetry() bool {
	return pe.shouldRetry
}
