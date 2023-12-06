package event

import (
	"log"
	"testing"
)

type testHandle struct {
	n      int
	state  bool
	result any
}

func (p *testHandle) Activate(*HandleContext) bool {
	p.n += 1
	return p.state
}

func (p *testHandle) Disconnect(*HandleContext) bool {
	p.n += 1
	return p.state
}

func (p *testHandle) Output(_ *HandleContext, data any) (any, bool) {
	p.n += 1
	return data, p.state
}

func (p *testHandle) Input(_ *HandleContext, data string) (string, bool) {
	p.n += 1
	return data, p.state
}

func TestPipeLine_createHandlerWrap(t *testing.T) {
	wrap, err := createHandlerWrap(new(testHandle))
	log.Println(wrap, err)
	if err != nil {
		t.Fatal("PipeLine.createHandlerWrap")
	}

	_, err = createHandlerWrap(newSession())
	log.Println(err)
	if err == nil {
		t.Fatal("PipeLine.createHandlerWrap")
	}
}

func TestPipeLine_inbound(t *testing.T) {
	handler := &testHandle{state: true, result: "data"}
	pipeLine := NewPipeLine()
	_ = pipeLine.AddHandler(handler)
	_ = pipeLine.AddHandler(handler)

	// 建立连接场景
	session := &Session{state: Accept}
	_, _ = pipeLine.inbound(session.InputContext())

	// 断开连接场景
	session = &Session{state: Disconnect}
	_, _ = pipeLine.inbound(session.InputContext())

	// 活动会话收到消息场景
	session = &Session{state: Active}
	session.InputContext().exchange = "test"
	_, _ = pipeLine.inbound(session.InputContext())

	// 交换数据 exchange == nil, 流水线不会触发执行
	session = &Session{state: Active}
	_, _ = pipeLine.inbound(session.InputContext())

	if handler.n != 6 {
		t.Errorf("pipeLine.inbound() error, %d", handler.n)
	}
}

func TestPipeLine_outbound(t *testing.T) {
	pipeLine := NewPipeLine()
	handler := &testHandle{state: true}
	_ = pipeLine.AddHandler(handler)
	_ = pipeLine.AddHandler(handler)

	session := newSession()
	session.state = Active
	session.OutputContext().exchange = "test"
	_, _ = pipeLine.outbound(session.OutputContext())

	session = newSession()
	session.state = Active
	// 交换数据 exchange == nil, 流水线不会触发执行
	_, _ = pipeLine.outbound(session.OutputContext())

	if handler.n != 2 {
		t.Errorf("pipeLine.outbound() error, %d", handler.n)
	}
}
