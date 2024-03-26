package simpleHttp

import (
	"bufio"
	"github.com/simple-set/simple.io/src/event"
	"github.com/sirupsen/logrus"
)

type ClientHandler struct{}

func (c ClientHandler) Input(context *event.HandleContext, reader *bufio.Reader) (any, bool) {
	response := NewResponseReader(reader)
	if err := NewResponseDecoded(response).Decoder(); err != nil {
		logrus.Errorln(err)
		_ = context.Session().Close()
		return nil, false
	}
	return response, true
}

func (c ClientHandler) Output(context *event.HandleContext, request *Request) (any, bool) {
	//request.bufWriter = context.Session().Sock().Writer
	//if err := NewRequestEncode(request).Encoder(); err != nil {
	//	logrus.Errorln(err)
	//	_ = context.Session().Close()
	//	return nil, false
	//}
	//context.Session().Flush()
	return request, true
}

func NewClientHandler() *ClientHandler {
	return &ClientHandler{}
}
