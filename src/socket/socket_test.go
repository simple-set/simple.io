package socket

import (
	"bytes"
	"github.com/golang/mock/gomock"
	"net"
	mockNet "simple-io/src/mocks/net"
	"testing"
)

func TestSocket_LocalAddr(t *testing.T) {
	mockConn := mockNet.NewMockConn(gomock.NewController(t))
	mockConn.EXPECT().LocalAddr().Return(&net.TCPAddr{}).Times(1)
	socket := NewSocket(mockConn, nil, nil)
	socket.LocalAddr()
}

func TestSocket_RemoteAddr(t *testing.T) {
	mockConn := mockNet.NewMockConn(gomock.NewController(t))
	mockConn.EXPECT().RemoteAddr().Return(&net.TCPAddr{}).Times(1)
	socket := NewSocket(mockConn, nil, nil)
	socket.RemoteAddr()
}

func TestSocket_Read(t *testing.T) {
	data := []byte("helloWorld")
	mockConn := mockNet.NewMockConn(gomock.NewController(t))
	mockConn.EXPECT().Read(gomock.Any()).Return(len(data), nil).Times(1)

	// run test
	socket := NewSocket(mockConn, nil, nil)
	socket.readBuf = data
	readData, err := socket.Read()

	// verify
	if err != nil {
		t.Fatal("socket.Read()")
	}
	if !bytes.Equal(readData, data) {
		t.Fatal("socket.Read()")
	}
}

func TestSocket_Write(t *testing.T) {
	data := []byte("helloWorld")
	socket := NewSocket(nil, nil, nil)
	socket.Write(data)

	if !bytes.Equal(data, socket.writeBuf) {
		t.Fatal("socket.Write()")
	}
}

func TestSocket_Flush(t *testing.T) {
	data := []byte("HelloWorld")
	mockConn := mockNet.NewMockConn(gomock.NewController(t))
	mockConn.EXPECT().Write(gomock.Any()).Return(len(data)/2, nil).Times(1)
	mockConn.EXPECT().Write(gomock.Any()).Return(len(data)/2, nil).Times(1)

	// run test
	socket := NewSocket(mockConn, nil, nil)
	socket.writeBuf = data
	err := socket.Flush()
	if err != nil {
		t.Fatal("socket.flush()")
	}
}

func TestSocket_Close(t *testing.T) {
	mockConn := mockNet.NewMockConn(gomock.NewController(t))
	mockConn.EXPECT().Close().Times(1)

	socket := NewSocket(mockConn, nil, nil)
	_ = socket.Close()
}
