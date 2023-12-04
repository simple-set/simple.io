package handle

import (
	"github.com/simple-set/simple.io/src/event"
)

type EchoHandler struct {
}

func NewEchoHandler() *EchoHandler {
	return &EchoHandler{}
}

func (p *EchoHandler) Input(context *event.HandleContext, data any) (any, bool) {
	context.Session().WriteAndFlush(data)
	return data, true
}
