package events

// Bus is an event bus
type Bus struct {
	listeners map[string][]chan interface{}
}

// New creates a new bus
func New() *Bus {
	return &Bus{
		listeners: make(map[string][]chan interface{}),
	}
}

// Subscribe subscribes to an event
func (e *Bus) Subscribe(event Events) chan interface{} {
	if _, ok := e.listeners[event.string()]; !ok {
		e.listeners[event.string()] = []chan interface{}{}
	}

	listener := make(chan interface{})
	e.listeners[event.string()] = append(e.listeners[event.string()], listener)
	return listener
}

// Emit emits an event
func (e *Bus) Emit(event Events, arg interface{}) {
	if _, ok := e.listeners[event.string()]; !ok {
		return
	}

	for _, v := range e.listeners[event.string()] {
		go func() {
			v <- arg
		}()
	}
}
