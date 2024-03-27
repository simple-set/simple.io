package event

import (
	"errors"
	"reflect"
)

// 处理器的包装类型
// 使用包装类型的主要原因是泛化处理器, 由于go的泛型无法断言, 包装器维护反射对象调用处理器
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
