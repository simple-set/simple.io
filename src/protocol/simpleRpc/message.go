package simpleRpc

import "strings"

type SerializerType byte

const (
	Json SerializerType = iota
)

type kind byte

const (
	Request kind = iota
	Response
)

// HeadLength 包头长度
func HeadLength() int {
	return 16
}

// Magic 魔术字符
func Magic() [2]byte {
	return [2]byte{0x23, 0x23}
}

// Version 协议版本
func Version() int {
	return 1
}

type Message struct {
	magic      [2]byte        // 协议魔数
	Version    int            // 协议版本号
	Length     int            // 消息正文长度
	Serializer SerializerType // 序列化算法
	Kind       kind           // 消息类型, 请求 or 响应
	method     [16]byte       // 方法
	Sequence   int            // 消息序号, 四个字节,请求和应答序号一致, 用于支持全双工通讯
	body       []byte         // 正文数据
}

// HeadLength 消息头长度
func (m *Message) HeadLength() int {
	return HeadLength()
}

// MessageLength 消息总长度
func (m *Message) MessageLength() int {
	return HeadLength() + len(m.body)
}

// SetBody 设置消息正文
func (m *Message) SetBody(data []byte) {
	if data != nil {
		m.body = data
		m.Length = len(m.body)
	}
}

// Body 消息正文
func (m *Message) Body() []byte {
	return m.body
}

// Method 消息方法
func (m *Message) Method() string {
	return strings.Trim(string(m.method[:]), string([]byte{0}))
}

// SetMethod 设置方法名
func (m *Message) SetMethod(name string) {
	_ = copy(m.method[:], []byte(name)[:16])
}

// NewEmpty 创建一个空的消息
func NewEmpty() *Message {
	return &Message{}
}

// NewRequest 请求消息
func NewRequest(data []byte) *Message {
	return NewMessage(data, Request)
}

// NewResponse 响应消息
func NewResponse(data []byte) *Message {
	return NewMessage(data, Response)
}

func NewMessage(data []byte, method kind) *Message {
	return &Message{
		magic:      Magic(),
		Version:    Version(),
		Length:     len(data),
		Serializer: Json,
		Kind:       method,
		Sequence:   1,
		body:       data,
	}
}
