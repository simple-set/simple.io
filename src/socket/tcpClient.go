package socket

import (
	"net"
)

type TcpClient struct {
	remoteAddr *net.TCPAddr
}

func NewTcpClient(addr string) (*TcpClient, error) {
	tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		return nil, err
	}
	return &TcpClient{remoteAddr: tcpAddr}, nil
}

func (t *TcpClient) Connect() (*Socket, error) {
	conn, err := net.DialTCP("tcp", nil, t.remoteAddr)
	if err != nil {
		return nil, err
	}
	return NewSocket(conn, nil, t), nil
}

func (t *TcpClient) Close() error {
	//TODO implement me
	panic("implement me")
}
