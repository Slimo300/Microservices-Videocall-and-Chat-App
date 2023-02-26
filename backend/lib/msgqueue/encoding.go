package msgqueue

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"reflect"
)

/////////////////////////////////////////////////// JSON

type jSONEncoder struct{}

// NewJSONEncoder creates json encoder
func NewJSONEncoder() Encoder {
	return jSONEncoder{}
}

// Encode encodes payload to its json representation
func (jSONEncoder) Encode(payload interface{}) ([]byte, error) {
	return json.Marshal(payload)
}

type jsonDecoder struct{}

// NewJSONDecoder creates new json decoder
func NewJSONDecoder() Decoder {
	return jsonDecoder{}
}

// Decode decodes json into given destination
func (jsonDecoder) Decode(byt []byte, dest interface{}) error {
	return json.Unmarshal(byt, dest)
}

///////////////////////////////////////////////// GOB

type gobEncoder struct {
	encoder gob.Encoder
	buffer  bytes.Buffer
}

// NewGobEncoder creates new gob encoder
func NewGobEncoder(types ...reflect.Type) Encoder {
	for _, typ := range types {
		gob.Register(typ)
	}
	var buffer bytes.Buffer

	return &gobEncoder{
		encoder: *gob.NewEncoder(&buffer),
		buffer:  buffer,
	}
}

// Encode encodes payload to its gob representation
func (e *gobEncoder) Encode(payload interface{}) ([]byte, error) {

	err := e.encoder.Encode(payload)

	return e.buffer.Bytes(), err
}

type gobDecoder struct {
	decoder gob.Decoder
	buffer  bytes.Buffer
}

// NewGobDecoder creates new gob decoder
func NewGobDecoder(types ...reflect.Type) Decoder {
	for _, typ := range types {
		gob.Register(typ)
	}
	var buffer bytes.Buffer
	return &gobDecoder{
		decoder: *gob.NewDecoder(&buffer),
		buffer:  buffer,
	}
}

// Decode decodes gob into given destination
func (d *gobDecoder) Decode(byt []byte, dest interface{}) error {

	_, err := d.buffer.Write(byt)
	if err != nil {
		return err
	}

	return d.decoder.Decode(dest)
}
