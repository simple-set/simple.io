package handle

import (
	"encoding/json"
	"github.com/simple-set/simple.io/src/event"
)

type StringEncoder struct{}

func NewStringEncoder() *StringEncoder {
	return &StringEncoder{}
}

func (s StringEncoder) Output(_ *event.HandleContext, data interface{}) (interface{}, bool) {
	if value, ok := data.([]byte); ok {
		return value, true
	}
	if value, ok := data.(string); ok {
		return []byte(value), true
	}
	if bytes, err := json.Marshal(data); err == nil {
		return bytes, true
	}
	return nil, true
}
