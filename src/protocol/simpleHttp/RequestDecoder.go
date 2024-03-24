package simpleHttp

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/simple-set/simple.io/src/protocol/codec"
	"golang.org/x/net/http/httpguts"
	"net/http"
	"net/textproto"
	"net/url"
	"strconv"
	"strings"
)

// RequestDecoder 请求解码器, 从io缓冲区读取字节流, 并解码为http请求对象, 缓冲区一般为socket连接, 也可以使字节数组等
type RequestDecoder struct{}

// Decoder 解码入口
func (r *RequestDecoder) Decoder(reader *bufio.Reader) (request *Request, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = errors.New(fmt.Sprintln("Error decoding request, ", e))
		}
	}()

	if reader == nil {
		return nil, errors.New("the buffer has no data to read")
	}
	request = NewRequest(codec.NewReadByteBuf(reader))
	r.line(request)
	r.header(request)
	r.uri(request)
	r.contentLength(request)
	r.body(request)
	return request, nil
}

// 解码请求行
func (r *RequestDecoder) line(request *Request) {
	line := request.readBuff.ReadLine()
	method, rest, ok1 := strings.Cut(line, " ")
	requestURI, proto, ok2 := strings.Cut(rest, " ")
	if ok1 && ok2 {
		request.Method = method
		request.RequestURI = requestURI
		request.Proto = proto
	} else {
		panic(errors.New("HTTP request line format error, " + string(line)))
	}

	if len(request.Method) < 0 || !ValidMethod(request.Method) || !ValidPath(request.RequestURI) {
		panic(errors.New("invalid request " + request.Method + " " + request.RequestURI))
	}

	var ok bool
	if request.ProtoMajor, request.ProtoMinor, ok = http.ParseHTTPVersion(request.Proto); !ok {
		panic(errors.New("malformed HTTP version " + request.Proto))
	}
}

// 解码请求头
func (r *RequestDecoder) header(request *Request) {
	header := make(http.Header, 4)

	for {
		line := request.readBuff.ReadLine()
		if len(line) == 0 {
			break
		}
		if line[0] == ' ' || line[0] == '\t' {
			panic(badRequestError("malformed MIME header initial line: " + line))
		}

		name, value, found := strings.Cut(line, ":")
		name = strings.TrimSpace(name)
		value = strings.TrimSpace(value)

		if name == "" || value == "" || !found {
			panic(badRequestError("Hearer format error： " + line))
		}
		if !httpguts.ValidHeaderFieldName(name) {
			panic(badRequestError("invalid header name: " + value))
		}
		if !httpguts.ValidHeaderFieldValue(value) {
			panic(badRequestError("invalid header value: " + value))
		}

		if header.Get(name) == "" {
			header.Add(name, value)
		} else {
			header.Set(name, value)
		}
	}

	// 调整缓存头
	PragmaCacheControl(header)
	request.Header = header
	request.Host = header.Get("Host")
	// 解析表单数据
	_ = request.ParseForm()
}

// 解码请求URL
func (r *RequestDecoder) uri(request *Request) {
	if uri, err := url.ParseRequestURI("http://" + request.Host + request.RequestURI); err != nil {
		panic(err)
	} else {
		request.URL = uri
	}
}

// 解码请求体长度
func (r *RequestDecoder) contentLength(request *Request) {
	contentLength := textproto.TrimString(request.Header.Get("Content-Length"))
	if contentLength == "" {
		request.ContentLength = 0
	} else if n, err := strconv.ParseUint(contentLength, 10, 63); err != nil {
		request.ContentLength = -1
		panic(errors.New("bad Content-Length: " + contentLength))
	} else {
		request.ContentLength = int64(n)
	}
}

// 解码请求体
func (r *RequestDecoder) body(request *Request) {
	if request.ContentLength > 0 {
		request.body = NewBody(request.ContentLength, request.readBuff)
	}
}

// NewRequestDecoder 构造函数
func NewRequestDecoder() *RequestDecoder {
	return &RequestDecoder{}
}
