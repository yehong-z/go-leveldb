package block

import "go-leveldb/internal"

type Iterator struct {
	block *Block
	index int
}

// Valid Returns true iff the iterator is positioned at a valid node.
func (it *Iterator) Valid() bool {
	return it.index >= 0 && it.index < len(it.block.items)
}

func (it *Iterator) InternalKey() *internal.InternalKey {
	return &it.block.items[it.index]
}

// Next Advances to the next position.
// REQUIRES: Valid()
func (it *Iterator) Next() {
	it.index++
}

// Prev Advances to the previous position.
// REQUIRES: Valid()
func (it *Iterator) Prev() {
	it.index--
}

// Seek Advance to the first entry with a key >= target
func (it *Iterator) Seek(target interface{}) {
	// 二分法查询
	left := 0
	right := len(it.block.items) - 1
	for left < right {
		mid := (left + right) / 2
		if internal.UserKeyComparator(it.block.items[mid].UserKey, target) < 0 {
			left = mid + 1
		} else {
			right = mid
		}
	}
	if left == len(it.block.items)-1 {
		if internal.UserKeyComparator(it.block.items[left].UserKey, target) < 0 {
			// not found
			left++
		}
	}
	it.index = left
}

// SeekToFirst Position at the first entry in list.
// Final state of iterator is Valid() iff list is not empty.
func (it *Iterator) SeekToFirst() {
	it.index = 0
}

// SeekToLast Position at the last entry in list.
// Final state of iterator is Valid() iff list is not empty.
func (it *Iterator) SeekToLast() {
	if len(it.block.items) > 0 {
		it.index = len(it.block.items) - 1
	}
}
