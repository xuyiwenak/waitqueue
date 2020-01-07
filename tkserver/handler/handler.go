package handler

import (
	"math/rand"
	"time"
	"waitqueue/proto/token"
	"waitqueue/tkserver/access"
	"log"
	"context"
	"google.golang.org/grpc"
	"waitqueue/utils/queue"
)
const (
	// token 推送的间隔
	tickerInterval time.Duration=5
)


func StartTicker(conn *grpc.ClientConn, q *queue.Queue)  {
	log.Println("StartTicker")

	ticker := time.NewTicker(time.Second*tickerInterval)
	defer ticker.Stop()

	for {
		select {
		case t := <-ticker.C:
			log.Println(t)
			// 如果队列里面没有token，则生成一部分
			if q.QLEN()<=0{
				tokenList := RandomToken(t)
				for i:=0; i<len(tokenList); i++ {
					q.QPUSH(tokenList[i])
					log.Printf("new token:%s\n", tokenList[i])
				}
			}
			log.Printf("cur q:%v, len(q):%d", q, q.QLEN())
			var resTokenList [] string
			// 先一次发一个
			pEle:=q.QPOP()
			if pEle==nil{
				continue
			}
			pbTokenList := &token.TokenRequest{
				Token:append(resTokenList,pEle.(string)),
			}
			c := token.NewTokenServiceClient(conn)

			r, err := c.SendTokenInfo(context.Background(), pbTokenList)
			if err != nil {
				log.Fatalf("could not send token: %v", err)
			}
			switch r.RetCode {
			case -1001:
				// 如果等待队列已经没有用户了则延迟一段时间再推送
			default:

			}

			log.Printf("bind RetCode:%d BindList:%v", r.RetCode, r.BindList)
		}
	}
}

func RandomToken(t time.Time) []string{
	rand.Seed(t.UnixNano())
	x := rand.Intn(2)   //生成0-10随机整数
	retToken := make([]string, x)
	for i:=0;i<x ;i++ {
		if curToken, err :=access.MakeAccessToken(int64(i)+t.UnixNano());err!=nil{
			log.Fatalf("生成token异常", err)
			continue
		} else{
			retToken=append(retToken, curToken)
		}
	}
	return retToken
}
