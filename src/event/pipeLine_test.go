package event

import (
	"github.com/simple-set/simple.io/src/collect"
	"testing"
)

type testHandle struct {
	n      int
	state  bool
	result any
}

func (p *testHandle) Activate(*HandleContext) (interface{}, bool) {
	p.n += 1
	return p.result, p.state
}

func (p *testHandle) Disconnect(*HandleContext) (interface{}, bool) {
	p.n += 1
	return p.result, p.state
}

func (p *testHandle) Output(_ *HandleContext, data any) (interface{}, bool) {
	p.n += 1
	return data, p.state
}

func (p *testHandle) Input(_ *HandleContext, data any) (interface{}, bool) {
	p.n += 1
	return data, p.state
}

func TestPipeLine_inbound(t *testing.T) {
	handles := collect.NewLinkedNode[any]()
	handler := &testHandle{state: true, result: "data"}
	handles.Add(handler)
	handles.Add(handler)

	session := newSession()
	context := session.OutputContext()

	pipeLine := NewPipeLine(handles)
	result, _ := pipeLine.inbound(context)

	if result != handler.result {
		t.Errorf("pipeLine.inbound() error, %s", result)
	}
	if handler.n != 2 {
		t.Errorf("pipeLine.inbound() error, %d", handler.n)
	}
}

func TestPipeLine_outbound(t *testing.T) {
	handles := collect.NewLinkedNode[any]()
	handler := &testHandle{state: true}
	handles.Add(handler)
	handles.Add(handler)

	session := newSession()
	session.state = Active
	context := session.InputContext()
	context.exchangeBuff = "data"

	pipeLine := NewPipeLine(handles)
	result, _ := pipeLine.outbound(context)

	if result != "data" {
		t.Errorf("pipeLine.outbound() error, %s", result)
	}

	if handler.n != 2 {
		t.Errorf("pipeLine.outbound() error, %d", handler.n)
	}
}
