package skiplist

// 跳表结点

type Node struct {
	key  any
	next []*Node
}

func newNode(key any, height int) *Node {
	x := &Node{key: key}
	x.next = make([]*Node, height)

	return x
}

func (node *Node) getNext(level int) *Node {
	return node.next[level]
}

func (node *Node) setNext(level int, x *Node) {
	node.next[level] = x
}
