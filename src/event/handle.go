package event

// ActivateHandle 连接建立时执行
type ActivateHandle interface {
	Activate(context *HandleContext) bool
}

// DisconnectHandle 连接断开时执行
type DisconnectHandle interface {
	Disconnect(context *HandleContext) bool
}

type InputHandle[T any] interface {
	Input(context *HandleContext, data T) (any, bool)
}

// OutputHandle 数据出站, 发送数据时执行
type OutputHandle[T any] interface {
	Output(context *HandleContext, data T) (any, bool)
}
