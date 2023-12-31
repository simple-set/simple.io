package simpleHttp

import (
	"errors"
	"golang.org/x/net/http/httpguts"
	"net/http"
	"net/textproto"
	"net/url"
	"strconv"
	"strings"
)

// RequestDecoder 请求解码器, 从io缓冲区读取字节流, 并解码为http请求对象, 缓冲区一般为socket连接, 也可以使字节数组等
type RequestDecoder struct {
	request *Request
}

// Decoder 解码入口
func (r *RequestDecoder) Decoder() error {
	if r.request == nil || r.request.bufReader == nil {
		return errors.New("the buffer has no data to read")
	}

	if err := r.line(); err != nil {
		return err
	}
	if err := r.header(); err != nil {
		return err
	}
	if err := r.uri(); err != nil {
		return err
	}
	if err := r.contentLength(); err != nil {
		return err
	}
	if err := r.body(); err != nil {
		return err
	}
	return nil
}

// 解码请求行
func (r *RequestDecoder) line() error {
	line, err := readLine(r.request.bufReader)
	if err != nil {
		return err
	}

	method, rest, ok1 := strings.Cut(string(line), " ")
	requestURI, proto, ok2 := strings.Cut(rest, " ")
	if ok1 && ok2 {
		r.request.Method = method
		r.request.RequestURI = requestURI
		r.request.Proto = proto
	} else {
		return errors.New("HTTP request line format error, " + string(line))
	}

	if len(r.request.Method) < 0 || !ValidMethod(r.request.Method) || !ValidPath(r.request.RequestURI) {
		return errors.New("invalid request " + r.request.Method + " " + r.request.RequestURI)
	}
	var ok bool
	if r.request.ProtoMajor, r.request.ProtoMinor, ok = http.ParseHTTPVersion(r.request.Proto); !ok {
		return errors.New("malformed HTTP version " + r.request.Proto)
	}
	return nil
}

// 解码请求头
func (r *RequestDecoder) header() error {
	header := make(http.Header, 4)

	for {
		line, err := readLine(r.request.bufReader)
		if err != nil {
			return err
		}
		if len(line) == 0 {
			break
		}
		if line[0] == ' ' || line[0] == '\t' {
			return badRequestError("malformed MIME header initial line: " + string(line))

		}

		name, value, found := strings.Cut(string(line), ":")
		name = strings.TrimSpace(name)
		value = strings.TrimSpace(value)

		if name == "" || value == "" || !found {
			return badRequestError("Hearer format error： " + string(line))
		}
		if !httpguts.ValidHeaderFieldName(name) {
			return badRequestError("invalid header name: " + value)
		}
		if !httpguts.ValidHeaderFieldValue(value) {
			return badRequestError("invalid header value: " + value)
		}

		if header.Get(name) == "" {
			header.Add(name, value)
		} else {
			header.Set(name, value)
		}
	}

	// 调整缓存头
	PragmaCacheControl(header)
	r.request.Header = header
	r.request.Host = header.Get("Host")
	// 解析表单数据
	_ = r.request.ParseForm()
	return nil
}

// 解码请求URL
func (r *RequestDecoder) uri() error {
	uri, err := url.ParseRequestURI("http://" + r.request.Host + r.request.RequestURI)
	if err != nil {
		return err
	}
	r.request.URL = uri
	return nil
}

// 解码请求体长度
func (r *RequestDecoder) contentLength() error {
	contentLength := textproto.TrimString(r.request.Header.Get("Content-Length"))
	if contentLength == "" {
		r.request.ContentLength = 0
		return nil
	}

	if n, err := strconv.ParseUint(contentLength, 10, 63); err != nil {
		r.request.ContentLength = -1
		return errors.New("bad Content-Length: " + contentLength)
	} else {
		r.request.ContentLength = int64(n)
	}
	return nil
}

// 解码请求体
func (r *RequestDecoder) body() error {
	if r.request.ContentLength <= 0 {
		r.request.Body = NewBody([]byte{})
		return nil
	}

	body := make([]byte, r.request.ContentLength)
	readLength, err := r.request.bufReader.Read(body)
	if err != nil {
		return err
	}
	if int64(readLength) != r.request.ContentLength {
		return errors.New("failed to read request Body")
	}
	r.request.Body = NewBody(body)
	return nil
}

// NewRequestDecoder 构造函数
func NewRequestDecoder(request *Request) *RequestDecoder {
	return &RequestDecoder{request: request}
}
