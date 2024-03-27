package simpleHttp

import (
	"bufio"
	"bytes"
	"errors"
	"github.com/simple-set/simple.io/src/protocol/codec"
	"io"
)

type Body struct {
	size int64
	buff *codec.ByteBuf
}

func (b *Body) Size() int64 {
	return b.size
}

func (b *Body) Write(p []byte) (int, error) {
	buffer := b.buff.WriteBuffer()
	if buffer == nil {
		return -1, errors.New("buffer cannot be written")
	}

	n, err := buffer.Write(p)
	_ = buffer.Flush()
	if err != nil {
		return 0, err
	}
	b.size += int64(n)
	return n, nil
}

func (b *Body) WriteString(data string) (int, error) {
	buffer := b.buff.WriteBuffer()
	if buffer == nil {
		return -1, errors.New("buffer cannot be written")
	}

	n, err := buffer.WriteString(data)
	_ = buffer.Flush()
	if err != nil {
		return 0, err
	}
	b.size += int64(n)
	return n, nil
}

func (b *Body) Read(p []byte) (n int, err error) {
	buffer := b.buff.ReadBuffer()
	if buffer == nil {
		return 0, errors.New("unable to read buffer")
	}
	if b.size <= 0 {
		return 0, io.EOF
	}
	if int64(cap(p)) > b.size {
		// 限制读取长度
		n, err = buffer.Read(p[:b.size])
	} else {
		n, err = buffer.Read(p)
	}

	if err != nil {
		return 0, err
	}
	b.size -= int64(n)
	return n, nil
}

func (b *Body) ReadBytes() ([]byte, error) {
	if b.size > 0 {
		buff := make([]byte, b.size)
		if _, err := b.Read(buff); err != nil {
			return nil, err
		}
		return buff, nil
	}
	return nil, errors.New("")
}

func (b *Body) ReadString() (string, error) {
	if readBytes, err := b.ReadBytes(); err != nil {
		return string(readBytes), nil
	} else {
		return "", err
	}
}

func NewReadBody(size int64, buff *codec.ByteBuf) *Body {
	if buff.ReadBuffer() != nil {
		return &Body{size: size, buff: buff}
	}
	return nil
}

func NewReaderWriteBody(buff *codec.ByteBuf) *Body {
	if buff.ReadBuffer() != nil && buff.WriteBuffer() != nil {
		return &Body{buff: buff}
	}
	return nil
}

func MakeReaderWriteBody() *Body {
	buffer := bytes.NewBuffer(make([]byte, 0))
	return NewReaderWriteBody(codec.NewReadWriteByteBuf(bufio.NewReader(buffer), bufio.NewWriter(buffer)))
}
