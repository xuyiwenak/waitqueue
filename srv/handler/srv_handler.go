package handler

import (
	pb "github.com/gogo/protobuf/proto"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"strconv"
	"sync"
	"waitqueue/proto"
	"waitqueue/srv/conn"
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
	// seqId生成队列
	seqQueue = make(chan uint64, 10000)
	// 无缓冲的channel负责挂起主线程
	msgChan  = make(chan int)
	// 注册正在排队的用户
	registMap sync.Map
	// 当前正在处理的序号
	curNum uint64
)

func InitSeqQueue(num uint64)  {
	for i:=num;i>0 ;i--{
		seqQueue<-i
	}
}
func GetOneSeqId()  uint64{
	return <-seqQueue
}
func QueryExist(userId uint64) bool {
	_, exists := registMap.Load(userId)
	log.Printf("exists:%v", exists)
	return exists
}
func WriteRecord(userId uint64, cc *clic.ClientConn, seqId uint64)  {
	registMap.Store(userId, cc)
	PUSHQ(userId)
}
func RemoveRecord(userId uint64)  {
	registMap.Delete(userId)
}
// 从channel读数据
func POPQ() uint64{
	for {
		r := <-waitQueue
		if QueryExist(r){
			log.Printf("read value: %d\n", r)
			return r
		}
		log.Printf("userId: %d not find\n", r)
	}
}
// 向channel写数据
func PUSHQ(i uint64) {
	waitQueue <- i
}
func Login(w http.ResponseWriter, r *http.Request) {

	conn, err := upgrader.Upgrade(w, r, nil)
	userId := r.URL.Query().Get("userId")
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	userIdInt,_:=strconv.Atoi(userId)

	clientConn := clic.NewClient(uint64(userIdInt), conn, GetOneSeqId())
	defer clientConn.Conn.Close()
	for {
		mt, buffer, err := clientConn.Conn.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			RemoveRecord(clientConn.UserId)
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
			}else {
				// 查询如果没有写过就写
				log.Println("empty map:%v", registMap)
				WriteRecord(userId, clientConn)
				PUSHQ(userId)
				log.Printf("insert map:%v", registMap)
			}
			break
		default:
			log.Println("not this msgId!")
			break
		}
		serverRes.UserId = userId
		serverRes.MsgId = revMsgId
		serverRes.RankNum = uint64(len(waitQueue))
		log.Printf("write to client: %s", serverRes.String())
		pbBuffer, _ := pb.Marshal(&serverRes)
		err = clientConn.Conn.WriteMessage(mt, pbBuffer)
		if err != nil {
			log.Println("write:", err)
			break
		}
	}
}