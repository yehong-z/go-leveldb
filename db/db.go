package db

import (
	"errors"
	"go-leveldb/internal"
	"go-leveldb/memtable"
	"go-leveldb/version"
	"sync"
	"time"
)

type DB struct {
	name                  string
	mu                    sync.Mutex
	cond                  *sync.Cond
	mem                   *memtable.MemTable
	imm                   *memtable.MemTable
	current               *version.Version
	bgCompactionScheduled bool
}

func Open(dbName string) *DB {
	db := DB{
		name:                  dbName,
		mem:                   memtable.NewMemTable(),
		imm:                   nil,
		bgCompactionScheduled: false,
	}

	db.cond = sync.NewCond(&db.mu)
	num := db.ReadCurrentFile()
	if num > 0 {
		v, err := version.Load(dbName, num)
		if err != nil {
			return nil
		}
		db.current = v
	} else {
		db.current = version.New(dbName)
	}

	return &db
}

func (db *DB) Put(key, value []byte) error {
	seq, err := db.makeRoomForWrite()
	if err != nil {
		return err
	}

	// todo : add log

	db.mem.Add(seq, internal.TypeValue, key, value)
	return nil
}

func (db *DB) Get(key []byte) ([]byte, error) {
	db.mu.Lock()
	mem := db.mem
	imm := db.imm
	current := db.current
	db.mu.Unlock()
	value, err := mem.Get(key)
	if !errors.Is(err, internal.ErrNotFound) {
		return value, err
	}

	if imm != nil {
		value, err := imm.Get(key)
		if !errors.Is(err, internal.ErrNotFound) {
			return value, err
		}
	}

	value, err = current.Get(key)
	return value, err
}

func (db *DB) Delete(key []byte) error {
	seq, err := db.makeRoomForWrite()
	if err != nil {
		return err
	}

	db.mem.Add(seq, internal.TypeDeletion, key, nil)
	return nil
}

func (db *DB) Close() {
	db.mu.Lock()
	defer db.mu.Unlock()
	for db.bgCompactionScheduled {
		db.cond.Wait()
	}
}

func (db *DB) makeRoomForWrite() (uint64, error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	for {
		if db.current.NumLevelFiles(0) >= internal.L0_SlowdownWritesTrigger {
			db.mu.Unlock()
			time.Sleep(time.Duration(1000) * time.Microsecond)
			db.mu.Lock()
		} else if db.mem.ApproximateMemoryUsage() <= internal.Write_buffer_size {
			return db.current.NextSeq(), nil
		} else if db.imm != nil {
			//  Current memtable full; waiting
			db.cond.Wait()
		} else {
			// Attempt to switch to a new memtable and trigger compaction of old
			// todo : switch log
			db.imm = db.mem
			db.mem = memtable.NewMemTable()
			db.maybeScheduleCompaction()
		}
	}
}
