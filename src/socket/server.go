package socket

type Server interface {
	Accept() (*Socket, error)
	Close() error
}

type Client interface {
	Connect() (*Socket, error)
	Close() error
}
