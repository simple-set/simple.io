package event

import (
	"github.com/simple-set/simple.io/src/collect"
	"github.com/sirupsen/logrus"
	"reflect"
)

// 处理的包装类型
type handlerWrap struct {
	name          string // 处理器名称
	isActivate    bool   // 是否支持Activate
	isDisconnect  bool   // 是否支持Disconnect
	isInput       bool   // 是否支持input
	isOutput      bool   // 是否支持output
	inputArgType  any    // 输入参数类型
	outputArgType any    // 输出参数类型
}

func parseArgs(handler any, funcName string) any {
	method, ok := reflect.TypeOf(handler).MethodByName(funcName)
	if ok {
		return method.Type.In(2)
	}
	return nil
}

// 生成包装对象
func createHandlerWrap(handler any) *handlerWrap {
	warp := new(handlerWrap)
	warp.name = reflect.TypeOf(handler).String()

	if _, ok := handler.(ActivateHandle[any]); ok {
		warp.isActivate = true
	}

	if _, ok := handler.(DisconnectHandle[any]); ok {
		warp.isDisconnect = true
	}

	if value, ok := handler.(InputHandle[any, any]); ok {
		warp.isInput = true
		warp.inputArgType = parseArgs(value, "Input")
	}

	if _, ok := handler.(OutputHandle[any, any]); ok {
		warp.isOutput = true
		//typeValue := reflect.TypeOf(value)
		//nameFunc, _ := typeValue.FieldByNameFunc(func(string) bool { return true })
		//warp.inputArgType = parseArgs(nameFunc)
	}
	return warp
}

type PipeLine struct {
	handlers  *collect.LinkedNode[any]
	handlerss *collect.LinkedNode[handlerWrap]
	//handlerInfos map[any][]*handlerInfo
}

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

func (p *PipeLine) execute(handler any, context *HandleContext, exchange any) (any, bool) {
	//defer func() {
	//	// 处理器执行异常
	//	if err := recover(); err != nil {
	//		logrus.Errorln("Exception in executing handler", reflect.TypeOf(handler).String(), err)
	//	}
	//}()
	//
	//if currentHandle, ok := handler.(ActivateHandle[any]); context.session.state == Accept && ok {
	//	logrus.Debugln("Execute the Activate method of the", reflect.TypeOf(handler).String())
	//	return currentHandle.Activate(context)
	//}
	//
	//if handler, ok := handler.(DisconnectHandle[any]); context.session.state == Disconnect && ok {
	//	logrus.Debugln("Execute the Disconnect method of the", reflect.TypeOf(handler).String())
	//	return handler.Disconnect(context)
	//}
	//
	//if context.session.state == Active {
	//	if handler, ok := handler.(InputHandle[any, string]); context.direction == inbound && ok {
	//		logrus.Debugln("Execute the inbound method of the", reflect.TypeOf(handler).String())
	//		return handler.Input(context, exchange)
	//	}
	//	if handler, ok := handler.(OutputHandle); context.direction == outbound && ok {
	//		logrus.Debugln("Execute the outbound method of the", reflect.TypeOf(handler).String())
	//		return handler.Output(context, exchange)
	//	}
	//}
	return exchange, true
}

func (p *PipeLine) doExecute() (any, bool) {
	return nil, true
}

func NewPipeLine(handlers *collect.LinkedNode[any]) *PipeLine {
	return &PipeLine{handlers: handlers}
}
