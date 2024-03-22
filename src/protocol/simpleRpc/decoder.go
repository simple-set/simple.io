package simpleRpc

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/simple-set/simple.io/src/protocol/codec"
)

// Decoder 解码器
type Decoder struct {
	LengthBased *codec.LengthBased
}

// Decoder 解码
func (d *Decoder) Decoder(reader *bufio.Reader) (*Message, error) {
	if verify, err := d.verifyMagic(reader); err != nil || verify != true {
		return nil, errors.New("Magic character verification failed, " + err.Error())
	}

	if _, err := d.LengthBased.Decode(reader); err != nil {
		return nil, errors.New("Data frame length error, " + err.Error())
	}

	return d.parserFrame(codec.NewReadByteBuf(reader))
}

// 解析数据帧
func (d *Decoder) parserFrame(byteBuf *codec.ByteBuf) (message *Message, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = errors.New(fmt.Sprintln("Parsing data frame error, ", e))
		}
	}()

	message = NewEmpty()
	message.magic = *(*[2]byte)(byteBuf.ReadBytes(2))
	message.Version = byteBuf.ReadInt()
	message.Length = byteBuf.ReadInt()
	message.Serializer = SerializerType(byteBuf.ReadByte())
	message.method = Method(byteBuf.ReadByte())
	message.Sequence = byteBuf.ReadInt()
	message.SetBody(byteBuf.ReadBytes(message.Length))
	return message, nil
}

// 校验魔术字符
func (d *Decoder) verifyMagic(reader *bufio.Reader) (bool, error) {
	magic, err := reader.Peek(2)
	if err != nil {
		return false, err
	}

	if magic[0] != Magic()[0] || magic[1] != Magic()[1] {
		return false, errors.New(fmt.Sprintln(magic[0], magic[1]))
	}
	return true, nil
}

func NewDecoder() *Decoder {
	lengthBased := codec.NewLengthBased(6, 4, 6, 0)
	return &Decoder{LengthBased: lengthBased}
}
