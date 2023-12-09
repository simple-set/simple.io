package simpleHttp

import (
	"strconv"
)

type ResponseEncode struct {
	response *Response
}

func (r *ResponseEncode) Codec() error {
	if err := r.line(); err != nil {
		return err
	}
	if err := r.header(); err != nil {
		return err
	}
	if err := r.body(); err != nil {
		return err
	}
	err := r.response.bufWriter.Flush()
	if err != nil {
		return err
	}
	return nil
}

func (r *ResponseEncode) line() error {
	if _, err := r.response.bufWriter.WriteString(r.response.Proto + " "); err != nil {
		return nil
	}
	if _, err := r.response.bufWriter.WriteString(strconv.Itoa(r.response.statusCode) + " "); err != nil {
		return err
	}
	if _, err := r.response.bufWriter.WriteString(r.response.statusText); err != nil {
		return err
	}
	if _, err := r.response.bufWriter.Write(crlf); err != nil {
		return err
	}
	return nil
}

func (r *ResponseEncode) header() error {
	if r.response.Server != "" {
		err := writeHeader(r.response.bufWriter, "Server", r.response.Server)
		if err != nil {
			return err
		}
	}
	if r.response.contentLength > 0 {
		err := writeHeader(r.response.bufWriter, "Content-Length", strconv.FormatInt(r.response.contentLength, 10))
		if err != nil {
			return err
		}
	}

	if r.response.Header == nil {
		return nil
	}
	for name, values := range r.response.Header {
		if values == nil || len(values) == 0 {
			continue
		}
		for i := 0; i < len(values); i++ {
			if name == "Cookie" {
				if err := r.cookie(); err != nil {
					return err
				}
				continue
			}
			err := writeHeader(r.response.bufWriter, name, values[i])
			if err != nil {
				return err
			}
		}
	}
	if _, err := r.response.bufWriter.Write(crlf); err != nil {
		return err
	}
	return nil
}

// 编码cookie
func (r *ResponseEncode) cookie() error {
	cookies := r.response.Cookies()
	for i := 0; i < len(cookies); i++ {
		if err := writeHeader(r.response.bufWriter, "Set-Cookie", cookies[i].Name+"="+cookies[i].Value); err != nil {
			return err
		}
	}
	return nil
}

func (r *ResponseEncode) body() error {
	if r.response.contentLength <= 0 {
		return nil
	}
	bytes, err := r.response.body.ReadBytes()
	if err != nil {
		return err
	}
	if _, err := r.response.bufWriter.Write(bytes); err != nil {
		return err
	}
	return nil
}

func NewResponseEncoded(response *Response) *ResponseEncode {
	return &ResponseEncode{response: response}
}
