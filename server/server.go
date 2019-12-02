package main

import (
	pb "github.com/gogo/protobuf/proto"
	"github.com/gorilla/websocket"
	"github.com/micro/go-micro/web"
	"log"
	"net/http"
	"sync"
	"time"
	"waitqueue/proto"
)

var upGrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}
var (
	clientReq login.Request
	serverRes login.Response
)
var (
	// 排队的队列长度
	waitQueue = make(chan uint64, 10000)
	// 无缓冲的channel负责挂起主线程
	msgChan  = make(chan int)
	// 查询是否存在过
	recordMap sync.Map
)

func QueryExist(userId uint64) bool {
	_, exists := recordMap.Load(userId)
	return exists
}
func WriteRecord(userId uint64)  {
	recordMap.Store(userId, time.Now())
}
func main() {
	// New web service

	service := web.NewService(
		web.Name("go.micro.web.login"),
	)

	if err := service.Init(); err != nil {
		log.Fatal("Init", err)
	}
	// websocket 连接接口 web.name注册根据.分割路由路径，所以注册的路径要和name对应上
	service.HandleFunc("/login", conn)

	if err := service.Run(); err != nil {
		log.Fatal("Run: ", err)
	}
}
func conn(w http.ResponseWriter, r *http.Request) {
	//
	conn, err := upGrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("upgrade: %s", err)
		return
	}

	defer conn.Close()

	go func() {
		for {
			_, buffer, err := conn.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				break
			}
			if err := pb.Unmarshal(buffer, &clientReq); err != nil {
				log.Printf("proto unmarshal: %s", err)
			}
			log.Printf("recv userId=%d MsgId=%d Data=%s", clientReq.UserId, clientReq.MsgId, clientReq.Data)
			revMsgId:=clientReq.MsgId
			switch revMsgId {
			case 1:
				if !QueryExist(clientReq.UserId){
					// 如果已经插入过就查询, 不写了
					WriteRecord(clientReq.UserId)
					
				}else {
					// 查询如果没有写过就写
					PUSHQ(clientReq.UserId)
				}
				break
			default:
				log.Println("not this msgId!")
			}

			serverRes.MsgId = revMsgId + 1
			serverRes.RankNum = uint64(len(waitQueue))
			pbBuffer, _ := pb.Marshal(&serverRes)
			if err:=conn.WriteMessage(websocket.BinaryMessage, pbBuffer);err != nil {
				log.Println("read:", err)
				break
			}
		}
	}()
}
// 从channel读数据
func POPQ(max int) {
	for {
		r := <-waitQueue
		log.Println("read value: %d\n", r)
	}

}
// 向channel写数据
func PUSHQ(i uint64) {
	waitQueue <- i
}


