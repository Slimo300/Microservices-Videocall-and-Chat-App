package msgqueue

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/mitchellh/mapstructure"
)

type DynamicEventMapper struct {
	typeMap map[string]reflect.Type
}

func NewDynamicEventMapper() DynamicEventMapper {
	return DynamicEventMapper{
		typeMap: make(map[string]reflect.Type),
	}
}

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

func (m *DynamicEventMapper) MapEvent(eventName string, eventPayload interface{}) (Event, error) {

	typ, ok := m.typeMap[eventName]
	if !ok {
		return nil, fmt.Errorf("No type with eventName: %s", eventName)
	}

	inter := reflect.New(typ).Interface()
	event, ok := inter.(Event)
	if !ok {
		return nil, fmt.Errorf("Type %s does not match Event interface", eventName)
	}

	cfg := mapstructure.DecoderConfig{
		Result:  event,
		TagName: "json",
	}

	dec, err := mapstructure.NewDecoder(&cfg)
	if err != nil {
		return nil, fmt.Errorf("Error when creating decoder: %s", err.Error())
	}

	if err = dec.Decode(eventPayload); err != nil {
		return nil, fmt.Errorf("Error when decoding: %s", err.Error())
	}

	return event, nil

}
