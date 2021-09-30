package main

import (
	"fmt"
	"log"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gorilla/websocket"
)

const (
	wsCount       = 2000 // 创建的websocket数量
	createConnGap = 100  // 创建连接的间隔时间
	// 每个连接发送的信息数量 = msgNums + 1
	msgNums = 10
)

var (
	origin = "http://127.0.0.1:8080"
	url    = "ws://127.0.0.1:8080/ws"
	start  = make(chan struct{})
	// 计数
	done int32

	wg sync.WaitGroup
)

func Worker(id int) {
	ws, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Fatal("dial:", err)
	}

	// 等待开始
	<-start

	send := 0

	defer func() {
		ws.Close()
		log.Printf("worker %3d done, send:%3d \n", id, send)
		wg.Done()
	}()
	go readFromServer(ws)

	for {
		msg := []byte("{\"msg\":\"worker\" " + strconv.Itoa(id) + "\"send\": " + strconv.Itoa(send) + "\"}")
		err := ws.WriteMessage(1, msg)
		if err != nil {
			log.Println(err)
		}
		send++
		// 自定义数量
		if send > msgNums {
			// 结束全部任务
			atomic.AddInt32(&done, 1)
			return
		}

		time.Sleep(time.Second)
	}

}
func readFromServer(ws *websocket.Conn) {
	for {
		// read in a message
		_, p, err := ws.ReadMessage()
		if err != nil {
			// log.Println("err in ReadMessage:\n", err)
			return
		}
		// print out that message for clarity
		fmt.Println("msg from server: ", string(p))
	}
}

func main() {
	// 建立连接
	for i := range [wsCount][]int{} {
		time.Sleep(time.Millisecond * createConnGap)
		go Worker(i)
		wg.Add(1)
	}
	// 开始发送
	close(start)
	// 等待发送任务完成
	wg.Wait()
	// 打印完成结果
	log.Println("done:", done)
}
