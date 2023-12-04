package simpleHttp

import (
	"bytes"
)

type Body struct {
	*bytes.Buffer
}

func NewBody(body []byte) *Body {
	return &Body{bytes.NewBuffer(body)}
}

func (b *Body) ReadBytes() ([]byte, error) {
	if b.Len() == 0 {
		return []byte{}, nil
	}
	body := make([]byte, b.Len())
	if _, err := b.Read(body); err != nil {
		return nil, err
	}
	return body, nil
}

func (b *Body) String() (string, error) {
	body, err := b.ReadBytes()
	if err != nil {
		return "", err
	}
	return string(body), nil
}
