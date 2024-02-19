package simpleHttp

import (
	"bufio"
	"errors"
	"golang.org/x/net/http/httpguts"
	"log"
	"net/http"
	"net/textproto"
	"strings"
)

var HttpMethods = []string{
	"GET",
	"HEAD",
	"POST",
	"PUT",
	"PATCH", // RFC 5789
	"DELETE",
	"CONNECT",
	"OPTIONS",
	"TRACE",
}

var (
	crlf       = []byte("\r\n")
	colonSpace = []byte(": ")
)

// ValidMethod 验证http方法
func ValidMethod(method string) bool {
	for _, name := range HttpMethods {
		if name == method {
			return true
		}
	}
	return false
}

func ValidPath(path string) bool {
	return true
}

func ValidHeader(header string) bool {
	return true
}

func ValidBody(body string) bool {
	return true
}

func ValidQuery(query string) bool {
	return true
}

// PragmaCacheControl RFC 7234, section 5.4: Pragma: no-cache && Cache-Control: no-cache
func PragmaCacheControl(header http.Header) {
	if value, ok := header["Pragma"]; ok && len(value) > 0 && value[0] == "no-cache" {
		if _, ok := header["Cache-Control"]; !ok {
			header["Cache-Control"] = []string{"no-cache"}
		}
	}
}

// 从缓冲区中读取一行, 注意:没有限制长度
func readLine(buf *bufio.Reader) ([]byte, error) {
	if buf == nil {
		return nil, errors.New("the buffer has no data to read")
	}

	var readData []byte
	for {
		line, more, err := buf.ReadLine()
		if err != nil {
			return nil, err
		}
		if line != nil {
			readData = append(readData, line...)
		}
		if more {
			continue
		}
		break
	}
	return readData, nil
}

// 向缓冲区写入一行
func writeHeader(writer *bufio.Writer, name, value string) error {
	if _, err := writer.WriteString(name); err != nil {
		return err
	}
	if _, err := writer.Write(colonSpace); err != nil {
		return err
	}
	if _, err := writer.WriteString(value); err != nil {
		return err
	}
	if _, err := writer.Write(crlf); err != nil {
		return err
	}
	return nil
}

func readCookies(h http.Header, filter string) []*http.Cookie {
	lines := h["Cookie"]
	if len(lines) == 0 {
		return []*http.Cookie{}
	}

	cookies := make([]*http.Cookie, 0, len(lines)+strings.Count(lines[0], ";"))
	for _, line := range lines {
		line = textproto.TrimString(line)

		var part string
		for len(line) > 0 { // continue since we have rest
			part, line, _ = strings.Cut(line, ";")
			part = textproto.TrimString(part)
			if part == "" {
				continue
			}
			name, val, _ := strings.Cut(part, "=")
			if !isCookieNameValid(name) {
				continue
			}
			if filter != "" && filter != name {
				continue
			}
			val, ok := parseCookieValue(val, true)
			if !ok {
				continue
			}
			cookies = append(cookies, &http.Cookie{Name: name, Value: val})
		}
	}
	return cookies
}

func isCookieNameValid(raw string) bool {
	if raw == "" {
		return false
	}
	return strings.IndexFunc(raw, isNotToken) < 0
}

func parseCookieValue(raw string, allowDoubleQuote bool) (string, bool) {
	// Strip the quotes, if present.
	if allowDoubleQuote && len(raw) > 1 && raw[0] == '"' && raw[len(raw)-1] == '"' {
		raw = raw[1 : len(raw)-1]
	}
	for i := 0; i < len(raw); i++ {
		if !validCookieValueByte(raw[i]) {
			return "", false
		}
	}
	return raw, true
}

func isNotToken(r rune) bool {
	return !httpguts.IsTokenRune(r)
}

func validCookieValueByte(b byte) bool {
	return 0x20 <= b && b < 0x7f && b != '"' && b != ';' && b != '\\'
}

func sanitizeCookieName(n string) string {
	return strings.NewReplacer("\n", "-", "\r", "-").Replace(n)
}

func sanitizeCookieValue(v string) string {
	v = sanitizeOrWarn("Cookie.Value", validCookieValueByte, v)
	if len(v) == 0 {
		return v
	}
	if strings.ContainsAny(v, " ,") {
		return `"` + v + `"`
	}
	return v
}

func sanitizeOrWarn(fieldName string, valid func(byte) bool, v string) string {
	ok := true
	for i := 0; i < len(v); i++ {
		if valid(v[i]) {
			continue
		}
		log.Printf("net/http: invalid byte %q in %s; dropping invalid bytes", v[i], fieldName)
		ok = false
		break
	}
	if ok {
		return v
	}
	buf := make([]byte, 0, len(v))
	for i := 0; i < len(v); i++ {
		if b := v[i]; valid(b) {
			buf = append(buf, b)
		}
	}
	return string(buf)
}
