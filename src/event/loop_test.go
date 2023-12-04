package event

import (
	"errors"
	"github.com/golang/mock/gomock"
	mockNet "github.com/simple-set/simple.io/src/mocks/net"
	mockSocket "github.com/simple-set/simple.io/src/mocks/socket"
	"github.com/simple-set/simple.io/src/socket"
	"github.com/sirupsen/logrus"
	"net"
	"testing"
	"time"
)

func init() {
	logrus.SetLevel(logrus.TraceLevel)
}

func TestLoop_Connect(t *testing.T) {
	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()

	// mocks conn
	mockConn := mockNet.NewMockConn(mockCtl)
	mockConn.EXPECT().RemoteAddr().Return(&net.TCPAddr{}).AnyTimes()
	mockConn.EXPECT().Close().Times(1)
	mockConn.EXPECT().Read(gomock.Any()).Return(1, nil).Times(1)
	mockConn.EXPECT().Read(gomock.Any()).Return(0, errors.New("close")).Times(1)
	// mock socket
	mockSock := socket.NewSocket(mockConn, nil, nil)
	// mock client
	mockClient := mockSocket.NewMockClient(mockCtl)
	mockClient.EXPECT().Connect().Return(mockSock, nil).AnyTimes()

	eventLoop := NewEventLoop(nil)
	session := eventLoop.Connect(mockClient)
	session.Wait()

	if session.client != mockClient {
		t.Fatal("eventLoop.Connect failed")
	}
	if session.Sock() != mockSock {
		t.Fatal("eventLoop.Connect failed")
	}
	if session.state != Disconnect {
		t.Fatal("eventLoop.Connect failed")
	}
}

func TestLoop_Bind(t *testing.T) {
	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()

	// mocks conn
	mockConn := mockNet.NewMockConn(mockCtl)
	mockConn.EXPECT().RemoteAddr().Return(&net.TCPAddr{}).AnyTimes()
	mockConn.EXPECT().Close().Times(1)
	mockConn.EXPECT().Read(gomock.Any()).Return(1, nil).Times(1)
	mockConn.EXPECT().Read(gomock.Any()).Return(0, errors.New("close")).Times(1)
	// mock socket
	mockSock := socket.NewSocket(mockConn, nil, nil)
	// mock server
	mockServer := mockSocket.NewMockServer(mockCtl)
	mockServer.EXPECT().Accept().Return(mockSock, nil).Times(1)
	mockServer.EXPECT().Accept().Return(nil, errors.New("close")).Times(1)

	eventLoop := NewEventLoop(nil)
	listenSession := eventLoop.Bind(mockServer)

	listenSession.Wait()
	time.Sleep(1 * time.Second)
	if listenSession.server != mockServer {
		t.Fatal("eventLoop.Bind failed")
	}
}
