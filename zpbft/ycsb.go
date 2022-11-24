package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"math/rand"
	"praft/zlog"
	"strconv"

	"github.com/google/uuid"
	"github.com/syndtr/goleveldb/leveldb"
)

var (
	db   *leveldb.DB
	kmap map[string]struct{}
)

func InitDB(dir string) {
	var err error
	db, err = leveldb.OpenFile(dir, nil)
	if err != nil {
		zlog.Error("open db dir %s failed!", dir)
	}
	kmap = make(map[string]struct{})
}

func Read(key string) string {
	val, _ := db.Get([]byte(key), nil)
	return string(val)
}

func Write(key, val string) {
	db.Put([]byte(key), []byte(val), nil)
}

func Update(key, val string) {
	db.Put([]byte(key), []byte(val), nil)
}

type OpType int

const (
	OpRead OpType = iota
	OpWrite
	OpUpdate
	NumOfOpType
)

type Op struct {
	Type OpType
	Key  string
	Val  string
}

type Tx struct {
	Ops []Op
}

func GenTxSet() []byte {
	n := KConfig.BatchTxNum
	m := KConfig.OpsPerTx
	valFormat := "%0" + strconv.Itoa(KConfig.ValueSize) + "%s"
	rates := KConfig.Rates
	txs := make([]Tx, n)
	for i := range txs {
		ops := make([]Op, m)
		for j := range ops {
			r := rand.Float64()
			if len(kmap) != 0 && r < rates[OpRead] {
				ops[j].Type = OpRead
				for key := range kmap {
					ops[j].Key = key
					break
				}
			} else if r < rates[OpWrite] {
				ops[j].Type = OpWrite
				ops[j].Key = uuid.NewString()
				ops[j].Val = fmt.Sprintf(valFormat, uuid.NewString())
			} else {
				ops[j].Type = OpUpdate
				for key := range kmap {
					ops[j].Key = key
					break
				}
				ops[j].Val = fmt.Sprintf(valFormat, uuid.NewString())
			}
		}
		txs[i].Ops = ops
	}

	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	enc.Encode(txs)
	return buf.Bytes()
}

func ExecTxSet(txSet []byte) int {
	var txs []Tx
	buf := bytes.NewBuffer(txSet)
	dec := gob.NewDecoder(buf)
	dec.Decode(&txs)

	for _, tx := range txs {
		for _, op := range tx.Ops {
			if op.Type == OpRead {
				Read(op.Key)
			} else if op.Type == OpWrite {
				Write(op.Key, op.Val)
			} else if op.Type == OpUpdate {
				Update(op.Key, op.Val)
			}
		}
	}
	return len(txs)
}
