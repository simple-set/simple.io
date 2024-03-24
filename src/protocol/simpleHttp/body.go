package simpleHttp

import (
	"github.com/simple-set/simple.io/src/protocol/codec"
)

type Body struct {
	r, w int64
	size int64
	buff *codec.ByteBuf
}

func (b *Body) Size() int64 {
	return b.size
}

func (b *Body) Read(p []byte) (n int, err error) {
	buffer := b.buff.ReadBuffer()
	return buffer.Read(p)
}

func (b *Body) ReadString() (string, error) {
	panic("implement me")
}

//func (b *Body) ReadBytes() ([]byte, error) {
//	if b.Len() == 0 {
//		return []byte{}, nil
//	}
//	body := make([]byte, b.Len())
//	if _, err := b.Read(body); err != nil {
//		return nil, err
//	}
//	return body, nil
//}
//
//func (b *Body) String() (string, error) {
//	body, err := b.ReadBytes()
//	if err != nil {
//		return "", err
//	}
//	return string(body), nil
//}

func NewBody(size int64, buff *codec.ByteBuf) *Body {
	return &Body{size: size, buff: buff}
}
