package simpleHttp

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"github.com/simple-set/simple.io/src/protocol/codec"
	"strings"
)

// RequestEncode http编码器, 把request对象编码成字节流, 根据 RFC2616 规范编码
// 适用于http客户端场景，把构建好的请求对象编码成字节流发送给服务器。
type RequestEncode struct{}

func (r *RequestEncode) Encoder(request *Request) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = errors.New(fmt.Sprintln("Error encoding request, ", e))
		}
	}()

	if request == nil {
		return errors.New("request cannot be nil")
	}
	if request.writerBuff == nil {
		buffer := bytes.NewBuffer(make([]byte, 0))
		request.writerBuff = codec.NewReadByteBuf(bufio.NewReader(buffer))
	}

	r.line(request)
	r.header(request)
	r.body(request)
	return nil
}

// 编码请求行
func (r *RequestEncode) line(request *Request) {
	request.writerBuff.WriteString(request.Method + " ")
	request.writerBuff.WriteString(request.RequestURI + " ")
	request.writerBuff.WriteString(request.Proto + " ")
	request.writerBuff.WriteBytes(crlf)
}

// 编码请求头
func (r *RequestEncode) header(request *Request) {
	if err := writeHeader(request.writerBuff.WriteBuffer(), "Host", request.Host); err != nil {
		panic(err)
	}
	if request.Header == nil || len(request.Header) == 0 {
		return
	}

	for name, values := range request.Header {
		if values == nil || len(values) == 0 {
			continue
		}
		if name == "Cookie" {
			r.cookie(request)
			continue
		}
		for _, value := range values {
			if err := writeHeader(request.writerBuff.WriteBuffer(), name, value); err != nil {
				panic(err)
			}
		}
	}
	request.writerBuff.WriteBytes(crlf)
}

// 编码cookie
func (r *RequestEncode) cookie(request *Request) {
	cookies := request.Cookies()
	if cookies == nil || len(cookies) == 0 {
		return
	}
	var cookieValues []string
	for i := 0; i < len(cookies); i++ {
		cookieValues = append(cookieValues, cookies[i].Name+"="+cookies[i].Value)
	}
	if err := writeHeader(request.writerBuff.WriteBuffer(), "Cookie", strings.Join(cookieValues, "; ")); err != nil {
		panic(err)
	}
}

// 编码请求体
func (r *RequestEncode) body(request *Request) {
	if request.Body != nil && request.body.size != request.ContentLength {
		panic("Body length error")
	}
}

func NewRequestEncode() *RequestEncode {
	return &RequestEncode{}
}
