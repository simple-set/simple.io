package socket

import (
	"bufio"
	"errors"
	"io"
	"net"
)

type Socket struct {
	*bufio.Reader
	*bufio.Writer
	server Server
	client Client
	conn   net.Conn
}

func NewSocket(conn net.Conn, server Server, client Client) *Socket {
	return &Socket{
		conn:   conn,
		server: server,
		client: client,
		Reader: bufio.NewReaderSize(conn, 1024),
		Writer: bufio.NewWriterSize(conn, 1024),
	}
}

func (p *Socket) Server() Server {
	return p.server
}

func (p *Socket) Client() Client {
	return p.client
}

func (p *Socket) LocalAddr() net.Addr {
	if p.conn != nil {
		return p.conn.LocalAddr()
	}
	return nil
}

func (p *Socket) RemoteAddr() net.Addr {
	if p.conn != nil {
		return p.conn.RemoteAddr()
	}
	return nil
}

func (p *Socket) Read() (io.Reader, error) {
	if p.conn == nil || p.Reader == nil {
		return nil, errors.New("socket connection not found")
	}

	if _, err := p.Reader.Peek(1); err != nil {
		return nil, err
	}
	return p.Reader, nil
}

func (p *Socket) Close() error {
	if p.conn != nil {
		return p.conn.Close()
	}
	return nil
}
