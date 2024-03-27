package codec

import (
	"bufio"
	"encoding/binary"
	"errors"
)

// LengthBased 长度解码器, 重试从缓冲区解析数据帧程度
type LengthBased struct {
	// 长度偏移，长度开始位置
	lengthFieldOffset int
	// 长度占用字节
	lengthFieldLength int
	// 长度调整，从长度字段开始计算，还有多少字符是内容，也就是长度和内容之间还有占位
	lengthAdjustment int
	// 剥离字节数，从头部开始剥离，也就是的丢弃头部几个字节，比如协议魔术字符
	initialBytesToStrip int
}

// Decode 长度解码器
func (l *LengthBased) Decode(reader *bufio.Reader) (*bufio.Reader, error) {
	frameLength, err := l.getFrameLength(reader)
	if err != nil {
		return nil, err
	}
	if frameLength <= 0 {
		return nil, errors.New("the data frame must be greater than zero")
	}
	return reader, nil
}

// 获取数据帧长度
func (l *LengthBased) getFrameLength(reader *bufio.Reader) (int, error) {
	data, err := reader.Peek(l.lengthFieldOffset + l.lengthFieldLength)
	if err != nil {
		return 0, err
	}
	frameLength := int(binary.BigEndian.Uint32(data[l.lengthFieldOffset : l.lengthFieldOffset+l.lengthFieldLength]))
	if frameLength <= 0 {
		return 0, errors.New("frame length parsing error")
	}
	return frameLength, nil
}

func NewLengthBased(lengthFieldOffset int, lengthFieldLength int, lengthAdjustment int, initialBytesToStrip int) *LengthBased {
	return &LengthBased{
		lengthFieldOffset:   lengthFieldOffset,
		lengthFieldLength:   lengthFieldLength,
		lengthAdjustment:    lengthAdjustment,
		initialBytesToStrip: initialBytesToStrip}
}
