package simpleHttp

import (
	"bufio"
	"errors"
	"net/http"
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

// 从缓冲区中读取一行, 注意没有限制长度
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
