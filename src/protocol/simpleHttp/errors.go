package simpleHttp

import "net/http"

// http错误
type statusError struct {
	code int
	text string
}

func (e *statusError) Error() string {
	return http.StatusText(e.code) + ": " + e.text
}

// 错误的请求 400
func badRequestError(e string) error { return &statusError{http.StatusBadRequest, e} }

// 服务器错误 500
func serverError(e string) error { return &statusError{http.StatusInternalServerError, e} }
