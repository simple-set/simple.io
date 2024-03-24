package simpleHttp

import (
	"bufio"
	"bytes"
	"errors"
	"github.com/simple-set/simple.io/src/event"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

type ServerHandler struct {
	requestDecoder  *RequestDecoder
	responseEncoder *ResponseEncode
}

func (s *ServerHandler) Input(context *event.HandleContext, reader *bufio.Reader) (any, bool) {
	if request, err := s.requestDecoder.Decoder(reader); err != nil {
		logrus.Errorln("Decoding HTTP protocol error, ", err)
		_ = context.Session().Close()
		return nil, false
	} else {
		return request, true
	}
}

func (s *ServerHandler) Output(context *event.HandleContext, response *Response) (any, bool) {
	if response != nil && response.Header.Get("Date") == "" {
		response.Header.Add("Date", time.Now().Format(http.TimeFormat))
	}

	if err := s.responseEncoder.Encode(response); err != nil {
		logrus.Errorln("Encoding response exception, ", err)
		_ = context.Session().Close()
		return nil, false
	}
	if err := s.response(context, response); err != nil {
		logrus.Errorln("Return response exception, ", err)
		_ = context.Session().Close()
	}
	return nil, false
}

func (s *ServerHandler) response(context *event.HandleContext, response *Response) error {
	var buffer any = response.bufWriter.WriteBuffer()
	if value, ok := buffer.(*bytes.Buffer); ok {
		if _, err := context.Session().WriteSocket(value); err != nil {
			return err
		}
		if response.body != nil {
			if _, err := context.Session().WriteSocket(response.body); err != nil {
				return err
			}
		}
		return nil
	}
	return errors.New("the bufWriter type of the response is incorrect and cannot be read")
}

func NewServerHandler() *ServerHandler {
	return &ServerHandler{requestDecoder: NewRequestDecoder(), responseEncoder: NewResponseEncoded()}
}
