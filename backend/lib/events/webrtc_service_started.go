package events

// ServiceStartedEvent holds information about new webrtc service being created
type ServiceStartedEvent struct {
	ServiceAddress string `json:"address" mapstructure:"address"`
}

// EventName fulfills msgqueue.Event interface
func (ServiceStartedEvent) EventName() string {
	return "webrtc.service-start"
}
