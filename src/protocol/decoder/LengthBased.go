package decoder

import (
	"encoding/binary"
	"errors"
	"github.com/simple-set/simple.io/src/event"
	"github.com/sirupsen/logrus"
)

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

func NewLengthBased(lengthFieldOffset int, lengthFieldLength int, lengthAdjustment int, initialBytesToStrip int) *LengthBased {
	return &LengthBased{
		lengthFieldOffset:   lengthFieldOffset,
		lengthFieldLength:   lengthFieldLength,
		lengthAdjustment:    lengthAdjustment,
		initialBytesToStrip: initialBytesToStrip}
}

func (l *LengthBased) Input(context *event.HandleContext, data interface{}) (interface{}, bool) {
	if value, ok := data.([]byte); ok {
		frameData, err := l.Decode(&value)
		if err != nil {
			logrus.Println(err)
			return nil, false
		}
		return frameData, true
	}
	return data, true
}

func (l *LengthBased) Decode(data *[]byte) ([]byte, error) {
	if l.lengthFieldOffset+l.lengthFieldLength > len(*data) {
		return nil, nil
	}

	dataLength := (*data)[l.lengthFieldOffset:l.lengthFieldLength]
	frameLength := int(binary.BigEndian.Uint32(dataLength))
	if frameLength <= 0 {
		return nil, errors.New("frame length parsing error")
	}

	if frameLength >= len(*data) {
		frameData := (*data)[l.initialBytesToStrip : frameLength-l.initialBytesToStrip]
		*data = (*data)[frameLength:]
		return frameData, nil
	}
	return *data, nil
}
