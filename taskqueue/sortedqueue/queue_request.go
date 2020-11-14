package sortedqueue

// Response is a request response
type Response struct {
	err     error
	payload interface{}
}

// NewResponse creates a new response
func NewResponse(err error, payload interface{}) *Response {
	return &Response{
		err:     err,
		payload: payload,
	}
}

// QueueRequest represents a request
type QueueRequest struct {
	method  method
	payload interface{}
	res     chan *Response
}

// NewRequest creates new request
func NewRequest(method method, payload interface{}) (*QueueRequest, chan *Response) {
	res := make(chan *Response)
	return &QueueRequest{
		method:  method,
		payload: payload,
		res:     res,
	}, res
}
