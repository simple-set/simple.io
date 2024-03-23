package example

import (
	"github.com/simple-set/simple.io/src/event"
	"github.com/simple-set/simple.io/src/protocol/simpleRpc"
	"log"
)

type RpcClient struct {
	addr    string
	session *event.Session
}

func (r *RpcClient) Input(_ *event.HandleContext, message *simpleRpc.Message) (any, bool) {
	log.Println(string(message.Body()))
	return nil, true
}

func (r *RpcClient) connect() {
	bootstrap := event.NewBootstrap()
	bootstrap.TcpClient(r.addr)
	bootstrap.AddHandler(simpleRpc.NewHandle())
	bootstrap.AddHandler(r)

	if session, err := bootstrap.Connect(); err != nil {
		panic(err)
	} else {
		r.session = session
	}
}

func (r *RpcClient) Call(name string, args []byte) (result any) {
	request := simpleRpc.NewRequest(args)
	request.SetMethod(name)
	r.session.Write(request)
	return
}

func NewRpcClient(addr string) *RpcClient {
	client := &RpcClient{addr: addr}
	client.connect()
	return client
}
