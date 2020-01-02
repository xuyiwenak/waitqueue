package main

import (
	"google.golang.org/grpc"
	"log"
	"waitqueue/tkserver/handler"
	"waitqueue/utils/queue"
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
	q:=queue.NewQueue(100)
	defer conn.Close()
	handler.StartTicker(conn, q)
}

