package handler

import (
	"log"
	"context"
	"google.golang.org/grpc"
	"net"
	pb "waitqueue/proto/token"
)

const (
	port = ":8012"
)
// server is used to implement helloworld.GreeterServer.
type tkServer struct{}

// token放入号牌队列中

func (s *tkServer) SendTokenInfo(ctx context.Context, in *pb.TokenRequest) (response *pb.TokenResponse, error) {
	log.Printf("Received: %v", in.Token)
	popUserId:=POPQ()
	return &pb.TokenResponse{RetCode:1, UserId:popUserId, BindToken:in.Token}, nil
}
func RunTokenServer()  {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterTokenServiceServer(s, &tkServer{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
