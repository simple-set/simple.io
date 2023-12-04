package decoder

import "github.com/simple-set/simple.io/src/event"

type FixedLengthDecoder struct {
	FrameLength   int
	decoderBuffer []byte
}

func NewFixedLengthDecoder(frameLength int) *FixedLengthDecoder {
	return &FixedLengthDecoder{FrameLength: frameLength}
}

func (d *FixedLengthDecoder) Input(context *event.HandleContext, data interface{}) (interface{}, bool) {
	panic("implement me")
}

func (d *FixedLengthDecoder) Decode(data *[]byte) ([]byte, error) {
	if len(*data) >= d.FrameLength {
		frameData := (*data)[:d.FrameLength]
		*data = (*data)[d.FrameLength:]
		return frameData, nil
	}
	return nil, nil
}
