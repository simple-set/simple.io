package decoder

// Decoder 解码器
type Decoder interface {
	Decode(data *[]byte) ([]byte, error)
}
