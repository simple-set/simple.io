package decoder

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"testing"
)

func TestFixedLengthDecoder_Decode(t *testing.T) {
	lengthDecoder := NewFixedLengthDecoder(2)

	data := []byte("helloWorld\n\r")
	for {
		decode, err := lengthDecoder.Decode(&data)
		if err != nil {
			logrus.Fatalf(err.Error())
		}
		if decode == nil {
			break
		}
		fmt.Printf("%v, %s\n", decode, decode)
	}
}
