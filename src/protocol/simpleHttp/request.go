package simpleHttp

import (
	"bufio"
	"errors"
	"net/http"
	"net/textproto"
	"net/url"
	"strconv"
	"strings"
)

type Request struct {
	http.Request
	Body     *Body
	Response *Response
}

// MakeRequest 从字节缓冲对象读取解析，并构建http请求对象
func MakeRequest(reader *bufio.Reader) (*Request, error) {
	if reader == nil {
		return nil, errors.New("the buffer has no data to read")
	}

	request := new(Request)
	if err := ParseRequestLine(request, reader); err != nil {
		return request, err
	}
	if err := ParseRequestHeader(request, reader); err != nil {
		return request, err
	}
	if err := ParseRequestURI(request); err != nil {
		return request, err
	}
	if err := ParseContentLength(request); err != nil {
		return request, err
	}
	if err := ParseRequestBody(request, reader); err != nil {
		return request, err
	}

	request.Response = bindResponse(request)
	return request, nil
}

// ParseRequestLine 解析http请求行
func ParseRequestLine(request *Request, buf *bufio.Reader) error {
	line, err := readLine(buf)
	if err != nil {
		return err
	}

	method, rest, ok1 := strings.Cut(string(line), " ")
	requestURI, proto, ok2 := strings.Cut(rest, " ")
	if ok1 && ok2 {
		request.Method = method
		request.RequestURI = requestURI
		request.Proto = proto
	} else {
		return errors.New("HTTP request line format error, " + string(line))
	}

	if len(request.Method) < 0 || !ValidMethod(request.Method) || !ValidPath(request.RequestURI) {
		return errors.New("invalid request " + request.Method + " " + request.RequestURI)
	}
	var ok bool
	if request.ProtoMajor, request.ProtoMinor, ok = http.ParseHTTPVersion(request.Proto); !ok {
		return errors.New("malformed HTTP version " + request.Proto)
	}
	return nil
}

// ParseRequestHeader 解析http请求头
func ParseRequestHeader(request *Request, buf *bufio.Reader) error {
	header := make(http.Header, 4)

	for {
		line, err := readLine(buf)
		if err != nil {
			return err
		}
		if len(line) == 0 {
			break
		}
		if line[0] == ' ' || line[0] == '\t' {
			return errors.New("malformed MIME header initial line: " + string(line))
		}

		name, value, found := strings.Cut(string(line), ":")
		name = strings.TrimSpace(name)
		value = strings.TrimSpace(value)
		if name == "" || value == "" || !found {
			return errors.New("Hearer format error： " + string(line))
		}

		if header.Get(name) == "" {
			header.Add(name, value)
		} else {
			header.Set(name, value)
		}
	}

	PragmaCacheControl(header)
	////!httpguts.ValidHostHeader(hosts[0])
	request.Header = header
	request.Host = header.Get("Host")
	_ = request.ParseForm()
	return nil
}

func encodeCookie(request *Request, cookie string) {
	cookies := strings.Split(cookie, ";")
	for i := 0; i < len(cookies); i++ {
		before, after, found := strings.Cut(strings.TrimSpace(cookies[i]), "=")
		if before != "" || after != "" || !found {
			request.AddCookie(&http.Cookie{Name: before, Value: after})
		}
	}
}

func ParseRequestURI(request *Request) error {
	if uri, err := url.ParseRequestURI("http://" + request.Host + request.RequestURI); err == nil {
		request.URL = uri
	}
	return nil
}

// ParseContentLength 解析http消息体长度
func ParseContentLength(request *Request) error {
	contentLength := textproto.TrimString(request.Header.Get("Content-Length"))
	if contentLength == "" {
		request.ContentLength = 0
		return nil
	}

	if n, err := strconv.ParseUint(contentLength, 10, 63); err != nil {
		request.ContentLength = -1
		return errors.New("bad Content-Length: " + contentLength)
	} else {
		request.ContentLength = int64(n)
	}
	return nil
}

// ParseRequestBody 解析http请求体
func ParseRequestBody(request *Request, buf *bufio.Reader) error {
	if request.ContentLength <= 0 {
		request.Body = NewBody([]byte{})
		return nil
	}

	body := make([]byte, request.ContentLength)
	readLength, err := buf.Read(body)
	if err != nil {
		return err
	}
	if int64(readLength) != request.ContentLength {
		return errors.New("failed to read request Body")
	}
	request.Body = NewBody(body)
	return nil
}

func bindResponse(request *Request) *Response {
	response := NewResponse()
	if request.Proto != "" {
		response.Proto = request.Proto
		response.ProtoMinor = request.ProtoMinor
		response.ProtoMajor = request.ProtoMajor
	}
	if request.Header.Get("Connection") != "" {
		response.AddHeader("Connection", request.Header.Get("Connection"))
	}

	return response
}
