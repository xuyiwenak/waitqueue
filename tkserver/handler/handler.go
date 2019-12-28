package handler

import (
	"math/rand"
	"time"
	"waitqueue/proto/token"
	"waitqueue/tkserver/access"
)

func startTicker(conn *ClientConn)  {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case t := <-ticker.C:
			tokenList := RandomToken(t)
			// Contact the server and print out its response.
			pbTokenList := &token.TokenRequest{
				Token:tokenList,
			}
			c := token.NewTokenServiceClient(conn)

			r, err := c.SendTokenInfo(context.Background(), pbTokenList)
			if err != nil {
				log.Fatalf("could not send token: %v", err)
			}

		}
	}
}
func RandomToken(t time.Time) []string{
	rand.Seed(t.UnixNano())
	x := rand.Intn(10)   //生成0-10随机整数
	tmpStr := make([]string, x)
	for i:=0;i<x ;i++ {
		if curToken, err :=access.MakeAccessToken(int64(i)+t.UnixNano());err!=nil{
			log.Fatalf("生成token异常", err)
			continue
		} else{
			tmpStr=append(tmpStr, curToken)
		}
	}
	return tmpStr
}
