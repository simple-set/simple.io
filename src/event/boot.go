package event

import (
	"github.com/simple-set/simple.io/src/socket"
	"github.com/sirupsen/logrus"
	"sync"
)

type Bootstrap struct {
	mutex    sync.Mutex
	server   socket.Server
	client   socket.Client
	pipeLine *PipeLine
}

func NewBootstrap() *Bootstrap {
	return &Bootstrap{pipeLine: NewPipeLine()}
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

func (p *Bootstrap) AddHandler(handler any) *Bootstrap {
	if err := p.pipeLine.AddHandler(handler); err != nil {
		logrus.Panic(err)
	}
	return p
}

func (p *Bootstrap) Connect() *Session {
	return NewEventLoop(p.pipeLine).Connect(p.client)
}

func (p *Bootstrap) Bind() *Session {
	return NewEventLoop(p.pipeLine).Bind(p.server)
}
