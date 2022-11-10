package msgqueue

type EventEmitter interface {
	Emit(event Event) error
}
