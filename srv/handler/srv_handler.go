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
	curNum uint64
)
// 基础数据的初始化
func Init()  {
	log.Println("Init waitQueue...")
	// 当前处理的初值
	atomic.StoreUint64(&curNum, 0)

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

func QueryExist(userId uint64) bool {
	_, exists := registMap.Load(userId)
	log.Printf("exists:%v", exists)
	return exists
}

func GetCurProcessNum()  uint64{
	return atomic.LoadUint64(&curNum)
}
func AddCurProcessNum()  {
	atomic.AddUint64(&curNum, 1)
}
func GetRankNum(userId uint64)  (uint64, bool){
	if cc, ok := registMap.Load(userId);ok == true{
		value:=cc.(*clic.ClientConn)
		return value.SeqId - GetCurProcessNum(), ok
	}else {
		return 0, ok
	}
}

func WriteRecord(userId uint64, cc *clic.ClientConn)  {
	registMap.Store(userId, cc)
	waitQueue.QPUSH(userId)
}

func RemoveRecord(userId uint64)  {
	registMap.Delete(userId)
}
// 从channel读数据
func POPQ() uint64{
	for {
		r := waitQueue.QPOP()
		AddCurProcessNum()
		field := reflect.ValueOf(r)
		switch field.Kind() {
		case reflect.Uint64:
			if QueryExist(field.Interface().(uint64)){
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
	log.Println("Login start...")
	conn, err := upgrader.Upgrade(w, r, nil)
	userId := r.URL.Query().Get("userId")
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	userIdInt,_:=strconv.Atoi(userId)
	curSeqId :=seqIdQueue.QPOP()
	clientConn := clic.NewClient(uint64(userIdInt), conn, curSeqId.(uint64))
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
				WriteRecord(userId, clientConn)
			}
			break
		default:
			log.Println("not this msgId!")
			break
		}
		serverRes.UserId = userId
		serverRes.MsgId = revMsgId
		serverRes.RankNum, _ = GetRankNum(userId)

		pbBuffer, _ := pb.Marshal(&serverRes)
		err = clientConn.Conn.WriteMessage(mt, pbBuffer)
		if err != nil {
			log.Println("write:", err)
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