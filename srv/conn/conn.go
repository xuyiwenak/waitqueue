package clic

import (
	pb "github.com/gogo/protobuf/proto"
	"github.com/gorilla/websocket"
	"log"
)
type ClientConn struct {
	UserId uint64 // 用户id
	Conn *websocket.Conn // 个人专属连接
	SeqId uint64 // 个人排号
}

func NewClient(userId uint64, con *websocket.Conn, seqId uint64)  *ClientConn{
	return &ClientConn{
		UserId:userId,
		Conn: con,
		SeqId:seqId,
	}
}
// 返回用户连接G
func(c *ClientConn) GetConnection()  websocket.Conn{
	return *c.Conn
}

func(c *ClientConn) SendPBMsg(userId uint64, buffer pb.Message)  error{
	pbBuffer, _ := pb.Marshal(buffer)
	err := c.Conn.WriteMessage(websocket.BinaryMessage, pbBuffer)
	if err != nil {
		log.Println("write:", err)
		return err
	}
	return nil
}
