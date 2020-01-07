package handler

import (
	"context"
	"google.golang.org/grpc"
	"log"
	"net"
	pb "waitqueue/proto/token"
)

const (
	port = ":8012"
)

type tkServer struct{}

// token放入号牌队列中

func (s *tkServer) SendTokenInfo(ctx context.Context, in *pb.TokenRequest) (out *pb.TokenResponse, err error) {
	log.Printf("Received: %v", in.Token)
	var outList [] *pb.UserBindList
	for i:=0; i<len(in.Token);i++ {
		if waitQueue.QLEN()<=0{
			return &pb.TokenResponse{RetCode:-1001, BindList:outList}, nil
		}
		popUserId:=POPQ()
		tmp:=&pb.UserBindList{
			UserId:popUserId,
			BindToken:in.Token[i],
		}
		outList =append(outList, tmp)

		// 移除在排队状态中的userId
		RemoveRecord(popUserId)
	}

	return &pb.TokenResponse{RetCode:1, BindList:outList}, nil
}
func RunTokenServer()  {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Printf("token server start on port:%s ...", port)
	s := grpc.NewServer()
	pb.RegisterTokenServiceServer(s, &tkServer{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
