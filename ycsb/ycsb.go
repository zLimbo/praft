package ycsb

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"math/rand"

	"github.com/google/uuid"
	"github.com/syndtr/goleveldb/leveldb"
)

var (
	db   *leveldb.DB
	kmap map[string]struct{}
)

func init() {
	db1, err := leveldb.OpenFile("levdb", nil)
	if err != nil {
		panic("open db failed!")
	}
	db = db1
	kmap = make(map[string]struct{})
}

func Read(key string) string {
	val, _ := db.Get([]byte(key), nil)
	return string(val)
}

func Write(key string, val string) {
	db.Put([]byte(key), []byte(val), nil)
}

type TxType int

const (
	R TxType = iota
	W
)

type Tx struct {
	Type TxType
	Key  string
	Val  string
}


const (
	Wrate = 0.5
)

func GenTxSet(wrate float64, num int) []byte {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	enc.Encode(num)
	for i := 0; i < num; i++ {
		r := rand.Float64()
		if len(kmap) == 0 || r < wrate {
			k := uuid.NewString()
			v := fmt.Sprintf("%0512s", k)
			enc.Encode(Tx{W, k, v})
			kmap[k] = struct{}{}
		} else {
			for k := range kmap {
				enc.Encode(Tx{R, k, ""})
				break
			}
		}
	}
	return buf.Bytes()
}

func ExecTxSet(txSet []byte) {
	var num int
	buf := bytes.NewBuffer(txSet)
	dec := gob.NewDecoder(buf)
	dec.Decode(&num)
	// fmt.Printf("n: %d\n", num)
	for i := 0; i < num; i++ {
		var tx Tx
		err := dec.Decode(&tx)
		if err != nil {
			panic("dec failed.")
		}
		// fmt.Println(tx.Type, tx.Key, tx.Val, "$")
		if tx.Type == R {
			Read(tx.Key)
		} else {
			Write(tx.Key, tx.Val)
		}
	}
}