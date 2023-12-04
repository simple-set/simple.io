package simpleHttp

import "net/http"

type Cookie struct {
	http.Cookie
}

func NewCookie(name, value string) *Cookie {
	return nil
}
