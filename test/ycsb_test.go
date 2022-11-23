package test

import (
	"fmt"
	"io/ioutil"
	"praft/ycsb"
	"testing"
	"time"
)

func TestYcsb(t *testing.T) {
	ycsb.InitDB("testdb")
	txSet := ycsb.GenTxSet(0.5, 10000)
	before := time.Now()
	ycsb.ExecTxSet(txSet)
	take := time.Since(before)
	ioutil.WriteFile("out", txSet, 0777)
	fmt.Println("take", take)
}
