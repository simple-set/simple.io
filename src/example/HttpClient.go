package example

import (
	"github.com/simple-set/simple.io/src/event"
	"github.com/simple-set/simple.io/src/protocol/simpleHttp"
	"github.com/sirupsen/logrus"
	"sync"
)

type SimpleHttpClient struct {
	wg       sync.WaitGroup
	response *simpleHttp.Response
}

func (s *SimpleHttpClient) Input(_ *event.HandleContext, response *simpleHttp.Response) (any, bool) {
	s.response = response
	logrus.Infoln(response)
	s.wg.Done()
	return response, true
}

func (s *SimpleHttpClient) makeRequest(url string) *simpleHttp.Request {
	return simpleHttp.NewRequestUrl(url)
}

func (s *SimpleHttpClient) connect(addr string) (*event.Session, error) {
	bootstrap := event.NewBootstrap()
	bootstrap.TcpClient(addr)
	bootstrap.AddHandler(simpleHttp.NewClientHandler())
	bootstrap.AddHandler(s)
	return bootstrap.Connect()
}

func (s *SimpleHttpClient) Get(url string) (*simpleHttp.Response, error) {
	request := s.makeRequest(url)
	session, err := s.connect(request.Host)
	if err != nil {
		return nil, err
	}

	s.wg.Add(1)
	session.Write(request)
	s.wg.Wait()
	return s.response, nil
}

func NewSimpleHttpClient() *SimpleHttpClient {
	return &SimpleHttpClient{}
}
