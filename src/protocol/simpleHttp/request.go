package simpleHttp

import (
	"github.com/simple-set/simple.io/src/protocol/codec"
	"github.com/simple-set/simple.io/src/version"
	"net/http"
)

// Request http请求体
type Request struct {
	http.Request
	body       *Body
	Response   *Response
	buffReadr  *codec.ByteBuf
	buffWriter *codec.ByteBuf
}

func (r *Request) Body() *Body {
	return r.body
}

func (r *Request) SetBody(body *Body) {
	r.body = body
}

func (r *Request) AddHeader(name, value string) {
	if r.Header == nil {
		r.Header = make(http.Header, 4)
	}
	if r.Header.Get("name") != "" {
		r.Header.Add(name, value)
	} else {
		r.Header.Set(name, value)
	}
}

func (r *Request) AddCookie(name, value string) {
	r.AddCookieEntity(&http.Cookie{Name: name, Value: value})
}

func (r *Request) AddCookieEntity(cookie *http.Cookie) {
	if r.Header == nil {
		r.Header = make(http.Header, 4)
	}
	r.Request.AddCookie(cookie)
}

func NewRequestUrl(url string) *Request {
	request := NewRequestBuild().
		Proto("HTTP/1.1").
		Method("GET").
		Agent(version.Name + "/" + version.Version).
		Uri(url).
		Build()
	request.ProtoMajor, request.ProtoMinor, _ = http.ParseHTTPVersion(request.Proto)
	return request
}

func NewRequest(readBuff *codec.ByteBuf) *Request {
	return &Request{buffReadr: readBuff}
}
