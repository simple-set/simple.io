package simpleHttp

import (
	"bufio"
	"github.com/simple-set/simple.io/src/event"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

type ServerHandler struct{}

func (s ServerHandler) Input(context *event.HandleContext, reader *bufio.Reader) (any, bool) {
	request := NewRequestReader(reader)
	if err := NewRequestDecoder(request).Decoder(); err != nil {
		logrus.Errorln("Decoding HTTP protocol error, ", err)
		_ = context.Session().Close()
		return nil, false
	}

	request.Response = NewReplyResponse(request)
	return request, true
}

func (s ServerHandler) Output(context *event.HandleContext, data any) (any, bool) {
	var response *Response
	if value, ok := data.(*Request); ok && value.Response != nil {
		response = value.Response
	} else if value, ok := data.(*Response); ok {
		response = value
	}

	if response == nil {
		logrus.Warnln("HTTP encoder execution failed with no available response")
		_ = context.Session().Close()
		return nil, false
	}
	if response.Header.Get("Date") == "" {
		response.Header.Add("Date", time.Now().Format(http.TimeFormat))
	}

	response.bufWriter = context.Session().Sock().Writer
	if err := NewResponseEncoded(response).Codec(); err != nil {
		logrus.Errorln("HTTP encoder execution failed, ", err)
	}
	context.Session().Flush()
	return nil, false
}

func NewServerHandler() *ServerHandler {
	return &ServerHandler{}
}
