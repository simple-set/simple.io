package simpleHttp

import (
	"bufio"
	"github.com/simple-set/simple.io/src/version"
	"net/http"
	"strconv"
)

var (
	crlf       = []byte("\r\n")
	colonSpace = []byte(": ")
)

type Response struct {
	Proto         string
	ProtoMajor    int
	ProtoMinor    int
	statusCode    int
	statusText    string
	Header        http.Header
	Cookie        []*Cookie
	Close         bool
	body          *Body
	responseBytes []byte
	contentLength int64
	request       *Request
}

func (r *Response) ResponseBytes() []byte {
	return r.responseBytes
}

func (r *Response) StatusCode() int {
	return r.statusCode
}

func (r *Response) SetStatusCode(statusCode int) {
	r.statusCode = statusCode
	if statusText := http.StatusText(statusCode); statusText != "" {
		r.statusText = statusText
	}
}

func (r *Response) Request() *Request {
	return r.request
}

func (r *Response) AddHeader(name, value string) {
	if r.Header == nil {
		r.Header = make(http.Header, 4)
	}
	if r.Header.Get("name") != "" {
		r.Header.Add(name, value)
	} else {
		r.Header.Set(name, value)
	}
}

func (r *Response) AddCookie(name, value string) {
	r.AddCookieEntity(NewCookie(name, value))
}

func (r *Response) AddCookieEntity(cookie *Cookie) {
	if r.Cookie == nil {
		r.Cookie = make([]*Cookie, 0)
	}
	r.Cookie = append(r.Cookie, cookie)
}

func (r *Response) Write(p []byte) (int, error) {
	if r.body == nil {
		r.body = NewBody(make([]byte, 0))
	}
	size, err := r.body.Write(p)
	r.contentLength = int64(r.body.Len())
	return size, err
}

func NewResponse() *Response {
	return NewResponseBuild().
		Header("Server", version.Name+"/"+version.Version).
		Status(http.StatusOK).
		Build()
}

// ResponseDecode Response编码器
type ResponseDecode struct {
	response   *Response
	buf        *bufio.Writer
	decodeData []byte
	decodeErr  error
}

func (r *ResponseDecode) Decode() ([]byte, error) {
	defer func() {
		// TODO 编码错误，服务器500
	}()
	r.responseLine()
	r.responseHeader()
	r.responseBody()
	r.buf.Flush()
	r.response.responseBytes = r.decodeData
	return r.response.ResponseBytes(), nil
}

func (r *ResponseDecode) responseLine() {
	r.buf.WriteString(r.response.Proto + " ")
	r.buf.WriteString(strconv.Itoa(r.response.statusCode) + " ")
	r.buf.WriteString(r.response.statusText)
	r.buf.Write(crlf)
}

func (r *ResponseDecode) responseHeader() {
	if r.response.Header == nil {
		return
	}

	for name, values := range r.response.Header {
		if values == nil || len(values) == 0 {
			continue
		}
		for i := 0; i < len(values); i++ {
			r.writeHeader(name, values[i])
		}
	}
	if r.response.contentLength > 0 {
		r.writeHeader("Content-Length", strconv.FormatInt(r.response.contentLength, 10))
	}
	r.encodeCookie()
	r.buf.Write(crlf)
}

func (r *ResponseDecode) writeHeader(name, value string) {
	r.buf.WriteString(name)
	r.buf.Write(colonSpace)
	r.buf.WriteString(value)
	r.buf.Write(crlf)
}

func (r *ResponseDecode) encodeCookie() {
	if r.response.Cookie == nil || len(r.response.Cookie) == 0 {
		return
	}
	for _, cookie := range r.response.Cookie {
		if v := cookie.String(); v != "" {
			r.writeHeader("Set-Cookie", v)
		}
	}
}
func (r *ResponseDecode) responseCookie() {

}

func (r *ResponseDecode) responseBody() {
	if r.response.contentLength <= 0 {
		return
	}
	bytes, err := r.response.body.ReadBytes()
	if err != nil {
		r.decodeErr = err
		return
	}
	r.buf.Write(bytes)
}

func (r *ResponseDecode) Write(data []byte) (int, error) {
	if r.decodeData == nil {
		r.decodeData = make([]byte, 0)
	}
	r.decodeData = append(r.decodeData, data...)
	return len(data), nil
}

func NewResponseDecode(response *Response) *ResponseDecode {
	responseDecode := &ResponseDecode{response: response}
	responseDecode.buf = bufio.NewWriter(responseDecode)
	return responseDecode
}
