package event

import (
	"github.com/simple-set/simple.io/src/collect"
	"github.com/sirupsen/logrus"
	"reflect"
)

type PipeLine struct {
	handlers *collect.LinkedNode[any]
}

// NewPipeLine 流水线构造函数
func NewPipeLine(handlers *collect.LinkedNode[any]) *PipeLine {
	return &PipeLine{handlers: handlers}
}

// AddHandler 为流水线添加处理器
func (p *PipeLine) AddHandler(handler any) {
	if !IsHandler(handler) {
		logrus.Errorln("Wrong type, only handler can be added")
	}
	p.handlers.Add(handler)
}

func (p *PipeLine) InsertHandler(index int, handler any) {
	if !IsHandler(handler) {
		logrus.Errorln("Wrong type, only handler can be added")
	}
	err := p.handlers.Insert(index, handler)
	if err != nil {
		logrus.Errorln(err)
	}
}

func (p *PipeLine) inbound(context *HandleContext) (interface{}, bool) {
	if p.handlers == nil || p.handlers.GetFirst() == nil {
		return context.exchangeBuff, true
	}
	pointHandle := p.handlers.GetFirst()
	var result = context.exchangeBuff
	var state = true

	for {
		// 执行处理器
		result, state = p.executeHandle(pointHandle.Value(), context, result)
		if state && p.handlers.Next(pointHandle) != nil {
			// 继续执行下个处理器
			pointHandle = p.handlers.Next(pointHandle)
			continue
		}
		return result, state
	}
}

func (p *PipeLine) outbound(context *HandleContext) (interface{}, bool) {
	if p.handlers == nil || p.handlers.GetLast() == nil {
		return context.exchangeBuff, true
	}
	pointHandle := p.handlers.GetLast()
	var result = context.exchangeBuff
	var state = true

	for {
		// 执行处理器
		result, state = p.executeHandle(pointHandle.Value(), context, result)
		if state && p.handlers.Prev(pointHandle) != nil {
			// 继续执行下个处理器
			pointHandle = p.handlers.Prev(pointHandle)
			continue
		}
		return result, state
	}
}

// 执行处理器
func (p *PipeLine) executeHandle(handler any, context *HandleContext, args any) (any, bool) {
	defer func() {
		// 处理器执行异常
		if err := recover(); err != nil {
			logrus.Errorln("pipeLine Exception handler", reflect.TypeOf(handler).String(), err)
		}
	}()

	if currentHandle, ok := handler.(ActivateHandle); context.session.state == Accept && ok {
		logrus.Debugln("Execute the Activate method of the", reflect.TypeOf(handler).String())
		return currentHandle.Activate(context)
	}

	if handler, ok := handler.(DisconnectHandle); context.session.state == Disconnect && ok {
		logrus.Debugln("Execute the Disconnect method of the", reflect.TypeOf(handler).String())
		return handler.Disconnect(context)
	}

	if context.session.state == Active {
		if handler, ok := handler.(InputHandle); context.direction == inbound && ok {
			logrus.Debugln("Execute the inbound method of the", reflect.TypeOf(handler).String())
			return handler.Input(context, args)
		}
		if handler, ok := handler.(OutputHandle); context.direction == outbound && ok {
			logrus.Debugln("Execute the outbound method of the", reflect.TypeOf(handler).String())
			return handler.Output(context, args)
		}
	}
	return args, true
}
