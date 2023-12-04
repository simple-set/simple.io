package handle

import (
	"bufio"
	"github.com/simple-set/simple.io/src/event"
	"github.com/sirupsen/logrus"
)

type StringDecoder struct{}

func (s *StringDecoder) Input(_ *event.HandleContext, reader bufio.Reader) (string, bool) {
	bytes := make([]byte, reader.Size())
	if _, err := reader.Read(bytes); err != nil {
		logrus.Errorln("Exception reading data from buffer", err)
	}
	return string(bytes), false
}

func NewStringDecoder() *StringDecoder {
	return &StringDecoder{}
}
