package simpleHttp

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"github.com/simple-set/simple.io/src/protocol/codec"
	"strconv"
)

type ResponseEncode struct{}

func (r *ResponseEncode) Encode(response *Response) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = errors.New(fmt.Sprintln("Error encoding request, ", e))
		}
	}()

	if response == nil {
		return errors.New("response cannot be nil")
	}
	if response.bufWriter == nil {
		buffer := bytes.NewBuffer(make([]byte, 0))
		reader := bufio.NewReaderSize(buffer, 32)
		writer := bufio.NewWriterSize(buffer, 32)
		response.bufWriter = codec.NewReadWriteByteBuf(reader, writer)
	}
	return r.encodeFrame(response)
}

func (r *ResponseEncode) encodeFrame(response *Response) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = errors.New(fmt.Sprintln("Error encoding request, ", e))
		}
	}()
	r.line(response)
	r.header(response)
	r.body(response)
	response.bufWriter.Flush()
	return nil
}

func (r *ResponseEncode) line(response *Response) {
	response.bufWriter.WriteString(response.Proto + " ")
	response.bufWriter.WriteString(strconv.Itoa(response.statusCode) + " ")
	response.bufWriter.WriteString(response.statusText)
	response.bufWriter.WriteBytes(crlf)
}

func (r *ResponseEncode) header(response *Response) {
	buffer := response.bufWriter.WriteBuffer()

	if response.Server != "" {
		if err := writeHeader(buffer, "Server", response.Server); err != nil {
			panic(err)
		}
	}
	if response.Body().Size() > 0 {
		if err := writeHeader(buffer, "Content-Length", strconv.FormatInt(response.Body().Size(), 10)); err != nil {
			panic(err)
		}
	}

	if response.Header != nil {
		for name, values := range response.Header {
			if values == nil || len(values) == 0 {
				continue
			}
			for i := 0; i < len(values); i++ {
				if name == "Cookie" {
					r.cookie(response)
					continue
				}
				if err := writeHeader(buffer, name, values[i]); err != nil {
					panic(err)
				}
			}
		}
	}
	response.bufWriter.WriteBytes(crlf)
}

// 编码cookie
func (r *ResponseEncode) cookie(response *Response) {
	cookies := response.Cookies()
	if cookies == nil {
		return
	}
	for _, cookie := range cookies[:] {
		if err := writeHeader(response.bufWriter.WriteBuffer(), "Set-Cookie", cookie.String()); err != nil {
			panic(err)
		}
	}
}

func (r *ResponseEncode) body(_ *Response) {
	// TODO
}

func NewResponseEncoded() *ResponseEncode {
	return &ResponseEncode{}
}
