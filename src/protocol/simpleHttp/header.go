package simpleHttp

import "net/http"

type Header struct {
	http.Header
}

//func NewHeader() *Header {
//	return &Header{Header: make(http.Header, 4)}
//}
//
//func (h *Header) Has(key string) bool {
//	value := h.Get(key)
//	if len(value) > 0 {
//		return true
//	}
//	return false
//}
