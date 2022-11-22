package ycsb

import "github.com/syndtr/goleveldb/leveldb"

var (
	db   *leveldb.DB
	keys map[string]struct{}
)

func init() {
	db1, err := leveldb.OpenFile("levdb", nil)
	if err != nil {
		panic("open db failed!")
	}
	db = db1
	keys = make(map[string]struct{})
}

func Read(key string) string {
	val, _ := db.Get([]byte(key), nil)
	return string(val)
}

func Write(key string, val string) {
	db.Put([]byte(key), []byte(val), nil)
}
