package simpleRpc

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/simple-set/simple.io/src/protocol/codec"
)

type Encoder struct {
}

func (e *Encoder) Encoder(message *Message, writer *bufio.Writer) error {
	return e.makeFrame(message, codec.NewWriteByteBuf(writer))
}

func (e *Encoder) makeFrame(message *Message, buf *codec.ByteBuf) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = errors.New(fmt.Sprintln("make data frame error, ", e))
		}
	}()

	buf.WriteBytes(message.magic[:])
	buf.WriteInt(message.Version)
	buf.WriteInt(message.Length)
	buf.WriteBytea(byte(message.Sequence))
	buf.WriteBytea(byte(message.Kind))
	buf.WriteBytes(message.method[:])
	buf.WriteInt(message.Sequence)
	buf.WriteBytes(message.Body())
	buf.Flush()
	return err
}

func NewEncoder() *Encoder {
	return &Encoder{}
}
