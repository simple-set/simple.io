package collect

// LinkedList 双向链表, list风格API
type LinkedList[T any] struct {
	headNode *NodeData[T]
	tailNode *NodeData[T]
	size     int
}

func NewLinkedList[T any]() *LinkedList[T] {
	return &LinkedList[T]{}
}

func (l *LinkedList[T]) Size() int {
	return l.size
}

func (l *LinkedList[T]) Add(element T) {
	if l.headNode == nil {
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

func (l *LinkedList[T]) AddFirst(element T) {
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

func (l *LinkedList[T]) Get(index ...int) T {
	position := 0
	if len(index) > 0 {
		position = index[0]
	}

	var node T
	if position >= 0 && position < l.size {
		i := 0
		nodeValue := l.headNode
		for {
			if i == position {
				node = nodeValue.value
			}
			nodeValue = nodeValue.nextNode
		}
	}
	return node
}

func (l *LinkedList[T]) GetFirst() T {
	var node T
	if l.headNode != nil {
		node = l.headNode.value
	}
	return node
}

func (l *LinkedList[T]) GetLast() T {
	var node T
	if l.tailNode != nil {
		node = l.tailNode.value
	}
	return node
}

func (l *LinkedList[T]) Pop(index ...int) T {
	position := 0
	if len(index) > 0 {
		position = index[0]
	}

	var node T
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
				return nodeValue.value
			}
			i = -1
			node = nodeValue.nextNode.value
		}
	}
	return node
}

func (l *LinkedList[T]) PopFirst() T {
	return l.Pop(0)
}

func (l *LinkedList[T]) PopLast() T {
	return l.Pop(l.size - 1)
}

func (l *LinkedList[T]) Remove(index ...int) {
	l.Pop(index...)
}

func (l *LinkedList[T]) RemoveAll() {
	l.headNode = nil
	l.tailNode = nil
	l.size = 0
}
