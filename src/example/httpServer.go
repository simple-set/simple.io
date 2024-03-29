package example

import (
	"github.com/simple-set/simple.io/src/event"
	"github.com/simple-set/simple.io/src/protocol/simpleHttp"
	"github.com/sirupsen/logrus"
	"net/http"
	"strings"
	"time"
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
	bootstrap.AddHandler(simpleHttp.NewServerHandler())
	bootstrap.AddHandler(h)

	h.session = bootstrap.Bind()
	h.session.Wait()
}

func (h *HttpServerDemo) Input(context *event.HandleContext, request *simpleHttp.Request) (any, bool) {
	request.Response = simpleHttp.NewReplyResponse(request)
	request.Response.AddCookie("sessionId", context.Session().Id())

	h.dispatch(request, request.Response)
	context.Session().Write(request.Response)

	logrus.Println("path="+request.URL.RequestURI(), ", method="+request.Method, ", status=", request.Response.StatusCode())
	return request, true
}

func (h *HttpServerDemo) dispatch(request *simpleHttp.Request, response *simpleHttp.Response) {
	for path := range h.mapping {
		if strings.HasPrefix(request.URL.Path, path) {
			h.mapping[path](request, response)
			return
		}
	}

	request.Response.SetStatusCode(404)
	body := simpleHttp.MakeReaderWriteBody()
	_, _ = body.WriteString("404 not found, " + time.Now().Format(http.TimeFormat))
	response.SetBody(body)
}

func (h *HttpServerDemo) AddController(path string, controller Controller) {
	if h.mapping == nil {
		h.mapping = make(map[string]Controller)
	}
	h.mapping[path] = controller
}

func NewHttpServerDemo(addr string) {
	httpServerDemo := &HttpServerDemo{addr: addr}
	httpServerDemo.AddController("/index", indexController)
	httpServerDemo.Start()
}

func indexController(request *simpleHttp.Request, response *simpleHttp.Response) {
	if request.Body() != nil && request.Body().Size() > 0 {
		bytes := make([]byte, request.Body().Size())
		if _, err := request.Body().Read(bytes); err != nil {
			return
		}
		if _, err := response.Body().Write(bytes); err == nil {
			return
		}
	}

	body := simpleHttp.MakeReaderWriteBody()
	_, _ = body.WriteString("hello world")
	response.SetBody(body)
}
