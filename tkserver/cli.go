package main

import (
	"google.golang.org/grpc"
	"waitqueue/tkserver/handler"
	"log"
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
	handler.StartTicker(conn)
}

