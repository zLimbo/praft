package ycsb

import (
	"fmt"
	"math/rand"

	"github.com/google/uuid"
)

var TmpV string

func Ycsb(writeRate float64, txNum int) {

	for i := 0; i < txNum; i++ {
		r := rand.Float64()
		// fmt.Printf("== r:%f\n", r)
		if len(keys) == 0 || r < writeRate {
			k := uuid.NewString()
			v := fmt.Sprintf("%0512s", k)
			Write(k, v)
			keys[k] = struct{}{}
			// fmt.Printf("write\t(%s, %s)\n", k, v)
		} else {
			for k := range keys {
				TmpV = Read(k)
				// fmt.Printf("read\t(%s, %s)\n", k, TmpV)
				break
			}
		}
	}
}
