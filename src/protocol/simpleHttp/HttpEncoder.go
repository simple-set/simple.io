package simpleHttp

import (
	"bufio"
	"github.com/simple-set/simple.io/src/event"
	"github.com/sirupsen/logrus"
)

// HttpEncoder http编码器, 适用于客户端模式
type HttpEncoder struct{}

func NewHttpEncoder() *HttpEncoder {
	return &HttpEncoder{}
}

func (h HttpEncoder) Input(context *event.HandleContext, reader *bufio.Reader) (*Response, bool) {
	return nil, false
}

func (h HttpEncoder) Output(context *event.HandleContext, request *Request) (any, bool) {
	request.bufWriter = context.Session().Sock().Writer
	if err := NewRequestEncode(request).Encoder(); err != nil {
		logrus.Errorln(err)
		_ = context.Session().Close()
		return nil, false
	}
	return nil, true
}
