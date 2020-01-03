package main

import (
	"flag"
	"fmt"
	pb "github.com/gogo/protobuf/proto"
	"github.com/gorilla/websocket"
	"log"
	"math/rand"
	"net/url"
	"os"
	"os/signal"
	"time"
	"waitqueue/proto/login"
)

var addr = flag.String("addr", "localhost:8083", "http service address")
var(
	clientReq login.Request
	serverRes login.Response
	userId uint64
)
func main() {
	flag.Parse()
	log.SetFlags(0)
	rand.Seed(time.Now().UnixNano())
	userId = uint64(rand.Int63n(10000))
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	urlRawPath:=fmt.Sprintf("userId=%d", userId)
	u := url.URL{Scheme: "ws", Host: *addr, Path: "/Login", RawQuery:urlRawPath}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			_, buffer, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			if err := pb.Unmarshal(buffer, &serverRes); err != nil {
				log.Printf("proto unmarshal: %s", err)
			}
			log.Printf("recv from server :%s", serverRes.String())
		}
	}()

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			return
		case t := <-ticker.C:
			clientReq.MsgId = 1
			clientReq.UserId = userId
			clientReq.Data = t.String()
			pbBuffer, _ := pb.Marshal(&clientReq)
			log.Println(t)
			err = c.WriteMessage(websocket.BinaryMessage, pbBuffer)
			if err != nil {
				log.Println("write:", err)
				break
			}
		case <-interrupt:
			// 断开的时候指明是哪个userId
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, string(userId)))
			if err != nil {
				log.Println("write close:", err)
				return
			}
			log.Fatalln("interrupt")
		}
	}
}
