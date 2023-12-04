package event

import (
	"github.com/google/uuid"
	"github.com/simple-set/simple.io/src/collect"
	"github.com/simple-set/simple.io/src/socket"
	"github.com/sirupsen/logrus"
	"sync"
)

type SessionType int

const (
	ClientType SessionType = iota
	ServerType
	ListenType
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

func newSession() *Session {
	session := &Session{id: uuid.New().String(), state: Accept}
	session.wg.Add(1)
	return session
}

func ListenSession(server socket.Server) *Session {
	session := newSession()
	session.server = server
	session.sessionType = ListenType
	return session
}

func ServerSession(sock *socket.Socket, server socket.Server, pipeline *PipeLine) *Session {
	session := newSession()
	session.sock = sock
	session.server = server
	session.pipeLine = pipeline
	session.sessionType = ServerType
	return session
}

func ClientSession(sock *socket.Socket, client socket.Client, pipeline *PipeLine) *Session {
	session := newSession()
	session.sock = sock
	session.client = client
	session.pipeLine = pipeline
	session.sessionType = ClientType
	return session
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
	p.OutputContext().exchangeBuff = data
	result := p.submitOutput(p.outputContext)

	if result == nil {
		return
	}
	if value, ok := result.(string); ok {
		p.sock.WriteString(value)
		p.Flush()
		return
	}
	if value, ok := result.([]byte); ok {
		p.sock.Write(value)
		p.Flush()
		return
	}
	logrus.Warnf("Sending data is discarded and must be a string or byte array")
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
		context.exchangeBuff = nil
	}()

	p.pipeLine.inbound(context)
}

func (p *Session) submitOutput(context *HandleContext) any {
	if p.pipeLine == nil || context == nil {
		return nil
	}
	if p.outputStack >= 2 {
		logrus.Errorln("Outbound processor call overflow")
		return nil
	}
	defer func() {
		p.outputStack = -1
		context.exchangeBuff = nil
	}()

	p.outputStack += 1
	if result, ok := p.pipeLine.outbound(context); ok {
		return result
	}
	return nil
}

func (p *Session) Wait() {
	p.wg.Wait()
}
