package msgqueue

type Event interface {
	EventName() string
}

type EventEmitter interface {
	Emit(event Event) error
}

type EventListener interface {
	Listen(eventNames ...string) (<-chan Event, <-chan error, error)
}
