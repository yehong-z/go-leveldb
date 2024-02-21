package skiplist

type Iterator struct {
	list *SkipList
	node *Node
}

func (it *Iterator) Valid() bool {
	return it.node != nil
}

func (it *Iterator) Key() any {
	return it.node.key
}

func (it *Iterator) Next() {
	it.list.mu.RLock()
	defer it.list.mu.RUnlock()

	it.node = it.node.getNext(0)
}

// Prev Advances to the previous position.
// REQUIRES: Valid()
func (it *Iterator) Prev() {
	it.list.mu.RLock()
	defer it.list.mu.RUnlock()

	it.node = it.list.findLessThan(it.node.key)
	if it.node == it.list.head {
		it.node = nil
	}
}

// Seek Advance to the first entry with a key >= target
func (it *Iterator) Seek(target interface{}) {
	it.list.mu.RLock()
	defer it.list.mu.RUnlock()

	it.node, _ = it.list.findGreaterOrEqual(target)
}

// SeekToFirst Position at the first entry in list.
// Final state of iterator is Valid() iff list is not empty.
func (it *Iterator) SeekToFirst() {
	it.list.mu.RLock()
	defer it.list.mu.RUnlock()

	it.node = it.list.head.getNext(0)
}

// SeekToLast Position at the last entry in list.
// Final state of iterator is Valid() iff list is not empty.
func (it *Iterator) SeekToLast() {
	it.list.mu.RLock()
	defer it.list.mu.RUnlock()

	it.node = it.list.findLast()
	if it.node == it.list.head {
		it.node = nil
	}
}
