package socket

import (
	"fmt"
	"github.com/golang/mock/gomock"
	mockNet "github.com/simple-set/simple.io/src/mocks/net"
	"net"
	//mockNet "simple-io/src/mocks/net"
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
	_, _ = socket.Read()
}

func TestSocket_Write(t *testing.T) {
	data := []byte("helloWorld")
	mockConn := mockNet.NewMockConn(gomock.NewController(t))
	mockConn.EXPECT().Write(data).Return(len(data), nil).AnyTimes()

	socket := NewSocket(mockConn, nil, nil)
	nn, err := socket.Write(data)
	fmt.Println(nn, err)

	if nn != len(data) {
		t.Fatal("socket.write()")
	}
}

func TestSocket_Flush(t *testing.T) {
	data := []byte("HelloWorld")
	mockConn := mockNet.NewMockConn(gomock.NewController(t))
	mockConn.EXPECT().Write(data).Return(len(data), nil).Times(1)

	// run test
	socket := NewSocket(mockConn, nil, nil)
	_, _ = socket.Write(data)
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
