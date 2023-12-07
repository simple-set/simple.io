package simpleHttp

import (
	"bufio"
	"log"
	"net/http"
	"net/url"
)

// RequestBuild 请求对象构造器
type RequestBuild struct {
	request *Request
}

func (r *RequestBuild) bufReader(bufReader *bufio.Reader) *RequestBuild {
	r.request.bufReader = bufReader
	return r
}

func (r *RequestBuild) Header(name, value string) *RequestBuild {
	r.request.AddHeader(name, value)
	return r
}

func (r *RequestBuild) Cookie(name, value string) *RequestBuild {
	r.request.AddCookie(name, value)
	return r
}

func (r *RequestBuild) Agent(agent string) *RequestBuild {
	r.request.AddHeader("User-Agent", agent)
	return r
}

func (r *RequestBuild) Uri(uri string) *RequestBuild {
	requestURI, err := url.ParseRequestURI(uri)
	if err != nil {
		log.Panic(err)
	}
	r.request.URL = requestURI
	r.request.Host = requestURI.Host
	return r
}

func (r *RequestBuild) Get() *RequestBuild {
	r.request.Method = "GET"
	return r
}

func (r *RequestBuild) Post() *RequestBuild {
	r.request.Method = "GET"
	return r
}

func (r *RequestBuild) method(method string) *RequestBuild {
	r.request.Method = method
	return r
}

func (r *RequestBuild) Build() *Request {
	return r.request
}

func NewRequestBuild() *RequestBuild {
	build := &RequestBuild{request: new(Request)}
	build.request.Proto = "HTTP/1.1"
	build.request.ProtoMajor, build.request.ProtoMinor, _ = http.ParseHTTPVersion(build.request.Proto)
	return build
}

// ByReaderRequest 根据io.Reader构建请求对象
func ByReaderRequest(reader *bufio.Reader) *Request {
	return NewRequestBuild().bufReader(reader).Build()
}
