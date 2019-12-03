package clic

import (
	"github.com/gorilla/websocket"
)
type ClientConn struct {
	UserId uint64 // 用户id
	Conn *websocket.Conn // 个人专属连接
}

func NewClient(userId uint64, con *websocket.Conn)  *ClientConn{
	return &ClientConn{
		UserId:userId,
		Conn: con,
	}
}
// 返回用户连接
func(c *ClientConn) getConnection()  websocket.Conn{
	return *c.Conn
}
