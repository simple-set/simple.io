package event

import (
	"github.com/simple-set/simple.io/src/collect"
	"github.com/sirupsen/logrus"
	"reflect"
)

type PipeLine struct {
	handlers *collect.LinkedNode[*handlerWrap]
}

func (p *PipeLine) AddHandler(handler any) error {
	if wrap, err := createHandlerWrap(handler); err == nil {
		p.handlers.Add(wrap)
		return nil
	} else {
		return err
	}
}

func (p *PipeLine) InsertHandler(index int, handler any) (err error) {
	wrap, err := createHandlerWrap(handler)
	if err == nil {
		_ = p.handlers.Insert(index, wrap)
	}
	return
}

func (p *PipeLine) inbound(context *HandleContext) (interface{}, bool) {
	if p.handlers == nil || p.handlers.GetFirst() == nil {
		return context.exchange, true
	}
	pointHandle := p.handlers.GetFirst()
	var exchange = context.exchange
	var state = true

	for {
		// 执行处理器
		exchange, state = p.execute(pointHandle.Value(), context, exchange)
		if state && p.handlers.Next(pointHandle) != nil {
			// 继续执行下个处理器
			pointHandle = p.handlers.Next(pointHandle)
			continue
		}
		return exchange, state
	}
}

func (p *PipeLine) outbound(context *HandleContext) (interface{}, bool) {
	if p.handlers == nil || p.handlers.GetLast() == nil {
		return context.exchange, true
	}
	pointHandle := p.handlers.GetLast()
	var exchange = context.exchange
	var state = true

	for {
		// 执行处理器
		exchange, state = p.execute(pointHandle.Value(), context, exchange)
		if state && p.handlers.Prev(pointHandle) != nil {
			// 继续执行下个处理器
			pointHandle = p.handlers.Prev(pointHandle)
			continue
		}
		return exchange, state
	}
}

func (p *PipeLine) execute(wrap *handlerWrap, context *HandleContext, exchange any) (result any, state bool) {
	result = exchange
	state = true

	defer func() {
		// 处理器执行预期外异常
		if err := recover(); err != nil {
			logrus.Errorln("Exception in executing handler", wrap.name, exchange, err)
			_ = context.session.Close()
		}
	}()

	if context.session.state == Accept && wrap.activateMethod != nil {
		value := wrap.activateMethod.Call([]reflect.Value{reflect.ValueOf(context)})
		state = value[0].Bool()
		return
	}
	if context.session.state == Disconnect && wrap.disconnectMethod != nil {
		value := wrap.disconnectMethod.Call([]reflect.Value{reflect.ValueOf(context)})
		state = value[0].Bool()
		return
	}

	if context.session.state == Active && exchange != nil {
		var results []reflect.Value
		if context.direction == inbound && wrap.inputMethod != nil {
			results = wrap.inputMethod.Call([]reflect.Value{reflect.ValueOf(context), reflect.ValueOf(exchange)})
		} else if context.direction == outbound && wrap.outputMethod != nil {
			results = wrap.outputMethod.Call([]reflect.Value{reflect.ValueOf(context), reflect.ValueOf(exchange)})
		}

		if results != nil {
			result = results[0].Interface()
			state = results[1].Bool()
		}
	}
	return
}

func NewPipeLine() *PipeLine {
	return &PipeLine{handlers: collect.NewLinkedNode[*handlerWrap]()}
}
