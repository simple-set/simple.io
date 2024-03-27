package simpleHttp

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"github.com/simple-set/simple.io/src/protocol/codec"
	"strconv"
	"strings"
)

// RequestEncode http编码器, 把request对象编码成字节流, 根据 RFC2616 规范编码
// 适用于http客户端场景，把构建好的请求对象编码成字节流发送给服务器。
type RequestEncode struct{}

func (r *RequestEncode) Encoder(request *Request) (err error) {
	if request == nil {
		return errors.New("request cannot be nil")
	}
	if request.buffWriter == nil {
		buffer := bytes.NewBuffer(make([]byte, 0))
		reader := bufio.NewReaderSize(buffer, 32)
		writer := bufio.NewWriterSize(buffer, 32)
		request.buffWriter = codec.NewReadWriteByteBuf(reader, writer)
	}
	return r.encodeFrame(request)
}

func (r *RequestEncode) encodeFrame(request *Request) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = errors.New(fmt.Sprintln("Error encoding request, ", e))
		}
	}()

	r.line(request)
	r.header(request)
	r.body(request)
	_ = request.buffWriter.Writer.Flush()
	return nil
}

// 编码请求行
func (r *RequestEncode) line(request *Request) {
	request.buffWriter.WriteString(request.Method + " ")
	request.buffWriter.WriteString(request.RequestURI + " ")
	request.buffWriter.WriteString(request.Proto)
	request.buffWriter.WriteBytes(crlf)
}

// 编码请求头
func (r *RequestEncode) header(request *Request) {
	writeBuffer := request.buffWriter.WriteBuffer()
	if err := writeHeader(writeBuffer, "Host", request.Host); err != nil {
		panic(err)
	}
	if request.Header != nil || len(request.Header) > 0 {
		for name, values := range request.Header {
			if values == nil || len(values) == 0 {
				continue
			}
			if name == "Cookie" {
				r.cookie(request)
				continue
			}
			for _, value := range values {
				if err := writeHeader(writeBuffer, name, value); err != nil {
					panic(err)
				}
			}
		}
	}

	if request.body != nil && request.body.size > 0 {
		if err := writeHeader(writeBuffer, "Content-Length", strconv.FormatInt(request.body.size, 10)); err != nil {
			panic(err)
		}
	}
	request.buffWriter.WriteBytes(crlf)
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
	if err := writeHeader(request.buffWriter.WriteBuffer(), "Cookie", strings.Join(cookieValues, "; ")); err != nil {
		panic(err)
	}
}

// 编码请求体
func (r *RequestEncode) body(request *Request) {
	if request.body != nil && request.body.size != request.ContentLength {
		panic("Body length error")
	}
}

func NewRequestEncode() *RequestEncode {
	return &RequestEncode{}
}
