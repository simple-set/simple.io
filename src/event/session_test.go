package event

import (
	"github.com/golang/mock/gomock"
	mockNet "github.com/simple-set/simple.io/src/mocks/net"
	"github.com/simple-set/simple.io/src/socket"
	"github.com/sirupsen/logrus"
	"net"
	"testing"
)

func init() {
	logrus.SetLevel(logrus.TraceLevel)
}

func TestSession_Id(t *testing.T) {
	session := newSession()
	if session.id != session.Id() {
		t.Fatalf("session.Id() error, %s", session.id)
	}
}

func TestSession_InputContext(t *testing.T) {
	session := newSession()
	context := session.InputContext()
	if context.direction != inbound || context.Session() != session {
		t.Fatal("session.inbound()")
	}
}

func TestSession_OutputContext(t *testing.T) {
	session := newSession()
	context := session.OutputContext()
	if context.direction != outbound || context.Session() != session {
		t.Fatal("session.OutputContext()")
	}
}

func TestSession_Close(t *testing.T) {
	mockConn := mockNet.NewMockConn(gomock.NewController(t))
	mockConn.EXPECT().Close().Return(nil).Times(1)
	mockConn.EXPECT().RemoteAddr().Return(&net.IPAddr{}).Times(1)
	newSocket := socket.NewSocket(mockConn, nil, nil)

	session := newSession()
	session.sock = newSocket
	err := session.Close()
	if err != nil || session.state != Disconnect {
		t.Fatal("session.Close()")
	}
}

func TestSession_Write(t *testing.T) {
	// mocks conn
	mockConn := mockNet.NewMockConn(gomock.NewController(t))
	mockConn.EXPECT().Write(gomock.Any()).Return(4, nil).AnyTimes()

	session := newSession()
	session.pipeLine = NewPipeLine()
	session.sock = socket.NewSocket(mockConn, nil, nil)
	session.Write("data")
}

func TestSession_Flush(t *testing.T) {
	data := []byte("data")
	mockConn := mockNet.NewMockConn(gomock.NewController(t))
	mockConn.EXPECT().Write(gomock.Any()).Return(4, nil).AnyTimes()

	session := newSession()
	session.sock = socket.NewSocket(mockConn, nil, nil)
	session.pipeLine = NewPipeLine()
	session.Write(data)
	session.Flush()
}

func TestSession_submitInput(t *testing.T) {
	data := []byte("data")
	mockConn := mockNet.NewMockConn(gomock.NewController(t))
	mockConn.EXPECT().Write(gomock.Any()).Return(len(data), nil).AnyTimes()

	handle := &testHandle{state: true, result: "data"}
	pipeLine := NewPipeLine()
	_ = pipeLine.AddHandler(handle)
	_ = pipeLine.AddHandler(handle)

	session := NewClientSession(nil, nil, pipeLine)
	session.sock = socket.NewSocket(mockConn, nil, nil)
	session.state = Active
	session.InputContext().exchange = "data"
	session.submitInput(session.InputContext())

	if handle.n != 4 || session.InputContext().exchange != nil {
		t.Fatal("session.submitInput()")
	}
}

func TestSession_submitOutput(t *testing.T) {
	data := []byte("data")
	mockConn := mockNet.NewMockConn(gomock.NewController(t))
	mockConn.EXPECT().Write(gomock.Any()).Return(len(data), nil).AnyTimes()

	handler := &testHandle{state: true, result: "data"}
	pipeLine := NewPipeLine()
	_ = pipeLine.AddHandler(handler)
	_ = pipeLine.AddHandler(handler)

	session := NewClientSession(nil, nil, pipeLine)
	session.sock = socket.NewSocket(mockConn, nil, nil)
	session.state = Disconnect
	context := session.InputContext()
	context.exchange = "data"
	session.submitOutput(context)

	if handler.n != 2 || context.exchange != nil {
		t.Fatal("session.submitOutput()")
	}
}
