package example

import (
	"github.com/simple-set/simple.io/src/event"
	"github.com/simple-set/simple.io/src/protocol/simpleHttp"
	"github.com/simple-set/simple.io/src/version"
	"github.com/sirupsen/logrus"
	"sync"
)

type SimpleHttpClient struct {
	wg       sync.WaitGroup
	response *simpleHttp.Response
}

func (s *SimpleHttpClient) Input(context *event.HandleContext, response *simpleHttp.Response) (any, bool) {
	logrus.Infoln(response)
	s.wg.Done()
	return response, true
}

func NewSimpleHttpClient() *SimpleHttpClient { return &SimpleHttpClient{} }

func (s *SimpleHttpClient) Connect() *simpleHttp.Response {
	bootstrap := event.NewBootstrap()
	bootstrap.TcpClient("153.3.238.110:80")
	bootstrap.AddHandler(simpleHttp.NewHttpEncoder())
	bootstrap.AddHandler(s)
	session, err := bootstrap.Connect()
	if err != nil {
		logrus.Fatal(err)
	}

	request := simpleHttp.NewRequestBuild().
		Uri("https://www.baidu.com/index?a=1&b=2&c").
		Agent(version.Name+"/"+version.Version).
		Cookie("simple.id", session.Id()).
		Proto("HTTP/1.1").
		Get().
		Build()

	s.wg.Add(1)
	session.Write(request)
	s.wg.Wait()
	return s.response
}
