package simpleRpc

import (
	"bufio"
	"bytes"
	"github.com/simple-set/simple.io/src/event"
	"github.com/sirupsen/logrus"
)

type Handle struct {
	encoder *Encoder
	decoder *Decoder
}

func (h Handle) Input(context *event.HandleContext, reader *bufio.Reader) (*Message, bool) {
	if message, err := h.decoder.Decoder(reader); err == nil {
		logrus.Error(err)
		_ = context.Session().Close()
		return nil, false
	} else {
		return message, true
	}
}

func (h Handle) Output(context *event.HandleContext, message *Message) (any, bool) {
	buffer := bytes.NewBuffer(make([]byte, 0, message.MessageLength()))

	if err := h.encoder.Encoder(message, bufio.NewWriter(buffer)); err != nil {
		logrus.Errorln(err)
		_ = context.Session().Close()
		return nil, false
	}
	return buffer, true
}

func NewHandle() *Handle {
	return &Handle{decoder: NewDecoder(), encoder: NewEncoder()}
}
