package example

import (
	"github.com/simple-set/simple.io/src/event"
	"github.com/simple-set/simple.io/src/protocol/simpleHttp"
	"github.com/sirupsen/logrus"
	"strings"
)

type Controller func(request *simpleHttp.Request, response *simpleHttp.Response)

type HttpServerDemo struct {
	addr    string
	session *event.Session
	mapping map[string]Controller
}

func (h *HttpServerDemo) Start() {
	bootstrap := event.NewBootstrap()
	bootstrap.TcpServer(h.addr)
	bootstrap.AddHandler(simpleHttp.NewHttpDecoder())
	bootstrap.AddHandler(h)

	h.session = bootstrap.Bind()
	h.session.Wait()
}

func (h *HttpServerDemo) Input(context *event.HandleContext, request *simpleHttp.Request) (*simpleHttp.Request, bool) {
	logrus.Println("path="+request.URL.RequestURI(), ", method="+request.Method, ", status=", request.Response.StatusCode())

	request.Response.AddCookie("sessionId", context.Session().Id())
	h.dispatch(request)
	return request, true
}

func (h *HttpServerDemo) dispatch(request *simpleHttp.Request) {
	for path := range h.mapping {
		if strings.HasPrefix(request.URL.Path, path) {
			h.mapping[path](request, request.Response)
			return
		}
	}

	request.Response.SetStatusCode(404)
	_, _ = request.Response.Write([]byte("404 not found"))
}

func (h *HttpServerDemo) AddController(path string, controller Controller) {
	if h.mapping == nil {
		h.mapping = make(map[string]Controller)
	}
	h.mapping[path] = controller
}

func NewHttpServerDemo(addr string) {
	httpServerDemo := &HttpServerDemo{addr: addr}
	httpServerDemo.AddController("/index", func(request *simpleHttp.Request, response *simpleHttp.Response) {
		_, _ = response.Write([]byte("hello "))
		_, _ = response.Write([]byte("world"))
	})
	httpServerDemo.Start()
}
