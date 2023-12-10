package handle

import (
	"bufio"
	"github.com/simple-set/simple.io/src/event"
	"github.com/sirupsen/logrus"
)

type StringDecoder struct{}

func NewStringDecoder() *StringDecoder { return &StringDecoder{} }

func (s *StringDecoder) Input(_ *event.HandleContext, reader *bufio.Reader) (any, bool) {
	size := reader.Buffered()
	if size <= 0 {
		return nil, false
	}

	bytes := make([]byte, reader.Buffered())
	if _, err := reader.Read(bytes); err == nil {
		return string(bytes), true
	} else {
		logrus.Errorln("Exception reading data from buffer", err)
		return nil, false
	}
}
