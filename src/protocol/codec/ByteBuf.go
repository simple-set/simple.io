package codec

import (
	"bufio"
)

type ByteBuf struct {
	*bufio.Reader
	*bufio.Writer
}

func (b *ByteBuf) ReadBuffer() *bufio.Reader {
	return b.Reader
}

func (b *ByteBuf) ReadInt() int {
	readInt, err := ReadInt(b.Reader)
	if err != nil {
		panic(err)
	}
	return readInt
}

func (b *ByteBuf) ReadBytea() byte {
	readByte, err := b.Reader.ReadByte()
	if err != nil {
		panic(err)
	}
	return readByte
}

func (b *ByteBuf) ReadBytes(size int) []byte {
	bytes, err := ReadBytes(b.Reader, size)
	if err != nil {
		panic(err)
	}
	return bytes
}

func (b *ByteBuf) ReadLine() string {
	var readData []byte
	for {
		line, more, err := b.Reader.ReadLine()
		if err != nil {
			panic(err)
		}
		if line != nil {
			readData = append(readData, line...)
		}
		if more {
			continue
		}
		break
	}
	return string(readData)
}

func (b *ByteBuf) WriteBuffer() *bufio.Writer {
	return b.Writer
}

func (b *ByteBuf) WriteInt(data int) {
	if n, err := b.Write(IntToByteBigEndian(data)); n != 4 || err != nil {
		panic(err)
	}
}

func (b *ByteBuf) WriteBytea(data byte) {
	if err := b.Writer.WriteByte(data); err != nil {
		panic(err)
	}
}

func (b *ByteBuf) WriteBytes(data []byte) {
	if n, err := b.Writer.Write(data); err != nil || n != len(data) {
		panic(err)
	}
}

func (b *ByteBuf) WriteString(data string) {
	if _, err := b.Writer.WriteString(data); err != nil {
		panic(err)
	}
}

func (b *ByteBuf) Flush() {
	if b.Writer != nil {
		if err := b.Writer.Flush(); err != nil {
			panic(err)
		}
	}
}

func NewReadByteBuf(reader *bufio.Reader) *ByteBuf {
	return &ByteBuf{Reader: reader}
}

func NewWriteByteBuf(writeBuffer *bufio.Writer) *ByteBuf {
	return &ByteBuf{Writer: writeBuffer}
}

func NewReadWriteByteBuf(reader *bufio.Reader, writer *bufio.Writer) *ByteBuf {
	return &ByteBuf{Reader: reader, Writer: writer}
}
