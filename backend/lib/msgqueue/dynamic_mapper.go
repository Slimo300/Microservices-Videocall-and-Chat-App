package msgqueue

import (
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/google/uuid"
	"github.com/mitchellh/mapstructure"
)

// DynamicEventMapper allows to define event types in runtime
type DynamicEventMapper struct {
	typeMap map[string]reflect.Type
}

// NewDynamicEventMapper creates new mapper
func NewDynamicEventMapper() *DynamicEventMapper {
	return &DynamicEventMapper{
		typeMap: make(map[string]reflect.Type),
	}
}

// RegisterEventType adds new type to list for later mapping
func (m *DynamicEventMapper) RegisterEventType(eventType reflect.Type) error {

	emptyEvent := reflect.New(eventType)
	inter := emptyEvent.Interface()

	newEvent, ok := inter.(Event)
	if !ok {
		return errors.New("eventType does not match Event interface")
	}

	m.typeMap[newEvent.EventName()] = eventType
	return nil
}

// RegisterTypes allows to register any number of types
func (m *DynamicEventMapper) RegisterTypes(eventTypes ...reflect.Type) error {
	for _, typ := range eventTypes {
		if err := m.RegisterEventType(typ); err != nil {
			return fmt.Errorf("error when registering type %s: %s", typ.Name(), err.Error())
		}
	}
	return nil
}

// MapEvent searches for event with given event name and if it finds it, it returns event payload
// decoded to this type, decoded value is returned as an instance of Event interface
func (m *DynamicEventMapper) MapEvent(eventName string, eventPayload interface{}) (Event, error) {

	typ, ok := m.typeMap[eventName]
	if !ok {
		return nil, fmt.Errorf("no type with eventName: %s", eventName)
	}

	inter := reflect.New(typ).Interface()
	event, ok := inter.(Event)
	if !ok {
		return nil, fmt.Errorf("type %s does not match Event interface", eventName)
	}

	cfg := mapstructure.DecoderConfig{
		DecodeHook: mapstructure.ComposeDecodeHookFunc(
			stringToUUIDHookFunc(),
			stringToTimeHookFunc(),
		),
		Result: event,
	}

	dec, err := mapstructure.NewDecoder(&cfg)
	if err != nil {
		return nil, fmt.Errorf("error when creating decoder: %s", err.Error())
	}

	if err = dec.Decode(eventPayload); err != nil {
		return nil, fmt.Errorf("error when decoding: %s", err.Error())
	}

	return event, nil

}

func stringToUUIDHookFunc() mapstructure.DecodeHookFunc {
	return func(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
		if f.Kind() != reflect.String {
			return data, nil
		}
		if t != reflect.TypeOf(uuid.UUID{}) {
			return data, nil
		}

		return uuid.Parse(data.(string))
	}
}

func stringToTimeHookFunc() mapstructure.DecodeHookFunc {
	return func(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
		if f.Kind() != reflect.String {
			return data, nil
		}
		if t != reflect.TypeOf(time.Time{}) {
			return data, nil
		}

		return time.Parse(time.RFC3339, data.(string))
	}
}
