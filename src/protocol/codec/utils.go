package codec

import (
	"encoding/binary"
	"io"
)

func IntToByteBigEndian(i int) []byte {
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b[:], uint32(i))
	return b
}

func ReadBytes(reader io.Reader, size int) ([]byte, error) {
	bytes := make([]byte, size)
	if n, err := reader.Read(bytes); err != nil || n != size {
		return nil, err
	}
	return bytes, nil
}

func ReadInt(reader io.Reader) (int, error) {
	bytes, err := ReadBytes(reader, 4)
	if err != nil {
		return -1, err
	}
	return int(binary.BigEndian.Uint32(bytes)), nil
}
