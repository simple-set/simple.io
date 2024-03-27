package simpleHttp

import (
	"fmt"
	"github.com/simple-set/simple.io/src/protocol/codec"
	"github.com/simple-set/simple.io/src/version"
	"net/http"
)

type Response struct {
	Proto         string
	ProtoMajor    int
	ProtoMinor    int
	statusCode    int
	statusText    string
	Header        http.Header
	Close         bool
	body          *Body
	contentLength int64
	Server        string
	request       *Request
	bufWriter     *codec.ByteBuf
	bufReader     *codec.ByteBuf
}

func (r *Response) SetBody(body *Body) {
	r.body = body
}

func (r *Response) Body() *Body {
	return r.body
}

func (r *Response) StatusCode() int {
	return r.statusCode
}

func (r *Response) SetStatusCode(statusCode int) {
	r.statusCode = statusCode
	if statusText := http.StatusText(statusCode); statusText != "" {
		r.statusText = statusText
	}
}

func (r *Response) Request() *Request {
	return r.request
}

func (r *Response) AddHeader(name, value string) {
	if r.Header == nil {
		r.Header = make(http.Header, 4)
	}
	if r.Header.Get("name") != "" {
		r.Header.Add(name, value)
	} else {
		r.Header.Set(name, value)
	}
}

func (r *Response) AddCookie(name, value string) {
	r.AddCookieEntity(&http.Cookie{Name: name, Value: value})
}

func (r *Response) AddCookieEntity(cookie *http.Cookie) {
	if r.Header == nil {
		r.Header = make(http.Header, 4)
	}
	s := fmt.Sprintf("%s=%s", sanitizeCookieName(cookie.Name), sanitizeCookieValue(cookie.Value))
	if c := r.Header.Get("Cookie"); c != "" {
		r.Header.Set("Cookie", c+"; "+s)
	} else {
		r.Header.Set("Cookie", s)
	}
}

func (r *Response) Cookies() []*http.Cookie {
	return readCookies(r.Header, "")
}

func (r *Response) cookie(name string) (*http.Cookie, error) {
	for _, c := range readCookies(r.Header, name) {
		return c, nil
	}
	return nil, http.ErrNoCookie
}

func NewResponse() *Response {
	return NewResponseBuild().
		Server(version.Name + "/" + version.Version).
		Status(http.StatusOK).
		Build()
}

func NewResponseReader(bufReader *codec.ByteBuf) *Response {
	return &Response{bufReader: bufReader}
}

// NewReplyResponse 创建请求响应体, 根据request创建Response, 用于编写HttpServer服务器时响应客户端请求
func NewReplyResponse(request *Request) *Response {
	response := NewResponse()
	response.request = request
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
