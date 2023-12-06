package event

import (
	"errors"
	"github.com/simple-set/simple.io/src/collect"
	"github.com/sirupsen/logrus"
	"reflect"
)

// 处理的包装类型
// 使用包装器的主要原因是泛化处理器, 由于go的泛型无法断言(像java一样多好啊), 包装器维护反射对象调用处理器
type handlerWrap struct {
	name             string
	activateMethod   *reflect.Value
	disconnectMethod *reflect.Value
	inputMethod      *reflect.Value
	outputMethod     *reflect.Value
}

// 验证处理
func verifyHandler(value reflect.Value, num int) (*reflect.Value, error) {
	if value.Kind() != reflect.Func {
		return nil, errors.New("invalid number of arguments")
	}

	if value.Type().NumIn() != num || value.Type().NumOut() != num {
		return nil, errors.New("invalid number of arguments")
	}

	if value.Type().In(0).Elem().Name() != "HandleContext" {
		return nil, errors.New("invalid number of arguments")
	}

	numOut := value.Type().NumOut()
	if value.Type().Out(numOut-1).Kind() != reflect.Bool {
		return nil, errors.New("invalid number of arguments")
	}

	return &value, nil
}

// 创建处理器的包装对象
func createHandlerWrap(handler any) (*handlerWrap, error) {
	typeValue := reflect.ValueOf(handler)
	warp := new(handlerWrap)
	warp.name = typeValue.Elem().Type().String()

	if method, err := verifyHandler(typeValue.MethodByName("Activate"), 1); err == nil {
		warp.activateMethod = method
	}
	if method, err := verifyHandler(typeValue.MethodByName("Disconnect"), 1); err == nil {
		warp.disconnectMethod = method
	}
	if method, err := verifyHandler(typeValue.MethodByName("Input"), 2); err == nil {
		warp.inputMethod = method
	}
	if method, err := verifyHandler(typeValue.MethodByName("Output"), 2); err == nil {
		warp.outputMethod = method
	}
	if warp.inputMethod == nil && warp.outputMethod == nil && warp.activateMethod == nil && warp.disconnectMethod == nil {
		return nil, errors.New("the handler must have one of the Activate, Disconnect, Input, or Output methods")
	}
	return warp, nil
}

type PipeLine struct {
	handlers *collect.LinkedNode[*handlerWrap]
}

func (p *PipeLine) AddHandler(handler any) (err error) {
	wrap, err := createHandlerWrap(handler)
	if err == nil {
		p.handlers.Add(wrap)
	}
	return
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
			logrus.Panicln("Exception in executing handler", wrap.name, exchange, err)
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

	if context.session.state == Active {
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
