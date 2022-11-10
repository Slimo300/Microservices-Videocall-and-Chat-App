package msgqueue

type EventMapper interface {
	MapEvent(eventName string, eventPayload interface{}) (Event, error)
}
