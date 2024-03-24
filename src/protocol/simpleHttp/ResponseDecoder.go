package simpleHttp

import (
	"errors"
	"net/http"
	"net/textproto"
	"strconv"
	"strings"
)

type ResponseDecode struct {
	response *Response
}

func (r *ResponseDecode) Decoder() error {
	if err := r.line(); err != nil {
		return err
	}
	if err := r.header(); err != nil {
		return err
	}
	if err := r.contentLength(); err != nil {
		return err
	}
	if err := r.body(); err != nil {
		return err
	}
	return nil
}

func (r *ResponseDecode) line() error {
	line, err := readLine(r.response.bufReader)
	if err != nil {
		return err
	}

	proto, rest, ok1 := strings.Cut(string(line), " ")
	statusCode, statusText, ok2 := strings.Cut(rest, " ")
	if ok1 && ok2 {
		r.response.Proto = proto
		r.response.statusText = statusText
	} else {
		return errors.New("HTTP request line format error, " + string(line))
	}

	if code, err := strconv.Atoi(statusCode); err != nil {
		return errors.New("malformed HTTP statusCode " + statusCode)
	} else if http.StatusText(code) == "" {
		return errors.New("unknown HTTP statusCode " + statusCode)
	} else {
		r.response.statusCode = code
	}

	var ok bool
	if r.response.ProtoMajor, r.response.ProtoMinor, ok = http.ParseHTTPVersion(r.response.Proto); !ok {
		return errors.New("malformed HTTP version " + r.response.Proto)
	}
	return nil
}

func (r *ResponseDecode) header() error {
	header := make(http.Header, 4)

	for {
		line, err := readLine(r.response.bufReader)
		if err != nil {
			return err
		}
		if len(line) == 0 {
			break
		}
		if line[0] == ' ' || line[0] == '\t' {
			return errors.New("malformed MIME header initial line: " + string(line))
		}

		name, value, found := strings.Cut(string(line), string(colonSpace))
		name = strings.TrimSpace(name)
		value = strings.TrimSpace(value)
		if name == "" || value == "" || !found {
			return errors.New("Hearer format error： " + string(line))
		}

		if header.Get(name) == "" {
			header.Add(name, value)
		} else {
			header.Set(name, value)
		}
	}

	PragmaCacheControl(header)
	r.response.Header = header
	r.response.Server = header.Get("Server")
	return nil
}

func (r *ResponseDecode) contentLength() error {
	contentLength := textproto.TrimString(r.response.Header.Get("Content-Length"))
	if contentLength == "" {
		r.response.contentLength = 0
		return nil
	}

	if n, err := strconv.ParseUint(contentLength, 10, 63); err != nil {
		r.response.contentLength = -1
		return errors.New("bad Content-Length: " + contentLength)
	} else {
		r.response.contentLength = int64(n)
	}
	return nil
}

func (r *ResponseDecode) body() error {
	//if r.response.contentLength <= 0 {
	//	r.response.body = NewBody([]byte{})
	//	return nil
	//}
	//
	//body := make([]byte, r.response.contentLength)
	//readLength, err := r.response.bufReader.Read(body)
	//if err != nil {
	//	return err
	//}
	//if int64(readLength) != r.response.contentLength {
	//	return errors.New("failed to read request Body")
	//}
	//r.response.body = NewBody(body)
	return nil
}

func NewResponseDecoded(response *Response) *ResponseDecode {
	return &ResponseDecode{response: response}
}
