package skiplist

import (
	"go-leveldb/utils"
	"math/rand"
	"sync"
)

// 跳表实现

const kMaxHeight = 12
const kBranching = 4

type SkipList struct {
	maxHeight  int
	head       *Node
	comparator utils.Comparator
	mu         sync.RWMutex
}

func New(cmp utils.Comparator) *SkipList {
	return &SkipList{
		head:       newNode(nil, kMaxHeight),
		maxHeight:  1,
		comparator: cmp,
	}
}

func (l *SkipList) keyIsAfterNode(key any, n *Node) bool {
	return (n != nil) && (l.comparator(n.key, key) < 0)
}

func (l *SkipList) findGreaterOrEqual(key any) (*Node, [kMaxHeight]*Node) {
	var prev [kMaxHeight]*Node
	x := l.head
	level := l.maxHeight - 1
	for {
		next := x.getNext(level)
		if l.keyIsAfterNode(key, next) {
			x = next
		} else {
			prev[level] = x
			if level == 0 {
				return next, prev
			} else {
				level--
			}
		}
	}
}

func (l *SkipList) randomHeight() int {
	height := 1
	for height < kMaxHeight && (rand.Intn(kBranching) == 0) {
		height++
	}

	return height
}

func (l *SkipList) Insert(key any) {
	l.mu.Lock()
	defer l.mu.Unlock()

	_, prev := l.findGreaterOrEqual(key)
	height := l.randomHeight()
	if height > l.maxHeight {
		for i := l.maxHeight; i < height; i++ {
			prev[i] = l.head
		}

		l.maxHeight = height
	}

	x := newNode(key, height)
	for i := 0; i < height; i++ {
		x.setNext(i, prev[i].getNext(i))
		prev[i].setNext(i, x)
	}
}

func (l *SkipList) Contains(key any) bool {
	l.mu.RLock()
	defer l.mu.RUnlock()
	x, _ := l.findGreaterOrEqual(key)
	if x != nil && l.comparator(x.key, key) == 0 {
		return true
	}

	return false
}

func (l *SkipList) NewIterator() *Iterator {
	var it Iterator
	it.list = l
	return &it
}

func (l *SkipList) findLessThan(key interface{}) *Node {
	x := l.head
	level := l.maxHeight - 1
	for {
		next := x.getNext(level)
		if next == nil || l.comparator(next.key, key) >= 0 {
			if level == 0 {
				return x
			} else {
				level--
			}
		} else {
			x = next
		}
	}
}

func (l *SkipList) findLast() *Node {
	x := l.head
	level := l.maxHeight - 1
	for {
		next := x.getNext(level)
		if next == nil {
			if level == 0 {
				return x
			} else {
				level--
			}
		} else {
			x = next
		}
	}
}
