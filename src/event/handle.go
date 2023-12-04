package event

type handlerInterface interface {
}

// ActivateHandle 连接建立时执行
type ActivateHandle[T any] interface {
	Activate(context *HandleContext) (T, bool)
}

// DisconnectHandle 连接断开时执行
type DisconnectHandle[T any] interface {
	Disconnect(context *HandleContext) (T, bool)
}

// InputHandle 数据入站, 接收到数据时执行
//type InputHandle interface {
//	Input(context *HandleContext, data interface{}) (interface{}, bool)
//}

type InputHandle[T, E any] interface {
	Input(context *HandleContext, data T) (E, bool)
}

// OutputHandle 数据出站, 发送数据时执行
type OutputHandle[T, E any] interface {
	Output(context *HandleContext, data T) (E, bool)
}

// IsHandler 检查是否为处理器类型
func IsHandler(handles ...any) bool {
	//for _, handle := range handles {
	//	switch handle.(type) {
	//	case OutputHandle:
	//		continue
	//	case InputHandle:
	//		continue
	//	case ActivateHandle:
	//		continue
	//	case DisconnectHandle:
	//		continue
	//	default:
	//		return false
	//	}
	//}
	return true
}
