package clic

import (
	"github.com/gorilla/websocket"
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
		seqId:seqId,
	}
}
// 返回用户连接
func(c *ClientConn) getConnection()  websocket.Conn{
	return *c.Conn
}
// 返回用户当前排名
func(c *ClientConn) getCurRank(curSeq uint64)  uint64{
	return c.seqId - curSeq
}
