package simpleHttp

import (
	"bufio"
	"github.com/simple-set/simple.io/src/event"
	"github.com/sirupsen/logrus"
)

type ClientHandler struct {
	requestEncode  *RequestEncode
	responseDecode *ResponseDecode
}

func (c *ClientHandler) Input(context *event.HandleContext, reader *bufio.Reader) (any, bool) {
	if response, err := c.responseDecode.Decoder(reader); err != nil {
		logrus.Errorln(err)
		_ = context.Session().Close()
		return nil, false
	} else {
		return response, true
	}
}

func (c *ClientHandler) Output(context *event.HandleContext, request *Request) (any, bool) {
	if err := c.requestEncode.Encoder(request); err != nil {
		logrus.Errorln(err)
		_ = context.Session().Close()
		return nil, false
	}

	if err := c.sendRequest(context, request); err != nil {
		logrus.Errorln(err)
		_ = context.Session().Close()
		return nil, false
	}
	return nil, true
}

func (c *ClientHandler) sendRequest(context *event.HandleContext, request *Request) error {
	if _, err := context.Session().WriteSocket(request.buffWriter); err != nil {
		return err
	}
	if request.body != nil && request.body.size > 0 {
		if _, err := context.Session().WriteSocket(request.body); err != nil {
			return err
		}
	}
	return nil
}

func NewClientHandler() *ClientHandler {
	return &ClientHandler{requestEncode: NewRequestEncode(), responseDecode: NewResponseDecoded()}
}
