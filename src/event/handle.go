package event

/*
	流水线处理器, 用于加工从socket收到的数据, 多个处理器是链式执行
*/

// ActivateHandle 连接建立时执行
type ActivateHandle interface {
	Activate(context *HandleContext) (interface{}, bool)
}

// DisconnectHandle 连接断开时执行
type DisconnectHandle interface {
	Disconnect(context *HandleContext) (interface{}, bool)
}

// InputHandle 数据入站, 接收到数据时执行
type InputHandle interface {
	Input(context *HandleContext, data interface{}) (interface{}, bool)
}

// OutputHandle 数据出站, 发送数据时执行
type OutputHandle interface {
	Output(context *HandleContext, data interface{}) (interface{}, bool)
}

// IsHandler 检查是否为处理器类型
func IsHandler(handles ...any) bool {
	for _, handle := range handles {
		switch handle.(type) {
		case OutputHandle:
			continue
		case InputHandle:
			continue
		case ActivateHandle:
			continue
		case DisconnectHandle:
			continue
		default:
			return false
		}
	}
	return true
}
