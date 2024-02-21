package db

import (
	"go-leveldb/internal"
	"go-leveldb/memtable"
	"sync/atomic"
)

type DB struct {
	seq uint64
	mem *memtable.MemTable
}

func Open(dbName string) *DB {
	return &DB{mem: memtable.NewMemTable()}
}

func (D *DB) Put(key, value []byte) error {
	seq := atomic.AddUint64(&D.seq, 1)
	D.mem.Add(seq, internal.TypeValue, key, value)
	return nil
}

func (D *DB) Get(key []byte) ([]byte, error) {
	value, err := D.mem.Get(key)
	if err != nil {

	}

	return value, nil
}

func (D *DB) Delete(key []byte) error {
	seq := atomic.AddUint64(&D.seq, 1)
	D.mem.Add(seq, internal.TypeDeletion, key, nil)
	return nil
}

func (D *DB) Close() {

}
