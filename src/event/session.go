package event

import (
	"github.com/google/uuid"
	"github.com/simple-set/simple.io/src/collect"
	"github.com/simple-set/simple.io/src/socket"
	"github.com/sirupsen/logrus"
	"io"
	"sync"
)

type SessionType int

const (
	ClientSession SessionType = iota
	ServerSession
	ListenSession
)

type SessionState int

const (
	Accept SessionState = iota
	Active
	Disconnect
)

type Session struct {
	id            string
	state         SessionState
	sessionType   SessionType
	server        socket.Server
	client        socket.Client
	sock          *socket.Socket
	pipeLine      *PipeLine
	inputContext  *HandleContext
	outputContext *HandleContext
	outputStack   int
	handles       *collect.LinkedNode[any]
	wg            sync.WaitGroup
}

func (p *Session) Id() string {
	return p.id
}

func (p *Session) State() SessionState {
	return p.state
}

func (p *Session) Sock() *socket.Socket {
	return p.sock
}

func (p *Session) PipeLine() *PipeLine {
	return p.pipeLine
}

func (p *Session) SessionType() SessionType {
	return p.sessionType
}

func (p *Session) InputContext() *HandleContext {
	if p.inputContext == nil {
		p.inputContext = NewContext(p)
		p.inputContext.direction = inbound
	}
	return p.inputContext
}

func (p *Session) OutputContext() *HandleContext {
	if p.outputContext == nil {
		p.outputContext = NewContext(p)
		p.outputContext.direction = outbound
	}
	return p.outputContext
}

func (p *Session) Close() error {
	if p.state == Disconnect {
		return nil
	}
	defer func() {
		p.wg.Done()
		p.state = Disconnect
	}()

	if p.sock != nil {
		err := p.sock.Close()
		logrus.Debugf("Close session, Id: %s, addr: %s, %v", p.Id(), p.sock.RemoteAddr().String(), err)
		return err
	}
	return nil
}

func (p *Session) Write(data any) {
	p.OutputContext().exchange = data
	p.submitOutput(p.outputContext)
}

func (p *Session) WriteSocket(data any) (n int, err error) {
	if value, ok := data.(string); ok {
		n, err = p.sock.WriteString(value)
	} else if value, ok := data.([]byte); ok {
		n, err = p.sock.Write(value)
	} else if value, ok := data.(io.Reader); ok {
		_, err = p.sock.ReadFrom(value)
	}

	if err != nil {
		_ = p.Close()
	}
	p.Flush()
	return
}

func (p *Session) Flush() {
	if p.sock == nil {
		return
	}
	if err := p.sock.Flush(); err != nil {
		logrus.Warnln(err)
		_ = p.Close()
	}
}

func (p *Session) WriteAndFlush(data any) {
	p.Write(data)
	p.Flush()
}

func (p *Session) submitInput(context *HandleContext) {
	if p.pipeLine == nil || context == nil {
		return
	}
	defer func() {
		context.exchange = nil
	}()

	result, state := p.pipeLine.inbound(context)
	if state && result != nil {
		p.OutputContext().exchange = result
		p.submitOutput(p.outputContext)
	}
}

func (p *Session) submitOutput(context *HandleContext) {
	if p.pipeLine == nil || context == nil {
		return
	}
	if p.outputStack >= 2 {
		logrus.Errorln("Outbound processor call overflow")
		return
	}

	p.outputStack += 1
	defer func() {
		p.outputStack = -1
		context.exchange = nil
	}()

	if result, state := p.pipeLine.outbound(context); state && result != nil {
		_, _ = p.WriteSocket(result)
	}
}

func (p *Session) Wait() {
	p.wg.Wait()
}

func newSession() *Session {
	session := &Session{id: uuid.New().String(), state: Accept}
	session.wg.Add(1)
	return session
}

func newListenSession(server socket.Server) *Session {
	session := newSession()
	session.server = server
	session.sessionType = ListenSession
	return session
}

func NewServerSession(sock *socket.Socket, server socket.Server, pipeline *PipeLine) *Session {
	session := newSession()
	session.sock = sock
	session.server = server
	session.pipeLine = pipeline
	session.sessionType = ServerSession
	return session
}

func NewClientSession(sock *socket.Socket, client socket.Client, pipeline *PipeLine) *Session {
	session := newSession()
	session.sock = sock
	session.client = client
	session.pipeLine = pipeline
	session.sessionType = ClientSession
	return session
}
