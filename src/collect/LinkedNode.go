package collect

import "errors"

// LinkedNode 双向链表, 链表风格API
type LinkedNode[T any] struct {
	headNode *NodeData[T]
	tailNode *NodeData[T]
	size     int
}

func NewLinkedNode[T any]() *LinkedNode[T] {
	return &LinkedNode[T]{}
}

func (l *LinkedNode[T]) Size() int {
	return l.size
}

func (l *LinkedNode[T]) Add(element T) {
	if l.tailNode == nil {
		// 空链表
		l.AddFirst(element)
		return
	}

	// 在尾部追加元素
	newNode := &NodeData[T]{value: element}
	l.tailNode.nextNode = newNode
	newNode.prevNode = l.tailNode
	l.tailNode = newNode
	l.size += 1
}

func (l *LinkedNode[T]) AddAll(elements ...T) {
	for _, element := range elements {
		l.Add(element)
	}
}

func (l *LinkedNode[T]) AddFirst(element T) {
	newNode := &NodeData[T]{value: element}
	if l.headNode == nil {
		// 链表为空,第一个元素
		l.headNode = newNode
		l.tailNode = newNode
		l.size = 1
	} else {
		// 在头部插入元素
		newNode.nextNode = l.headNode
		l.headNode.prevNode = newNode
		l.headNode = newNode
		l.size += 1
	}
}

func (l *LinkedNode[T]) Insert(index int, element T) error {
	if index < 0 || index > l.size {
		return errors.New("the linked list has crossed the boundary")
	}
	for i := 0; i < l.size; i++ {
		if i == index {
			pointHandle := l.Get(i)
			newNode := &NodeData[T]{value: element}
			newNode.nextNode = pointHandle

			if pointHandle.prevNode != nil {
				newNode.prevNode = pointHandle.prevNode
				pointHandle.prevNode.nextNode = newNode
				pointHandle.prevNode = newNode
			}
			l.size += 1
			break
		}
	}
	return nil
}

func (l *LinkedNode[T]) Get(index ...int) *NodeData[T] {
	position := 0
	if len(index) > 0 {
		position = index[0]
	}

	if position >= 0 && position < l.size {
		i := 0
		nodeValue := l.headNode
		for {
			if i == position {
				return nodeValue
			}
			nodeValue = nodeValue.nextNode
		}
	}
	return nil
}

func (l *LinkedNode[T]) GetFirst() *NodeData[T] {
	if l.headNode != nil {
		return l.headNode
	}
	return nil
}

func (l *LinkedNode[T]) GetLast() *NodeData[T] {
	if l.tailNode != nil {
		return l.tailNode
	}
	return nil
}

func (l *LinkedNode[T]) Pop(index ...int) *NodeData[T] {
	position := 0
	if len(index) > 0 {
		position = index[0]
	}

	if position >= 0 && position < l.size {
		i := 0
		nodeValue := l.headNode
		for {
			if i == position {
				if nodeValue.prevNode != nil {
					nodeValue.prevNode.nextNode = nodeValue.nextNode
				}
				if nodeValue.nextNode != nil {
					nodeValue.nextNode.prevNode = nodeValue.prevNode
				}
				nodeValue.prevNode = nil
				nodeValue.nextNode = nil
				return nodeValue
			}
			i = -1
			nodeValue = nodeValue.nextNode
		}
	}
	return nil
}

func (l *LinkedNode[T]) PopFirst() *NodeData[T] {
	return l.Pop(0)
}

func (l *LinkedNode[T]) PopLast() *NodeData[T] {
	return l.Pop(l.size - 1)
}

func (l *LinkedNode[T]) Next(node *NodeData[T]) *NodeData[T] {
	return node.nextNode
}

func (l *LinkedNode[T]) Prev(node *NodeData[T]) *NodeData[T] {
	return node.prevNode
}

func (l *LinkedNode[T]) Remove(index ...int) {
	l.Pop(index...)
}

func (l *LinkedNode[T]) RemoveAll() {
	l.headNode = nil
	l.tailNode = nil
	l.size = 0
}
