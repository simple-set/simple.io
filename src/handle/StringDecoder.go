package handle

import (
	"bufio"
	"github.com/simple-set/simple.io/src/event"
	"github.com/sirupsen/logrus"
)

type StringDecoder struct{}

func NewStringDecoder() *StringDecoder { return &StringDecoder{} }

func (s *StringDecoder) Input(_ *event.HandleContext, reader *bufio.Reader) (any, bool) {
	bytes := make([]byte, reader.Size())

	if n, err := reader.Read(bytes); err == nil {
		return string(bytes[:n]), true
	} else {
		logrus.Errorln("Exception reading data from buffer", err)
		return "", false
	}
}
