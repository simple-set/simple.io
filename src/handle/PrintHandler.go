package handle

import (
	"github.com/simple-set/simple.io/src/event"
	"github.com/sirupsen/logrus"
)

type PrintHandler struct{}

func NewPrintHandler() *PrintHandler {
	return &PrintHandler{}
}

func (p PrintHandler) Input(context *event.HandleContext, data interface{}) (interface{}, bool) {
	session := context.Session()
	logrus.Printf("RemoteAddr: %s, %s, received data: %v", session.Sock().RemoteAddr(), session.Id(), data)
	return data, true
}
