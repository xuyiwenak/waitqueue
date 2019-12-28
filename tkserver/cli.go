package main

import (
	"google.golang.org/grpc"
	"math/rand"
	"time"
	"waitqueue/proto/token"
	"log"
	"context"
	"waitqueue/tkserver/access"
)


const (
	address     = "localhost:8012"
)

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	startTicker(conn)

	log.Printf("bind userId:%d token:%s", r.UserId, r.BindToken)
}

