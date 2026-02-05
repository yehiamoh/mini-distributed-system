package service

import (
	"context"
	pb "service1/internal/gen/message"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)
type MessageService struct{
	pb.UnimplementedMessageServiceServer
}

func NewMessageService() *MessageService{
	return &MessageService{}
}

func(s *MessageService) AddNumber(ctx context.Context, req *pb.Request) (*pb.Response,error){
	if req.Number_1<0 || req.Number_2<0{
		return nil,status.Error(codes.InvalidArgument,"both numbers can't be zero")
	}
	return &pb.Response{Result: req.Number_1+req.Number_2},nil
}