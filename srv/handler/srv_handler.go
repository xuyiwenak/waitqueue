package handler

import (
	"flag"
	pb "github.com/gogo/protobuf/proto"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"reflect"
	"strconv"
	"sync"
	"sync/atomic"
	"waitqueue/proto/login"
	"waitqueue/srv/conn"
	"waitqueue/utils/msg"
	"waitqueue/utils/queue"
)

const loginPort string  = "localhost:8083"
var addr = flag.String("addr", loginPort, "http service address")

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
} // use default options

var (
	clientReq login.Request
	serverRes login.Response
)
var (
	// 排队的队列
	waitQueue *queue.Queue
	// 座位号队列
	seqIdQueue *queue.Queue
	// 注册正在排队的用户
	registMap sync.Map
	// 当前正在处理的序号
	curNum int64
)
// 基础数据的初始化
func Init()  {
	log.Println("Init waitQueue...")
	// 当前处理的初值
	atomic.StoreInt64(&curNum, 0)
	// 队列初始化
	waitQueue = queue.NewQueue(10000)
	seqIdQueue = queue.NewQueue(10000)
}
//
func InsertSeqId(num int)  {
	log.Printf("insert %d Ids", num)
	for i:=0; i<num;i++  {
		seqIdQueue.QPUSH(uint64(i))
	}
}

func QueryUserExist(userId uint64) bool {
	_, exists := registMap.Load(userId)
	return exists
}

func GetCurProcessNum()  int64{
	return atomic.LoadInt64(&curNum)
}
func AddCurProcessNum(delta int64)  {
	atomic.AddInt64(&curNum, delta)
}
func GetRankNum(userId uint64)  (int64, bool){
	if cc, ok := registMap.Load(userId);ok == true{
		value:=cc.(*clic.ClientConn)
		log.Printf("userId:%d seqId:%d hasProcess:%d", userId, value.SeqId, GetCurProcessNum())
		return int64(value.SeqId) - GetCurProcessNum()+1, ok
	}else {
		return 0, ok
	}
}

func WriteRecord(userId uint64, cc *clic.ClientConn)  {
	registMap.Store(userId, cc)
	waitQueue.QPUSH(userId)
}

func RemoveRecord(userId uint64)  {
	conn, exists := registMap.Load(userId)
	if exists{
		value:=conn.(*clic.ClientConn)
		var cancel login.Response
		cancel.MsgId=msg.CANCEL
		cancel.UserId=userId
		if err:=value.SendPBMsg(userId, &cancel);err!=nil{
			log.Println(err)
		}
		if err:=value.Conn.Close();err!=nil{
			log.Printf("close client err: err=%s",err)
		} else {
			log.Printf("userId:%d close successfully!",userId)
		}
	}
	registMap.Delete(userId)
}
// 从channel读数据
func POPQ() uint64{
	for {
		r := waitQueue.QPOP()
		AddCurProcessNum(1)
		field := reflect.ValueOf(r)
		switch field.Kind() {
		case reflect.Uint64:
			if QueryUserExist(field.Interface().(uint64)){
				log.Printf("read value: %d\n", r)
				return field.Interface().(uint64)
			}
		default:
			log.Printf("type:%s not in list",field.Kind())
			return 0
		}
		log.Printf("userId: %d not find\n", r)
	}
}

func Login(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	RevUserId := r.URL.Query().Get("userId")
	var userId uint64
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	if userIdInt,err:=strconv.Atoi(RevUserId);err!=nil{
		log.Println(err)
		return
	}else {
		userId =uint64(userIdInt)
	}
	log.Printf("user[%d] start login...", userId)
	if QueryUserExist(userId){
		log.Printf("QueryUserExist userId=%d", userId)
		return
	}
	curSeqId :=seqIdQueue.QPOP()
	clientConn := clic.NewClient(userId, conn, curSeqId.(uint64))
    WriteRecord(userId, clientConn)
	defer conn.Close()
	for {
		mt, buffer, err := clientConn.Conn.ReadMessage()
		if err != nil {
			log.Printf("read conntion buffer:%s:", err)
			RemoveRecord(clientConn.UserId)
			AddCurProcessNum(1)
			break
		}
		if err := pb.Unmarshal(buffer, &clientReq); err != nil {
			log.Printf("proto unmarshal: %s", err)
		}
		log.Printf("recv from client=%s, msgtype=%d", clientReq.String(), mt)
		revMsgId:=clientReq.MsgId
		userId:=clientReq.UserId
		switch revMsgId {
		case msg.QUERY:
			serverRes.UserId = userId
			serverRes.MsgId = revMsgId
			serverRes.RankNum, _ = GetRankNum(userId)
			pbBuffer, _ := pb.Marshal(&serverRes)
			err = clientConn.Conn.WriteMessage(mt, pbBuffer)
			if err != nil {
				log.Println("write:", err)
				break
			}
		default:
			log.Println("not this msgId!")
			break
		}
	}
}

func RunLoginServer()  {
	flag.Parse()
	log.SetFlags(0)
	Init()
	InsertSeqId(10000)
	http.HandleFunc("/Login", Login)
	log.Printf("login server start on port:%s ...", loginPort)
	if err:=http.ListenAndServe(*addr, nil);err!=nil{
		log.Fatal(err)
	}
}