package memtable

import (
	"fmt"
	"go-leveldb/internal"
	"testing"
)

func Test_MemTable(t *testing.T) {
	memTable := NewMemTable()
	memTable.Add(1234567, internal.TypeValue, []byte("aadsa34a"), []byte("bb23b3423"))
	value, _ := memTable.Get([]byte("aadsa34a"))
	fmt.Println(string(value))
	fmt.Println(memTable.ApproximateMemoryUsage())
}
