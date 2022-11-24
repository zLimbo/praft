package main

import (
	"flag"
	"fmt"
	"praft/zlog"
	"runtime"
)

func main() {
	//var virtualNodeNum int
	//var itemNum int
	//flag.IntVar(&virtualNodeNum, "v", 1, "单个节点虚拟节点数量")
	//flag.IntVar(&itemNum, "n", 1, "交易总数目")
	//flag.Parse()
	//zlog.Debug("itemNum = %d", itemNum)
	//ExampleNew(virtualNodeNum,itemNum)
	var processIdx int
	var delayRange int64
	flag.Int64Var(&delayRange, "d", 0, "延迟范围, 默认为0")
	flag.IntVar(&processIdx, "i", 1, "进程号")
	flag.Parse()

	InitConfig(processIdx)
	zlog.Info("start server...")
	runtime.GOMAXPROCS(KConfig.GoMaxProcs)
	nodeId := GetId(KConfig.LocalIp, KConfig.PortBase+processIdx)
	dbDir := fmt.Sprintf("levdb/%d", nodeId)
	InitDB(dbDir)
	RunServer(nodeId, delayRange)
}
