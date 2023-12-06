package event

type Direction int

const (
	// 入站
	inbound Direction = iota
	// 出站
	outbound
)

type HandleContext struct {
	// 流水线方向
	direction Direction
	// 会话对象
	session *Session
	// 交换缓冲区, handle之间传递数据
	exchange interface{}
}

func (h HandleContext) Direction() Direction {
	return h.direction
}

func (h HandleContext) Session() *Session {
	return h.session
}

func NewContext(session *Session) *HandleContext {
	return &HandleContext{session: session}
}
