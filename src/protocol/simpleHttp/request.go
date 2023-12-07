package simpleHttp

import (
	"bufio"
	"net/http"
)

type Request struct {
	http.Request
	Body *Body
	//cookies   *[]Cookie
	Response  *Response
	bufWriter *bufio.Writer
	bufReader *bufio.Reader
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
