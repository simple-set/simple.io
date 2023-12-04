package event

import (
	"github.com/simple-set/simple.io/src/collect"
	"github.com/simple-set/simple.io/src/socket"
	"github.com/sirupsen/logrus"
	"reflect"
	"sync"
)

type Bootstrap struct {
	mutex    sync.Mutex
	server   socket.Server
	client   socket.Client
	handlers *collect.LinkedNode[any]
}

func NewBootstrap() *Bootstrap {
	return &Bootstrap{handlers: collect.NewLinkedNode[any]()}
}

// TcpServer 设置tcp服务器的地址和端口, 如 localhost:8080、127.0.0.1:8080、0.0.0.0:8080、:8080
func (p *Bootstrap) TcpServer(addr string) *Bootstrap {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if p.server != nil {
		logrus.Panicln("Repeatedly setting the socket server")
	}

	// 设置socket服务器
	socketServer, err := socket.NewTcpServer(addr)
	if err != nil {
		logrus.Panicln(err)
	}
	p.server = socketServer
	return p
}

func (p *Bootstrap) TcpClient(addr string) *Bootstrap {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	if p.server != nil {
		logrus.Panicln("Repeatedly setting the socket server")
	}

	client, err := socket.NewTcpClient(addr)
	if err != nil {
		logrus.Panicln(err)
	}
	p.client = client
	return p
}

// AddHandler 添加事件处理器
func (p *Bootstrap) AddHandler(handler any) *Bootstrap {
	if !IsHandler(handler) {
		logrus.Panicln("Wrong type, only handler can be added")
	}
	logrus.Debugf("Add handler: %s", reflect.TypeOf(handler).String())
	p.handlers.Add(handler)
	return p
}

func (p *Bootstrap) Connect() *Session {
	return NewEventLoop(p.handlers).Connect(p.client)
}

func (p *Bootstrap) Bind() *Session {
	return NewEventLoop(p.handlers).Bind(p.server)
}
