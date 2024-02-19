package simpleHttp

import (
	"bufio"
	"github.com/simple-set/simple.io/src/event"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

// HttpDecoder http解码器, 适用于服务端模式
type HttpDecoder struct{}

func (h *HttpDecoder) Input(context *event.HandleContext, reader *bufio.Reader) (*Request, bool) {
	request := NewRequestReader(reader)
	err := NewRequestDecoder(request).Decoder()
	if err != nil {
		logrus.Errorln("Decoding HTTP protocol error, ", err)
		_ = context.Session().Close()
		return nil, false
	}

	request.Response = NewReplyResponse(request)
	return request, true
}

func (h *HttpDecoder) Output(context *event.HandleContext, data any) (any, bool) {
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

// HttpEncoder http编码器, 适用于客户端模式
type HttpEncoder struct{}

func (h HttpEncoder) Input(context *event.HandleContext, reader *bufio.Reader) (*Response, bool) {
	response := NewResponseReader(reader)
	err := NewResponseDecoded(response).Decoder()
	if err != nil {
		logrus.Errorln(err)
		_ = context.Session().Close()
	}
	return response, true
}

func (h HttpEncoder) Output(context *event.HandleContext, request *Request) (any, bool) {
	request.bufWriter = context.Session().Sock().Writer
	if err := NewRequestEncode(request).Encoder(); err != nil {
		logrus.Errorln(err)
		_ = context.Session().Close()
		return nil, false
	}
	return request, true
}

func NewHttpEncoder() *HttpEncoder { return &HttpEncoder{} }

func NewHttpDecoder() *HttpDecoder { return &HttpDecoder{} }
