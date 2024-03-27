package event

import (
	"github.com/simple-set/simple.io/src/socket"
	"github.com/sirupsen/logrus"
)

type Loop struct {
	pipeLine *PipeLine
}

func NewEventLoop(pipeLine *PipeLine) *Loop {
	return &Loop{pipeLine: pipeLine}
}

func (p *Loop) Connect(client socket.Client) (*Session, error) {
	sock, err := client.Connect()
	if err != nil {
		return nil, err
	}
	clientSession := NewClientSession(sock, client, p.pipeLine)
	logrus.Debugf("Create a new clientSession, Id: %s, addr: %s", clientSession.Id(), sock.RemoteAddr())
	clientSession.submitInput(clientSession.InputContext())
	clientSession.state = Active
	go p.pool(clientSession)
	return clientSession, nil
}

func (p *Loop) Bind(server socket.Server) *Session {
	listenSession := newListenSession(server)
	go p.accept(listenSession)
	return listenSession
}

func (p *Loop) accept(listenSession *Session) {
	for {
		if sock, err := listenSession.server.Accept(); err == nil {
			go func(sock *socket.Socket) {
				serverSession := NewServerSession(sock, listenSession.server, p.pipeLine)
				logrus.Debugf("Create a new serverSession, Id: %s, addr: %s", serverSession.Id(), sock.RemoteAddr())

				serverSession.submitInput(serverSession.InputContext())
				serverSession.state = Active
				p.pool(serverSession)
			}(sock)
		} else {
			logrus.Errorln("connection Accept exception", err)
			_ = listenSession.Close()
			return
		}
	}
}

func (p *Loop) pool(session *Session) {
	for {
		if reader, err := session.sock.Read(); err == nil {
			context := session.InputContext()
			context.exchange = reader
			session.submitInput(context)
		} else {
			logrus.Debugln("Exception reading data from socket", err)
			_ = session.Close()
			context := session.InputContext()
			context.exchange = nil
			session.submitInput(context)
			return
		}
	}
}
