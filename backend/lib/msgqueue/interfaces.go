package msgqueue

// Event is any type to be send by message broker
type Event interface {
	EventName() string
}

// EventEmiter is type able to send events to message broker
type EventEmiter interface {
	Emit(event Event) error
}

// EventListener is able to forward events with given names from message broker
type EventListener interface {
	Listen(eventNames ...string) (<-chan Event, <-chan error, error)
}

// EventMapper is able to map serialized event to its type representation
type EventMapper interface {
	MapEvent(eventName string, eventPayload interface{}) (Event, error)
}

// Encoder encodes data to its format, currently there are json and gob representations
type Encoder interface {
	Encode(payload interface{}) ([]byte, error)
}

// Decoder decodes data from its format, currently there are json and gob representations
type Decoder interface {
	Decode([]byte, interface{}) error
}
