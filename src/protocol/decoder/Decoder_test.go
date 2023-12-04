package decoder

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"testing"
)

func decoder(data *[]byte) {
	*data = (*data)[len(*data)-1:]
}

// 字节转换成整形
func BytesToInt(b []byte) int {
	bytesBuffer := bytes.NewBuffer(b)

	var x int32
	binary.Read(bytesBuffer, binary.BigEndian, &x)

	return int(x)
}

// 整形转换成字节
func IntToBytes(n int) []byte {
	x := int32(n)

	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.BigEndian, x)
	return bytesBuffer.Bytes()
}

func TestName(t *testing.T) {
	data := []byte{0, 0, 0, 65}
	//toInt := BytesToInt(data)
	fmt.Printf("%v\n", data)

	toBytes := IntToBytes(100)
	fmt.Printf("%v\n", toBytes)
}
