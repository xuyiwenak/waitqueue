package handler

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"time"
	"waitqueue/proto"
	"sync"
	pb "github.com/gogo/protobuf/proto"
)
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
} // use default options

var (
	clientReq login.Request
	serverRes login.Response
)
var (
	// 排队的队列长度
	waitQueue = make(chan uint64, 10000)
	// 无缓冲的channel负责挂起主线程
	msgChan  = make(chan int)
	// 注册正在排队的用户
	registMap sync.Map
)
func QueryExist(userId uint64) bool {
	_, exists := registMap.Load(userId)
	return exists
}
func WriteRecord(userId uint64)  {
	registMap.Store(userId, time.Now())
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
func Login(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()
	for {
		mt, buffer, err := c.ReadMessage()
		log.Println("xxxxxx", mt, buffer, err)
		// 如果客户端主动关闭tcp
		if mt==websocket.CloseMessage{
			log.Println("client close!",buffer, buffer[2:], err)
		}
		if err != nil {
			log.Println("read:", err)
			break
		}
		if err := pb.Unmarshal(buffer, &clientReq); err != nil {
			log.Printf("proto unmarshal: %s", err)
		}
		log.Printf("recv from client=%s, msgtype=%d", clientReq.String(), mt)
		revMsgId:=clientReq.MsgId
		userId:=clientReq.UserId
		switch revMsgId {
		case 1:
			if QueryExist(userId){
				// 如果已经插入过就查询, 不写了
				log.Printf("QueryExist userId=%d", userId)
				WriteRecord(userId)

			}else {
				// 查询如果没有写过就写
				PUSHQ(userId)
				log.Printf("map:%v", registMap)
			}
			break
		default:
			log.Println("not this msgId!")
		}
		serverRes.UserId = userId
		serverRes.MsgId = revMsgId
		serverRes.RankNum = uint64(len(waitQueue))
		log.Printf("write to client: %s", serverRes.String())
		pbBuffer, _ := pb.Marshal(&serverRes)
		err = c.WriteMessage(mt, pbBuffer)
		if err != nil {
			log.Println("write:", err)
			break
		}
	}
}