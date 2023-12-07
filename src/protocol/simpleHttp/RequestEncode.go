package simpleHttp

import (
	"errors"
	"strings"
)

// RequestEncode http编码器, 把request对象编码成字节流, 根据 RFC2616 规范编码
// 适用于编写http客户端场景，把构建好的请求对象编码成字节流发送给服务器。
type RequestEncode struct {
	request *Request
}

func (r *RequestEncode) Encoder() error {
	return nil
}

// 编码请求行
func (r *RequestEncode) line() error {
	if _, err := r.request.bufWriter.WriteString(r.request.Method + " "); err != nil {
		return err
	}
	if _, err := r.request.bufWriter.WriteString(r.request.RequestURI + " "); err != nil {
		return err
	}
	if _, err := r.request.bufWriter.WriteString(r.request.Proto); err != nil {
		return err
	}
	if _, err := r.request.bufWriter.Write(crlf); err != nil {
		return err
	}
	return nil
}

// 编码请求头
func (r *RequestEncode) header() error {
	if r.request.Header == nil || len(r.request.Header) == 0 {
		return nil
	}
	for name, values := range r.request.Header {
		if values == nil || len(values) == 0 {
			continue
		}
		if name == "Cookie" {
			if err := r.cookie(); err != nil {
				return err
			}
			continue
		}
		for _, value := range values {
			if err := writeHeader(r.request.bufWriter, name, value); err != nil {
				return err
			}
		}
	}

	if _, err := r.request.bufWriter.Write(crlf); err != nil {
		return err
	}
	return nil
}

// 编码cookie
func (r *RequestEncode) cookie() error {
	cookies := r.request.Cookies()
	if cookies == nil || len(cookies) == 0 {
		return nil
	}
	cookieValue := strings.Builder{}
	for i := 0; i < len(cookies); i++ {
		cookieValue.WriteString(cookies[i].Name)
		cookieValue.WriteString("=")
		cookieValue.WriteString(cookies[i].Value)
		if i < len(cookies) {
			cookieValue.WriteString("; ")
		}
	}
	if err := writeHeader(r.request.bufWriter, "Cookie", cookieValue.String()); err != nil {
		return err
	}
	return nil
}

// 编码请求题
func (r *RequestEncode) body() error {
	if (r.request.Method == "POST" || r.request.Method == "PUT") && r.request.ContentLength > 0 {
		if n, err := r.request.bufWriter.Write(r.request.Body.Bytes()); err != nil {
			return err
		} else if int64(n) != r.request.ContentLength {
			return errors.New("xxx")
		}
	}
	return nil
}

func NewRequestEncode(request *Request) *RequestEncode {
	return &RequestEncode{request: request}
}
