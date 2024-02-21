package memtable

import (
	"go-leveldb/internal"
	"go-leveldb/skiplist"
)

type Iterator struct {
	listIter *skiplist.Iterator
}

// Valid Returns true iff the iterator is positioned at a valid node.
func (it *Iterator) Valid() bool {
	return it.listIter.Valid()
}

func (it *Iterator) InternalKey() *internal.InternalKey {
	return it.listIter.Key().(*internal.InternalKey)
}

// Next Advances to the next position.
// REQUIRES: Valid()
func (it *Iterator) Next() {
	it.listIter.Next()
}

// Prev Advances to the previous position.
// REQUIRES: Valid()
func (it *Iterator) Prev() {
	it.listIter.Prev()
}

// Seek Advance to the first entry with a key >= target
func (it *Iterator) Seek(target interface{}) {
	it.listIter.Seek(target)
}

// SeekToFirst Position at the first entry in list.
// Final state of iterator is Valid() iff list is not empty.
func (it *Iterator) SeekToFirst() {
	it.listIter.SeekToFirst()
}

// SeekToLast Position at the last entry in list.
// Final state of iterator is Valid() iff list is not empty.
func (it *Iterator) SeekToLast() {
	it.listIter.SeekToLast()
}
