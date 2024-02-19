package example

import (
	"github.com/simple-set/simple.io/src/event"
	"github.com/simple-set/simple.io/src/protocol/simpleHttp"
	"github.com/sirupsen/logrus"
	"time"
)

type SimpleHttpServer struct{}

func (h *SimpleHttpServer) Start() {
	bootstrap := event.NewBootstrap()
	bootstrap.TcpServer("localhost:8000")
	bootstrap.AddHandler(simpleHttp.NewHttpDecoder())
	bootstrap.AddHandler(h)
	bootstrap.Bind().Wait()
}

func (h *SimpleHttpServer) Input(context *event.HandleContext, request *simpleHttp.Request) (*simpleHttp.Request, bool) {
	logrus.Println("path=", request.URL.Path, ", method=", request.Method, ", status=", request.Response.StatusCode())
	request.Response.AddCookie("data", time.Now().Format(time.RFC3339))
	if _, err := request.Response.Write([]byte("Hello, world!")); err != nil {
		logrus.Errorln(err)
		return nil, false
	}
	return request, true
}
