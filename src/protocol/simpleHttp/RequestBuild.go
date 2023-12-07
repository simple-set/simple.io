package simpleHttp

import (
	"log"
	"net/url"
)

// RequestBuild 请求对象构造器
type RequestBuild struct {
	request *Request
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

func (r *RequestBuild) Method(method string) *RequestBuild {
	r.request.Method = method
	return r
}

func (r *RequestBuild) Proto(proto string) *RequestBuild {
	r.request.Proto = proto
	return r
}

func (r *RequestBuild) Build() *Request {
	return r.request
}

func NewRequestBuild() *RequestBuild {
	return &RequestBuild{request: new(Request)}
}
