package msgqueue

type EventListener interface {
	Listen(eventNames ...string) (<-chan Event, <-chan error, error)
}
