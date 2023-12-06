package socket

import (
	"github.com/sirupsen/logrus"
	"net"
)

type TcpServer struct {
	listenAddr *net.TCPAddr
	listenTcp  *net.TCPListener
}

func NewTcpServer(addr string) (*TcpServer, error) {
	// 创建tcp地址
	tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		return nil, err
	}

	// 创建TCP监听器
	listenTCP, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		return nil, err
	}

	logrus.Infoln("Start the TCP socket at address", tcpAddr)
	return &TcpServer{tcpAddr, listenTCP}, nil
}

func (p *TcpServer) Accept() (*Socket, error) {
	conn, err := p.listenTcp.Accept()
	if err != nil {
		return nil, err
	}
	return NewSocket(conn, p, nil), nil
}

func (p *TcpServer) Close() error {
	err := p.listenTcp.Close()
	return err
}
