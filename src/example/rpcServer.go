package example

import (
	"errors"
	"fmt"
	"github.com/simple-set/simple.io/src/event"
	"github.com/simple-set/simple.io/src/protocol/simpleRpc"
	"github.com/sirupsen/logrus"
	"time"
)

type Service func(args []byte) any

type RpcServer struct {
	addr    string
	session *event.Session
	mapping map[string]Service
}

// Input 请求入口
func (r *RpcServer) Input(context *event.HandleContext, message *simpleRpc.Message) (any, bool) {
	result, err := r.dispatch(message)
	if err != nil {
		logrus.Errorln("rpc error", err)
		_ = context.Session().Close()
		return nil, false
	}
	if result == nil {
		return nil, false
	}

	if value, ok := result.(string); ok {
		return simpleRpc.NewResponse([]byte(value)), true
	}
	return nil, false
}

// 路由分发器
func (r *RpcServer) dispatch(message *simpleRpc.Message) (any, error) {
	if len(message.Method()) <= 0 {
		return nil, errors.New("method name cannot be empty")
	}

	if service := r.mapping[message.Method()]; service == nil {
		return nil, errors.New(message.Method() + ", method not found")
	} else {
		return r.callService(service, message.Body())
	}
}

// 执行业务函数
func (r *RpcServer) callService(service Service, args []byte) (result any, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = errors.New(fmt.Sprintln("Method execution exception, ", e))
		}
	}()

	result = service(args)
	return
}

func (r *RpcServer) Start() {
	boot := event.NewBootstrap()
	boot.TcpServer(r.addr)
	boot.AddHandler(simpleRpc.NewHandle())
	boot.AddHandler(r)

	r.session = boot.Bind()
	r.session.Wait()
}

func (r *RpcServer) registerServer(name string, service Service) {
	if r.mapping == nil {
		r.mapping = make(map[string]Service)
	}
	r.mapping[name] = service
}

func NewRpcServer(addr string) *RpcServer {
	server := &RpcServer{addr: addr}
	server.registerServer("hello", hello)
	return server
}

func hello(_ []byte) any {
	return time.Now().Format(time.RFC3339) + " hello World"
}
