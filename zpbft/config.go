package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"praft/consistent"
	"praft/zlog"
	"strconv"
	"time"
)

var strByte = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")
var strByteLen = len(strByte)

const (
	MBSize      = 1024 * 1024
	ChanSize    = 10000
	KConfigFile = "./config/config.json"
	KCertsDir   = "./certs"
	KPeersFile  = "./config/peers.json"
	// KLocalIpFile = "./config/local_ip.txt"
)

type Config struct {
	PeerIps          []string      `json:"PeerIps"`
	IpNum            int           `json:"IpNum"`
	PortBase         int           `json:"PortBase"`
	ProcessNum       int           `json:"ProcessNum"`
	ReqNum           int           `json:"ReqNum"`
	BoostNum         int           `json:"BoostNum"`
	StartDelay       int           `json:"StartDelay"`
	RecvBufSize      int           `json:"RecvBufSize"`
	LogStdout        bool          `json:"LogStdout"`
	LogLevel         zlog.LogLevel `json:"LogLevel"`
	GoMaxProcs       int           `json:"GoMaxProcs"`
	BatchTxNum       int           `json:"BatchTxNum"`
	TxSize           int           `json:"TxSize"`
	GossipNum        int           `json:"GossipNum"`
	EnableGossip     bool          `json:"EnableGossip"`
	ExecNum          int           `json:"ExecNum"`
	ProposerNum      int           `json:"ProposerNum"`
	Load             int           `json:"Load"`
	Delay            int           `json:"Delay"`
	Delays           []int         `json:"Delays"`
	RotateOrNot      bool          `json:"RotateOrNot"`
	RandomDelayOrNot bool          `json:"RandomDelayOrNot"`
	ProcessNumArray  []int         `json:"ProcessNumArray"`
	DuplicateMode    int           `json:"DuplicateMode"`

	// for ycsb
	OpsPerTx   int     `json:"OpsPerTx"`
	ReadRate   float64 `json:"ReadRate"`
	WriteRate  float64 `json:"WriteRate"`
	UpdateRate float64 `json:"UpdateRate"`
	KeySize    int     `json:"KeySize"`
	ValueSize  int     `json:"ValueSize"`
	Rates      [NumOfOpType]float64

	Id2Node     map[int64]*Node
	ClientNode  *Node
	PeerIds     []int64
	LocalIp     string
	FaultNum    int
	RouteMap    map[int64][]int64
	ProposerIds []int64
	IsProposer  bool
}

var KConfig Config

func InitConfig(processId int) {

	// ?????? json
	context, err := ioutil.ReadFile(KConfigFile)
	if err != nil {
		zlog.Error("read %s failed.", KConfigFile)
	}
	err = json.Unmarshal(context, &KConfig)
	if err != nil {
		zlog.Error("json.Unmarshal() err: %v", err)
	}

	zlog.Debug("config file: ", string(context))
	zlog.Info("config: ", KConfig)

	// ????????????ip, port, ?????????
	KConfig.PeerIps = KConfig.PeerIps[:KConfig.IpNum]
	KConfig.Id2Node = make(map[int64]*Node)
	for i, ip := range KConfig.PeerIps {
		for j := 0; j < KConfig.ProcessNumArray[i]; j++ {
			port := KConfig.PortBase + 1 + j
			nodeId := GetId(ip, port)
			priKey, pubKey := ReadKeyPairDefault() // ????????????????????????????????????
			KConfig.Id2Node[nodeId] = NewNode(ip, port, priKey, pubKey)
			KConfig.PeerIds = append(KConfig.PeerIds, nodeId)
		}
	}
	KConfig.LocalIp = GetLocalIp()
	//??????Proposer???peerIps???KConfig.ProposerNum??????Proposer
	KConfig.ProposerIds = make([]int64, KConfig.ProposerNum)
	KConfig.IsProposer = false
	for i := 0; i < KConfig.ProposerNum; i++ {
		KConfig.ProposerIds[i] = KConfig.PeerIds[i]
		//KConfig.ProposerIds[i] = GetId(KConfig.PeerIps[i % len(KConfig.PeerIps)], KConfig.PortBase+ i / len(KConfig.PeerIps)+1)
		if KConfig.ProposerIds[i] == GetId(KConfig.LocalIp, KConfig.PortBase+processId) {
			KConfig.IsProposer = true
		}
		zlog.Debug("proposer id = %d", KConfig.ProposerIds[i])
		zlog.Debug("local mode is proposer : %v", KConfig.IsProposer)
	}
	zlog.Debug("Duplicate mode = %d", KConfig.DuplicateMode)

	// ???????????????
	KConfig.FaultNum = (len(KConfig.Id2Node) - 1) / 3

	// ?????? Rates[]
	rateSum := KConfig.ReadRate + KConfig.WriteRate + KConfig.UpdateRate
	KConfig.Rates[OpRead] = KConfig.ReadRate / rateSum
	KConfig.Rates[OpWrite] = KConfig.WriteRate/rateSum + KConfig.Rates[OpRead]
	KConfig.Rates[OpUpdate] = KConfig.UpdateRate/rateSum + KConfig.Rates[OpWrite]
}

func GetNode(id int64) *Node {
	node, ok := KConfig.Id2Node[id]
	if !ok {
		zlog.Error("The node of this ID(%d) does not exist!", id)
	}
	return node
}

func GetIndex(nodeId int64) int {
	for idx, id := range KConfig.PeerIds {
		if nodeId == id {
			return idx
		}
	}
	return -1
}

func RandString(length int) []byte {

	bytes := make([]byte, length)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < length; i++ {
		bytes[i] = strByte[r.Intn(strByteLen)]
	}

	return bytes
}

func ExampleNew(virtualNodeNum int, itemNum int) {
	c := consistent.New()
	for i := 0; i < virtualNodeNum*7; i++ {
		//fmt.Print("cache"+strconv.FormatInt(int64(i),10))
		c.Add("cache" + strconv.FormatInt(int64(i), 10))
	}

	var cacheANum int
	var cacheBNum int
	var cacheCNum int
	var cacheDNum int
	var cacheENum int
	var cacheFNum int
	var cacheGNum int
	//var cacheCNum int

	var randomString string
	for i := 0; i < itemNum; i++ {
		randomString = string(RandString(60))
		//fmt.Printf("rand sring %d: %s\n", i, randomString)
		server, err := c.Get(randomString)
		if err != nil {
			log.Fatal(err)
		}

		num, _ := strconv.Atoi(server[5:])
		//zlog.Debug("%s", server)
		//zlog.Debug("%d", num)
		if num%7 == 1 {
			cacheANum++
		} else if num%7 == 2 {
			cacheBNum++
		} else if num%7 == 3 {
			cacheCNum++
		} else if num%7 == 4 {
			cacheDNum++
		} else if num%7 == 5 {
			cacheENum++
		} else if num%7 == 6 {
			cacheFNum++
		} else if num%7 == 0 {
			cacheGNum++
		}

		//if server == "cacheA" || server == "cacheH" || server == "cacheO" || server == "cacheV" || server == "cacheAC" || server == "cacheAJ" || server == "cacheAQ"|| server == "cacheAX"|| server == "cacheBE"|| server == "cacheBL"|| server == "cacheBS"|| server == "cacheCA"|| server == "cacheCH"|| server == "cacheCO"|| server == "cacheCV"|| server == "cacheDC"|| server == "cacheDJ"|| server == "cacheDQ"|| server == "cacheEA"|| server == "cacheEH"|| server == "cacheEO"|| server == "cacheFA"|| server == "cacheFH"|| server == "cacheFO"|| server == "cacheGA"|| server == "cacheGH"|| server == "cacheGO"|| server == "cacheHA"|| server == "cacheHH"|| server == "cacheHO"|| server == "cacheIA"|| server == "cacheIH"|| server == "cacheIO"|| server == "cacheJA"|| server == "cacheJH"|| server == "cacheJO"|| server == "cacheKA"|| server == "cacheKH"|| server == "cacheKO"|| server == "cacheLA"|| server == "cacheLH"|| server == "cacheLO"|| server == "cacheMA"|| server == "cacheMH"|| server == "cacheMO"|| server == "cacheNA"|| server == "cacheNH"|| server == "cacheNO"|| server == "cacheOA"|| server == "cacheOH"{
		//	cacheANum++
		//}
		//if server == "cacheB" || server == "cacheI" || server == "cacheP" || server == "cacheW" || server == "cacheAD" || server == "cacheAK" || server == "cacheAR"|| server == "cacheAY"|| server == "cacheBF"|| server == "cacheBM"|| server == "cacheBT"|| server == "cacheCB"|| server == "cacheCI"|| server == "cacheCP"|| server == "cacheCW"|| server == "cacheDD"|| server == "cacheDK"|| server == "cacheDR"|| server == "cacheEB"|| server == "cacheEI"|| server == "cacheEP"|| server == "cacheFB"|| server == "cacheFI"|| server == "cacheFP"|| server == "cacheGB"|| server == "cacheGI"|| server == "cacheGP"|| server == "cacheHB"|| server == "cacheHI"|| server == "cacheHP"|| server == "cacheIB"|| server == "cacheII"|| server == "cacheIP"|| server == "cacheJB"|| server == "cacheJI"|| server == "cacheJP"|| server == "cacheKB"|| server == "cacheKI"|| server == "cacheKP"|| server == "cacheLB"|| server == "cacheLI"|| server == "cacheLP"|| server == "cacheMB"|| server == "cacheMI"|| server == "cacheMP"|| server == "cacheNB"|| server == "cacheNI"|| server == "cacheNP"|| server == "cacheOB"|| server == "cacheOI"{
		//	cacheBNum++
		//}
		//if server == "cacheC" || server == "cacheJ" || server == "cacheQ" || server == "cacheX" || server == "cacheAE" || server == "cacheAL" || server == "cacheAS"|| server == "cacheAZ"|| server == "cacheBG"|| server == "cacheBN"|| server == "cacheBU"|| server == "cacheCC"|| server == "cacheCJ"|| server == "cacheCQ"|| server == "cacheCX"|| server == "cacheDE"|| server == "cacheDL"|| server == "cacheDS"|| server == "cacheEC"|| server == "cacheEJ"|| server == "cacheEQ"|| server == "cacheFC"|| server == "cacheFJ"|| server == "cacheFQ"|| server == "cacheGC"|| server == "cacheGJ"|| server == "cacheGQ"|| server == "cacheHC"|| server == "cacheHJ"|| server == "cacheHQ"|| server == "cacheIC"|| server == "cacheIJ"|| server == "cacheIQ"|| server == "cacheJC"|| server == "cacheJJ"|| server == "cacheJQ"|| server == "cacheKC"|| server == "cacheKJ"|| server == "cacheKQ"|| server == "cacheLC"|| server == "cacheLJ"|| server == "cacheLQ"|| server == "cacheMC"|| server == "cacheMJ"|| server == "cacheMQ"|| server == "cacheNC"|| server == "cacheNJ"|| server == "cacheNQ"|| server == "cacheOC"|| server == "cacheOJ"{
		//	cacheCNum++
		//}
		//if server == "cacheD" || server == "cacheK" || server == "cacheR" || server == "cacheY" || server == "cacheAF" || server == "cacheAM" || server == "cacheAT"|| server == "cacheBA"|| server == "cacheBH"|| server == "cacheBO"|| server == "cacheBV"|| server == "cacheCD"|| server == "cacheCK"|| server == "cacheCR"|| server == "cacheCY"|| server == "cacheDF"|| server == "cacheDM"|| server == "cacheDT"|| server == "cacheED"|| server == "cacheEK"|| server == "cacheER"|| server == "cacheFD"|| server == "cacheFK"|| server == "cacheFR"|| server == "cacheGD"|| server == "cacheGK"|| server == "cacheGR"|| server == "cacheHD"|| server == "cacheHK"|| server == "cacheHR"|| server == "cacheID"|| server == "cacheIK"|| server == "cacheIR"|| server == "cacheJD"|| server == "cacheJK"|| server == "cacheJR"|| server == "cacheKD"|| server == "cacheKK"|| server == "cacheKR"|| server == "cacheLD"|| server == "cacheLK"|| server == "cacheLR"|| server == "cacheMD"|| server == "cacheMK"|| server == "cacheMR"|| server == "cacheND"|| server == "cacheNK"|| server == "cacheNR"|| server == "cacheOD"|| server == "cacheOK"{
		//	cacheDNum++
		//}
		//if server == "cacheE" || server == "cacheL" || server == "cacheS" || server == "cacheZ" || server == "cacheAG" || server == "cacheAN" || server == "cacheAU"|| server == "cacheBB"|| server == "cacheBI"|| server == "cacheBP"|| server == "cacheBW"|| server == "cacheCE"|| server == "cacheCL"|| server == "cacheCS"|| server == "cacheCZ"|| server == "cacheDG"|| server == "cacheDN"|| server == "cacheDU"|| server == "cacheEE"|| server == "cacheEL"|| server == "cacheES"|| server == "cacheFE"|| server == "cacheFL"|| server == "cacheFS"|| server == "cacheGE"|| server == "cacheGL"|| server == "cacheGS"|| server == "cacheHE"|| server == "cacheHL"|| server == "cacheHS"|| server == "cacheIE"|| server == "cacheIL"|| server == "cacheIS"|| server == "cacheJE"|| server == "cacheJL"|| server == "cacheJS"|| server == "cacheKE"|| server == "cacheKL"|| server == "cacheKS"|| server == "cacheLE"|| server == "cacheLL"|| server == "cacheLS"|| server == "cacheME"|| server == "cacheML"|| server == "cacheMS"|| server == "cacheNE"|| server == "cacheNL"|| server == "cacheNS"|| server == "cacheOE"|| server == "cacheOL"{
		//	cacheENum++
		//}
		//if server == "cacheF" || server == "cacheM" || server == "cacheT" || server == "cacheAA" || server == "cacheAH" || server == "cacheAO" || server == "cacheAV"|| server == "cacheBC"|| server == "cacheBJ"|| server == "cacheBQ"|| server == "cacheBX"|| server == "cacheCF"|| server == "cacheCM"|| server == "cacheCT"|| server == "cacheDA"|| server == "cacheDH"|| server == "cacheDO"|| server == "cacheDV"|| server == "cacheEF"|| server == "cacheEM"|| server == "cacheET"|| server == "cacheFF"|| server == "cacheFM"|| server == "cacheFT"|| server == "cacheGF"|| server == "cacheGM"|| server == "cacheGT"|| server == "cacheHF"|| server == "cacheHM"|| server == "cacheHT"|| server == "cacheIF"|| server == "cacheIM"|| server == "cacheIT"|| server == "cacheJF"|| server == "cacheJM"|| server == "cacheJT"|| server == "cacheKF"|| server == "cacheKM"|| server == "cacheKT"|| server == "cacheLF"|| server == "cacheLM"|| server == "cacheLT"|| server == "cacheMF"|| server == "cacheMM"|| server == "cacheMT"|| server == "cacheNF"|| server == "cacheNM"|| server == "cacheNT"|| server == "cacheOF"|| server == "cacheOM"{
		//	cacheFNum++
		//}
		//if server == "cacheG" || server == "cacheN" || server == "cacheU" || server == "cacheAB" || server == "cacheAI" || server == "cacheAP" || server == "cacheAW"|| server == "cacheBD"|| server == "cacheBK"|| server == "cacheBR"|| server == "cacheBY"|| server == "cacheCG"|| server == "cacheCN"|| server == "cacheCU"|| server == "cacheDB"|| server == "cacheDI"|| server == "cacheDP"|| server == "cacheDW"|| server == "cacheEG"|| server == "cacheEN"|| server == "cacheEU"|| server == "cacheFG"|| server == "cacheFN"|| server == "cacheFU"|| server == "cacheGG"|| server == "cacheGN"|| server == "cacheGU"|| server == "cacheHG"|| server == "cacheHN"|| server == "cacheHU"|| server == "cacheIG"|| server == "cacheIN"|| server == "cacheIU"|| server == "cacheJG"|| server == "cacheJN"|| server == "cacheJU"|| server == "cacheKG"|| server == "cacheKN"|| server == "cacheKU"|| server == "cacheLG"|| server == "cacheLN"|| server == "cacheLU"|| server == "cacheMG"|| server == "cacheMN"|| server == "cacheMU"|| server == "cacheNG"|| server == "cacheNN"|| server == "cacheNU"|| server == "cacheOG"|| server == "cacheON"{
		//	cacheGNum++
		//}
	}
	zlog.Debug("%d %d %d %d %d %d %d", cacheANum, cacheBNum, cacheCNum, cacheDNum, cacheENum, cacheFNum, cacheGNum)
	//users := []string{"user_mcnulty", "user_bunk", "user_omar", "user_bunny", "user_stringer","user_mcnulty1", "user_bunk2", "user_omar3", "user_bunny4", "user_stringer5"}
	//	fmt.Printf("%s => %s\n", u, server)

	// Output:
	// user_mcnulty => cacheA
	// user_bunk => cacheA
	// user_omar => cacheA
	// user_bunny => cacheC
	// user_stringer => cacheC
}

func ExampleAdd() {
	c := consistent.New()
	c.Add("cacheA")
	c.Add("cacheB")
	c.Add("cacheC")
	users := []string{"user_mcnulty", "user_bunk", "user_omar", "user_bunny", "user_stringer"}
	fmt.Println("initial state [A, B, C]")
	for _, u := range users {
		server, err := c.Get(u)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s => %s\n", u, server)
	}
	c.Add("cacheD")
	c.Add("cacheE")
	fmt.Println("\nwith cacheD, cacheE [A, B, C, D, E]")
	for _, u := range users {
		server, err := c.Get(u)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s => %s\n", u, server)
	}
	// Output:
	// initial state [A, B, C]
	// user_mcnulty => cacheA
	// user_bunk => cacheA
	// user_omar => cacheA
	// user_bunny => cacheC
	// user_stringer => cacheC
	//
	// with cacheD, cacheE [A, B, C, D, E]
	// user_mcnulty => cacheE
	// user_bunk => cacheA
	// user_omar => cacheA
	// user_bunny => cacheE
	// user_stringer => cacheE
}

func ExampleRemove() {
	c := consistent.New()
	c.Add("cacheA")
	c.Add("cacheB")
	c.Add("cacheC")
	users := []string{"user_mcnulty", "user_bunk", "user_omar", "user_bunny", "user_stringer"}
	fmt.Println("initial state [A, B, C]")
	for _, u := range users {
		server, err := c.Get(u)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s => %s\n", u, server)
	}
	c.Remove("cacheC")
	fmt.Println("\ncacheC removed [A, B]")
	for _, u := range users {
		server, err := c.Get(u)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s => %s\n", u, server)
	}
	// Output:
	// initial state [A, B, C]
	// user_mcnulty => cacheA
	// user_bunk => cacheA
	// user_omar => cacheA
	// user_bunny => cacheC
	// user_stringer => cacheC
	//
	// cacheC removed [A, B]
	// user_mcnulty => cacheA
	// user_bunk => cacheA
	// user_omar => cacheA
	// user_bunny => cacheB
	// user_stringer => cacheB
}
