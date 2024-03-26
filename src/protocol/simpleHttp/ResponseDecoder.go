package simpleHttp

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/simple-set/simple.io/src/protocol/codec"
	"net/http"
	"net/textproto"
	"strconv"
	"strings"
)

type ResponseDecode struct{}

func (r *ResponseDecode) Decoder(reader *bufio.Reader) (response *Response, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = errors.New(fmt.Sprintln("Error decoding request, ", e))
		}
	}()
	if reader == nil {
		return nil, errors.New("the buffer has no data to read")
	}
	response = NewResponseReader(codec.NewReadByteBuf(reader))
	r.line(response)
	r.header(response)
	r.contentLength(response)
	r.body(response)
	return response, nil
}

func (r *ResponseDecode) line(response *Response) {
	line := response.bufReader.ReadLine()

	proto, rest, ok1 := strings.Cut(line, " ")
	statusCode, statusText, ok2 := strings.Cut(rest, " ")
	if ok1 && ok2 {
		response.Proto = proto
		response.statusText = statusText
	} else {
		panic(errors.New("HTTP request line format error, " + line))
	}

	if code, err := strconv.Atoi(statusCode); err != nil {
		panic(errors.New("malformed HTTP statusCode " + statusCode))
	} else if http.StatusText(code) == "" {
		panic(errors.New("unknown HTTP statusCode " + statusCode))
	} else {
		response.statusCode = code
	}

	var ok bool
	if response.ProtoMajor, response.ProtoMinor, ok = http.ParseHTTPVersion(response.Proto); !ok {
		panic(errors.New("malformed HTTP version " + response.Proto))
	}
}

func (r *ResponseDecode) header(response *Response) {
	header := make(http.Header, 4)

	for {
		line := response.bufReader.ReadLine()
		if len(line) == 0 {
			break
		}
		if line[0] == ' ' || line[0] == '\t' {
			panic(errors.New("malformed MIME header initial line: " + line))
		}

		name, value, found := strings.Cut(line, string(colonSpace))
		name = strings.TrimSpace(name)
		value = strings.TrimSpace(value)
		if name == "" || value == "" || !found {
			panic(errors.New("Hearer format error： " + line))
		}

		if header.Get(name) == "" {
			header.Add(name, value)
		} else {
			header.Set(name, value)
		}
	}

	PragmaCacheControl(header)
	response.Header = header
	if header.Get("Server") != "" {
		response.Server = header.Get("Server")
	}
}

func (r *ResponseDecode) contentLength(response *Response) {
	contentLength := textproto.TrimString(response.Header.Get("Content-Length"))
	if contentLength == "" {
		response.contentLength = 0
		return
	}

	if n, err := strconv.ParseUint(contentLength, 10, 63); err != nil {
		response.contentLength = -1
		panic(errors.New("bad Content-Length: " + contentLength))
	} else {
		response.contentLength = int64(n)
	}
}

func (r *ResponseDecode) body(response *Response) {
	if response.contentLength > 0 {
		response.body = NewReadBody(response.contentLength, response.bufReader)
		return
	}
}

func NewResponseDecoded() *ResponseDecode {
	return &ResponseDecode{}
}
