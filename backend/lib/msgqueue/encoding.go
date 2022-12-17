package msgqueue

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"reflect"
)

/////////////////////////////////////////////////// JSON

type jSONEncoder struct{}

func NewJSONEncoder() jSONEncoder {
	return jSONEncoder{}
}

func (jSONEncoder) Encode(payload interface{}) ([]byte, error) {
	return json.Marshal(payload)
}

type jsonDecoder struct{}

func NewJSONDecoder() jsonDecoder {
	return jsonDecoder{}
}

func (jsonDecoder) Decode(byt []byte, dest interface{}) error {
	return json.Unmarshal(byt, dest)
}

///////////////////////////////////////////////// GOB

type gobEncoder struct {
	encoder gob.Encoder
	buffer  bytes.Buffer
}

// Gob encoder
func NewGobEncoder(types ...reflect.Type) *gobEncoder {
	for _, typ := range types {
		gob.Register(typ)
	}
	var buffer bytes.Buffer

	return &gobEncoder{
		encoder: *gob.NewEncoder(&buffer),
		buffer:  buffer,
	}
}

func (e *gobEncoder) Encode(payload interface{}) ([]byte, error) {

	err := e.encoder.Encode(payload)

	return e.buffer.Bytes(), err
}

// Gob decoder
type gobDecoder struct {
	decoder gob.Decoder
	buffer  bytes.Buffer
}

func NewGobDecoder(types ...reflect.Type) *gobDecoder {
	for _, typ := range types {
		gob.Register(typ)
	}
	var buffer bytes.Buffer
	return &gobDecoder{
		decoder: *gob.NewDecoder(&buffer),
		buffer:  buffer,
	}
}

func (d *gobDecoder) Decode(byt []byte, dest interface{}) error {

	_, err := d.buffer.Write(byt)
	if err != nil {
		return err
	}

	return d.decoder.Decode(dest)
}
