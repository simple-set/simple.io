package collect

type NodeData[T any] struct {
	nextNode *NodeData[T]
	prevNode *NodeData[T]
	value    T
}

func (n *NodeData[T]) Value() T {
	return n.value
}
