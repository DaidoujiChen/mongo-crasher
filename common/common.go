package common

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"sync"
	"time"

	mgo "gopkg.in/mgo.v2"
)

// CoreRunLoop =w=
func CoreRunLoop(mongoDBClient *mgo.Session, wait *sync.WaitGroup, done chan bool, fetch func(mongoDBClient *mgo.Session, wait *sync.WaitGroup)) {
	totalRounds := 0
	var totalTimeSpend time.Duration
	for {

		// 每輪會在 2000 次到 4000 次之間
		perRound := rand.Intn(2000) + 2000
		totalRounds += perRound

		// 等待這次的每個結果都做完
		wait.Add(perRound)
		start := time.Now()
		for i := 0; i < perRound; i++ {
			go fetch(mongoDBClient, wait)
		}
		wait.Wait()

		// 總計時間
		totalTimeSpend += (time.Now().Sub(start))

		// mgo 開啟 connection 數量
		eachCount := []string{}
		totalSocketCount := 0
		servers := mongoDBClient.Cluster().Servers()
		for _, server := range servers.Slice() {
			eachCount = append(eachCount, strconv.Itoa(len(server.LiveSockets())))
			totalSocketCount += len(server.LiveSockets())
		}
		eachCountString := "(" + strings.Join(eachCount[:], "/") + ")"

		// 打印一下結果
		fmt.Printf("\r[ %fsec ] Rounds : %d, Cost : %f ns/op, MGO Connections : %d %s          ", totalTimeSpend.Seconds(), totalRounds, float32(totalTimeSpend)/float32(totalRounds), totalSocketCount, eachCountString)

		// 睡一下, 好像不睡用 freestyle 會爆
		time.Sleep(500 * time.Millisecond)
	}
}
