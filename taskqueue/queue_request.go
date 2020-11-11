package taskqueue

type method string

const (
	add    method = "add"
	remove method = "remove"
	cancel method = "cancel"
)

// QueueRequest represents a request
type QueueRequest struct {
	method  method
	payload interface{}
	err     chan error
}

// NewQueueRequest creates new request
func NewQueueRequest(method method, payload interface{}) (*QueueRequest, chan error) {
	err := make(chan error)
	return &QueueRequest{
		method:  method,
		payload: payload,
		err:     err,
	}, err
}
