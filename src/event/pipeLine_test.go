package event

import (
	"fmt"
	"reflect"
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

type MyStruct struct {
}

func (s MyStruct) MyMethod(a int, b string) {
	// 方法实现
}

func Benchmark_reflectCall(b *testing.B) {
	for i := 0; i < b.N; i++ {
		reflectCall()
	}
}

func Benchmark_methodCall(b *testing.B) {
	for i := 0; i < b.N; i++ {
		methodCall()
	}
}
func methodCall() {
	//handler := new(testHandle)
	//input, _ := handler.Input(NewContext(nil), "name")
	//_, _ = handler.Output(NewContext(nil), input)
}

func reflectCall() {
	handler := new(testHandle)
	value := reflect.ValueOf(handler)
	methodInput := value.MethodByName("Input")
	inputResult := methodInput.Call([]reflect.Value{reflect.ValueOf(NewContext(nil)), reflect.ValueOf("name")})

	methodOutput := reflect.ValueOf(new(testHandle)).MethodByName("Output")
	_ = methodOutput.Call([]reflect.Value{reflect.ValueOf(NewContext(nil)), inputResult[0]})
}

func TestPipeLine_createHandlerWrap(t *testing.T) {
	wrap, err := createHandlerWrap(new(testHandle))
	fmt.Println(wrap, err)

	//typeValue := reflect.ValueOf(new(testHandle))
	////warp := new(handlerWrap)
	//
	//activateMethod := typeValue.MethodByName("Input")
	//if activateMethod.Kind() == reflect.Func && activateMethod.Type().NumIn() == 2 {
	//	in1 := activateMethod.Type().In(0)
	//	fmt.Println(in1)
	//	//warp.activateMethod = &activateMethod
	//}

	//println(outputResult)

	//var handler InputHandle := new(testHandle)
	//handler := new(testHandle)
	//wrap := createHandlerWrap(handler)
	//fmt.Println(wrap)

	// 获取类型信息
	//typeValue := reflect.TypeOf(MyStruct{})
	//typeValue := reflect.TypeOf(handler)
	//inputType := reflect.TypeOf((*InputHandle[string, any])(nil))
	//fmt.Println(typeValue, inputType)

	// 获取方法信息
	//method, ok := typeValue.MethodByName("MyMethod")
	//method, ok := typeValue.MethodByName("Input")
	//fmt.Println(method, ok)

	// 获取方法参数类型
	//paramTypes, ok := method.Type.In(0).Elem().FieldByName("a").Type() // 获取第一个参数类型  a
	//paramTypes, ok := method.Type.In(0).Elem().FieldByName("a") // 获取第一个参数类型  a
	//fmt.Println(paramTypes, ok)                                 // 输出: int
}

//func TestPipeLine_inbound(t *testing.T) {
//	handles := collect.NewLinkedNode[any]()
//	handler := &testHandle{state: true, result: "data"}
//	handles.Add(handler)
//	handles.Add(handler)
//
//	session := newSession()
//	context := session.OutputContext()
//
//	pipeLine := NewPipeLine(handles)
//	result, _ := pipeLine.inbound(context)
//
//	if result != handler.result {
//		t.Errorf("pipeLine.inbound() error, %s", result)
//	}
//	if handler.n != 2 {
//		t.Errorf("pipeLine.inbound() error, %d", handler.n)
//	}
//}
//
//func TestPipeLine_outbound(t *testing.T) {
//	handles := collect.NewLinkedNode[any]()
//	handler := &testHandle{state: true}
//	handles.Add(handler)
//	handles.Add(handler)
//
//	session := newSession()
//	session.state = Active
//	context := session.InputContext()
//	context.exchange = "data"
//
//	pipeLine := NewPipeLine(handles)
//	result, _ := pipeLine.outbound(context)
//
//	if result != "data" {
//		t.Errorf("pipeLine.outbound() error, %s", result)
//	}
//
//	if handler.n != 2 {
//		t.Errorf("pipeLine.outbound() error, %d", handler.n)
//	}
//}
