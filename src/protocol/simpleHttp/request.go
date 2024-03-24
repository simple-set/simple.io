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
	readBuff   *codec.ByteBuf
	writerBuff *codec.ByteBuf
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

func DefaultRequest() *Request {
	request := NewRequestBuild().Proto("HTTP/1.1").Agent(version.Name + "/" + version.Version).Build()
	request.ProtoMajor, request.ProtoMinor, _ = http.ParseHTTPVersion(request.Proto)
	return request
}

//func NewRequestReader(bufReader *bufio.Reader) *Request {
//	return &Request{bufReader: bufReader}
//}

func NewRequest(readBuff *codec.ByteBuf) *Request {
	return &Request{readBuff: readBuff}
}
