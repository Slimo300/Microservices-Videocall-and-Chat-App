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

type EventMapper interface {
	MapEvent(eventName string, eventPayload interface{}) (Event, error)
}

type Encoder interface {
	Encode(payload interface{}) ([]byte, error)
}

type Decoder interface {
	Decode([]byte, interface{}) error
}
