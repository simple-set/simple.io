package codec

import (
	"bufio"
)

type ByteBuf struct {
	readBuffer  *bufio.Reader
	writeBuffer *bufio.Writer
}

func (b *ByteBuf) ReadBuffer() *bufio.Reader {
	return b.readBuffer
}

func (b *ByteBuf) ReadInt() int {
	readInt, err := ReadInt(b.readBuffer)
	if err != nil {
		panic(err)
	}
	return readInt
}

//goland:noinspection GoStandardMethods
func (b *ByteBuf) ReadByte() byte {
	readByte, err := b.readBuffer.ReadByte()
	if err != nil {
		panic(err)
	}
	return readByte
}

func (b *ByteBuf) ReadBytes(size int) []byte {
	bytes, err := ReadBytes(b.readBuffer, size)
	if err != nil {
		panic(err)
	}
	return bytes
}

func (b *ByteBuf) ReadLine() string {
	var readData []byte
	for {
		line, more, err := b.readBuffer.ReadLine()
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
	return b.writeBuffer
}

func (b *ByteBuf) WriteInt(data int) {
	if n, err := b.writeBuffer.Write(IntToByteBigEndian(data)); n != 4 || err != nil {
		panic(err)
	}
}

//goland:noinspection GoStandardMethods
func (b *ByteBuf) WriteByte(data byte) {
	if err := b.writeBuffer.WriteByte(data); err != nil {
		panic(err)
	}
}

func (b *ByteBuf) WriteBytes(data []byte) {
	if n, err := b.writeBuffer.Write(data); err != nil || n != len(data) {
		panic(err)
	}
}

func (b *ByteBuf) WriteString(data string) {
	if _, err := b.writeBuffer.WriteString(data); err != nil {
		panic(err)
	}
}

func (b *ByteBuf) Flush() {
	if b.writeBuffer != nil {
		if err := b.writeBuffer.Flush(); err != nil {
			panic(err)
		}
	}
}

func NewReadByteBuf(buffer *bufio.Reader) *ByteBuf {
	return &ByteBuf{readBuffer: buffer}
}

func NewWriteByteBuf(writeBuffer *bufio.Writer) *ByteBuf {
	return &ByteBuf{writeBuffer: writeBuffer}
}
