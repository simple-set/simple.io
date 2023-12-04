package event

import (
	"github.com/simple-set/simple.io/src/collect"
	"github.com/simple-set/simple.io/src/socket"
	"github.com/sirupsen/logrus"
)

type Loop struct {
	handlers *collect.LinkedNode[any]
}

func NewEventLoop(handles *collect.LinkedNode[any]) *Loop {
	return &Loop{handlers: handles}
}

func (p *Loop) Connect(client socket.Client) *Session {
	sock, err := client.Connect()
	if err != nil {
		logrus.Errorln(err)
		return nil
	}
	clientSession := ClientSession(sock, client, NewPipeLine(p.handlers))
	logrus.Debugf("Create a new clientSession, Id: %s, addr: %s", clientSession.Id(), sock.RemoteAddr())
	clientSession.submitInput(clientSession.InputContext())
	clientSession.state = Active
	go p.pool(clientSession)
	return clientSession
}

func (p *Loop) Bind(server socket.Server) *Session {
	listenSession := ListenSession(server)
	go p.accept(listenSession)
	return listenSession
}

func (p *Loop) accept(listenSession *Session) {
	for {
		if sock, err := listenSession.server.Accept(); err == nil {
			go func(sock *socket.Socket) {
				serverSession := ServerSession(sock, listenSession.server, NewPipeLine(p.handlers))
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
			context.exchangeBuff = reader
			session.submitInput(context)
		} else {
			logrus.Debugln("Exception reading data from socket", err)
			_ = session.Close()
			context := session.InputContext()
			context.exchangeBuff = nil
			session.submitInput(context)
			return
		}
	}
}
