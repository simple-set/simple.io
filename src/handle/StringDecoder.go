package handle

import (
	"github.com/simple-set/simple.io/src/event"
	"github.com/sirupsen/logrus"
)

type StringDecoder struct{}

func NewStringDecoder() *StringDecoder {
	return &StringDecoder{}
}

func (s *StringDecoder) Input(_ *event.HandleContext, data interface{}) (interface{}, bool) {
	if value, ok := data.([]byte); ok {
		return string(value), true
	} else if value, ok := data.(string); ok {
		return value, true
	}
	logrus.Warnf("Convert to Replace String")
	return nil, false
}
